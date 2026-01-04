import {
  HttpInterceptorFn,
  HttpRequest,
  HttpHandlerFn,
  HttpErrorResponse,
} from '@angular/common/http';
import { inject } from '@angular/core';
import { catchError, switchMap, throwError } from 'rxjs';
import { AuthService } from '../services/auth.service';
import { NotificationService } from '../services/notification.service';

const AUTH_ENDPOINTS = ['/api/v1/auth/login', '/api/v1/auth/register', '/api/v1/auth/refresh'];

function isAuthEndpoint(url: string): boolean {
  return AUTH_ENDPOINTS.some((endpoint) => url.includes(endpoint));
}

function addAuthHeader(request: HttpRequest<unknown>, token: string): HttpRequest<unknown> {
  return request.clone({
    setHeaders: {
      Authorization: `Bearer ${token}`,
    },
  });
}

function extractErrorMessage(error: HttpErrorResponse): string {
  if (error.error?.errors?.length > 0) {
    return error.error.errors[0].message;
  }
  return error.message || 'An unexpected error occurred';
}

export const authInterceptor: HttpInterceptorFn = (
  request: HttpRequest<unknown>,
  next: HttpHandlerFn
) => {
  const authService = inject(AuthService);
  const notificationService = inject(NotificationService);

  // Skip auth header for auth endpoints
  if (isAuthEndpoint(request.url)) {
    return next(request);
  }

  const accessToken = authService.getAccessToken();

  if (accessToken) {
    request = addAuthHeader(request, accessToken);
  }

  return next(request).pipe(
    catchError((error: HttpErrorResponse) => {
      // Handle 401 Unauthorized - Token expired/invalid
      if (error.status === 401 && !isAuthEndpoint(request.url)) {
        return authService.refreshToken().pipe(
          switchMap((response) => {
            // Retry original request with new token
            const newRequest = addAuthHeader(request, response.access_token);
            return next(newRequest);
          }),
          catchError((refreshError) => {
            // Refresh failed, logout user
            authService.logout();
            return throwError(() => refreshError);
          })
        );
      }

      return throwError(() => error);
    })
  );
};
