import { Hono } from "hono";
import { UserRoleService } from "../domain";
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

export function newUserRoleHandler(
  app: Hono,
  routerResources: RouterResources,
  service: UserRoleService
): void {
  const userRolesRouter = new Hono();

  userRolesRouter.post(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.UserRole, PermissionName.Create)
    ),
    async (c) => {
      const body = await c.req.json<{ userId: number; roleId: number }>();
      const errors = await service.createUserRole(body);

      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 400 | 409 | 500);
      }

      return c.json(successResponse({ message: "UserRole created successfully" }), 201);
    }
  );

  userRolesRouter.get(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.UserRole, PermissionName.Read)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid user role ID")]),
          400
        );
      }

      const { userRole, errors } = await service.getUserRole(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse(userRole));
    }
  );

  userRolesRouter.get(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.UserRole, PermissionName.List)
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

      const result = await service.getUserRoles(pagination, search);
      if (result.errors.length > 0) {
        return c.json(errorResponse(result.errors), result.errors[0].code as 500);
      }

      return c.json(
        successResponse(result.userRoles, {
          pagination: result.pagination,
          search: result.search,
        })
      );
    }
  );

  userRolesRouter.delete(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.UserRole, PermissionName.Delete)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid user role ID")]),
          400
        );
      }

      const errors = await service.deleteUserRole(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse({ message: "UserRole deleted successfully" }));
    }
  );

  app.route("/user/roles", userRolesRouter);
}
