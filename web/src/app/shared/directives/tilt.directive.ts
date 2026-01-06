import { Directive, ElementRef, HostListener, Input, Renderer2 } from '@angular/core';

@Directive({
  selector: '[appTilt]',
  standalone: true
})
export class TiltDirective {
  @Input() maxTilt = 10; // Maximum tilt angle
  @Input() perspective = 1000; // Perspective depth
  @Input() scale = 1.05; // Scale on hover

  private el: HTMLElement;

  constructor(private elementRef: ElementRef, private renderer: Renderer2) {
    this.el = this.elementRef.nativeElement;
    this.renderer.setStyle(this.el, 'transition', 'transform 0.1s ease-out');
    this.renderer.setStyle(this.el, 'transform-style', 'preserve-3d');
  }

  @HostListener('mousemove', ['$event'])
  onMouseMove(event: MouseEvent) {
    const rect = this.el.getBoundingClientRect();
    const width = rect.width;
    const height = rect.height;
    
    // Calculate mouse position relative to element center
    const centerX = rect.left + width / 2;
    const centerY = rect.top + height / 2;
    const mouseX = event.clientX - centerX;
    const mouseY = event.clientY - centerY;

    // Calculate rotation (-1 to 1) * maxTilt
    const rotateX = -1 * (mouseY / (height / 2)) * this.maxTilt;
    const rotateY = (mouseX / (width / 2)) * this.maxTilt;

    this.renderer.setStyle(this.el, 'transform', 
      `perspective(${this.perspective}px) rotateX(${rotateX}deg) rotateY(${rotateY}deg) scale3d(${this.scale}, ${this.scale}, ${this.scale})`
    );
    
    // Optional: Add a glare effect via CSS variable if wanted, 
    // or just rely on the 3D rotation for now.
    // Let's add interaction classes if needed.
  }

  @HostListener('mouseleave')
  onMouseLeave() {
    this.renderer.setStyle(this.el, 'transition', 'transform 0.5s ease-out');
    this.renderer.setStyle(this.el, 'transform', 
      `perspective(${this.perspective}px) rotateX(0deg) rotateY(0deg) scale3d(1, 1, 1)`
    );
  }
  
  @HostListener('mouseenter')
  onMouseEnter() {
    // Reset transition for quick response
    this.renderer.setStyle(this.el, 'transition', 'transform 0.1s ease-out');
  }
}
