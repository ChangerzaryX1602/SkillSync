import { Injectable, inject, PLATFORM_ID } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { isPlatformBrowser } from '@angular/common';
import { Observable, of } from 'rxjs';
import { environment } from '../../../environments/environment';

// Interfaces
export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface CreateUserDto {
  username: string;
  email: string;
  password: string;
}

export interface UpdateUserDto {
  username?: string;
  email?: string;
}

export interface Role {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface UserRole {
  id: number;
  user_id: number;
  role_id: number;
  created_at: string;
  updated_at: string;
}

export interface Pagination {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface Search {
  keyword: string;
  column: string;
}

export interface PaginationParams {
  page?: number;
  per_page?: number;
  order_by?: string;
  sort_by?: 'asc' | 'desc';
  keyword?: string;
  column?: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  errors?: Array<{
    code: number;
    source: string;
    title: string;
    message: string;
  }>;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T[];
  result: {
    pagination: Pagination;
    search: Search;
  };
}

@Injectable({
  providedIn: 'root',
})
export class AdminApiService {
  private readonly http = inject(HttpClient);
  private readonly platformId = inject(PLATFORM_ID);
  private readonly apiUrl = environment.apiUrl;

  // ===== USER APIs =====

  getUsers(params: PaginationParams = {}): Observable<PaginatedResponse<User>> {
    if (!isPlatformBrowser(this.platformId)) {
      return of({ success: false, data: [], result: { pagination: {} as Pagination, search: {} as Search } });
    }

    let httpParams = new HttpParams();
    if (params.page) httpParams = httpParams.set('page', params.page.toString());
    if (params.per_page) httpParams = httpParams.set('per_page', params.per_page.toString());
    if (params.order_by) httpParams = httpParams.set('order_by', params.order_by);
    if (params.sort_by) httpParams = httpParams.set('sort_by', params.sort_by);
    if (params.keyword) httpParams = httpParams.set('keyword', params.keyword);
    if (params.column) httpParams = httpParams.set('column', params.column);

    return this.http.get<PaginatedResponse<User>>(`${this.apiUrl}/api/v1/users`, { params: httpParams });
  }

  getUser(id: number): Observable<ApiResponse<User>> {
    return this.http.get<ApiResponse<User>>(`${this.apiUrl}/api/v1/users/${id}`);
  }

  createUser(data: CreateUserDto): Observable<ApiResponse<User>> {
    return this.http.post<ApiResponse<User>>(`${this.apiUrl}/api/v1/users`, data);
  }

  updateUser(id: number, data: UpdateUserDto): Observable<ApiResponse<null>> {
    return this.http.patch<ApiResponse<null>>(`${this.apiUrl}/api/v1/users/${id}`, data);
  }

  deleteUser(id: number): Observable<ApiResponse<null>> {
    return this.http.delete<ApiResponse<null>>(`${this.apiUrl}/api/v1/users/${id}`);
  }

  // ===== ROLE APIs =====

  getRoles(params: PaginationParams = {}): Observable<PaginatedResponse<Role>> {
    if (!isPlatformBrowser(this.platformId)) {
      return of({ success: false, data: [], result: { pagination: {} as Pagination, search: {} as Search } });
    }

    let httpParams = new HttpParams();
    if (params.page) httpParams = httpParams.set('page', params.page.toString());
    if (params.per_page) httpParams = httpParams.set('per_page', params.per_page.toString());

    return this.http.get<PaginatedResponse<Role>>(`${this.apiUrl}/api/v1/roles`, { params: httpParams });
  }

  getRole(id: number): Observable<ApiResponse<Role>> {
    return this.http.get<ApiResponse<Role>>(`${this.apiUrl}/api/v1/roles/${id}`);
  }

  // ===== USER ROLE APIs =====

  getUserRoles(userId: number): Observable<ApiResponse<UserRole[]>> {
    return this.http.get<ApiResponse<UserRole[]>>(`${this.apiUrl}/api/v1/user/roles/user/${userId}`);
  }

  assignRole(userId: number, roleId: number): Observable<ApiResponse<UserRole>> {
    return this.http.post<ApiResponse<UserRole>>(`${this.apiUrl}/api/v1/user/roles`, {
      user_id: userId,
      role_id: roleId,
    });
  }

  removeUserRole(userRoleId: number): Observable<ApiResponse<null>> {
    return this.http.delete<ApiResponse<null>>(`${this.apiUrl}/api/v1/user/roles/${userRoleId}`);
  }

  removeAllUserRoles(userId: number): Observable<ApiResponse<null>> {
    return this.http.delete<ApiResponse<null>>(`${this.apiUrl}/api/v1/user/roles/user/${userId}`);
  }
}
