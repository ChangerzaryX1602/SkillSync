import { NodePgDatabase } from "drizzle-orm/node-postgres";
import Redis from "ioredis";
import { JwtResources } from "../../internal/datasources/jwt";

export interface Resources {
  mainDbConn: NodePgDatabase;
  redisStorage: Redis | null;
  jwtResources: JwtResources;
}

export interface RouterResources {
  jwtResources: JwtResources;
  mainDbConn: NodePgDatabase;
  redisStorage: Redis | null;
}
