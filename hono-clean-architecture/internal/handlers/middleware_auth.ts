import { Context, Next } from "hono";
import { RouterResources } from "../../pkg/models";
import { errorResponse, createError, whereAmI } from "../infrastructure/custom_error";
import { verifyJwt } from "../datasources/jwt";

export interface AuthContext {
  userId: number;
  roles: string[];
  permissions: string[];
}

export function authMiddleware(routerResources: RouterResources, ...requiredPermissions: string[]) {
  return async (c: Context, next: Next): Promise<Response | void> => {
    const authHeader = c.req.header("Authorization");

    if (!authHeader) {
      return c.json(
        errorResponse([
          createError(401, whereAmI(), "Unauthorized", "Authorization header is required"),
        ]),
        401
      );
    }

    const parts = authHeader.split(" ");
    if (parts.length !== 2 || parts[0] !== "Bearer") {
      return c.json(
        errorResponse([
          createError(401, whereAmI(), "Unauthorized", "Invalid authorization header format"),
        ]),
        401
      );
    }

    const token = parts[1];

    try {
      const payload = await verifyJwt(routerResources.jwtResources, token);

      const userId = Number(payload.sub);
      const permissions = (payload.permissions as string[]) || [];
      const roles = (payload.roles as string[]) || [];

      if (requiredPermissions.length > 0) {
        const hasPermission = requiredPermissions.some((perm) => permissions.includes(perm));

        if (!hasPermission) {
          return c.json(
            errorResponse([createError(403, whereAmI(), "Forbidden", "Insufficient permissions")]),
            403
          );
        }
      }

      c.set("auth", {
        userId,
        roles,
        permissions,
      } as AuthContext);

      await next();
    } catch (error) {
      const message = error instanceof Error ? error.message : "Authentication failed";
      return c.json(errorResponse([createError(401, whereAmI(), "Unauthorized", message)]), 401);
    }
  };
}

export function optionalAuthMiddleware(routerResources: RouterResources) {
  return async (c: Context, next: Next) => {
    const authHeader = c.req.header("Authorization");

    if (authHeader) {
      const parts = authHeader.split(" ");
      if (parts.length === 2 && parts[0] === "Bearer") {
        const token = parts[1];

        try {
          const payload = await verifyJwt(routerResources.jwtResources, token);

          const userId = Number(payload.sub);
          const permissions = (payload.permissions as string[]) || [];
          const roles = (payload.roles as string[]) || [];

          c.set("auth", {
            userId,
            roles,
            permissions,
          } as AuthContext);
        } catch {}
      }
    }

    await next();
  };
}

export function getAuthContext(c: Context): AuthContext | null {
  return c.get("auth") as AuthContext | null;
}

export function requireAuth(c: Context): AuthContext {
  const auth = c.get("auth") as AuthContext | null;
  if (!auth) {
    throw new Error("Auth context not found");
  }
  return auth;
}
