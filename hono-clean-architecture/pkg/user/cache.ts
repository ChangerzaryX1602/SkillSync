import Redis from "ioredis";
import { UserResponse, Pagination, Search, CacheKeys } from "../models";
import { randomJitter } from "../utils";

const CACHE_TTL = 15 * 60 * 1000;
const CACHE_TTL_LIST = 1 * 60 * 1000;

export class UserCache {
  private store: Redis | null;

  constructor(store: Redis | null) {
    this.store = store;
  }

  async get(id: number): Promise<UserResponse | null> {
    if (!this.store) return null;

    const key = `${CacheKeys.PkgUserGetUser}:${id}`;
    const data = await this.store.get(key);

    if (!data) return null;

    try {
      return JSON.parse(data) as UserResponse;
    } catch {
      return null;
    }
  }

  async set(user: UserResponse): Promise<void> {
    if (!this.store || !user) return;

    const key = `${CacheKeys.PkgUserGetUser}:${user.id}`;
    const ttl = randomJitter(CACHE_TTL);

    await this.store.set(key, JSON.stringify(user), "PX", ttl);
  }

  async getByEmail(email: string): Promise<UserResponse | null> {
    if (!this.store) return null;

    const key = `${CacheKeys.PkgUserGetUserByEmail}:${email}`;
    const data = await this.store.get(key);

    if (!data) return null;

    try {
      return JSON.parse(data) as UserResponse;
    } catch {
      return null;
    }
  }

  async setByEmail(email: string, user: UserResponse): Promise<void> {
    if (!this.store || !user) return;

    const key = `${CacheKeys.PkgUserGetUserByEmail}:${email}`;
    const ttl = randomJitter(CACHE_TTL);

    await this.store.set(key, JSON.stringify(user), "PX", ttl);
  }

  async getList(
    pagination: Pagination,
    search: Search
  ): Promise<{ users: UserResponse[]; pagination: Pagination; search: Search } | null> {
    if (!this.store) return null;

    const cacheKey = this.buildListCacheKey(pagination, search);
    const data = await this.store.get(cacheKey);

    if (!data) return null;

    try {
      return JSON.parse(data);
    } catch {
      return null;
    }
  }

  async setList(users: UserResponse[], pagination: Pagination, search: Search): Promise<void> {
    if (!this.store) return;

    const cacheKey = this.buildListCacheKey(pagination, search);
    const ttl = randomJitter(CACHE_TTL_LIST);

    await this.store.set(cacheKey, JSON.stringify({ users, pagination, search }), "PX", ttl);
  }

  async invalidate(id: number, email?: string): Promise<void> {
    if (!this.store) return;

    const keys: string[] = [`${CacheKeys.PkgUserGetUser}:${id}`];
    if (email) {
      keys.push(`${CacheKeys.PkgUserGetUserByEmail}:${email}`);
    }

    await this.store.del(...keys);
  }

  async invalidateAllLists(): Promise<void> {
    if (!this.store) return;

    const pattern = `${CacheKeys.PkgUserGetList}:*`;
    const keys = await this.store.keys(pattern);

    if (keys.length > 0) {
      await this.store.del(...keys);
    }
  }

  private buildListCacheKey(pagination: Pagination, search: Search): string {
    const paginationStr = `p${pagination.page}_pp${pagination.perPage}_ob${pagination.orderBy || ""}_o${pagination.order || ""}`;
    const searchStr = `k${search.keyword || ""}_c${search.column || ""}`;
    return `${CacheKeys.PkgUserGetList}:${paginationStr}:${searchStr}`;
  }
}

export function newUserCache(store: Redis | null): UserCache {
  return new UserCache(store);
}
