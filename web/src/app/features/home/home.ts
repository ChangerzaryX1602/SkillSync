import { ChangeDetectionStrategy, Component, inject, signal, OnInit } from '@angular/core';
import { RouterLink } from '@angular/router';
import { AuthService } from '../../core/services/auth.service';
import { UserService } from '../../core/services/user.service';
import { NotificationService } from '../../core/services/notification.service';

@Component({
  selector: 'app-home',
  imports: [RouterLink],
  templateUrl: './home.html',
  styleUrl: './home.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class HomeComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly userService = inject(UserService);
  protected readonly notificationService = inject(NotificationService);


  protected readonly isSidebarOpen = signal(true);
  protected readonly isAdminMenuOpen = signal(false);

  // Expose user service signals
  protected readonly currentUser = this.userService.currentUser;
  protected readonly isAdmin = this.userService.isAdmin;
  protected readonly isLoading = this.userService.isLoading;

  protected readonly skills = [
    { name: 'TypeScript', level: 85, color: '#3178c6' },
    { name: 'Angular', level: 78, color: '#dd0031' },
    { name: 'Go', level: 72, color: '#00add8' },
    { name: 'PostgreSQL', level: 68, color: '#336791' },
    { name: 'Docker', level: 65, color: '#2496ed' },
  ];

  protected readonly courses = [
    {
      id: 1,
      title: 'Advanced Angular Patterns',
      instructor: 'John Doe',
      progress: 65,
      thumbnail: 'https://picsum.photos/seed/angular/300/200',
      category: 'Frontend',
    },
    {
      id: 2,
      title: 'Go Microservices',
      instructor: 'Jane Smith',
      progress: 40,
      thumbnail: 'https://picsum.photos/seed/golang/300/200',
      category: 'Backend',
    },
    {
      id: 3,
      title: 'System Design Fundamentals',
      instructor: 'Alex Johnson',
      progress: 20,
      thumbnail: 'https://picsum.photos/seed/system/300/200',
      category: 'Architecture',
    },
  ];

  protected readonly stats = [
    { label: 'Courses Enrolled', value: 12, icon: '📚' },
    { label: 'Skills Mastered', value: 8, icon: '🎯' },
    { label: 'Hours Learned', value: 156, icon: '⏱️' },
    { label: 'Certificates', value: 5, icon: '🏆' },
  ];

  ngOnInit(): void {
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
