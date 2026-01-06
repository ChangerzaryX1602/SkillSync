import { Hono } from "hono";
import { RolePermissionService } from "../domain";
import {
  RouterResources,
  createPagination,
  createSearch,
  PermissionGroupType,
  PermissionName,
  permissionGroupName,
} from "../models";
import {
  successResponse,
  errorResponse,
  createError,
  whereAmI,
} from "../../internal/infrastructure/custom_error";
import { authMiddleware } from "../../internal/handlers/middleware_auth";

export function newRolePermissionHandler(
  app: Hono,
  routerResources: RouterResources,
  service: RolePermissionService
): void {
  const rolePermissionsRouter = new Hono();

  rolePermissionsRouter.post(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.RolePermission, PermissionName.Create)
    ),
    async (c) => {
      const body = await c.req.json<{ roleId: number; permissionId: number }>();
      const errors = await service.createRolePermission(body);

      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 400 | 409 | 500);
      }

      return c.json(successResponse({ message: "RolePermission created successfully" }), 201);
    }
  );

  rolePermissionsRouter.get(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.RolePermission, PermissionName.Read)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([
            createError(400, whereAmI(), "Invalid Request", "Invalid role permission ID"),
          ]),
          400
        );
      }

      const { rolePermission, errors } = await service.getRolePermission(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse(rolePermission));
    }
  );

  rolePermissionsRouter.get(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.RolePermission, PermissionName.List)
    ),
    async (c) => {
      const query = c.req.query();
      const pagination = createPagination(
        query.page ? parseInt(query.page, 10) : undefined,
        query.per_page ? parseInt(query.per_page, 10) : undefined,
        query.order_by,
        query.sort_by
      );
      const search = createSearch(query.keyword, query.column, query.text);

      const result = await service.getRolePermissions(pagination, search);
      if (result.errors.length > 0) {
        return c.json(errorResponse(result.errors), result.errors[0].code as 500);
      }

      return c.json(
        successResponse(result.rolePermissions, {
          pagination: result.pagination,
          search: result.search,
        })
      );
    }
  );

  rolePermissionsRouter.delete(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.RolePermission, PermissionName.Delete)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([
            createError(400, whereAmI(), "Invalid Request", "Invalid role permission ID"),
          ]),
          400
        );
      }

      const errors = await service.deleteRolePermission(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse({ message: "RolePermission deleted successfully" }));
    }
  );

  app.route("/role/permissions", rolePermissionsRouter);
}
