import { RolePermission, Pagination, Search } from "../models";
import {
  RolePermissionRepository,
  RolePermissionService as IRolePermissionService,
} from "../domain";
import { ResponseError } from "../../internal/infrastructure/custom_error";

export class RolePermissionService implements IRolePermissionService {
  private repository: RolePermissionRepository;

  constructor(repository: RolePermissionRepository) {
    this.repository = repository;
  }

  async createRolePermission(data: {
    roleId: number;
    permissionId: number;
  }): Promise<ResponseError[]> {
    const error = await this.repository.createRolePermission(data);
    if (error) {
      return [error];
    }
    return [];
  }

  async getRolePermission(
    id: number
  ): Promise<{ rolePermission: RolePermission | null; errors: ResponseError[] }> {
    const { rolePermission, error } = await this.repository.getRolePermission(id);
    if (error) {
      return { rolePermission: null, errors: [error] };
    }
    return { rolePermission, errors: [] };
  }

  async getRolePermissions(
    pagination: Pagination,
    search: Search
  ): Promise<{
    rolePermissions: RolePermission[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }> {
    const result = await this.repository.getRolePermissions(pagination, search);
    if (result.error) {
      return {
        rolePermissions: [],
        pagination,
        search,
        errors: [result.error],
      };
    }
    return {
      rolePermissions: result.rolePermissions,
      pagination: result.pagination,
      search: result.search,
      errors: [],
    };
  }

  async deleteRolePermission(id: number): Promise<ResponseError[]> {
    const error = await this.repository.deleteRolePermission(id);
    if (error) {
      return [error];
    }
    return [];
  }
}

export function newRolePermissionService(
  repository: RolePermissionRepository
): IRolePermissionService {
  return new RolePermissionService(repository);
}
