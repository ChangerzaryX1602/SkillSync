import { Hono } from "hono";
import { cors } from "hono/cors";
import { logger } from "hono/logger";
import { prettyJSON } from "hono/pretty-json";
import { Resources, RouterResources } from "../../pkg/models";
import { setupRoutes } from "./router";

export async function createServer(
  resources: Resources,
  routerResources: RouterResources
): Promise<Hono> {
  const app = new Hono();

  app.use("*", logger());
  app.use("*", prettyJSON());
  app.use(
    "*",
    cors({
      origin: "*",
      allowMethods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"],
      allowHeaders: ["Content-Type", "Authorization"],
      exposeHeaders: ["Content-Length"],
      maxAge: 86400,
    })
  );

  await setupRoutes(app, resources, routerResources);

  return app;
}

export async function startServer(app: Hono, port: number, host: string): Promise<void> {
  console.log(`🚀 Server starting on http://${host}:${port}`);

  const server = Bun.serve({
    port,
    hostname: host,
    fetch: app.fetch,
  });

  console.log(`✅ Server running on http://${host}:${port}`);

  process.on("SIGINT", () => {
    console.log("\n🛑 Shutting down gracefully...");
    server.stop();
    process.exit(0);
  });

  process.on("SIGTERM", () => {
    console.log("\n🛑 Shutting down gracefully...");
    server.stop();
    process.exit(0);
  });
}
