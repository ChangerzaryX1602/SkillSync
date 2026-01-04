import { ChangeDetectionStrategy, Component, inject, signal, NgZone } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router, RouterLink } from '@angular/router';
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
  selector: 'app-register',
  imports: [ReactiveFormsModule, RouterLink],
  templateUrl: './register.html',
  styleUrl: './register.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RegisterComponent {
  private readonly fb = inject(FormBuilder);
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly ngZone = inject(NgZone);

  protected readonly isLoading = signal(false);
  protected readonly errorMessage = signal<string | null>(null);
  protected readonly successMessage = signal<string | null>(null);

  protected readonly registerForm = this.fb.nonNullable.group({
    username: ['', [Validators.required, Validators.minLength(3)]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
  });

  protected onSubmit(): void {
    if (this.registerForm.invalid) {
      this.registerForm.markAllAsTouched();
      return;
    }

    this.isLoading.set(true);
    this.errorMessage.set(null);
    this.successMessage.set(null);

    const { email, password, username } = this.registerForm.getRawValue();

    this.authService.register(email, password, username).subscribe({
      next: () => {
        this.isLoading.set(false);
        this.successMessage.set('Account created successfully! Redirecting to login...');
        setTimeout(() => {
          this.ngZone.run(() => {
            this.router.navigate(['/login']);
          });
        }, 1500);
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

      // Handle 409 Conflict - duplicate email/username
      if (firstError.code === 409 || firstError.title === 'Conflict') {
        return firstError.message;
      }

      // Handle Bad Request
      if (firstError.title === 'Bad Request') {
        return firstError.message;
      }

      return `${firstError.title}: ${firstError.message}`;
    }

    // Fallback for non-structured errors
    switch (error.status) {
      case 400:
        return 'Please check your input and try again.';
      case 409:
        return 'Email or username already exists.';
      case 500:
        return 'Server error. Please try again later.';
      default:
        return 'An unexpected error occurred. Please try again.';
    }
  }
}
