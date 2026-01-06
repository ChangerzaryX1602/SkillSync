import { User, UserResponse } from "../models";
import { ResponseError } from "../../internal/infrastructure/custom_error";
import { JWTPayload } from "jose";

export interface AuthRepository {
  signToken(
    user: User | UserResponse,
    host: string,
    ttlSeconds: number
  ): Promise<{ token: string | null; error: ResponseError | null }>;
  saveRefreshToken(
    userId: number,
    token: string,
    ttlSeconds: number
  ): Promise<ResponseError | null>;
  getRefreshToken(userId: number): Promise<{ token: string | null; error: ResponseError | null }>;
  deleteRefreshToken(userId: number): Promise<ResponseError | null>;
  parseToken(
    tokenString: string
  ): Promise<{ payload: JWTPayload | null; valid: boolean; error: ResponseError | null }>;
}

export interface AuthService {
  register(user: { username: string; email: string; password: string }): Promise<ResponseError[]>;
  login(
    user: { email: string; password: string },
    host: string
  ): Promise<{
    accessToken: string | null;
    refreshToken: string | null;
    errors: ResponseError[];
  }>;
  refreshToken(refreshToken: string): Promise<{
    accessToken: string | null;
    refreshToken: string | null;
    errors: ResponseError[];
  }>;
  getUserById(userId: number): Promise<{ user: UserResponse | null; errors: ResponseError[] }>;
}
