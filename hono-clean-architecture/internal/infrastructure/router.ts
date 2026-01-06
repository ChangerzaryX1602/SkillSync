import { Hono } from "hono";
import { Resources, RouterResources } from "../../pkg/models";
import { newIndexHandler } from "../handlers/handler_index";
import { getPool } from "../datasources/drizzle";
import { newTxManager } from "../../pkg/tx_context/tx_manager";
import { migratePermissions } from "../../pkg/utils/migrate-permission";

import { newUserRepository } from "../../pkg/user/repository";
import { newRoleRepository } from "../../pkg/role/repository";
import { newPermissionRepository } from "../../pkg/permission/repository";
import { newRolePermissionRepository } from "../../pkg/role_permission/repository";
import { newUserRoleRepository } from "../../pkg/user_role/repository";
import { newAuthRepository } from "../../pkg/auth/repository";

import { newUserService } from "../../pkg/user/service";
import { newRoleService } from "../../pkg/role/service";
import { newPermissionService } from "../../pkg/permission/service";
import { newRolePermissionService } from "../../pkg/role_permission/service";
import { newUserRoleService } from "../../pkg/user_role/service";
import { newAuthService } from "../../pkg/auth/service";

import { newUserHandler } from "../../pkg/user/handler";
import { newRoleHandler } from "../../pkg/role/handler";
import { newPermissionHandler } from "../../pkg/permission/handler";
import { newRolePermissionHandler } from "../../pkg/role_permission/handler";
import { newUserRoleHandler } from "../../pkg/user_role/handler";
import { newAuthHandler } from "../../pkg/auth/handler";

export async function setupRoutes(
  app: Hono,
  resources: Resources,
  routerResources: RouterResources
): Promise<void> {
  newIndexHandler(app);

  // Migrate default permissions and roles on startup
  try {
    await migratePermissions(resources.mainDbConn);
    console.log("✅ Permissions migrated successfully");
  } catch (err) {
    console.error("❌ Failed to migrate permissions:", err);
  }

  const apiV1 = new Hono();

  apiV1.get("/", (c) => {
    return c.json({ message: "SkillSync API v1", status: "ok" });
  });

  const txManager = newTxManager(getPool());

  const userRepository = newUserRepository(resources);
  const roleRepository = newRoleRepository(resources);
  const permissionRepository = newPermissionRepository(resources);
  const rolePermissionRepository = newRolePermissionRepository(resources);
  const userRoleRepository = newUserRoleRepository(resources);
  const authRepository = newAuthRepository(resources);

  const userService = newUserService(userRepository, roleRepository, userRoleRepository, txManager);
  const roleService = newRoleService(roleRepository);
  const permissionService = newPermissionService(permissionRepository);
  const rolePermissionService = newRolePermissionService(rolePermissionRepository);
  const userRoleService = newUserRoleService(userRoleRepository);
  const authService = newAuthService(authRepository, userService);

  newUserHandler(apiV1, routerResources, userService);
  newRoleHandler(apiV1, routerResources, roleService);
  newPermissionHandler(apiV1, routerResources, permissionService);
  newRolePermissionHandler(apiV1, routerResources, rolePermissionService);
  newUserRoleHandler(apiV1, routerResources, userRoleService);
  newAuthHandler(apiV1, authService);

  app.route("/api/v1", apiV1);
}
