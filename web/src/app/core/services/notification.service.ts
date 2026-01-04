import { Injectable, signal, computed } from '@angular/core';

export interface Notification {
  id: number;
  type: 'success' | 'error' | 'warning' | 'info';
  message: string;
  duration?: number;
}

@Injectable({
  providedIn: 'root',
})
export class NotificationService {
  private readonly _notifications = signal<Notification[]>([]);
  private notificationId = 0;

  readonly notifications = this._notifications.asReadonly();

  showSuccess(message: string, duration = 5000): void {
    this.addNotification({ type: 'success', message, duration });
  }

  showError(message: string, duration = 7000): void {
    this.addNotification({ type: 'error', message, duration });
  }

  showWarning(message: string, duration = 5000): void {
    this.addNotification({ type: 'warning', message, duration });
  }

  showInfo(message: string, duration = 5000): void {
    this.addNotification({ type: 'info', message, duration });
  }

  private addNotification(notification: Omit<Notification, 'id'>): void {
    const id = ++this.notificationId;
    const newNotification: Notification = { ...notification, id };

    this._notifications.update((list) => [...list, newNotification]);

    // Auto-remove after duration
    if (notification.duration && notification.duration > 0) {
      setTimeout(() => {
        this.removeNotification(id);
      }, notification.duration);
    }
  }

  removeNotification(id: number): void {
    this._notifications.update((list) => list.filter((n) => n.id !== id));
  }

  clearAll(): void {
    this._notifications.set([]);
  }
}
