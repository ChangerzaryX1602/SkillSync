import { User, UserResponse, Pagination, Search } from "../models";
import { ResponseError } from "../../internal/infrastructure/custom_error";

export interface UserRepository {
  migrate(): Promise<void>;
  createUser(user: {
    username: string;
    email: string;
    password: string;
  }): Promise<{ user: UserResponse | null; error: ResponseError | null }>;
  getUser(id: number): Promise<{ user: UserResponse | null; error: ResponseError | null }>;
  getUsers(
    pagination: Pagination,
    search: Search
  ): Promise<{
    users: UserResponse[];
    pagination: Pagination;
    search: Search;
    error: ResponseError | null;
  }>;
  updateUser(id: number, user: Partial<User>): Promise<ResponseError | null>;
  deleteUser(id: number): Promise<ResponseError | null>;
  getUserByEmail(email: string): Promise<{ user: User | null; error: ResponseError | null }>;
}

export interface UserService {
  createUser(user: { username: string; email: string; password: string }): Promise<ResponseError[]>;
  getUser(id: number): Promise<{ user: UserResponse | null; errors: ResponseError[] }>;
  getUsers(
    pagination: Pagination,
    search: Search
  ): Promise<{
    users: UserResponse[];
    pagination: Pagination;
    search: Search;
    errors: ResponseError[];
  }>;
  updateUser(id: number, user: Partial<User>): Promise<ResponseError[]>;
  deleteUser(id: number): Promise<ResponseError[]>;
  getUserByEmail(email: string): Promise<{ user: User | null; errors: ResponseError[] }>;
}
