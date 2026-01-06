import { Hono } from "hono";
import { RoleService } from "../domain";
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

export function newRoleHandler(
  app: Hono,
  routerResources: RouterResources,
  service: RoleService
): void {
  const rolesRouter = new Hono();

  rolesRouter.post(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Role, PermissionName.Create)
    ),
    async (c) => {
      const body = await c.req.json<{ name: string }>();
      const errors = await service.createRole(body);

      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 400 | 409 | 500);
      }

      return c.json(successResponse({ message: "Role created successfully" }), 201);
    }
  );

  rolesRouter.get(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Role, PermissionName.Read)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid role ID")]),
          400
        );
      }

      const { role, errors } = await service.getRole(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse(role));
    }
  );

  rolesRouter.get(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Role, PermissionName.List)
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

      const result = await service.getRoles(pagination, search);
      if (result.errors.length > 0) {
        return c.json(errorResponse(result.errors), result.errors[0].code as 500);
      }

      return c.json(
        successResponse(result.roles, {
          pagination: result.pagination,
          search: result.search,
        })
      );
    }
  );

  rolesRouter.patch(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Role, PermissionName.Update)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid role ID")]),
          400
        );
      }

      const body = await c.req.json();
      const errors = await service.updateRole(id, body);

      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 400 | 404 | 500);
      }

      return c.json(successResponse({ message: "Role updated successfully" }));
    }
  );

  rolesRouter.delete(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Role, PermissionName.Delete)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid role ID")]),
          400
        );
      }

      const errors = await service.deleteRole(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse({ message: "Role deleted successfully" }));
    }
  );

  app.route("/roles", rolesRouter);
}
