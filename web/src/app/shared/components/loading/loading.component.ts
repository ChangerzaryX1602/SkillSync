import { Component, Input, ViewEncapsulation } from '@angular/core';

@Component({
  selector: 'app-loading',
  standalone: true,
  template: `
    <div class="loader-container" [class.overlay]="overlay">
      <div class="reactor">
        <div class="reactor-core"></div>
        <div class="ring ring-outer"></div>
        <div class="ring ring-inner"></div>
        <div class="particles">
          <span></span><span></span><span></span><span></span>
        </div>
      </div>
      @if (text) {
        <div class="loading-text" [attr.data-text]="text">{{ text }}</div>
      }
    </div>
  `,
  styleUrls: ['./loading.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class LoadingComponent {
  @Input() overlay = false;
  @Input() text = 'INITIALIZING...';
}
