import { ChangeDetectionStrategy, Component, inject, signal, OnInit } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink, ActivatedRoute } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { AuthService } from '../../../core/services/auth.service';

interface ApiError {
  code: number;
  source: string;
  title: string;
  message: string;
}

interface ApiErrorResponse {
  success: boolean;
  errors: ApiError[];
}

@Component({
  selector: 'app-login',
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './login.html',
  styleUrl: './login.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginComponent implements OnInit {
  private readonly fb = inject(FormBuilder);
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);

  protected readonly isLoading = signal(false);
  protected readonly errorMessage = signal<string | null>(null);
  private returnUrl = '/';

  protected readonly loginForm = this.fb.nonNullable.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
  });

  ngOnInit(): void {
    // Get return url from query params or default to '/'
    this.returnUrl = this.route.snapshot.queryParams['returnUrl'] || '/';
  }

  protected onSubmit(): void {
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    this.isLoading.set(true);
    this.errorMessage.set(null);

    const { email, password } = this.loginForm.getRawValue();

    this.authService.login(email, password).subscribe({
      next: () => {
        this.router.navigateByUrl(this.returnUrl);
      },
      error: (error: HttpErrorResponse) => {
        this.isLoading.set(false);
        this.errorMessage.set(this.parseApiError(error));
      },
    });
  }

  private parseApiError(error: HttpErrorResponse): string {
    // Check if error response matches our API error format
    if (error.error?.errors?.length > 0) {
      const apiErrors = error.error as ApiErrorResponse;
      const firstError = apiErrors.errors[0];

      // Return user-friendly messages based on error title/code
      if (firstError.title === 'Unauthorized' || firstError.code === 401) {
        return 'Invalid email or password. Please try again.';
      }

      if (firstError.title === 'Bad Request') {
        return firstError.message;
      }

      if (firstError.title === 'Not Found') {
        return 'No account found with this email. Please register first.';
      }

      return `${firstError.title}: ${firstError.message}`;
    }

    // Fallback for non-structured errors
    switch (error.status) {
      case 400:
        return 'Invalid request. Please check your input.';
      case 401:
        return 'Invalid email or password.';
      case 404:
        return 'No account found with this email.';
      case 429:
        return 'Too many login attempts. Please try again later.';
      case 500:
        return 'Server error. Please try again later.';
      default:
        return 'An unexpected error occurred. Please try again.';
    }
  }
}
