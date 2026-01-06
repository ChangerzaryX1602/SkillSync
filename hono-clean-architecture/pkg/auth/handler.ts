import { Hono } from "hono";
import { AuthService } from "../domain";
import {
  successResponse,
  errorResponse,
  createError,
  whereAmI,
} from "../../internal/infrastructure/custom_error";

export function newAuthHandler(app: Hono, service: AuthService): void {
  const auth = new Hono();

  auth.post("/login", async (c) => {
    const body = await c.req.json<{ email: string; password: string }>();

    const host = c.req.header("host") || "localhost";
    const result = await service.login(body, host);

    if (result.errors.length > 0) {
      return c.json(errorResponse(result.errors), result.errors[0].code as 400 | 401 | 500);
    }

    return c.json({
      success: true,
      access_token: result.accessToken,
      refresh_token: result.refreshToken,
    });
  });

  auth.post("/register", async (c) => {
    const body = await c.req.json<{
      username: string;
      email: string;
      password: string;
    }>();

    const errors = await service.register(body);
    if (errors.length > 0) {
      return c.json(errorResponse(errors), errors[0].code as 400 | 409 | 500);
    }

    return c.json(successResponse({ message: "User registered successfully" }), 201);
  });

  auth.post("/refresh", async (c) => {
    const body = await c.req.json<{ refresh_token: string }>();

    if (!body.refresh_token) {
      return c.json(
        errorResponse([
          createError(400, whereAmI(), "Invalid Request", "refresh_token is required"),
        ]),
        400
      );
    }

    const result = await service.refreshToken(body.refresh_token);

    if (result.errors.length > 0) {
      return c.json(errorResponse(result.errors), result.errors[0].code as 400 | 401 | 500);
    }

    return c.json({
      success: true,
      access_token: result.accessToken,
      refresh_token: result.refreshToken,
    });
  });

  app.route("/auth", auth);
}
