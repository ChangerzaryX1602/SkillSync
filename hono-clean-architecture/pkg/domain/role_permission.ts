import { RolePermission, Pagination, Search } from "../models";
import { ResponseError } from "../../internal/infrastructure/custom_error";

export interface RolePermissionRepository {
  migrate(): Promise<void>;
  createRolePermission(rolePermission: {
    roleId: number;
    permissionId: number;
  }): Promise<ResponseError | null>;
  getRolePermission(
    id: number
  ): Promise<{ rolePermission: RolePermission | null; error: ResponseError | null }>;
  getRolePermissions(
    pagination: Pagination,
    search: Search
  ): Promise<{
    rolePermissions: RolePermission[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }>;
  deleteRolePermission(id: number): Promise<ResponseError | null>;
  getRolePermissionsByRoleId(
    roleId: number
  ): Promise<{ rolePermissions: RolePermission[]; error: ResponseError | null }>;
}

export interface RolePermissionService {
  createRolePermission(rolePermission: {
    roleId: number;
    permissionId: number;
  }): Promise<ResponseError[]>;
  getRolePermission(
    id: number
  ): Promise<{ rolePermission: RolePermission | null; errors: ResponseError[] }>;
  getRolePermissions(
    pagination: Pagination,
    search: Search
  ): Promise<{
    rolePermissions: RolePermission[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }>;
  deleteRolePermission(id: number): Promise<ResponseError[]>;
}
