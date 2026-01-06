import { Hono } from "hono";
import { PermissionService } from "../domain";
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

export function newPermissionHandler(
  app: Hono,
  routerResources: RouterResources,
  service: PermissionService
): void {
  const permissionsRouter = new Hono();

  permissionsRouter.post(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Permission, PermissionName.Create)
    ),
    async (c) => {
      const body = await c.req.json<{ group: string; name: string }>();
      const errors = await service.createPermission(body);

      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 400 | 409 | 500);
      }

      return c.json(successResponse({ message: "Permission created successfully" }), 201);
    }
  );

  permissionsRouter.get(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Permission, PermissionName.Read)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid permission ID")]),
          400
        );
      }

      const { permission, errors } = await service.getPermission(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse(permission));
    }
  );

  permissionsRouter.get(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Permission, PermissionName.List)
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

      const result = await service.getPermissions(pagination, search);
      if (result.errors.length > 0) {
        return c.json(errorResponse(result.errors), result.errors[0].code as 500);
      }

      return c.json(
        successResponse(result.permissions, {
          pagination: result.pagination,
          search: result.search,
        })
      );
    }
  );

  permissionsRouter.patch(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Permission, PermissionName.Update)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid permission ID")]),
          400
        );
      }

      const body = await c.req.json();
      const errors = await service.updatePermission(id, body);

      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 400 | 404 | 500);
      }

      return c.json(successResponse({ message: "Permission updated successfully" }));
    }
  );

  permissionsRouter.delete(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.Permission, PermissionName.Delete)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid permission ID")]),
          400
        );
      }

      const errors = await service.deletePermission(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse({ message: "Permission deleted successfully" }));
    }
  );

  app.route("/permissions", permissionsRouter);
}
