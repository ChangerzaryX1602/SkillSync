import { NodePgDatabase } from "drizzle-orm/node-postgres";
import { eq, and, sql, count } from "drizzle-orm";
import { permissions, Permission, Pagination, Search, Resources } from "../models";
import { PermissionRepository as IPermissionRepository } from "../domain";
import { ResponseError, createError, whereAmI } from "../../internal/infrastructure/custom_error";
import { applySearch, getOffset, getLimit } from "../utils";
import { PermissionCache, newPermissionCache } from "./cache";

export class PermissionRepository implements IPermissionRepository {
  private db: NodePgDatabase;
  private cache: PermissionCache;

  constructor(resources: Resources) {
    this.db = resources.mainDbConn;
    this.cache = newPermissionCache(resources.redisStorage);
  }

  async migrate(): Promise<void> {
    console.log("Permission table migration handled by Drizzle Kit");
  }

  async createPermission(permissionData: {
    group: string;
    name: string;
  }): Promise<{ permission: Permission | null; error: ResponseError | null }> {
    const existing = await this.db
      .select()
      .from(permissions)
      .where(
        and(eq(permissions.group, permissionData.group), eq(permissions.name, permissionData.name))
      )
      .limit(1);

    if (existing.length > 0) {
      return {
        permission: null,
        error: createError(409, whereAmI(), "Conflict", "Permission already exists."),
      };
    }

    try {
      const result = await this.db
        .insert(permissions)
        .values({ group: permissionData.group, name: permissionData.name })
        .returning();

      await this.cache.invalidateAllLists();
      return { permission: result[0], error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        permission: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async getPermission(
    id: number
  ): Promise<{ permission: Permission | null; error: ResponseError | null }> {
    const cached = await this.cache.get(id);
    if (cached) {
      return { permission: cached, error: null };
    }

    try {
      const result = await this.db
        .select()
        .from(permissions)
        .where(eq(permissions.id, id))
        .limit(1);

      if (result.length === 0) {
        return {
          permission: null,
          error: createError(404, whereAmI(), "Not Found", "Permission not found"),
        };
      }

      await this.cache.set(result[0]);
      return { permission: result[0], error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        permission: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async getPermissions(
    pagination: Pagination,
    search: Search
  ): Promise<{
    permissions: Permission[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }> {
    try {
      const searchCondition = applySearch(search, [permissions.group, permissions.name]);

      const countQuery = searchCondition
        ? this.db.select({ count: count() }).from(permissions).where(searchCondition)
        : this.db.select({ count: count() }).from(permissions);

      const countResult = await countQuery;
      const total = Number(countResult[0]?.count ?? 0);

      let query = this.db.select().from(permissions);

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
        permissions: result,
        pagination: { ...pagination, total },
        search,
        error: null,
      };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        permissions: [],
        pagination,
        search,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async updatePermission(
    id: number,
    permissionData: Partial<Permission>
  ): Promise<ResponseError | null> {
    try {
      const current = await this.db
        .select()
        .from(permissions)
        .where(eq(permissions.id, id))
        .limit(1);

      if (current.length === 0) {
        return createError(404, whereAmI(), "Not Found", "Permission not found");
      }

      await this.db
        .update(permissions)
        .set({ ...permissionData, updatedAt: new Date() })
        .where(eq(permissions.id, id));

      await this.cache.invalidate(id);
      await this.cache.invalidateAllLists();

      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async deletePermission(id: number): Promise<ResponseError | null> {
    try {
      const permission = await this.db
        .select()
        .from(permissions)
        .where(eq(permissions.id, id))
        .limit(1);

      if (permission.length === 0) {
        return createError(404, whereAmI(), "Not Found", "Permission not found");
      }

      await this.db.delete(permissions).where(eq(permissions.id, id));
      await this.cache.invalidate(id);
      await this.cache.invalidateAllLists();

      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async getPermissionByGroupAndName(
    group: string,
    name: string
  ): Promise<{ permission: Permission | null; error: ResponseError | null }> {
    try {
      const result = await this.db
        .select()
        .from(permissions)
        .where(and(eq(permissions.group, group), eq(permissions.name, name)))
        .limit(1);

      if (result.length === 0) {
        return {
          permission: null,
          error: createError(404, whereAmI(), "Not Found", "Permission not found"),
        };
      }

      return { permission: result[0], error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        permission: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }
}

export function newPermissionRepository(resources: Resources): IPermissionRepository {
  return new PermissionRepository(resources);
}
