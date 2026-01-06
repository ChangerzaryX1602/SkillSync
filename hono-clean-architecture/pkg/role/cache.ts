import Redis from "ioredis";
import { Role } from "../models";
import { randomJitter } from "../utils";

const CACHE_TTL = 15 * 60 * 1000;

export class RoleCache {
  private store: Redis | null;

  constructor(store: Redis | null) {
    this.store = store;
  }

  async get(id: number): Promise<Role | null> {
    if (!this.store) return null;

    const key = `pkg:role:get:${id}`;
    const data = await this.store.get(key);

    if (!data) return null;

    try {
      return JSON.parse(data) as Role;
    } catch {
      return null;
    }
  }

  async set(role: Role): Promise<void> {
    if (!this.store || !role) return;

    const key = `pkg:role:get:${role.id}`;
    const ttl = randomJitter(CACHE_TTL);

    await this.store.set(key, JSON.stringify(role), "PX", ttl);
  }

  async getByName(name: string): Promise<Role | null> {
    if (!this.store) return null;

    const key = `pkg:role:name:${name}`;
    const data = await this.store.get(key);

    if (!data) return null;

    try {
      return JSON.parse(data) as Role;
    } catch {
      return null;
    }
  }

  async setByName(name: string, role: Role): Promise<void> {
    if (!this.store || !role) return;

    const key = `pkg:role:name:${name}`;
    const ttl = randomJitter(CACHE_TTL);

    await this.store.set(key, JSON.stringify(role), "PX", ttl);
  }

  async invalidate(id: number, name?: string): Promise<void> {
    if (!this.store) return;

    const keys: string[] = [`pkg:role:get:${id}`];
    if (name) {
      keys.push(`pkg:role:name:${name}`);
    }

    await this.store.del(...keys);
  }

  async invalidateAllLists(): Promise<void> {
    if (!this.store) return;

    const pattern = `pkg:role:list:*`;
    const keys = await this.store.keys(pattern);

    if (keys.length > 0) {
      await this.store.del(...keys);
    }
  }
}

export function newRoleCache(store: Redis | null): RoleCache {
  return new RoleCache(store);
}
