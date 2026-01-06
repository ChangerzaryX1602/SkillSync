import { readFileSync } from "fs";
import { parse } from "yaml";
import { join, dirname } from "path";
import { fileURLToPath } from "url";

export interface AppConfig {
  app: {
    name: string;
    env: string;
    port: {
      http: number;
      https: number;
    };
    path: {
      cert: string;
      priv: string;
      ca: string;
      log: string;
    };
  };
  db: {
    postgres: {
      sock: string;
      host: string;
      port: number;
      username: string;
      password: string;
      db_name: string;
      conn: {
        min: number;
        max: number;
      };
    };
    redis: {
      host: string;
      port: number;
      username: string;
      password: string;
      db_name: number;
    };
  };
  jwt: {
    private: string;
    public: string;
  };
}

let currentConfig: AppConfig | null = null;

function getProjectRoot(): string {
  const currentFile = fileURLToPath(import.meta.url);
  return join(dirname(currentFile), "..", "..");
}

export function loadConfig(configName: string = "dev"): AppConfig {
  const projectRoot = getProjectRoot();
  const configPath = join(projectRoot, "configs", `${configName}.yaml`);
  const configFile = readFileSync(configPath, "utf-8");
  currentConfig = parse(configFile) as AppConfig;
  return currentConfig;
}

export function getConfig(): AppConfig {
  if (!currentConfig) {
    throw new Error("Config not loaded. Call loadConfig() first.");
  }
  return currentConfig;
}
