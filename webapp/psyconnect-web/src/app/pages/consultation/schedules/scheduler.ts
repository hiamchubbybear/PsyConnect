import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { SidebarService } from '../../../services/sidebar/sidebar';

interface Therapist {
  id: string;
  name: string;
  img: string;
}

@Component({
  selector: 'app-scheduler',
  templateUrl: './scheduler.html',
  styleUrls: ['./scheduler.scss'],
  imports: [CommonModule],
  standalone: true,
})
export class SchedulerComponent implements OnInit, OnDestroy {
  today: string = new Date().toDateString();
  isSidebarCollapsed = false;
  private sub?: Subscription;
  therapists: Therapist[] = [
    { id: 'gabi', name: 'Gabi Guimaraes', img: 'assets/img/t1.jpg' },
    { id: 'yuki', name: 'Yuki Ishikawa', img: 'assets/img/t5.jpg' },
    { id: 'fornal', name: 'Tomasz Fornal', img: 'assets/img/t4.jpg' },
    { id: 'sun', name: 'Sun YingSha', img: 'assets/img/t2.jpg' },
    { id: 'wang', name: 'Wang ChuQin', img: 'assets/img/t6.jpg' },
  ];

  viewDate: Date = new Date(new Date().getFullYear(), new Date().getMonth(), 1);
  selectedDate?: Date;
  selectedTime?: string;
  selectedTherapist?: Therapist;

  calendarDays: { date: Date; label: number; out?: boolean }[] = [];

  morningSlots = [
    '8:00 AM',
    '8:20 AM',
    '8:40 AM',
    '9:00 AM',
    '9:20 AM',
    '9:40 AM',
    '10:00 AM',
    '10:20 AM',
    '10:40 AM',
    '11:00 AM',
  ];
  afternoonSlots = [
    '1:00 PM',
    '1:30 PM',
    '2:00 PM',
    '2:30 PM',
    '3:00 PM',
    '3:30 PM',
    '4:00 PM',
    '4:30 PM',
    '5:00 PM',
    '5:30 PM',
  ];
  ngOnDestroy() {
    this.sub?.unsubscribe();
  }
  constructor(private sidebarService: SidebarService) {}
  ngOnInit() {
    this.sub = this.sidebarService.collapsed$.subscribe((collapsed) => {
      this.isSidebarCollapsed = collapsed;
    });
    this.generateCalendar();
  }

  monthLabel(date: Date) {
    return date
      .toLocaleString('en-US', { month: 'short', year: 'numeric' })
      .toUpperCase();
  }

  generateCalendar() {
    const y = this.viewDate.getFullYear();
    const m = this.viewDate.getMonth();
    const firstDay = new Date(y, m, 1).getDay();
    const daysInMonth = new Date(y, m + 1, 0).getDate();
    const prevDays = new Date(y, m, 0).getDate();
    const totalCells = 42;

    this.calendarDays = [];

    for (let i = 0; i < totalCells; i++) {
      let day: number;
      let date: Date;
      let out = false;

      if (i < firstDay) {
        day = prevDays - firstDay + 1 + i;
        date = new Date(y, m - 1, day);
        out = true;
      } else if (i >= firstDay + daysInMonth) {
        day = i - (firstDay + daysInMonth) + 1;
        date = new Date(y, m + 1, day);
        out = true;
      } else {
        day = i - firstDay + 1;
        date = new Date(y, m, day);
      }

      this.calendarDays.push({ date, label: day, out });
    }
  }

  prevMonth() {
    this.viewDate.setMonth(this.viewDate.getMonth() - 1);
    this.generateCalendar();
  }

  nextMonth() {
    this.viewDate.setMonth(this.viewDate.getMonth() + 1);
    this.generateCalendar();
  }

  selectDate(day: { date: Date; out?: boolean }) {
    if (!day.out) this.selectedDate = day.date;
  }

  selectTime(time: string) {
    this.selectedTime = time;
  }

  selectTherapist(therapist: Therapist) {
    this.selectedTherapist = therapist;
  }

  get footerText() {
    if (!this.selectedDate)
      return 'Chọn ngày, khung giờ và therapist để tiếp tục.';
    const prettyDate = this.selectedDate.toLocaleDateString('en-US', {
      weekday: 'short',
      day: '2-digit',
      month: 'short',
      year: 'numeric',
    });
    if (!this.selectedTime)
      return `Ngày: ${prettyDate}. Chọn khung giờ và therapist để tiếp tục.`;
    if (!this.selectedTherapist)
      return `Ngày: ${prettyDate} – Khung giờ: ${this.selectedTime}. Chọn therapist để tiếp tục.`;
    return `Ngày: ${prettyDate} – Khung giờ: ${this.selectedTime} – Therapist: ${this.selectedTherapist.name}`;
  }

  canNext() {
    return (
      !!this.selectedDate && !!this.selectedTime && !!this.selectedTherapist
    );
  }
}
