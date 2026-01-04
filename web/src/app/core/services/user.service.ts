import { Injectable, inject, signal, computed, PLATFORM_ID } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { isPlatformBrowser } from '@angular/common';
import { Observable, tap, catchError, of } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface UserRole {
  id: number;
  user_id: number;
  role_id: number;
  created_at: string;
  updated_at: string;
}

export interface Role {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
  updated_at: string;
  roles?: Role[];
}

interface ApiResponse<T> {
  success: boolean;
  data: T;
}

interface UserRolesResponse {
  success: boolean;
  data: Array<{
    id: number;
    user_id: number;
    role_id: number;
  }>;
}

@Injectable({
  providedIn: 'root',
})
export class UserService {
  private readonly http = inject(HttpClient);
  private readonly platformId = inject(PLATFORM_ID);

  private readonly _currentUser = signal<User | null>(null);
  private readonly _userRoles = signal<Role[]>([]);
  private readonly _isLoading = signal(false);

  readonly currentUser = this._currentUser.asReadonly();
  readonly userRoles = this._userRoles.asReadonly();
  readonly isLoading = this._isLoading.asReadonly();

  readonly isAdmin = computed(() => {
    const roles = this._userRoles();
    return roles.some((role) => role.name === 'Admin');
  });

  fetchCurrentUser(): Observable<User | null> {
    if (!isPlatformBrowser(this.platformId)) {
      return of(null);
    }

    this._isLoading.set(true);

    return this.http.get<ApiResponse<User>>(`${environment.apiUrl}/api/v1/users/me`).pipe(
      tap((response) => {
        if (response.success && response.data) {
          this._currentUser.set(response.data);
          this.fetchUserRoles(response.data.id);
        }
        this._isLoading.set(false);
      }),
      catchError(() => {
        this._isLoading.set(false);
        this._currentUser.set(null);
        return of(null);
      })
    ) as Observable<User | null>;
  }

  private fetchUserRoles(userId: number): void {
    this.http
      .get<UserRolesResponse>(`${environment.apiUrl}/api/v1/user/roles/user/${userId}`)
      .subscribe({
        next: (response) => {
          if (response.success && response.data && response.data.length > 0) {
            const roleIds = response.data.map((ur) => ur.role_id);
            this.fetchRoleDetails(roleIds);
          }
        },
        error: () => {
          this._userRoles.set([]);
        },
      });
  }

  private fetchRoleDetails(roleIds: number[]): void {
    const roles: Role[] = [];
    let completed = 0;

    for (const roleId of roleIds) {
      this.http
        .get<{ success: boolean; data: Role }>(`${environment.apiUrl}/api/v1/roles/${roleId}`)
        .subscribe({
          next: (response) => {
            if (response.success && response.data) {
              roles.push(response.data);
            }
            completed++;
            if (completed === roleIds.length) {
              this._userRoles.set(roles);
            }
          },
          error: () => {
            completed++;
            if (completed === roleIds.length) {
              this._userRoles.set(roles);
            }
          },
        });
    }
  }

  clearUser(): void {
    this._currentUser.set(null);
    this._userRoles.set([]);
  }
}
