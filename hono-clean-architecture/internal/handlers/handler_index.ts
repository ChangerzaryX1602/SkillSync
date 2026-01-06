import { Hono } from "hono";

export function newIndexHandler(app: Hono): void {
  app.get("/", (c) => {
    return c.json({
      message: "SkillSync API",
      status: "healthy",
      version: "1.0.0",
    });
  });

  app.get("/health", (c) => {
    return c.json({
      status: "ok",
      timestamp: new Date().toISOString(),
    });
  });
}
