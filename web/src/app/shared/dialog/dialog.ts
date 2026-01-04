import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';

@Component({
  selector: 'app-dialog',
  template: `
    @if (isOpen()) {
      <div class="dialog-overlay" (click)="onOverlayClick($event)">
        <div class="dialog" role="dialog" [attr.aria-labelledby]="'dialog-title'">
          <div class="dialog-header">
            <h2 id="dialog-title">{{ title() }}</h2>
            <button
              type="button"
              class="close-btn"
              (click)="close.emit()"
              aria-label="Close dialog"
            >
              ✕
            </button>
          </div>
          <div class="dialog-content">
            <ng-content></ng-content>
          </div>
        </div>
      </div>
    }
  `,
  styles: `
    .dialog-overlay {
      position: fixed;
      inset: 0;
      background: rgba(0, 0, 0, 0.6);
      backdrop-filter: blur(4px);
      display: flex;
      align-items: center;
      justify-content: center;
      z-index: 1000;
      animation: fadeIn 0.2s ease;
    }

    @keyframes fadeIn {
      from { opacity: 0; }
      to { opacity: 1; }
    }

    .dialog {
      background: rgba(30, 30, 60, 0.95);
      border: 1px solid rgba(255, 255, 255, 0.1);
      border-radius: 16px;
      min-width: 400px;
      max-width: 90vw;
      max-height: 90vh;
      overflow: hidden;
      box-shadow: 0 24px 48px rgba(0, 0, 0, 0.4);
      animation: slideUp 0.2s ease;
    }

    @keyframes slideUp {
      from {
        opacity: 0;
        transform: translateY(20px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }

    .dialog-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 1.25rem 1.5rem;
      border-bottom: 1px solid rgba(255, 255, 255, 0.08);

      h2 {
        margin: 0;
        font-size: 1.125rem;
        font-weight: 600;
        color: #ffffff;
      }
    }

    .close-btn {
      width: 32px;
      height: 32px;
      border-radius: 8px;
      border: 1px solid rgba(255, 255, 255, 0.1);
      background: transparent;
      color: rgba(255, 255, 255, 0.6);
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 0.875rem;
      transition: all 0.2s ease;

      &:hover {
        background: rgba(255, 255, 255, 0.1);
        color: #ffffff;
      }
    }

    .dialog-content {
      padding: 1.5rem;
      overflow-y: auto;
      max-height: calc(90vh - 80px);
    }
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DialogComponent {
  readonly isOpen = input.required<boolean>();
  readonly title = input<string>('');
  readonly close = output<void>();

  protected onOverlayClick(event: MouseEvent): void {
    if ((event.target as HTMLElement).classList.contains('dialog-overlay')) {
      this.close.emit();
    }
  }
}
