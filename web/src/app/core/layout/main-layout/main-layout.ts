import { ChangeDetectionStrategy, Component, inject, signal, OnInit } from '@angular/core';
import { RouterLink, RouterOutlet, RouterLinkActive } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { UserService } from '../../services/user.service';
import { NotificationService } from '../../services/notification.service';

@Component({
  selector: 'app-main-layout',
  imports: [RouterLink, RouterOutlet, RouterLinkActive],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MainLayoutComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly userService = inject(UserService);
  protected readonly notificationService = inject(NotificationService);

  protected readonly isSidebarOpen = signal(true);
  protected readonly isAdminMenuOpen = signal(false);

  // Expose user service signals
  protected readonly currentUser = this.userService.currentUser;
  protected readonly isAdmin = this.userService.isAdmin;

  ngOnInit(): void {
    // Ensure user data is fetched if not already valid
    this.userService.fetchCurrentUser().subscribe();
  }

  protected toggleSidebar(): void {
    this.isSidebarOpen.update((open) => !open);
  }

  protected toggleAdminMenu(): void {
    this.isAdminMenuOpen.update((open) => !open);
  }

  protected logout(): void {
    this.userService.clearUser();
    this.authService.logout();
  }

  protected dismissNotification(id: number): void {
    this.notificationService.removeNotification(id);
  }
}
