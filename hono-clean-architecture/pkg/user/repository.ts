import { NodePgDatabase } from "drizzle-orm/node-postgres";
import { eq, sql, count } from "drizzle-orm";
import bcrypt from "bcryptjs";
import {
  users,
  User,
  UserResponse,
  Pagination,
  Search,
  Resources,
  toUserResponse,
} from "../models";
import { UserRepository as IUserRepository } from "../domain";
import { ResponseError, createError, whereAmI } from "../../internal/infrastructure/custom_error";
import { applySearch, getOffset, getLimit } from "../utils";
import { UserCache, newUserCache } from "./cache";

export class UserRepository implements IUserRepository {
  private db: NodePgDatabase;
  private cache: UserCache;

  constructor(resources: Resources) {
    this.db = resources.mainDbConn;
    this.cache = newUserCache(resources.redisStorage);
  }

  async migrate(): Promise<void> {
    console.log("User table migration handled by Drizzle Kit");
  }

  async createUser(userData: {
    username: string;
    email: string;
    password: string;
  }): Promise<{ user: UserResponse | null; error: ResponseError | null }> {
    const existingEmail = await this.db
      .select()
      .from(users)
      .where(eq(users.email, userData.email))
      .limit(1);

    if (existingEmail.length > 0) {
      return {
        user: null,
        error: createError(409, whereAmI(), "Conflict", "User with this email already exists."),
      };
    }

    const existingUsername = await this.db
      .select()
      .from(users)
      .where(eq(users.username, userData.username))
      .limit(1);

    if (existingUsername.length > 0) {
      return {
        user: null,
        error: createError(409, whereAmI(), "Conflict", "User with this username already exists."),
      };
    }

    const hashedPassword = await bcrypt.hash(userData.password, 10);

    try {
      const result = await this.db
        .insert(users)
        .values({
          username: userData.username,
          email: userData.email,
          password: hashedPassword,
        })
        .returning();

      const user = result[0];
      await this.cache.invalidateAllLists();

      return { user: toUserResponse(user), error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        user: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async getUser(id: number): Promise<{ user: UserResponse | null; error: ResponseError | null }> {
    const cachedUser = await this.cache.get(id);
    if (cachedUser) {
      return { user: cachedUser, error: null };
    }

    try {
      const result = await this.db.select().from(users).where(eq(users.id, id)).limit(1);

      if (result.length === 0) {
        return {
          user: null,
          error: createError(404, whereAmI(), "Not Found", "User not found"),
        };
      }

      const userResponse = toUserResponse(result[0]);
      await this.cache.set(userResponse);

      return { user: userResponse, error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        user: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async getUsers(
    pagination: Pagination,
    search: Search
  ): Promise<{
    users: UserResponse[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }> {
    const cached = await this.cache.getList(pagination, search);
    if (cached) {
      return { ...cached, error: null };
    }

    try {
      const searchCondition = applySearch(search, [users.username, users.email]);

      const countQuery = searchCondition
        ? this.db.select({ count: count() }).from(users).where(searchCondition)
        : this.db.select({ count: count() }).from(users);

      const countResult = await countQuery;
      const total = Number(countResult[0]?.count ?? 0);

      let query = this.db.select().from(users);

      if (searchCondition) {
        query = query.where(searchCondition) as typeof query;
      }

      if (pagination.orderBy) {
        const orderDirection = pagination.order === "desc" ? sql`DESC` : sql`ASC`;
        query = query.orderBy(
          sql`${sql.identifier(pagination.orderBy)} ${orderDirection}`
        ) as typeof query;
      }

      const result = await query.limit(getLimit(pagination)).offset(getOffset(pagination));

      const userResponses = result.map(toUserResponse);
      const updatedPagination = { ...pagination, total };

      await this.cache.setList(userResponses, updatedPagination, search);

      return {
        users: userResponses,
        pagination: updatedPagination,
        search,
        error: null,
      };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        users: [],
        pagination,
        search,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async updateUser(id: number, userData: Partial<User>): Promise<ResponseError | null> {
    try {
      const currentUser = await this.db.select().from(users).where(eq(users.id, id)).limit(1);

      if (currentUser.length === 0) {
        return createError(404, whereAmI(), "Not Found", "User not found");
      }

      const updateData: Partial<User> = { ...userData };

      if (userData.password) {
        updateData.password = await bcrypt.hash(userData.password, 10);
      }

      await this.db
        .update(users)
        .set({ ...updateData, updatedAt: new Date() })
        .where(eq(users.id, id));

      await this.cache.invalidate(id, currentUser[0].email);
      await this.cache.invalidateAllLists();

      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async deleteUser(id: number): Promise<ResponseError | null> {
    try {
      const user = await this.db.select().from(users).where(eq(users.id, id)).limit(1);

      if (user.length === 0) {
        return createError(404, whereAmI(), "Not Found", "User not found");
      }

      await this.db.delete(users).where(eq(users.id, id));

      await this.cache.invalidate(id, user[0].email);
      await this.cache.invalidateAllLists();

      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async getUserByEmail(email: string): Promise<{ user: User | null; error: ResponseError | null }> {
    try {
      const result = await this.db.select().from(users).where(eq(users.email, email)).limit(1);

      if (result.length === 0) {
        return {
          user: null,
          error: createError(404, whereAmI(), "Not Found", "User not found"),
        };
      }

      return { user: result[0], error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        user: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }
}

export function newUserRepository(resources: Resources): IUserRepository {
  return new UserRepository(resources);
}
