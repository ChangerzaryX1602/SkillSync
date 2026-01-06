import { NodePgDatabase } from "drizzle-orm/node-postgres";
import { eq, sql, count } from "drizzle-orm";
import { roles, Role, Pagination, Search, Resources } from "../models";
import { RoleRepository as IRoleRepository } from "../domain";
import { ResponseError, createError, whereAmI } from "../../internal/infrastructure/custom_error";
import { applySearch, getOffset, getLimit } from "../utils";
import { RoleCache, newRoleCache } from "./cache";

export class RoleRepository implements IRoleRepository {
  private db: NodePgDatabase;
  private cache: RoleCache;

  constructor(resources: Resources) {
    this.db = resources.mainDbConn;
    this.cache = newRoleCache(resources.redisStorage);
  }

  async migrate(): Promise<void> {
    console.log("Role table migration handled by Drizzle Kit");
  }

  async createRole(roleData: {
    name: string;
  }): Promise<{ role: Role | null; error: ResponseError | null }> {
    const existing = await this.db
      .select()
      .from(roles)
      .where(eq(roles.name, roleData.name))
      .limit(1);

    if (existing.length > 0) {
      return {
        role: null,
        error: createError(409, whereAmI(), "Conflict", "Role with this name already exists."),
      };
    }

    try {
      const result = await this.db.insert(roles).values({ name: roleData.name }).returning();

      await this.cache.invalidateAllLists();
      return { role: result[0], error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        role: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async getRole(id: number): Promise<{ role: Role | null; error: ResponseError | null }> {
    const cachedRole = await this.cache.get(id);
    if (cachedRole) {
      return { role: cachedRole, error: null };
    }

    try {
      const result = await this.db.select().from(roles).where(eq(roles.id, id)).limit(1);

      if (result.length === 0) {
        return {
          role: null,
          error: createError(404, whereAmI(), "Not Found", "Role not found"),
        };
      }

      await this.cache.set(result[0]);
      return { role: result[0], error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        role: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async getRoles(
    pagination: Pagination,
    search: Search
  ): Promise<{
    roles: Role[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }> {
    try {
      const searchCondition = applySearch(search, [roles.name]);

      const countQuery = searchCondition
        ? this.db.select({ count: count() }).from(roles).where(searchCondition)
        : this.db.select({ count: count() }).from(roles);

      const countResult = await countQuery;
      const total = Number(countResult[0]?.count ?? 0);

      let query = this.db.select().from(roles);

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

      return {
        roles: result,
        pagination: { ...pagination, total },
        search,
        error: null,
      };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        roles: [],
        pagination,
        search,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async updateRole(id: number, roleData: Partial<Role>): Promise<ResponseError | null> {
    try {
      const currentRole = await this.db.select().from(roles).where(eq(roles.id, id)).limit(1);

      if (currentRole.length === 0) {
        return createError(404, whereAmI(), "Not Found", "Role not found");
      }

      await this.db
        .update(roles)
        .set({ ...roleData, updatedAt: new Date() })
        .where(eq(roles.id, id));

      await this.cache.invalidate(id, currentRole[0].name);
      await this.cache.invalidateAllLists();

      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async deleteRole(id: number): Promise<ResponseError | null> {
    try {
      const role = await this.db.select().from(roles).where(eq(roles.id, id)).limit(1);

      if (role.length === 0) {
        return createError(404, whereAmI(), "Not Found", "Role not found");
      }

      await this.db.delete(roles).where(eq(roles.id, id));
      await this.cache.invalidate(id, role[0].name);
      await this.cache.invalidateAllLists();

      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async getRoleByName(name: string): Promise<{ role: Role | null; error: ResponseError | null }> {
    const cachedRole = await this.cache.getByName(name);
    if (cachedRole) {
      return { role: cachedRole, error: null };
    }

    try {
      const result = await this.db.select().from(roles).where(eq(roles.name, name)).limit(1);

      if (result.length === 0) {
        return {
          role: null,
          error: createError(404, whereAmI(), "Not Found", "Role not found"),
        };
      }

      await this.cache.setByName(name, result[0]);
      return { role: result[0], error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        role: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }
}

export function newRoleRepository(resources: Resources): IRoleRepository {
  return new RoleRepository(resources);
}
