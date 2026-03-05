import { Directive, ElementRef, HostListener, Input } from '@angular/core';

@Directive({
  selector: 'img[appImgFallback]',
  standalone: true,
})
export class ImgFallbackDirective {
  @Input() appImgFallback = '/default_avatar/7249ed41cf4ab50c5d37c169786ff768.jpg';

  private hasFailed = false;

  constructor(private el: ElementRef<HTMLImageElement>) {}

  @HostListener('error')
  onError() {
    if (this.hasFailed) return; 
    this.hasFailed = true;
    this.el.nativeElement.src = this.appImgFallback;
  }
}

