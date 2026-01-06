import Redis from "ioredis";
import { Permission } from "../models";
import { randomJitter } from "../utils";

const CACHE_TTL = 15 * 60 * 1000;

export class PermissionCache {
  private store: Redis | null;

  constructor(store: Redis | null) {
    this.store = store;
  }

  async get(id: number): Promise<Permission | null> {
    if (!this.store) return null;

    const key = `pkg:permission:get:${id}`;
    const data = await this.store.get(key);

    if (!data) return null;

    try {
      return JSON.parse(data) as Permission;
    } catch {
      return null;
    }
  }

  async set(permission: Permission): Promise<void> {
    if (!this.store || !permission) return;

    const key = `pkg:permission:get:${permission.id}`;
    const ttl = randomJitter(CACHE_TTL);

    await this.store.set(key, JSON.stringify(permission), "PX", ttl);
  }

  async invalidate(id: number): Promise<void> {
    if (!this.store) return;
    await this.store.del(`pkg:permission:get:${id}`);
  }

  async invalidateAllLists(): Promise<void> {
    if (!this.store) return;

    const pattern = `pkg:permission:list:*`;
    const keys = await this.store.keys(pattern);

    if (keys.length > 0) {
      await this.store.del(...keys);
    }
  }
}

export function newPermissionCache(store: Redis | null): PermissionCache {
  return new PermissionCache(store);
}
