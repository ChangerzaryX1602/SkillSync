import { NodePgDatabase } from "drizzle-orm/node-postgres";
import { eq, and, sql, count } from "drizzle-orm";
import { userRoles, UserRole, Pagination, Search, Resources } from "../models";
import { UserRoleRepository as IUserRoleRepository } from "../domain";
import { ResponseError, createError, whereAmI } from "../../internal/infrastructure/custom_error";
import { getOffset, getLimit } from "../utils";
import { UserRoleCache, newUserRoleCache } from "./cache";

export class UserRoleRepository implements IUserRoleRepository {
  private db: NodePgDatabase;
  private cache: UserRoleCache;

  constructor(resources: Resources) {
    this.db = resources.mainDbConn;
    this.cache = newUserRoleCache(resources.redisStorage);
  }

  async migrate(): Promise<void> {
    console.log("UserRole table migration handled by Drizzle Kit");
  }

  async createUserRole(data: { userId: number; roleId: number }): Promise<ResponseError | null> {
    const existing = await this.db
      .select()
      .from(userRoles)
      .where(and(eq(userRoles.userId, data.userId), eq(userRoles.roleId, data.roleId)))
      .limit(1);

    if (existing.length > 0) {
      return createError(409, whereAmI(), "Conflict", "UserRole already exists.");
    }

    try {
      await this.db.insert(userRoles).values({
        userId: data.userId,
        roleId: data.roleId,
      });

      await this.cache.invalidateAllLists();
      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async getUserRole(
    id: number
  ): Promise<{ userRole: UserRole | null; error: ResponseError | null }> {
    const cached = await this.cache.get(id);
    if (cached) {
      return { userRole: cached, error: null };
    }

    try {
      const result = await this.db.select().from(userRoles).where(eq(userRoles.id, id)).limit(1);

      if (result.length === 0) {
        return {
          userRole: null,
          error: createError(404, whereAmI(), "Not Found", "UserRole not found"),
        };
      }

      await this.cache.set(result[0]);
      return { userRole: result[0], error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        userRole: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async getUserRoles(
    pagination: Pagination,
    search: Search
  ): Promise<{
    userRoles: UserRole[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }> {
    try {
      const countResult = await this.db.select({ count: count() }).from(userRoles);
      const total = Number(countResult[0]?.count ?? 0);

      let query = this.db.select().from(userRoles);

      if (pagination.orderBy) {
        const orderDirection = pagination.order === "desc" ? sql`DESC` : sql`ASC`;
        query = query.orderBy(
          sql`${sql.identifier(pagination.orderBy)} ${orderDirection}`
        ) as typeof query;
      }

      const result = await query.limit(getLimit(pagination)).offset(getOffset(pagination));

      return {
        userRoles: result,
        pagination: { ...pagination, total },
        search,
        error: null,
      };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        userRoles: [],
        pagination,
        search,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async deleteUserRole(id: number): Promise<ResponseError | null> {
    try {
      const userRole = await this.db.select().from(userRoles).where(eq(userRoles.id, id)).limit(1);

      if (userRole.length === 0) {
        return createError(404, whereAmI(), "Not Found", "UserRole not found");
      }

      await this.db.delete(userRoles).where(eq(userRoles.id, id));
      await this.cache.invalidate(id);
      await this.cache.invalidateAllLists();

      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async getUserRolesByUserId(
    userId: number
  ): Promise<{ userRoles: UserRole[]; error: ResponseError | null }> {
    try {
      const result = await this.db.select().from(userRoles).where(eq(userRoles.userId, userId));

      return { userRoles: result, error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        userRoles: [],
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }
}

export function newUserRoleRepository(resources: Resources): IUserRoleRepository {
  return new UserRoleRepository(resources);
}
