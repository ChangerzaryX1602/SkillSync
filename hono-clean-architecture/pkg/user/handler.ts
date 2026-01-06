import { Hono } from "hono";
import { UserService } from "../domain";
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
import { authMiddleware, getAuthContext } from "../../internal/handlers/middleware_auth";

export function newUserHandler(
  app: Hono,
  routerResources: RouterResources,
  service: UserService
): void {
  const users = new Hono();

  users.post(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.User, PermissionName.Create)
    ),
    async (c) => {
      const body = await c.req.json<{
        username: string;
        email: string;
        password: string;
      }>();

      const errors = await service.createUser(body);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 400 | 409 | 500);
      }

      return c.json(successResponse({ message: "User created successfully" }), 201);
    }
  );

  users.get(
    "/me",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.User, PermissionName.Me)
    ),
    async (c) => {
      const auth = getAuthContext(c);
      const userId = auth?.userId ?? 0;

      const { user, errors } = await service.getUser(userId);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse(user));
    }
  );

  users.get(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.User, PermissionName.Read)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid user ID")]),
          400
        );
      }

      const { user, errors } = await service.getUser(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse(user));
    }
  );

  users.get(
    "/",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.User, PermissionName.List)
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

      const result = await service.getUsers(pagination, search);
      if (result.errors.length > 0) {
        return c.json(errorResponse(result.errors), result.errors[0].code as 500);
      }

      return c.json(
        successResponse(result.users, {
          pagination: result.pagination,
          search: result.search,
        })
      );
    }
  );

  users.patch(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.User, PermissionName.Update)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid user ID")]),
          400
        );
      }

      const body = await c.req.json();
      const errors = await service.updateUser(id, body);

      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 400 | 404 | 500);
      }

      return c.json(successResponse({ message: "User updated successfully" }));
    }
  );

  users.delete(
    "/:id",
    authMiddleware(
      routerResources,
      permissionGroupName(PermissionGroupType.User, PermissionName.Delete)
    ),
    async (c) => {
      const id = parseInt(c.req.param("id"), 10);

      if (isNaN(id)) {
        return c.json(
          errorResponse([createError(400, whereAmI(), "Invalid Request", "Invalid user ID")]),
          400
        );
      }

      const errors = await service.deleteUser(id);
      if (errors.length > 0) {
        return c.json(errorResponse(errors), errors[0].code as 404 | 500);
      }

      return c.json(successResponse({ message: "User deleted successfully" }));
    }
  );

  app.route("/users", users);
}
