import { Role, Pagination, Search } from "../models";
import { ResponseError } from "../../internal/infrastructure/custom_error";

export interface RoleRepository {
  migrate(): Promise<void>;
  createRole(role: { name: string }): Promise<{ role: Role | null; error: ResponseError | null }>;
  getRole(id: number): Promise<{ role: Role | null; error: ResponseError | null }>;
  getRoles(
    pagination: Pagination,
    search: Search
  ): Promise<{
    roles: Role[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }>;
  updateRole(id: number, role: Partial<Role>): Promise<ResponseError | null>;
  deleteRole(id: number): Promise<ResponseError | null>;
  getRoleByName(name: string): Promise<{ role: Role | null; error: ResponseError | null }>;
}

export interface RoleService {
  createRole(role: { name: string }): Promise<ResponseError[]>;
  getRole(id: number): Promise<{ role: Role | null; errors: ResponseError[] }>;
  getRoles(
    pagination: Pagination,
    search: Search
  ): Promise<{
    roles: Role[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }>;
  updateRole(id: number, role: Partial<Role>): Promise<ResponseError[]>;
  deleteRole(id: number): Promise<ResponseError[]>;
  getRoleByName(name: string): Promise<{ role: Role | null; errors: ResponseError[] }>;
}
