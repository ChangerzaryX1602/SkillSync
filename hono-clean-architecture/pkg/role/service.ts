import { Role, Pagination, Search } from "../models";
import { RoleRepository, RoleService as IRoleService } from "../domain";
import { ResponseError } from "../../internal/infrastructure/custom_error";

export class RoleService implements IRoleService {
  private repository: RoleRepository;

  constructor(repository: RoleRepository) {
    this.repository = repository;
  }

  async createRole(roleData: { name: string }): Promise<ResponseError[]> {
    const { error } = await this.repository.createRole(roleData);
    if (error) {
      return [error];
    }
    return [];
  }

  async getRole(id: number): Promise<{ role: Role | null; errors: ResponseError[] }> {
    const { role, error } = await this.repository.getRole(id);
    if (error) {
      return { role: null, errors: [error] };
    }
    return { role, errors: [] };
  }

  async getRoles(
    pagination: Pagination,
    search: Search
  ): Promise<{
    roles: Role[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }> {
    const result = await this.repository.getRoles(pagination, search);
    if (result.error) {
      return { roles: [], pagination, search, errors: [result.error] };
    }
    return {
      roles: result.roles,
      pagination: result.pagination,
      search: result.search,
      errors: [],
    };
  }

  async updateRole(id: number, roleData: Partial<Role>): Promise<ResponseError[]> {
    const error = await this.repository.updateRole(id, roleData);
    if (error) {
      return [error];
    }
    return [];
  }

  async deleteRole(id: number): Promise<ResponseError[]> {
    const error = await this.repository.deleteRole(id);
    if (error) {
      return [error];
    }
    return [];
  }

  async getRoleByName(name: string): Promise<{ role: Role | null; errors: ResponseError[] }> {
    const { role, error } = await this.repository.getRoleByName(name);
    if (error) {
      return { role: null, errors: [error] };
    }
    return { role, errors: [] };
  }
}

export function newRoleService(repository: RoleRepository): IRoleService {
  return new RoleService(repository);
}
