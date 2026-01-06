import { ResponseError } from "./../../internal/infrastructure/custom_error";
import { User, UserResponse, Pagination, Search } from "../models";
import {
  UserRepository,
  UserService as IUserService,
  RoleRepository,
  UserRoleRepository,
} from "../domain";
import { createError, whereAmI } from "../../internal/infrastructure/custom_error";
import { isValidEmail, RoleType } from "../utils";
import { TxManager } from "../tx_context";

export class UserService implements IUserService {
  private repository: UserRepository;
  private roleRepository: RoleRepository;
  private userRoleRepository: UserRoleRepository;
  private txManager: TxManager;

  constructor(
    repository: UserRepository,
    roleRepository: RoleRepository,
    userRoleRepository: UserRoleRepository,
    txManager: TxManager
  ) {
    this.repository = repository;
    this.roleRepository = roleRepository;
    this.userRoleRepository = userRoleRepository;
    this.txManager = txManager;
  }

  async createUser(userData: {
    username: string;
    email: string;
    password: string;
  }): Promise<ResponseError[]> {
    if (!isValidEmail(userData.email) || !userData.password) {
      return [
        createError(400, whereAmI(), "Bad Request", "Invalid email format or password is empty"),
      ];
    }

    const { errors } = await this.txManager.withTransaction(async () => {
      const { user, error } = await this.repository.createUser(userData);
      if (error) {
        return {
          result: null,
          errors: [
            error,
            createError(
              500,
              whereAmI(),
              "Internal Server Error",
              "Something went wrong with createUser repository"
            ),
          ],
        };
      }

      if (!user) {
        return {
          result: null,
          errors: [createError(500, whereAmI(), "Internal Error", "Failed to create user")],
        };
      }

      const { role, error: roleError } = await this.roleRepository.getRoleByName(RoleType.User);
      if (roleError) {
        return {
          result: null,
          errors: [
            roleError,
            createError(
              500,
              whereAmI(),
              "Internal Server Error",
              "Something went wrong with getRoleByName repository"
            ),
          ],
        };
      }

      if (!role) {
        return {
          result: null,
          errors: [createError(404, whereAmI(), "Not Found", "Default role not found")],
        };
      }

      const userRoleError = await this.userRoleRepository.createUserRole({
        userId: user.id,
        roleId: role.id,
      });

      if (userRoleError) {
        return {
          result: null,
          errors: [
            userRoleError,
            createError(
              500,
              whereAmI(),
              "Internal Server Error",
              "Something went wrong with createUserRole repository"
            ),
          ],
        };
      }

      return { result: user, errors: [] };
    });

    return errors;
  }

  async getUser(id: number): Promise<{ user: UserResponse | null; errors: ResponseError[] }> {
    const { user, error } = await this.repository.getUser(id);
    if (error) {
      return {
        user: null,
        errors: [
          error,
          createError(
            500,
            whereAmI(),
            "Internal Server Error",
            "Something went wrong with getUser repository"
          ),
        ],
      };
    }
    return { user, errors: [] };
  }

  async getUsers(
    pagination: Pagination,
    search: Search
  ): Promise<{
    users: UserResponse[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }> {
    const result = await this.repository.getUsers(pagination, search);
    if (result.error) {
      return {
        users: [],
        pagination,
        search,
        errors: [
          result.error,
          createError(
            500,
            whereAmI(),
            "Internal Server Error",
            "Something went wrong with getUsers repository"
          ),
        ],
      };
    }
    return {
      users: result.users,
      pagination: result.pagination,
      search: result.search,
      errors: [],
    };
  }

  async updateUser(id: number, userData: Partial<User>): Promise<ResponseError[]> {
    if (userData.email && !isValidEmail(userData.email)) {
      return [createError(400, whereAmI(), "Bad Request", "Invalid email format")];
    }

    const error = await this.repository.updateUser(id, userData);
    if (error) {
      return [
        error,
        createError(
          500,
          whereAmI(),
          "Internal Server Error",
          "Something went wrong with updateUser repository"
        ),
      ];
    }
    return [];
  }

  async deleteUser(id: number): Promise<ResponseError[]> {
    const error = await this.repository.deleteUser(id);
    if (error) {
      return [
        error,
        createError(
          500,
          whereAmI(),
          "Internal Server Error",
          "Something went wrong with deleteUser repository"
        ),
      ];
    }
    return [];
  }

  async getUserByEmail(email: string): Promise<{ user: User | null; errors: ResponseError[] }> {
    if (!isValidEmail(email)) {
      return {
        user: null,
        errors: [createError(400, whereAmI(), "Bad Request", "Invalid email format")],
      };
    }

    const { user, error } = await this.repository.getUserByEmail(email);
    if (error) {
      return {
        user: null,
        errors: [
          error,
          createError(
            500,
            whereAmI(),
            "Internal Server Error",
            "Something went wrong with getUserByEmail repository"
          ),
        ],
      };
    }
    return { user, errors: [] };
  }
}

export function newUserService(
  repository: UserRepository,
  roleRepository: RoleRepository,
  userRoleRepository: UserRoleRepository,
  txManager: TxManager
): IUserService {
  return new UserService(repository, roleRepository, userRoleRepository, txManager);
}
