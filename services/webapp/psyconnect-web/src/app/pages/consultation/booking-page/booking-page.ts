import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import {
  CONSULTATION_MODES,
  ConsultationSession,
  SessionRequest,
} from '../../../models/consultation.model';
import { AuthService } from '../../../services/auth/auth.service';
import { ChatService } from '../../../services/chat/chat.service';
import { SessionService } from '../../../services/consultation/session.service';
import { TherapistService } from '../../../services/consultation/therapist.service';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';

@Component({
  selector: 'app-booking-page',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    TranslateModule,
    AvatarFallbackPipe,
  ],
  templateUrl: './booking-page.html',
  styleUrls: ['./booking-page.scss'],
})
export class BookingPageComponent implements OnInit {
  bookingForm!: FormGroup;
  minDate = new Date().toISOString().split('T')[0];
  isLoading = false;
  isFetchingTherapist = true;
  modes = CONSULTATION_MODES;
  therapist: any;
  therapistIdFromRoute = '';
  timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private therapistService: TherapistService,
    private sessionService: SessionService,
    private authService: AuthService,
    private chatService: ChatService,
    private translateService: TranslateService,
  ) {}

  ngOnInit() {
    const therapistId = this.route.snapshot.paramMap.get('therapistId');
    if (!therapistId) {
      this.router.navigate(['/feature/consultation/discover']);
      return;
    }

    this.therapistIdFromRoute = therapistId;
    this.fetchTherapist(therapistId);

    this.bookingForm = this.fb.group({
      scheduledDate: ['', Validators.required],
      startTime: ['', Validators.required],
      endTime: ['', Validators.required],
      mode: ['online', Validators.required],
      acceptResponsibility: [false, Validators.requiredTrue],
    });
  }

  fetchTherapist(id: string) {
    this.isFetchingTherapist = true;
    this.therapistService.getTherapistById(id).subscribe({
      next: (res) => {
        this.therapist = (res as any)?.data ?? res;
        this.isFetchingTherapist = false;
      },
      error: (err) => {
        console.error('Failed to fetch therapist', err);
        this.router.navigate(['/feature/consultation/discover']);
      },
    });
  }

  goBack() {
    window.history.back();
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
      therapist_id:
        this.therapist.profile_id ||
        this.therapist.profileId ||
        this.therapist.id ||
        this.therapistIdFromRoute ||
        '',
      start_time: startDateTime.toISOString(),
      end_time: endDateTime.toISOString(),
      time_zone: this.timezone,
      scheduled_date: formVal.scheduledDate,
      mode: formVal.mode,
      price: Number(this.therapist.rage_price || this.therapist.ragePrice || 0),
      location_info:
        formVal.mode === 'online'
          ? { link: 'To be generated' }
          : { address_line: this.therapist.address || '' },
    };

    this.sessionService.createSession(req).subscribe({
      next: (res) => {
        this.isLoading = false;
        this.sendBookingMessage(res);
        this.router.navigate(['/feature/consultation/sessions']);
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
        if (!conv || !conv.id) return;

        const modeKey =
          session.mode === 'online'
            ? 'CONSULTATION.Modes.Online'
            : 'CONSULTATION.Modes.InPerson';
        const modeLabel = this.translateService.instant(modeKey);

        const content = this.translateService.instant(
          'CONSULTATION.Booking.ConfirmMsg',
          {
            mode: modeLabel,
            date: session.scheduled_date,
          },
        );

        this.chatService.sendMessage({
          conversationId: conv.id,
          senderId: currentUser.profileId,
          content: content,
          sessionData: session,
        });
      });
  }
}
