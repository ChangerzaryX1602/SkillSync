import bcrypt from "bcryptjs";
import { UserResponse } from "../models";
import {
  AuthRepository,
  AuthService as IAuthService,
  UserService,
  UserRoleRepository,
  RolePermissionRepository,
  RoleRepository,
  PermissionRepository,
} from "../domain";
import { ResponseError, createError, whereAmI } from "../../internal/infrastructure/custom_error";
import { isValidEmail } from "../utils";

const ACCESS_TOKEN_TTL = 15 * 60;
const REFRESH_TOKEN_TTL = 7 * 24 * 60 * 60;

export class AuthService implements IAuthService {
  private repository: AuthRepository;
  private userService: UserService;
  private userRoleRepository: UserRoleRepository;
  private rolePermissionRepository: RolePermissionRepository;
  private roleRepository: RoleRepository;
  private permissionRepository: PermissionRepository;

  constructor(
    repository: AuthRepository,
    userService: UserService,
    userRoleRepository: UserRoleRepository,
    rolePermissionRepository: RolePermissionRepository,
    roleRepository: RoleRepository,
    permissionRepository: PermissionRepository
  ) {
    this.repository = repository;
    this.userService = userService;
    this.userRoleRepository = userRoleRepository;
    this.rolePermissionRepository = rolePermissionRepository;
    this.roleRepository = roleRepository;
    this.permissionRepository = permissionRepository;
  }

  private async getUserRolesAndPermissions(
    userId: number
  ): Promise<{ roles: string[]; permissions: string[]; errors: ResponseError[] }> {
    const { userRoles, error: userRolesError } =
      await this.userRoleRepository.getUserRolesByUserId(userId);

    if (userRolesError) {
      return { roles: [], permissions: [], errors: [userRolesError] };
    }

    const roles: string[] = [];
    const permissionSet = new Set<string>();

    for (const userRole of userRoles) {
      const { role, error: roleError } = await this.roleRepository.getRole(userRole.roleId);
      if (roleError || !role) {
        continue;
      }
      roles.push(role.name);

      const { rolePermissions, error: rpError } =
        await this.rolePermissionRepository.getRolePermissionsByRoleId(role.id);
      if (rpError) {
        continue;
      }

      for (const rp of rolePermissions) {
        const { permission, error: permError } = await this.permissionRepository.getPermission(
          rp.permissionId
        );
        if (!permError && permission) {
          permissionSet.add(`${permission.group}:${permission.name}`);
        }
      }
    }

    return { roles, permissions: Array.from(permissionSet), errors: [] };
  }

  async register(userData: {
    username: string;
    email: string;
    password: string;
  }): Promise<ResponseError[]> {
    return this.userService.createUser(userData);
  }

  async login(
    credentials: { email: string; password: string },
    host: string
  ): Promise<{
    accessToken: string | null;
    refreshToken: string | null;
    errors: ResponseError[];
  }> {
    if (!isValidEmail(credentials.email)) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [createError(400, whereAmI(), "Bad Request", "Invalid email format")],
      };
    }

    const { user, errors: userErrors } = await this.userService.getUserByEmail(credentials.email);
    if (userErrors.length > 0) {
      return { accessToken: null, refreshToken: null, errors: userErrors };
    }

    if (!user) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [createError(401, whereAmI(), "Unauthorized", "Invalid credentials")],
      };
    }

    const isPasswordValid = await bcrypt.compare(credentials.password, user.password);
    if (!isPasswordValid) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [createError(401, whereAmI(), "Unauthorized", "Invalid credentials")],
      };
    }

    const { roles, permissions, errors: rbacErrors } = await this.getUserRolesAndPermissions(
      user.id
    );
    if (rbacErrors.length > 0) {
      return { accessToken: null, refreshToken: null, errors: rbacErrors };
    }

    const { token: accessToken, error: accessTokenError } = await this.repository.signToken(
      user,
      host,
      ACCESS_TOKEN_TTL,
      roles,
      permissions
    );
    if (accessTokenError) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [accessTokenError],
      };
    }

    const { token: refreshToken, error: refreshTokenError } = await this.repository.signToken(
      user,
      host,
      REFRESH_TOKEN_TTL,
      roles,
      permissions
    );
    if (refreshTokenError) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [refreshTokenError],
      };
    }

    const saveError = await this.repository.saveRefreshToken(
      user.id,
      refreshToken!,
      REFRESH_TOKEN_TTL
    );
    if (saveError) {
      return { accessToken: null, refreshToken: null, errors: [saveError] };
    }

    return { accessToken, refreshToken, errors: [] };
  }

  async refreshToken(refreshTokenStr: string): Promise<{
    accessToken: string | null;
    refreshToken: string | null;
    errors: ResponseError[];
  }> {
    const { payload, valid, error: parseError } = await this.repository.parseToken(refreshTokenStr);

    if (parseError) {
      return { accessToken: null, refreshToken: null, errors: [parseError] };
    }

    if (!valid || !payload) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [createError(401, whereAmI(), "Unauthorized", "Invalid token")],
      };
    }

    const userId = parseInt(payload.sub || "0", 10);
    if (!userId) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [createError(401, whereAmI(), "Unauthorized", "Invalid user ID in token")],
      };
    }

    const { token: storedToken, error: getError } = await this.repository.getRefreshToken(userId);
    if (getError) {
      return { accessToken: null, refreshToken: null, errors: [getError] };
    }

    if (storedToken !== refreshTokenStr) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [
          createError(401, whereAmI(), "Unauthorized", "Refresh token mismatch (Reuse detected?)"),
        ],
      };
    }

    const { user, errors: userErrors } = await this.userService.getUser(userId);
    if (userErrors.length > 0) {
      return { accessToken: null, refreshToken: null, errors: userErrors };
    }

    if (!user) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [createError(404, whereAmI(), "Not Found", "User not found")],
      };
    }

    const host = payload.iss || "";

    const { roles, permissions, errors: rbacErrors } = await this.getUserRolesAndPermissions(
      user.id
    );
    if (rbacErrors.length > 0) {
      return { accessToken: null, refreshToken: null, errors: rbacErrors };
    }

    const { token: newAccessToken, error: accessTokenError } = await this.repository.signToken(
      user,
      host,
      ACCESS_TOKEN_TTL,
      roles,
      permissions
    );
    if (accessTokenError) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [accessTokenError],
      };
    }

    const { token: newRefreshToken, error: refreshTokenError } = await this.repository.signToken(
      user,
      host,
      REFRESH_TOKEN_TTL,
      roles,
      permissions
    );
    if (refreshTokenError) {
      return {
        accessToken: null,
        refreshToken: null,
        errors: [refreshTokenError],
      };
    }

    const saveError = await this.repository.saveRefreshToken(
      user.id,
      newRefreshToken!,
      REFRESH_TOKEN_TTL
    );
    if (saveError) {
      return { accessToken: null, refreshToken: null, errors: [saveError] };
    }

    return {
      accessToken: newAccessToken,
      refreshToken: newRefreshToken,
      errors: [],
    };
  }

  async getUserById(
    userId: number
  ): Promise<{ user: UserResponse | null; errors: ResponseError[] }> {
    return this.userService.getUser(userId);
  }
}

export function newAuthService(
  repository: AuthRepository,
  userService: UserService,
  userRoleRepository: UserRoleRepository,
  rolePermissionRepository: RolePermissionRepository,
  roleRepository: RoleRepository,
  permissionRepository: PermissionRepository
): IAuthService {
  return new AuthService(
    repository,
    userService,
    userRoleRepository,
    rolePermissionRepository,
    roleRepository,
    permissionRepository
  );
}
