import Redis from "ioredis";
import { UserRole } from "../models";
import { randomJitter } from "../utils";

const CACHE_TTL = 15 * 60 * 1000;

export class UserRoleCache {
  private store: Redis | null;

  constructor(store: Redis | null) {
    this.store = store;
  }

  async get(id: number): Promise<UserRole | null> {
    if (!this.store) return null;

    const key = `pkg:userrole:get:${id}`;
    const data = await this.store.get(key);

    if (!data) return null;

    try {
      return JSON.parse(data) as UserRole;
    } catch {
      return null;
    }
  }

  async set(userRole: UserRole): Promise<void> {
    if (!this.store || !userRole) return;

    const key = `pkg:userrole:get:${userRole.id}`;
    const ttl = randomJitter(CACHE_TTL);

    await this.store.set(key, JSON.stringify(userRole), "PX", ttl);
  }

  async invalidate(id: number): Promise<void> {
    if (!this.store) return;
    await this.store.del(`pkg:userrole:get:${id}`);
  }

  async invalidateAllLists(): Promise<void> {
    if (!this.store) return;

    const pattern = `pkg:userrole:list:*`;
    const keys = await this.store.keys(pattern);

    if (keys.length > 0) {
      await this.store.del(...keys);
    }
  }
}

export function newUserRoleCache(store: Redis | null): UserRoleCache {
  return new UserRoleCache(store);
}
