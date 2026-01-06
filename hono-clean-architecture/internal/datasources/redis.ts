import Redis from "ioredis";
import { getConfig } from "../../pkg/config/config";

export interface RedisConfig {
  host: string;
  port: number;
  username?: string;
  password?: string;
  dbName: number;
}

let redisClient: Redis | null = null;

export function buildRedisConfig(): RedisConfig {
  const config = getConfig();
  return {
    host: config.db.redis.host,
    port: config.db.redis.port,
    username: config.db.redis.username || undefined,
    password: config.db.redis.password || undefined,
    dbName: config.db.redis.db_name,
  };
}

export async function connectToRedis(config?: RedisConfig): Promise<Redis> {
  const redisConfig = config ?? buildRedisConfig();

  redisClient = new Redis({
    host: redisConfig.host,
    port: redisConfig.port,
    username: redisConfig.username,
    password: redisConfig.password,
    db: redisConfig.dbName,
    lazyConnect: true,
    maxRetriesPerRequest: 3,
  });

  await redisClient.connect();

  await redisClient.ping();

  return redisClient;
}

export function getRedis(): Redis {
  if (!redisClient) {
    throw new Error("Redis not connected. Call connectToRedis() first.");
  }
  return redisClient;
}

export async function closeRedis(): Promise<void> {
  if (redisClient) {
    await redisClient.quit();
    redisClient = null;
  }
}
