import { ChangeDetectionStrategy, Component, inject, signal, OnInit, computed } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { RouterLink } from '@angular/router';
import {
  AdminApiService,
  Role,
  Permission,
  RolePermission,
  PaginationParams,
  Pagination,
} from '../../../core/services/admin-api.service';
import { NotificationService } from '../../../core/services/notification.service';
import { DialogComponent } from '../../../shared/dialog/dialog';
import { SearchBoxComponent } from '../../../shared/search-box/search-box';
import { PaginationComponent } from '../../../shared/pagination/pagination';

@Component({
  selector: 'app-roles',
  imports: [ReactiveFormsModule, RouterLink, DialogComponent, SearchBoxComponent, PaginationComponent],
  templateUrl: './roles.html',
  styleUrl: './roles.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RolesComponent implements OnInit {
  private readonly adminApi = inject(AdminApiService);
  private readonly notification = inject(NotificationService);
  private readonly fb = inject(FormBuilder);

  // State
  protected readonly roles = signal<Role[]>([]);
  protected readonly pagination = signal<Pagination | null>(null);
  protected readonly isLoading = signal(false);
  protected readonly searchKeyword = signal('');
  protected readonly perPage = signal(10);
  protected readonly perPageOptions = [10, 20, 50, 100];

  // Dialog states
  protected readonly isCreateDialogOpen = signal(false);
  protected readonly isEditDialogOpen = signal(false);
  protected readonly isDeleteDialogOpen = signal(false);
  protected readonly isPermissionDialogOpen = signal(false);
  protected readonly selectedRole = signal<Role | null>(null);
  protected readonly isSubmitting = signal(false);

  // Permission states
  protected readonly permissions = signal<Permission[]>([]);
  protected readonly selectedRolePermissions = signal<RolePermission[]>([]);

  // Forms
  protected readonly createForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
  });

  protected readonly editForm = this.fb.nonNullable.group({
    name: ['', [Validators.required, Validators.minLength(2)]],
  });

  // Computed
  protected readonly currentPage = computed(() => this.pagination()?.page ?? 1);
  protected readonly totalPages = computed(() => {
    const p = this.pagination();
    if (!p || !p.total || !p.per_page) return 1;
    return Math.ceil(p.total / p.per_page);
  });

  ngOnInit(): void {
    this.loadRoles();
    this.loadPermissions();
  }

  protected loadRoles(params: PaginationParams = {}): void {
    this.isLoading.set(true);
    const searchParams: PaginationParams = {
      page: params.page ?? 1,
      per_page: params.per_page ?? this.perPage(),
      ...params,
    };

    if (this.searchKeyword()) {
      searchParams.keyword = this.searchKeyword();
    }

    this.adminApi.getRoles(searchParams).subscribe({
      next: (response) => {
        if (response.success) {
          this.roles.set(response.data);
          this.pagination.set(response.result.pagination);
        }
        this.isLoading.set(false);
      },
      error: (err) => {
        this.isLoading.set(false);
        if (err.status === 403) {
          this.notification.showError('คุณไม่มี Permission เพื่อทำสิ่งนี้');
        } else {
          this.notification.showError('ไม่สามารถโหลดข้อมูลบทบาทได้');
        }
      },
    });
  }

  protected loadPermissions(): void {
    this.adminApi.getPermissions({ per_page: 999 }).subscribe({
      next: (response) => {
        if (response.success) {
          this.permissions.set(response.data);
        }
      },
      error: (err) => {
        if (err.status === 403) {
          this.notification.showError('คุณไม่มี Permission เพื่อทำสิ่งนี้');
        }
      },
    });
  }

  protected loadRolePermissions(roleId: number): void {
    this.adminApi.getRolePermissions(roleId).subscribe({
      next: (response) => {
        if (response.success) {
          this.selectedRolePermissions.set(response.data || []);
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
    this.loadRoles({ page: 1 });
  }

  protected onPageChange(page: number): void {
    this.loadRoles({ page });
  }

  protected onPerPageChange(event: Event): void {
    const select = event.target as HTMLSelectElement;
    const newPerPage = parseInt(select.value, 10);
    this.perPage.set(newPerPage);
    this.loadRoles({ page: 1, per_page: newPerPage });
  }

  // Create Role
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

    this.adminApi.createRole(data).subscribe({
      next: (response) => {
        if (response.success) {
          this.notification.showSuccess('สร้างบทบาทสำเร็จ');
          this.closeCreateDialog();
          this.loadRoles();
        }
        this.isSubmitting.set(false);
      },
      error: (err) => {
        this.isSubmitting.set(false);
        if (err.status === 403) {
          this.notification.showError('คุณไม่มี Permission เพื่อทำสิ่งนี้');
        }
      },
    });
  }

  // Edit Role
  protected openEditDialog(role: Role): void {
    this.selectedRole.set(role);
    this.editForm.patchValue({
      name: role.name,
    });
    this.isEditDialogOpen.set(true);
  }

  protected closeEditDialog(): void {
    this.isEditDialogOpen.set(false);
    this.selectedRole.set(null);
  }

  protected submitEdit(): void {
    if (this.editForm.invalid || !this.selectedRole()) {
      this.editForm.markAllAsTouched();
      return;
    }

    this.isSubmitting.set(true);
    const data = this.editForm.getRawValue();
    const roleId = this.selectedRole()!.id;

    this.adminApi.updateRole(roleId, data).subscribe({
      next: (response) => {
        if (response.success) {
          this.notification.showSuccess('แก้ไขบทบาทสำเร็จ');
          this.closeEditDialog();
          this.loadRoles({ page: this.currentPage() });
        }
        this.isSubmitting.set(false);
      },
      error: (err) => {
        this.isSubmitting.set(false);
        if (err.status === 403) {
          this.notification.showError('คุณไม่มี Permission เพื่อทำสิ่งนี้');
        }
      },
    });
  }

  // Delete Role
  protected openDeleteDialog(role: Role): void {
    this.selectedRole.set(role);
    this.isDeleteDialogOpen.set(true);
  }

  protected closeDeleteDialog(): void {
    this.isDeleteDialogOpen.set(false);
    this.selectedRole.set(null);
  }

  protected confirmDelete(): void {
    if (!this.selectedRole()) return;

    this.isSubmitting.set(true);
    const roleId = this.selectedRole()!.id;

    this.adminApi.deleteRole(roleId).subscribe({
      next: (response) => {
        if (response.success) {
          this.notification.showSuccess('ลบบทบาทสำเร็จ');
          this.closeDeleteDialog();
          this.loadRoles({ page: this.currentPage() });
        }
        this.isSubmitting.set(false);
      },
      error: (err) => {
        this.isSubmitting.set(false);
        if (err.status === 403) {
          this.notification.showError('คุณไม่มี Permission เพื่อทำสิ่งนี้');
        }
      },
    });
  }

  // Permission Dialog
  protected openPermissionDialog(role: Role): void {
    this.selectedRole.set(role);
    this.loadRolePermissions(role.id);
    this.isPermissionDialogOpen.set(true);
  }

  protected closePermissionDialog(): void {
    this.isPermissionDialogOpen.set(false);
    this.selectedRole.set(null);
    this.selectedRolePermissions.set([]);
  }

  protected hasPermission(permissionId: number): boolean {
    return this.selectedRolePermissions().some((rp) => rp.permission_id === permissionId);
  }

  protected togglePermission(permissionId: number): void {
    const roleId = this.selectedRole()?.id;
    if (!roleId) return;

    const existingRolePermission = this.selectedRolePermissions().find((rp) => rp.permission_id === permissionId);

    if (existingRolePermission) {
      this.adminApi.removeRolePermission(existingRolePermission.id).subscribe({
        next: (response) => {
          if (response.success) {
            this.notification.showSuccess('ลบสิทธิสำเร็จ');
            this.loadRolePermissions(roleId);
          }
        },
        error: (err) => {
          if (err.status === 403) {
            this.notification.showError('คุณไม่มี Permission เพื่อทำสิ่งนี้');
          }
        },
      });
    } else {
      this.adminApi.assignPermission(roleId, permissionId).subscribe({
        next: (response) => {
          if (response.success) {
            this.notification.showSuccess('เพิ่มสิทธิสำเร็จ');
            this.loadRolePermissions(roleId);
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

  protected formatPermission(permission: Permission): string {
    return `${permission.group}:${permission.name}`;
  }

  protected formatDate(dateStr: string): string {
    return new Date(dateStr).toLocaleDateString('th-TH', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  }
}
