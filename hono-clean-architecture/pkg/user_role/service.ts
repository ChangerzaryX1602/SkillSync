import { UserRole, Pagination, Search } from "../models";
import { UserRoleRepository, UserRoleService as IUserRoleService } from "../domain";
import { ResponseError } from "../../internal/infrastructure/custom_error";

export class UserRoleService implements IUserRoleService {
  private repository: UserRoleRepository;

  constructor(repository: UserRoleRepository) {
    this.repository = repository;
  }

  async createUserRole(data: { userId: number; roleId: number }): Promise<ResponseError[]> {
    const error = await this.repository.createUserRole(data);
    if (error) {
      return [error];
    }
    return [];
  }

  async getUserRole(id: number): Promise<{ userRole: UserRole | null; errors: ResponseError[] }> {
    const { userRole, error } = await this.repository.getUserRole(id);
    if (error) {
      return { userRole: null, errors: [error] };
    }
    return { userRole, errors: [] };
  }

  async getUserRoles(
    pagination: Pagination,
    search: Search
  ): Promise<{
    userRoles: UserRole[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }> {
    const result = await this.repository.getUserRoles(pagination, search);
    if (result.error) {
      return { userRoles: [], pagination, search, errors: [result.error] };
    }
    return {
      userRoles: result.userRoles,
      pagination: result.pagination,
      search: result.search,
      errors: [],
    };
  }

  async deleteUserRole(id: number): Promise<ResponseError[]> {
    const error = await this.repository.deleteUserRole(id);
    if (error) {
      return [error];
    }
    return [];
  }
}

export function newUserRoleService(repository: UserRoleRepository): IUserRoleService {
  return new UserRoleService(repository);
}
