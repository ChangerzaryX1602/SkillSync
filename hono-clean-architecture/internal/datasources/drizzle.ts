import { drizzle, NodePgDatabase } from "drizzle-orm/node-postgres";
import { Pool, PoolConfig } from "pg";
import { getConfig } from "../../pkg/config/config";

export interface DbConfig {
  host: string;
  port: number;
  username: string;
  password: string;
  dbName: string;
  minConnections: number;
  maxConnections: number;
}

let dbPool: Pool | null = null;
let dbInstance: NodePgDatabase | null = null;

export function buildDbConfig(): DbConfig {
  const config = getConfig();
  return {
    host: config.db.postgres.host,
    port: config.db.postgres.port,
    username: config.db.postgres.username,
    password: config.db.postgres.password,
    dbName: config.db.postgres.db_name,
    minConnections: config.db.postgres.conn.min,
    maxConnections: config.db.postgres.conn.max,
  };
}

export async function connectDb(config?: DbConfig): Promise<NodePgDatabase> {
  const dbConfig = config ?? buildDbConfig();

  const poolConfig: PoolConfig = {
    host: dbConfig.host,
    port: dbConfig.port,
    user: dbConfig.username,
    password: dbConfig.password,
    database: dbConfig.dbName,
    min: dbConfig.minConnections,
    max: dbConfig.maxConnections,
  };

  dbPool = new Pool(poolConfig);

  const client = await dbPool.connect();
  client.release();

  dbInstance = drizzle(dbPool);
  return dbInstance;
}

export function getDb(): NodePgDatabase {
  if (!dbInstance) {
    throw new Error("Database not connected. Call connectDb() first.");
  }
  return dbInstance;
}

export function getPool(): Pool {
  if (!dbPool) {
    throw new Error("Database pool not initialized. Call connectDb() first.");
  }
  return dbPool;
}

export async function closeDb(): Promise<void> {
  if (dbPool) {
    await dbPool.end();
    dbPool = null;
    dbInstance = null;
  }
}
