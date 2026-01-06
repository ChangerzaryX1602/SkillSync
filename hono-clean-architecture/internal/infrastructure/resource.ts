import { NodePgDatabase } from "drizzle-orm/node-postgres";
import Redis from "ioredis";
import { JwtResources } from "../datasources/jwt";
import { Resources, RouterResources } from "../../pkg/models";

export function newResources(
  mainDbConn: NodePgDatabase,
  redisStorage: Redis | null,
  jwtResources: JwtResources
): Resources {
  return {
    mainDbConn,
    redisStorage,
    jwtResources,
  };
}

export function newRouterResources(
  jwtResources: JwtResources,
  mainDbConn: NodePgDatabase,
  redisStorage: Redis | null
): RouterResources {
  return {
    jwtResources,
    mainDbConn,
    redisStorage,
  };
}
