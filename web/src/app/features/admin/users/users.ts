import { ChangeDetectionStrategy, Component, inject, signal, OnInit, computed } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterLink } from '@angular/router';
import {
  AdminApiService,
  User,
  Role,
  UserRole,
  PaginationParams,
  Pagination,
} from '../../../core/services/admin-api.service';
import { NotificationService } from '../../../core/services/notification.service';
import { DialogComponent } from '../../../shared/dialog/dialog';
import { SearchBoxComponent } from '../../../shared/search-box/search-box';
import { PaginationComponent } from '../../../shared/pagination/pagination';

@Component({
  selector: 'app-users',
  imports: [ReactiveFormsModule, RouterLink, DialogComponent, SearchBoxComponent, PaginationComponent],
  templateUrl: './users.html',
  styleUrl: './users.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UsersComponent implements OnInit {
  private readonly adminApi = inject(AdminApiService);
  private readonly notification = inject(NotificationService);
  private readonly fb = inject(FormBuilder);

  // State
  protected readonly users = signal<User[]>([]);
  protected readonly roles = signal<Role[]>([]);
  protected readonly pagination = signal<Pagination | null>(null);
  protected readonly isLoading = signal(false);
  protected readonly searchKeyword = signal('');
  protected readonly perPage = signal(10);
  protected readonly perPageOptions = [5, 10, 20, 50];

  // Dialog states
  protected readonly isCreateDialogOpen = signal(false);
  protected readonly isEditDialogOpen = signal(false);
  protected readonly isDeleteDialogOpen = signal(false);
  protected readonly isRoleDialogOpen = signal(false);
  protected readonly selectedUser = signal<User | null>(null);
  protected readonly selectedUserRoles = signal<UserRole[]>([]);
  protected readonly isSubmitting = signal(false);

  // Forms
  protected readonly createForm = this.fb.nonNullable.group({
    username: ['', [Validators.required, Validators.minLength(3)]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
  });

  protected readonly editForm = this.fb.nonNullable.group({
    username: ['', [Validators.required, Validators.minLength(3)]],
    email: ['', [Validators.required, Validators.email]],
  });

  // Computed
  protected readonly currentPage = computed(() => this.pagination()?.page ?? 1);
  protected readonly totalPages = computed(() => {
    const p = this.pagination();
    if (!p || !p.total || !p.per_page) return 1;
    return Math.ceil(p.total / p.per_page);
  });

  ngOnInit(): void {
    this.loadUsers();
    this.loadRoles();
  }

  protected loadUsers(params: PaginationParams = {}): void {
    this.isLoading.set(true);
    const searchParams: PaginationParams = {
      page: params.page ?? 1,
      per_page: params.per_page ?? this.perPage(),
      ...params,
    };

    if (this.searchKeyword()) {
      searchParams.keyword = this.searchKeyword();
    }

    this.adminApi.getUsers(searchParams).subscribe({
      next: (response) => {
        if (response.success) {
          this.users.set(response.data);
          this.pagination.set(response.result.pagination);
        }
        this.isLoading.set(false);
      },
      error: () => {
        this.isLoading.set(false);
        this.notification.showError('ไม่สามารถโหลดข้อมูลผู้ใช้งานได้');
      },
    });
  }

  protected loadRoles(): void {
    this.adminApi.getRoles({ per_page: 999999 }).subscribe({
      next: (response) => {
        if (response.success) {
          this.roles.set(response.data);
        }
      },
      error: (err) => {
        if (err.status === 403) {
          this.notification.showError('คุณไม่มี Permission เพื่อทำสิ่งนี้');
        }
      },
    });
  }

  protected onSearch(keyword: string): void {
    this.searchKeyword.set(keyword);
    this.loadUsers({ page: 1 });
  }

  protected onPageChange(page: number): void {
    this.loadUsers({ page });
  }

  protected onPerPageChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    const newPerPage = parseInt(select.value, 10);
    this.perPage.set(newPerPage);
    this.loadUsers({ page: 1, per_page: newPerPage });
  }

  // Create User
  protected openCreateDialog(): void {
    this.createForm.reset();
    this.isCreateDialogOpen.set(true);
  }

  protected closeCreateDialog(): void {
    this.isCreateDialogOpen.set(false);
  }

  protected submitCreate(): void {
    if (this.createForm.invalid) {
      this.createForm.markAllAsTouched();
      return;
    }

    this.isSubmitting.set(true);
    const data = this.createForm.getRawValue();

    this.adminApi.createUser(data).subscribe({
      next: (response) => {
        if (response.success) {
          this.notification.showSuccess('สร้างผู้ใช้งานสำเร็จ');
          this.closeCreateDialog();
          this.loadUsers();
        }
        this.isSubmitting.set(false);
      },
      error: () => {
        this.isSubmitting.set(false);
      },
    });
  }

  // Edit User
  protected openEditDialog(user: User): void {
    this.selectedUser.set(user);
    this.editForm.patchValue({
      username: user.username,
      email: user.email,
    });
    this.isEditDialogOpen.set(true);
  }

  protected closeEditDialog(): void {
    this.isEditDialogOpen.set(false);
    this.selectedUser.set(null);
  }

  protected submitEdit(): void {
    if (this.editForm.invalid || !this.selectedUser()) {
      this.editForm.markAllAsTouched();
      return;
    }

    this.isSubmitting.set(true);
    const data = this.editForm.getRawValue();
    const userId = this.selectedUser()!.id;

    this.adminApi.updateUser(userId, data).subscribe({
      next: (response) => {
        if (response.success) {
          this.notification.showSuccess('แก้ไขผู้ใช้งานสำเร็จ');
          this.closeEditDialog();
          this.loadUsers({ page: this.currentPage() });
        }
        this.isSubmitting.set(false);
      },
      error: () => {
        this.isSubmitting.set(false);
      },
    });
  }

  // Delete User
  protected openDeleteDialog(user: User): void {
    this.selectedUser.set(user);
    this.isDeleteDialogOpen.set(true);
  }

  protected closeDeleteDialog(): void {
    this.isDeleteDialogOpen.set(false);
    this.selectedUser.set(null);
  }

  protected confirmDelete(): void {
    if (!this.selectedUser()) return;

    this.isSubmitting.set(true);
    const userId = this.selectedUser()!.id;

    this.adminApi.deleteUser(userId).subscribe({
      next: (response) => {
        if (response.success) {
          this.notification.showSuccess('ลบผู้ใช้งานสำเร็จ');
          this.closeDeleteDialog();
          this.loadUsers({ page: this.currentPage() });
        }
        this.isSubmitting.set(false);
      },
      error: () => {
        this.isSubmitting.set(false);
      },
    });
  }

  // Role Management
  protected openRoleDialog(user: User): void {
    this.selectedUser.set(user);
    this.loadUserRoles(user.id);
    this.isRoleDialogOpen.set(true);
  }

  protected closeRoleDialog(): void {
    this.isRoleDialogOpen.set(false);
    this.selectedUser.set(null);
    this.selectedUserRoles.set([]);
  }

  protected loadUserRoles(userId: number): void {
    this.adminApi.getUserRoles(userId).subscribe({
      next: (response) => {
        if (response.success) {
          this.selectedUserRoles.set(response.data);
        }
      },
    });
  }

  protected hasRole(roleId: number): boolean {
    return this.selectedUserRoles().some((ur) => ur.role_id === roleId);
  }

  protected toggleRole(roleId: number): void {
    const userId = this.selectedUser()?.id;
    if (!userId) return;

    const existingUserRole = this.selectedUserRoles().find((ur) => ur.role_id === roleId);

    if (existingUserRole) {
      this.adminApi.removeUserRole(existingUserRole.id).subscribe({
        next: (response) => {
          if (response.success) {
            this.notification.showSuccess('ลบบทบาทสำเร็จ');
            // Reload to ensure correct state
            this.loadUserRoles(userId);
          }
        },
        error: (err) => {
          if (err.status === 403) {
            this.notification.showError('คุณไม่มี Permission เพื่อทำสิ่งนี้');
          }
        },
      });
    } else {
      this.adminApi.assignRole(userId, roleId).subscribe({
        next: (response) => {
          if (response.success) {
            this.notification.showSuccess('เพิ่มบทบาทสำเร็จ');
            // Reload to get the new UserRole with correct ID
            this.loadUserRoles(userId);
          }
        },
        error: (err) => {
          if (err.status === 403) {
            this.notification.showError('คุณไม่มี Permission เพื่อทำสิ่งนี้');
          }
        },
      });
    }
  }

  protected formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString('th-TH', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  }
}
