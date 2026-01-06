import { UserRole, Pagination, Search } from "../models";
import { ResponseError } from "../../internal/infrastructure/custom_error";

export interface UserRoleRepository {
  migrate(): Promise<void>;
  createUserRole(userRole: { userId: number; roleId: number }): Promise<ResponseError | null>;
  getUserRole(id: number): Promise<{ userRole: UserRole | null; error: ResponseError | null }>;
  getUserRoles(
    pagination: Pagination,
    search: Search
  ): Promise<{
    userRoles: UserRole[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }>;
  deleteUserRole(id: number): Promise<ResponseError | null>;
  getUserRolesByUserId(
    userId: number
  ): Promise<{ userRoles: UserRole[]; error: ResponseError | null }>;
}

export interface UserRoleService {
  createUserRole(userRole: { userId: number; roleId: number }): Promise<ResponseError[]>;
  getUserRole(id: number): Promise<{ userRole: UserRole | null; errors: ResponseError[] }>;
  getUserRoles(
    pagination: Pagination,
    search: Search
  ): Promise<{
    userRoles: UserRole[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }>;
  deleteUserRole(id: number): Promise<ResponseError[]>;
}
