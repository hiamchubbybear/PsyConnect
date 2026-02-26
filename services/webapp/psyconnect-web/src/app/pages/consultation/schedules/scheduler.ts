import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component } from '@angular/core';
import { CalendarEvent, CalendarModule, CalendarView } from 'angular-calendar';
import { addDays, addHours, startOfDay, subDays } from 'date-fns';

const colors: any = {
  blue: { primary: '#1e90ff', secondary: '#D1E8FF' },
  green: { primary: '#48bb78', secondary: '#F0FDF4' },
  orange: { primary: '#ed8936', secondary: '#FFEDD5' },
};

@Component({
  selector: 'app-scheduler',
  standalone: true,
  imports: [CommonModule, CalendarModule],
  templateUrl: './scheduler.html',
  styleUrls: ['./scheduler.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SchedulerComponent {
  view: CalendarView = CalendarView.Week;
  CalendarView = CalendarView;

  viewDate: Date = new Date();

  events: CalendarEvent[] = [
    {
      start: addHours(startOfDay(new Date()), 9),
      end: addHours(startOfDay(new Date()), 10.5),
      title: 'Daily Meeting (Zoom)',
      color: colors.blue,
      cssClass: 'event-blue',
    },
    {
      start: addHours(startOfDay(new Date()), 11.5),
      end: addHours(startOfDay(new Date()), 13),
      title: 'Family Catchup',
      color: colors.green,
      cssClass: 'event-green',
    },
    {
      start: addHours(startOfDay(addDays(new Date(), 1)), 9),
      end: addHours(startOfDay(addDays(new Date(), 1)), 10),
      title: 'Gym Session',
      color: colors.green,
      cssClass: 'event-green',
    },
    {
      start: addHours(startOfDay(addDays(new Date(), 2)), 11.5),
      end: addHours(startOfDay(addDays(new Date(), 2)), 13),
      title: 'Products Sync',
      color: colors.blue,
      cssClass: 'event-blue',
    },
    {
      start: addHours(startOfDay(subDays(new Date(), 1)), 12),
      end: addHours(startOfDay(subDays(new Date(), 1)), 12.5),
      title: 'Side Hustle Review',
      color: colors.orange,
      cssClass: 'event-orange',
    },
  ];

  setView(view: CalendarView) {
    this.view = view;
  }
}
