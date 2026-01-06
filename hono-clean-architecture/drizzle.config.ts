import { defineConfig } from "drizzle-kit";
import { loadConfig } from "./pkg/config/config";

const config = loadConfig(process.env.ENV || "dev");

export default defineConfig({
  schema: "./pkg/models/*.ts",
  out: "./drizzle",
  dialect: "postgresql",
  dbCredentials: {
    host: config.db.postgres.host,
    port: config.db.postgres.port,
    user: config.db.postgres.username,
    password: config.db.postgres.password,
    database: config.db.postgres.db_name,
    ssl: false,
  },
});
