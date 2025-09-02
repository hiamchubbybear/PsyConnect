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
    SimpleChanges,
    ViewChild,
} from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { TherapistProfile } from '../profile-overlay/profile-overlay';

@Component({
  standalone: true,
  selector: 'app-swipe-card',
  templateUrl: './swipe-card.html',
  styleUrls: ['./swipe-card.scss'],
  imports: [CommonModule, TranslateModule],
})
export class SwipeCardComponent implements OnInit, OnDestroy {
  @Input() therapist: TherapistProfile = {} as TherapistProfile;
  @Input() index = 0;
  @Output() openOverlay = new EventEmitter<TherapistProfile>();
  @Input() isOverlayOpen: boolean = false;

  @Output() swiped = new EventEmitter<{
    direction: 'left' | 'right' | 'up';
    therapist: TherapistProfile;
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

  animateUp() {
    this.isOverlayOpen = !this.isOverlayOpen;

    this.transition = 'transform 300ms cubic-bezier(.2,.9,.2,1), opacity 300ms';

    if (this.isOverlayOpen) {
      const target = document.querySelector('.therapist-card__actions');
      if (target) {
        target.scrollIntoView({ behavior: 'smooth', block: 'end' });
      } else {
        const observer = new MutationObserver(() => {
          const targetNew = document.querySelector('.therapist-card__actions');
          if (targetNew) {
            targetNew.scrollIntoView({ behavior: 'smooth', block: 'end' });
            observer.disconnect();
          }
        });
        observer.observe(document.body, { childList: true, subtree: true });
      }

      this.transition =
        'transform 300ms cubic-bezier(.2,.9,.2,1), opacity 300ms';
      this.transform = `translate(0px, -50%)`;
    } else {
      const card = document.querySelector('.card-wrapper'); // hoặc '.therapist-card'
      if (card) {
        card.scrollIntoView({ behavior: 'smooth', block: 'start' });
      }
      this.transform = `translate(0px, 0px)`;
    }

    setTimeout(() => {
      this.swiped.emit({ direction: 'up', therapist: this.therapist });
      this.resetPosition();
    }, 300);
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
    const dy = this.currentY - this.startY;

    if (Math.abs(dx) >= this.swipeThreshold) {
      const dir: 'left' | 'right' = dx > 0 ? 'right' : 'left';
      this.animateOffScreen(dir);
    } else if (dy <= -this.swipeThreshold) {
      this.animateUp();
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
    } else if (ev.key === 'ArrowUp') {
      this.animateUp();
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

  match() {
    this.programmaticSwipe('right');
  }

  pass() {
    this.programmaticSwipe('left');
  }
  viewProfile() {
    this.openOverlay.emit(this.therapist);
  }
  ngOnChanges(changes: SimpleChanges) {
    if (changes['isOverlayOpen'] && this.isOverlayOpen) {
      const target = document.querySelector('.therapist-card__actions');
      if (target) {
        target.scrollIntoView({ behavior: 'smooth', block: 'end' });
      }
    }
  }
}
