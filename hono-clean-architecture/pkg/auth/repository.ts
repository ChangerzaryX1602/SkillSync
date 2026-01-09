import Redis from "ioredis";
import { User, UserResponse, Resources, CacheKeys } from "../models";
import { AuthRepository as IAuthRepository } from "../domain";
import { ResponseError, createError, whereAmI } from "../../internal/infrastructure/custom_error";
import { signJwt, parseJwt, JwtResources } from "../../internal/datasources/jwt";
import { JWTPayload } from "jose";

export class AuthRepository implements IAuthRepository {
  private jwtResources: JwtResources;
  private redisStorage: Redis | null;

  constructor(resources: Resources) {
    this.jwtResources = resources.jwtResources;
    this.redisStorage = resources.redisStorage;
  }

  async signToken(
    user: User | UserResponse,
    host: string,
    ttlSeconds: number,
    roles: string[] = [],
    permissions: string[] = []
  ): Promise<{ token: string | null; error: ResponseError | null }> {
    try {
      const now = Math.floor(Date.now() / 1000);
      const payload: JWTPayload = {
        sub: String(user.id),
        iss: host,
        exp: now + ttlSeconds,
        iat: now,
        roles,
        permissions,
      };

      const token = await signJwt(this.jwtResources, payload);
      return { token, error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to sign token";
      return {
        token: null,
        error: createError(500, whereAmI(), "Failed to sign json web token", message),
      };
    }
  }

  async saveRefreshToken(
    userId: number,
    token: string,
    ttlSeconds: number
  ): Promise<ResponseError | null> {
    if (!this.redisStorage) {
      return createError(500, whereAmI(), "Redis Error", "Redis storage not available");
    }

    try {
      const key = `${CacheKeys.RefreshToken}:${userId}`;
      await this.redisStorage.set(key, token, "EX", ttlSeconds);
      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to save refresh token";
      return createError(500, whereAmI(), "Failed to save refresh token", message);
    }
  }

  async getRefreshToken(
    userId: number
  ): Promise<{ token: string | null; error: ResponseError | null }> {
    if (!this.redisStorage) {
      return {
        token: null,
        error: createError(500, whereAmI(), "Redis Error", "Redis storage not available"),
      };
    }

    try {
      const key = `${CacheKeys.RefreshToken}:${userId}`;
      const token = await this.redisStorage.get(key);

      if (!token) {
        return {
          token: null,
          error: createError(401, whereAmI(), "Unauthorized", "Refresh token not found"),
        };
      }

      return { token, error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to get refresh token";
      return {
        token: null,
        error: createError(500, whereAmI(), "Failed to get refresh token", message),
      };
    }
  }

  async deleteRefreshToken(userId: number): Promise<ResponseError | null> {
    if (!this.redisStorage) {
      return createError(500, whereAmI(), "Redis Error", "Redis storage not available");
    }

    try {
      const key = `${CacheKeys.RefreshToken}:${userId}`;
      await this.redisStorage.del(key);
      return null;
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to delete refresh token";
      return createError(500, whereAmI(), "Failed to delete refresh token", message);
    }
  }

  async parseToken(tokenString: string): Promise<{
    payload: JWTPayload | null;
    valid: boolean;
    error: ResponseError | null;
  }> {
    try {
      const { payload, valid } = await parseJwt(this.jwtResources, tokenString);
      return { payload, valid, error: null };
    } catch (error) {
      const message = error instanceof Error ? error.message : "Invalid token";
      return {
        payload: null,
        valid: false,
        error: createError(401, whereAmI(), "Invalid token", message),
      };
    }
  }
}

export function newAuthRepository(resources: Resources): IAuthRepository {
  return new AuthRepository(resources);
}
