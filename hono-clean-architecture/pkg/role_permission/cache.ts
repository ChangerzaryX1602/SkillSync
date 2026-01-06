import Redis from "ioredis";
import { RolePermission } from "../models";
import { randomJitter } from "../utils";

const CACHE_TTL = 15 * 60 * 1000;

export class RolePermissionCache {
  private store: Redis | null;

  constructor(store: Redis | null) {
    this.store = store;
  }

  async get(id: number): Promise<RolePermission | null> {
    if (!this.store) return null;

    const key = `pkg:rolepermission:get:${id}`;
    const data = await this.store.get(key);

    if (!data) return null;

    try {
      return JSON.parse(data) as RolePermission;
    } catch {
      return null;
    }
  }

  async set(rolePermission: RolePermission): Promise<void> {
    if (!this.store || !rolePermission) return;

    const key = `pkg:rolepermission:get:${rolePermission.id}`;
    const ttl = randomJitter(CACHE_TTL);

    await this.store.set(key, JSON.stringify(rolePermission), "PX", ttl);
  }

  async invalidate(id: number): Promise<void> {
    if (!this.store) return;
    await this.store.del(`pkg:rolepermission:get:${id}`);
  }

  async invalidateAllLists(): Promise<void> {
    if (!this.store) return;

    const pattern = `pkg:rolepermission:list:*`;
    const keys = await this.store.keys(pattern);

    if (keys.length > 0) {
      await this.store.del(...keys);
    }
  }
}

export function newRolePermissionCache(store: Redis | null): RolePermissionCache {
  return new RolePermissionCache(store);
}
