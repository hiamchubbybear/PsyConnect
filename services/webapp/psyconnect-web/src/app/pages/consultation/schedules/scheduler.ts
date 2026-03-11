import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnInit,
} from '@angular/core';
import { CalendarEvent, CalendarModule, CalendarView } from 'angular-calendar';
import {
  addDays,
  addMonths,
  endOfMonth,
  startOfMonth,
  subDays,
  subMonths,
} from 'date-fns';
import { SessionService } from '../../../services/consultation/session.service';
import { TherapistService } from '../../../services/consultation/therapist.service';
import { SessionDetailModalComponent } from '../../../components/consultation/session-detail-modal/session-detail-modal';
import { forkJoin, map, of } from 'rxjs';
import { catchError } from 'rxjs/operators';

@Component({
  selector: 'app-scheduler',
  standalone: true,
  imports: [CommonModule, CalendarModule, SessionDetailModalComponent],
  templateUrl: './scheduler.html',
  styleUrls: ['./scheduler.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SchedulerComponent implements OnInit {
  view: CalendarView = CalendarView.Month;
  CalendarView = CalendarView;
  viewDate: Date = new Date();
  events: CalendarEvent[] = [];
  isLoading = true;

  // Interaction state
  selectedDate: Date | null = null;
  selectedDayEvents: CalendarEvent[] = [];
  therapistNames: Map<string, string> = new Map();

  // Modal state
  selectedDetailedEvent: CalendarEvent | null = null;

  constructor(
    private sessionService: SessionService,
    private therapistService: TherapistService,
    private cdr: ChangeDetectorRef,
  ) {}

  ngOnInit(): void {
    this.loadCalendar();
  }

  loadCalendar(): void {
    this.isLoading = true;
    const from = startOfMonth(this.viewDate).toISOString();
    const to = endOfMonth(this.viewDate).toISOString();

    this.sessionService.getCalendar(from, to).subscribe({
      next: (res) => {
        const rawEvents = res.events || [];
        const therapistIds = [...new Set(rawEvents.map(e => e.therapist_id).filter((id): id is string => !!id && !this.therapistNames.has(id)))];
        
        const therapistRequests = therapistIds.map(id => 
          this.therapistService.getTherapistById(id).pipe(
            map(t => ({ id, name: t.name || 'Hệ thống' })),
            catchError(() => of({ id, name: 'Chuyên gia' }))
          )
        );

        if (therapistRequests.length > 0) {
          forkJoin(therapistRequests).subscribe(results => {
            results.forEach(r => this.therapistNames.set(r.id, r.name));
            this.processEvents(rawEvents);
          });
        } else {
          this.processEvents(rawEvents);
        }
      },
      error: () => {
        this.events = [];
        this.isLoading = false;
        this.cdr.markForCheck();
      },
    });
  }

  private processEvents(rawEvents: any[]): void {
    this.events = rawEvents.map((e) => {
      const mode = e.mode?.replace('-', '_');
      return {
        id: e.id,
        title: mode === 'online' ? 'Online' : mode === 'in_person' ? 'Trực tiếp' : mode === 'phone' ? 'Điện thoại' : 'Tư vấn',
        start: new Date(e.start),
        end: new Date(e.end),
        color: { primary: e.color, secondary: e.color + '22' },
        meta: { 
          ...e, 
          mode,
          partnerName: e.therapist_id ? this.therapistNames.get(e.therapist_id) || 'Chuyên gia' : 'Bản thân',
          locationInfo: e.location_info
        },
      };
    });

    this.isLoading = false;
    if (this.selectedDate) {
      this.showDetailsForDate(this.selectedDate);
    }
    this.cdr.markForCheck();
  }

  handleDayClick(day: any): void {
    if (!day || !day.date) return;
    this.showDetailsForDate(day.date);
  }

  handleHourClick(date: Date): void {
    this.showDetailsForDate(date);
  }

  handleEventClick(event: CalendarEvent): void {
    this.openDetailedModal(event);
  }

  openDetailedModal(event: CalendarEvent): void {
    this.selectedDetailedEvent = event;
    this.cdr.markForCheck();
  }

  closeDetailedModal(): void {
    this.selectedDetailedEvent = null;
    this.cdr.markForCheck();
  }

  private showDetailsForDate(date: Date): void {
    this.selectedDate = date;
    const startOfDay = new Date(date).setHours(0, 0, 0, 0);
    const endOfDay = new Date(date).setHours(23, 59, 59, 999);

    this.selectedDayEvents = this.events.filter((event) => {
      const eventTime = event.start.getTime();
      return eventTime >= startOfDay && eventTime <= endOfDay;
    });
    
    this.cdr.markForCheck();
  }

  openMap(event: CalendarEvent): void {
    const meta = event.meta;
    if (!meta) return;

    if (meta.locationInfo?.coordinates?.length >= 2) {
      const lng = meta.locationInfo.coordinates[0];
      const lat = meta.locationInfo.coordinates[1];
      window.open(`https://www.google.com/maps/search/?api=1&query=${lat},${lng}`, '_blank');
      return;
    }
    
    if (meta.address) {
      const address = encodeURIComponent(meta.address);
      const url = `https://www.google.com/maps/search/?api=1&query=${address}`;
      window.open(url, '_blank');
    }
  }

  closeDetails(): void {
    this.selectedDate = null;
    this.selectedDayEvents = [];
    this.cdr.markForCheck();
  }

  setView(view: CalendarView) {
    this.view = view;
    this.loadCalendar();
  }

  previousPeriod() {
    if (this.view === CalendarView.Month) {
      this.viewDate = subMonths(this.viewDate, 1);
    } else {
      this.viewDate = subDays(this.viewDate, 7);
    }
    this.loadCalendar();
  }

  nextPeriod() {
    if (this.view === CalendarView.Month) {
      this.viewDate = addMonths(this.viewDate, 1);
    } else {
      this.viewDate = addDays(this.viewDate, 7);
    }
    this.loadCalendar();
  }

  today() {
    this.viewDate = new Date();
    this.loadCalendar();
  }
}
