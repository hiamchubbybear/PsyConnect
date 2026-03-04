import { CommonModule } from '@angular/common';
import { Component, Inject, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import {
  MAT_DIALOG_DATA,
  MatDialogModule,
  MatDialogRef,
} from '@angular/material/dialog';
import { TranslateModule } from '@ngx-translate/core';
import {
  CONSULTATION_MODES,
  ConsultationSession,
  SessionRequest,
} from '../../../models/consultation.model';
import { AuthService } from '../../../services/auth/auth.service';
import { ChatService } from '../../../services/chat/chat.service';
import { SessionService } from '../../../services/consultation/session.service';

export interface BookingDialogData {
  therapist: any;
}

@Component({
  selector: 'app-booking-dialog',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    TranslateModule,
  ],
  templateUrl: './booking-dialog.html',
  styleUrls: ['./booking-dialog.scss'],
})
export class BookingDialogComponent implements OnInit {
  bookingForm!: FormGroup;
  minDate = new Date().toISOString().split('T')[0];
  isLoading = false;
  modes = CONSULTATION_MODES;
  therapist: any;
  timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
  bookingSuccess = false;

  constructor(
    private fb: FormBuilder,
    private sessionService: SessionService,
    private authService: AuthService,
    private chatService: ChatService,
    public dialogRef: MatDialogRef<BookingDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: BookingDialogData,
  ) {
    this.therapist = data.therapist;
  }

  ngOnInit() {
    this.bookingForm = this.fb.group({
      scheduledDate: ['', Validators.required],
      startTime: ['', Validators.required],
      endTime: ['', Validators.required],
      mode: ['online', Validators.required],
      acceptResponsibility: [false, Validators.requiredTrue],
    });
  }

  onSubmit() {
    if (this.bookingForm.invalid) return;

    this.isLoading = true;
    const formVal = this.bookingForm.value;

    const startDateTime = new Date(
      `${formVal.scheduledDate}T${formVal.startTime}`,
    );
    const endDateTime = new Date(`${formVal.scheduledDate}T${formVal.endTime}`);

    const currentUser = this.authService.getCurrentUser();
    if (!currentUser) {
      this.isLoading = false;
      return;
    }

    const req: SessionRequest = {
      client_id: currentUser.profileId,
      therapist_id: this.therapist.profile_id || this.therapist.profileId || '',
      start_time: startDateTime.toISOString(),
      end_time: endDateTime.toISOString(),
      time_zone: this.timezone,
      scheduled_date: formVal.scheduledDate,
      mode: formVal.mode,
      price: this.therapist.rage_price || this.therapist.ragePrice || 0,
      location_info:
        formVal.mode === 'online'
          ? { link: 'To be generated' }
          : { address_line: this.therapist.address },
    };

    this.sessionService.createSession(req).subscribe({
      next: (res) => {
        this.isLoading = false;
        this.sendBookingMessage(res);
        this.dialogRef.close(res);
      },
      error: (err) => {
        this.isLoading = false;
        console.error('Booking failed', err);
      },
    });
  }

  private sendBookingMessage(session: ConsultationSession) {
    const currentUser = this.authService.getCurrentUser();
    if (!currentUser) return;

    const therapistId = this.therapist.profile_id || this.therapist.profileId;

    this.chatService
      .getOrCreateConversation(currentUser.profileId, therapistId)
      .subscribe((conv) => {
        if (!conv || !conv.id) {
          console.warn('Could not find conversation for booking message');
          return;
        }

        this.chatService.sendMessage({
          conversationId: conv.id,
          senderId: currentUser.profileId,
          content: `Lịch hẹn ${
            session.mode === 'online' ? 'Trực tuyến' : 'Trực tiếp'
          } vào lúc ${session.scheduled_date} đã được xác nhận.`,
          sessionData: session,
        });
      });
  }
}
