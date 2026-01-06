import { Permission, Pagination, Search } from "../models";
import { ResponseError } from "../../internal/infrastructure/custom_error";

export interface PermissionRepository {
  migrate(): Promise<void>;
  createPermission(permission: {
    group: string;
    name: string;
  }): Promise<{ permission: Permission | null; error: ResponseError | null }>;
  getPermission(
    id: number
  ): Promise<{ permission: Permission | null; error: ResponseError | null }>;
  getPermissions(
    pagination: Pagination,
    search: Search
  ): Promise<{
    permissions: Permission[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }>;
  updatePermission(id: number, permission: Partial<Permission>): Promise<ResponseError | null>;
  deletePermission(id: number): Promise<ResponseError | null>;
  getPermissionByGroupAndName(
    group: string,
    name: string
  ): Promise<{ permission: Permission | null; error: ResponseError | null }>;
}

export interface PermissionService {
  createPermission(permission: { group: string; name: string }): Promise<ResponseError[]>;
  getPermission(id: number): Promise<{ permission: Permission | null; errors: ResponseError[] }>;
  getPermissions(
    pagination: Pagination,
    search: Search
  ): Promise<{
    permissions: Permission[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }>;
  updatePermission(id: number, permission: Partial<Permission>): Promise<ResponseError[]>;
  deletePermission(id: number): Promise<ResponseError[]>;
}
