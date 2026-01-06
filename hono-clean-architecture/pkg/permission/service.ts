import { Permission, Pagination, Search } from "../models";
import { PermissionRepository, PermissionService as IPermissionService } from "../domain";
import { ResponseError } from "../../internal/infrastructure/custom_error";

export class PermissionService implements IPermissionService {
  private repository: PermissionRepository;

  constructor(repository: PermissionRepository) {
    this.repository = repository;
  }

  async createPermission(permissionData: {
    group: string;
    name: string;
  }): Promise<ResponseError[]> {
    const { error } = await this.repository.createPermission(permissionData);
    if (error) {
      return [error];
    }
    return [];
  }

  async getPermission(
    id: number
  ): Promise<{ permission: Permission | null; errors: ResponseError[] }> {
    const { permission, error } = await this.repository.getPermission(id);
    if (error) {
      return { permission: null, errors: [error] };
    }
    return { permission, errors: [] };
  }

  async getPermissions(
    pagination: Pagination,
    search: Search
  ): Promise<{
    permissions: Permission[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }> {
    const result = await this.repository.getPermissions(pagination, search);
    if (result.error) {
      return { permissions: [], pagination, search, errors: [result.error] };
    }
    return {
      permissions: result.permissions,
      pagination: result.pagination,
      search: result.search,
      errors: [],
    };
  }

  async updatePermission(
    id: number,
    permissionData: Partial<Permission>
  ): Promise<ResponseError[]> {
    const error = await this.repository.updatePermission(id, permissionData);
    if (error) {
      return [error];
    }
    return [];
  }

  async deletePermission(id: number): Promise<ResponseError[]> {
    const error = await this.repository.deletePermission(id);
    if (error) {
      return [error];
    }
    return [];
  }
}

export function newPermissionService(repository: PermissionRepository): IPermissionService {
  return new PermissionService(repository);
}
