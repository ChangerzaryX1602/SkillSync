import { loadConfig } from "../../pkg/config/config";
import {
  newResources,
  newRouterResources,
  createServer,
  startServer,
} from "../../internal/infrastructure";
import { connectDb } from "../../internal/datasources/drizzle";
import { connectToRedis } from "../../internal/datasources/redis";
import { loadJwtKeys } from "../../internal/datasources/jwt";
import { join, dirname } from "path";
import { fileURLToPath } from "url";

function getProjectRoot(): string {
  const currentFile = fileURLToPath(import.meta.url);
  return join(dirname(currentFile), "..", "..");
}

async function main(): Promise<void> {
  const env = process.env.ENV || "dev";
  console.log(`🔧 Loading configuration for environment: ${env}`);

  const config = loadConfig(env);
  console.log(`✅ Configuration loaded successfully`);

  const db = await connectDb();
  console.log(`✅ Database connected successfully`);

  let redis = null;
  try {
    redis = await connectToRedis();
    console.log(`✅ Redis connected successfully`);
  } catch (err) {
    console.warn(`⚠️ Redis connection failed, continuing without cache: ${err}`);
  }

  const projectRoot = getProjectRoot();
  const jwtKeyPath = join(projectRoot, config.jwt.private.replace("./", ""));
  const jwtResources = loadJwtKeys(jwtKeyPath);
  console.log(`✅ JWT configured successfully (algorithm: ${jwtResources.algorithm})`);

  const resources = newResources(db, redis, jwtResources);
  const routerResources = newRouterResources(jwtResources, db, redis);

  const app = await createServer(resources, routerResources);
  await startServer(app, config.app.port.http, "0.0.0.0");
}

main().catch((err) => {
  console.error("❌ Failed to start server:", err);
  process.exit(1);
});
