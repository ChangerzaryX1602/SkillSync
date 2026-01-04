import { ChangeDetectionStrategy, Component, input, output, computed } from '@angular/core';

@Component({
  selector: 'app-pagination',
  template: `
    <div class="pagination">
      <button
        type="button"
        class="btn-page btn-arrow"
        [disabled]="currentPage() === 1"
        (click)="pageChange.emit(currentPage() - 1)"
        aria-label="Previous page"
      >
        ←
      </button>

      @for (page of visiblePages(); track page) {
        @if (page === -1) {
          <span class="ellipsis">...</span>
        } @else {
          <button
            type="button"
            class="btn-page"
            [class.active]="page === currentPage()"
            (click)="pageChange.emit(page)"
            [attr.aria-label]="'Page ' + page"
            [attr.aria-current]="page === currentPage() ? 'page' : null"
          >
            {{ page }}
          </button>
        }
      }

      <button
        type="button"
        class="btn-page btn-arrow"
        [disabled]="currentPage() === totalPages()"
        (click)="pageChange.emit(currentPage() + 1)"
        aria-label="Next page"
      >
        →
      </button>
    </div>
  `,
  styles: `
    .pagination {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 0.5rem;
    }

    .btn-page {
      min-width: 36px;
      height: 36px;
      padding: 0 0.5rem;
      border-radius: 8px;
      border: 1px solid rgba(255, 255, 255, 0.08);
      background: rgba(255, 255, 255, 0.03);
      color: rgba(255, 255, 255, 0.8);
      font-size: 0.875rem;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s ease;

      &:hover:not(:disabled):not(.active) {
        background: rgba(255, 255, 255, 0.08);
        color: #ffffff;
      }

      &:disabled {
        opacity: 0.4;
        cursor: not-allowed;
      }

      &.active {
        background: linear-gradient(135deg, #6366f1, #8b5cf6);
        border-color: transparent;
        color: #ffffff;
      }

      &.btn-arrow {
        font-weight: 600;
      }
    }

    .ellipsis {
      color: rgba(255, 255, 255, 0.5);
      padding: 0 0.25rem;
    }
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class PaginationComponent {
  readonly currentPage = input.required<number>();
  readonly totalPages = input.required<number>();
  readonly pageChange = output<number>();

  protected readonly visiblePages = computed(() => {
    const current = this.currentPage();
    const total = this.totalPages();
    const pages: number[] = [];

    if (total <= 7) {
      // Show all pages if 7 or less
      for (let i = 1; i <= total; i++) {
        pages.push(i);
      }
    } else {
      // Always show first page
      pages.push(1);

      if (current > 3) {
        pages.push(-1); // ellipsis
      }

      // Show pages around current
      const start = Math.max(2, current - 1);
      const end = Math.min(total - 1, current + 1);

      for (let i = start; i <= end; i++) {
        pages.push(i);
      }

      if (current < total - 2) {
        pages.push(-1); // ellipsis
      }

      // Always show last page
      pages.push(total);
    }

    return pages;
  });
}
