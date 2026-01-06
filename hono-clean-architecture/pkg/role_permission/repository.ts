import { NodePgDatabase } from "drizzle-orm/node-postgres";
import { eq, and, sql, count } from "drizzle-orm";
import { rolePermissions, RolePermission, Pagination, Search, Resources } from "../models";
import { RolePermissionRepository as IRolePermissionRepository } from "../domain";
import { ResponseError, createError, whereAmI } from "../../internal/infrastructure/custom_error";
import { getOffset, getLimit } from "../utils";
import { RolePermissionCache, newRolePermissionCache } from "./cache";

export class RolePermissionRepository implements IRolePermissionRepository {
  private db: NodePgDatabase;
  private cache: RolePermissionCache;

  constructor(resources: Resources) {
    this.db = resources.mainDbConn;
    this.cache = newRolePermissionCache(resources.redisStorage);
  }

  async migrate(): Promise<void> {
    console.log("RolePermission table migration handled by Drizzle Kit");
  }

  async createRolePermission(data: {
    roleId: number;
    permissionId: number;
  }): Promise<ResponseError | null> {
    const existing = await this.db
      .select()
      .from(rolePermissions)
      .where(
        and(
          eq(rolePermissions.roleId, data.roleId),
          eq(rolePermissions.permissionId, data.permissionId)
        )
      )
      .limit(1);

    if (existing.length > 0) {
      return createError(409, whereAmI(), "Conflict", "RolePermission already exists.");
    }

    try {
      await this.db.insert(rolePermissions).values({
        roleId: data.roleId,
        permissionId: data.permissionId,
      });

      await this.cache.invalidateAllLists();
      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async getRolePermission(
    id: number
  ): Promise<{ rolePermission: RolePermission | null; error: ResponseError | null }> {
    const cached = await this.cache.get(id);
    if (cached) {
      return { rolePermission: cached, error: null };
    }

    try {
      const result = await this.db
        .select()
        .from(rolePermissions)
        .where(eq(rolePermissions.id, id))
        .limit(1);

      if (result.length === 0) {
        return {
          rolePermission: null,
          error: createError(404, whereAmI(), "Not Found", "RolePermission not found"),
        };
      }

      await this.cache.set(result[0]);
      return { rolePermission: result[0], error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        rolePermission: null,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async getRolePermissions(
    pagination: Pagination,
    search: Search
  ): Promise<{
    rolePermissions: RolePermission[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }> {
    try {
      const countResult = await this.db.select({ count: count() }).from(rolePermissions);
      const total = Number(countResult[0]?.count ?? 0);

      let query = this.db.select().from(rolePermissions);

      if (pagination.orderBy) {
        const orderDirection = pagination.order === "desc" ? sql`DESC` : sql`ASC`;
        query = query.orderBy(
          sql`${sql.identifier(pagination.orderBy)} ${orderDirection}`
        ) as typeof query;
      }

      const result = await query.limit(getLimit(pagination)).offset(getOffset(pagination));

      return {
        rolePermissions: result,
        pagination: { ...pagination, total },
        search,
        error: null,
      };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        rolePermissions: [],
        pagination,
        search,
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }

  async deleteRolePermission(id: number): Promise<ResponseError | null> {
    try {
      const rolePermission = await this.db
        .select()
        .from(rolePermissions)
        .where(eq(rolePermissions.id, id))
        .limit(1);

      if (rolePermission.length === 0) {
        return createError(404, whereAmI(), "Not Found", "RolePermission not found");
      }

      await this.db.delete(rolePermissions).where(eq(rolePermissions.id, id));
      await this.cache.invalidate(id);
      await this.cache.invalidateAllLists();

      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return createError(500, whereAmI(), "Database Error", message);
    }
  }

  async getRolePermissionsByRoleId(
    roleId: number
  ): Promise<{ rolePermissions: RolePermission[]; error: ResponseError | null }> {
    try {
      const result = await this.db
        .select()
        .from(rolePermissions)
        .where(eq(rolePermissions.roleId, roleId));

      return { rolePermissions: result, error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Database error";
      return {
        rolePermissions: [],
        error: createError(500, whereAmI(), "Database Error", message),
      };
    }
  }
}

export function newRolePermissionRepository(resources: Resources): IRolePermissionRepository {
  return new RolePermissionRepository(resources);
}
