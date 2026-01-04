import { ChangeDetectionStrategy, Component, input, output, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-search-box',
  imports: [FormsModule],
  template: `
    <div class="search-box">
      <input
        type="text"
        [placeholder]="placeholder()"
        [(ngModel)]="searchValue"
        (keydown.enter)="onSearch()"
        [attr.aria-label]="placeholder()"
      />
      <button type="button" class="search-btn" (click)="onSearch()" aria-label="Search">
        🔍
      </button>
    </div>
  `,
  styles: `
    .search-box {
      display: flex;
      align-items: center;
      max-width: 400px;
      background: rgba(255, 255, 255, 0.03);
      border: 1px solid rgba(255, 255, 255, 0.08);
      border-radius: 10px;
      overflow: hidden;
      transition: border-color 0.2s ease;
      margin-bottom: 10px;
      &:focus-within {
        border-color: #6366f1;
      }

      input {
        flex: 1;
        padding: 0.75rem 1rem;
        border: none;
        background: transparent;
        color: #ffffff;
        font-size: 0.875rem;

        &::placeholder {
          color: rgba(255, 255, 255, 0.5);
        }

        &:focus {
          outline: none;
        }
      }

      .search-btn {
        padding: 0.75rem 1rem;
        border: none;
        background: linear-gradient(135deg, #6366f1, #8b5cf6);
        color: white;
        font-size: 1rem;
        cursor: pointer;
        transition: opacity 0.2s ease;

        &:hover {
          opacity: 0.9;
        }
      }
    }
  `,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SearchBoxComponent {
  readonly placeholder = input<string>('ค้นหา...');
  readonly search = output<string>();

  protected searchValue = '';

  protected onSearch(): void {
    this.search.emit(this.searchValue);
  }

  // Allow parent to clear search
  clear(): void {
    this.searchValue = '';
  }
}
