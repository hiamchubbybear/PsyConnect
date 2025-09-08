// therapist-profile-overlay.component.ts
import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { Therapist } from '../../models/swipe-card';

@Component({
  selector: 'app-therapist-profile-overlay',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: 'profile-overlay.html',
  styleUrls: ['./profile-overlay.scss'],
})
export class TherapistProfileOverlayComponent implements OnInit {
  ngOnInit(): void {
    throw new Error('Method not implemented.');
  }
  @Input() therapist: Therapist | null = null;
  @Input() isVisible: boolean = false;

  @Output() closeOverlay = new EventEmitter<void>();
  @Output() bookAppointment = new EventEmitter<Therapist>();
  @Output() sendMessage = new EventEmitter<Therapist>();

  onClose(): void {
    this.closeOverlay.emit();
  }

  onBackdropClick(event: Event): void {
    if (event.target === event.currentTarget) {
      this.onClose();
    }
  }

  onBookAppointment() {
    if (this.therapist) this.bookAppointment.emit(this.therapist);
    this.onClose();
  }

  onSendMessage() {
    if (this.therapist) this.sendMessage.emit(this.therapist);
    this.onClose();
  }

  hasQualifications(): boolean {
    return !!(
      this.therapist?.professionalInfo?.degrees?.length ||
      this.therapist?.professionalInfo?.certifications?.length
    );
  }
}
