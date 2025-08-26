import { CommonModule } from '@angular/common';
import {
    Component,
    ElementRef,
    EventEmitter,
    HostListener,
    Input,
    OnDestroy,
    OnInit,
    Output,
    ViewChild,
} from '@angular/core';
import { Therapist } from '../../models/swipe-card';

@Component({
  selector: 'app-swipe-card',
  templateUrl: './swipe-card.html',
  styleUrls: ['./swipe-card.scss'],
  imports: [CommonModule],
})
export class SwipeCardComponent implements OnInit, OnDestroy {
  @Input() therapist: any = {};
  @Input() index = 0;
  @Output() swiped = new EventEmitter<{
    direction: 'left' | 'right';
    therapist: Therapist;
  }>();

  @ViewChild('card', { static: true }) cardRef!: ElementRef<HTMLElement>;

  tabs = [
    { key: 'specialties', label: 'Specialties' },
    { key: 'pricing', label: 'Pricing' },
    { key: 'availability', label: 'Availability' },
    { key: 'languages', label: 'Languages' },
    { key: 'modes', label: 'Modes' },
  ];
  activeTab = 'specialties';

  private pointerId: number | null = null;
  private startX = 0;
  private startY = 0;
  private currentX = 0;
  private currentY = 0;
  private dragging = false;
  transform = '';
  transition = '';
  likeOpacity = 0;
  nopeOpacity = 0;

  private readonly swipeThreshold = 120;
  private rafId = 0;

  ngOnInit() {
    this.activeTab = 'specialties';
  }

  ngOnDestroy() {
    cancelAnimationFrame(this.rafId);
  }

  onPointerDown(ev: PointerEvent) {
    if (this.pointerId !== null) return;
    this.pointerId = ev.pointerId;
    (ev.target as Element).setPointerCapture(this.pointerId);
    this.startX = ev.clientX;
    this.startY = ev.clientY;
    this.currentX = this.startX;
    this.currentY = this.startY;
    this.dragging = true;
    this.transition = '';
  }

  onPointerMove(ev: PointerEvent) {
    if (!this.dragging || ev.pointerId !== this.pointerId) return;
    this.currentX = ev.clientX;
    this.currentY = ev.clientY;
    this.scheduleUpdate();
  }

  onPointerUp(ev: PointerEvent) {
    if (!this.dragging || ev.pointerId !== this.pointerId) return;
    this.dragging = false;
    if (this.pointerId !== null) {
      try {
        (ev.target as Element).releasePointerCapture(this.pointerId);
      } catch {}
      this.pointerId = null;
    }

    const dx = this.currentX - this.startX;
    if (Math.abs(dx) >= this.swipeThreshold) {
      const dir: 'left' | 'right' = dx > 0 ? 'right' : 'left';
      this.animateOffScreen(dir);
    } else {
      this.resetPosition();
    }
  }

  @HostListener('keydown', ['$event'])
  onKeydown(ev: KeyboardEvent) {
    if (ev.key === 'ArrowLeft') {
      this.programmaticSwipe('left');
    } else if (ev.key === 'ArrowRight') {
      this.programmaticSwipe('right');
    }
  }

  private scheduleUpdate() {
    if (this.rafId) return;
    this.rafId = requestAnimationFrame(() => {
      this.rafId = 0;
      this.updateTransformFromPointer();
    });
  }

  private updateTransformFromPointer() {
    const dx = this.currentX - this.startX;
    const dy = this.currentY - this.startY;
    const rotate = Math.sign(dx) * Math.min(Math.abs(dx) / 12, 20);
    this.transform = `translate(${dx}px, ${dy * 0.4}px) rotate(${rotate}deg)`;

    const norm = Math.min(Math.abs(dx) / this.swipeThreshold, 1);
    if (dx > 0) {
      this.likeOpacity = norm;
      this.nopeOpacity = 0;
    } else {
      this.nopeOpacity = norm;
      this.likeOpacity = 0;
    }
  }

  private resetPosition() {
    this.transition = 'transform 250ms cubic-bezier(.2,.9,.2,1), opacity 200ms';
    this.transform = '';
    this.likeOpacity = 0;
    this.nopeOpacity = 0;
  }

  private animateOffScreen(direction: 'left' | 'right') {
    const width = window.innerWidth || document.documentElement.clientWidth;
    const offX = (direction === 'right' ? 1 : -1) * (width * 1.1);
    const rotate = direction === 'right' ? 20 : -20;
    this.transition = 'transform 300ms cubic-bezier(.2,.9,.2,1), opacity 300ms';
    this.transform = `translate(${offX}px, 0px) rotate(${rotate}deg)`;
    this.likeOpacity = direction === 'right' ? 1 : 0;
    this.nopeOpacity = direction === 'left' ? 1 : 0;

    setTimeout(() => {
      this.swiped.emit({
        direction,
        therapist: this.therapist,
      });

      this.transform = '';
      this.transition = '';
      this.likeOpacity = 0;
      this.nopeOpacity = 0;
    }, 320);
  }

  programmaticSwipe(dir: 'left' | 'right') {
    this.animateOffScreen(dir);
  }

  pass() {
    this.programmaticSwipe('left');
  }

  match() {
    this.programmaticSwipe('right');
  }
}
