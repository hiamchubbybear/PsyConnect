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

@Component({
  selector: 'app-scheduler',
  standalone: true,
  imports: [CommonModule, CalendarModule],
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

  constructor(
    private sessionService: SessionService,
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
        this.events = (res.events || []).map((e) => ({
          id: e.id,
          title:
            e.mode === 'online'
              ? '💻 Online'
              : e.mode === 'in_person'
                ? '🏥 Trực tiếp'
                : e.mode === 'phone'
                  ? '📞 Điện thoại'
                  : '💬 Tư vấn',
          start: new Date(e.start),
          end: new Date(e.end),
          color: { primary: e.color, secondary: e.color + '22' },
          meta: e,
        }));
        this.isLoading = false;
        this.cdr.markForCheck();
      },
      error: () => {
        this.events = [];
        this.isLoading = false;
        this.cdr.markForCheck();
      },
    });
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
