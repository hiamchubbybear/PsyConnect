import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

interface MoodOption {
  value: string;
  emoji: string;
  label: string;
  gradient: string;
}

interface MoodEntry {
  date: Date;
  mood: string;
  note?: string;
}

interface WeekDay {
  name: string;
  date: Date;
  mood?: string;
}

interface MoodStat {
  mood: string;
  emoji: string;
  label: string;
  count: number;
  percentage: number;
  gradient: string;
}

@Component({
  selector: 'app-mood-tracker',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: './mood-tracker.html',
  styleUrls: ['./mood-tracker.scss'],
})
export class MoodTrackerComponent implements OnInit {
  
  moodOptions: MoodOption[] = [
    {
      value: 'awful',
      emoji: '/assets/mood/awful.svg',
      label: 'MOOD.Awful',
      gradient: '#5B8DEE', 
    },
    {
      value: 'bad',
      emoji: '/assets/mood/bad.svg',
      label: 'MOOD.Bad',
      gradient: '#6DBEDB', 
    },
    {
      value: 'okay',
      emoji: '/assets/mood/okay.svg',
      label: 'MOOD.Okay',
      gradient: '#F4C430', 
    },
    {
      value: 'good',
      emoji: '/assets/mood/good.svg',
      label: 'MOOD.Good',
      gradient: '#A8D08D', 
    },
    {
      value: 'great',
      emoji: '/assets/mood/great.svg',
      label: 'MOOD.Great',
      gradient: '#6EBF8B', 
    },
  ];

  selectedMood: string | null = null;
  moodNote = '';
  todayMoodLogged = false;
  currentStreak = 0;
  currentWeekOffset = 0;

  weekDays: WeekDay[] = [];
  moodEntries: MoodEntry[] = [];
  moodStats: MoodStat[] = [];
  moodLinePath = '';

  averageMoodEmoji = '/assets/mood/good.svg';
  averageMoodLabel = 'Good';
  averageMoodGradient = '#A8D08D';

  patterns: any[] = [
    {
      icon: '/assets/icon/trending-up.svg',
      text: 'MOOD.Patterns.ImprovingTrend',
    },
    { icon: '/assets/icon/moon.svg', text: 'MOOD.Patterns.BetterMornings' },
    {
      icon: '/assets/icon/activity.svg',
      text: 'MOOD.Patterns.ConsistentLogging',
    },
  ];

  
  influenceFactors = [
    {
      id: 'exercise',
      icon: '/assets/icon/activity.svg',
      label: 'MOOD.Influences.Exercise',
      type: 'positive',
      active: false,
    },
    {
      id: 'sleep',
      icon: '/assets/icon/moon.svg',
      label: 'MOOD.Influences.Sleep',
      type: 'positive',
      active: false,
    },
    {
      id: 'social',
      icon: '/assets/icon/users.svg',
      label: 'MOOD.Influences.Social',
      type: 'positive',
      active: false,
    },
    {
      id: 'work',
      icon: '/assets/icon/briefcase.svg',
      label: 'MOOD.Influences.Work',
      type: 'negative',
      active: false,
    },
    {
      id: 'stress',
      icon: '/assets/icon/alert-circle.svg',
      label: 'MOOD.Influences.Stress',
      type: 'negative',
      active: false,
    },
  ];

  ngOnInit() {
    this.loadMoodData();
    this.generateWeekDays();
    this.calculateStats();
    this.checkTodayMood();
    this.calculateStreak();
    this.moodLinePath = this.calculateLinePath();
  }

  loadMoodData() {
    
    
    const today = new Date();
    this.moodEntries = [
      {
        date: new Date(
          today.getFullYear(),
          today.getMonth(),
          today.getDate() - 6
        ),
        mood: 'okay',
        note: 'Had a productive day!',
      },
      {
        date: new Date(
          today.getFullYear(),
          today.getMonth(),
          today.getDate() - 5
        ),
        mood: 'good',
      },
      {
        date: new Date(
          today.getFullYear(),
          today.getMonth(),
          today.getDate() - 4
        ),
        mood: 'bad',
      },
      {
        date: new Date(
          today.getFullYear(),
          today.getMonth(),
          today.getDate() - 3
        ),
        mood: 'bad',
      },
      {
        date: new Date(
          today.getFullYear(),
          today.getMonth(),
          today.getDate() - 2
        ),
        mood: 'good',
      },
      {
        date: new Date(
          today.getFullYear(),
          today.getMonth(),
          today.getDate() - 1
        ),
        mood: 'great',
      },
    ];
  }

  generateWeekDays() {
    const today = new Date();
    this.weekDays = [];

    for (let i = 6; i >= 0; i--) {
      const date = new Date(today);
      date.setDate(date.getDate() - i);

      const entry = this.moodEntries.find((e) => this.isSameDay(e.date, date));

      this.weekDays.push({
        name: date.toLocaleDateString('en-US', { weekday: 'short' }),
        date: date,
        mood: entry?.mood,
      });
    }
  }

  calculateStats() {
    const moodCounts: { [key: string]: number } = {};
    const total = this.moodEntries.length;

    this.moodEntries.forEach((entry) => {
      moodCounts[entry.mood] = (moodCounts[entry.mood] || 0) + 1;
    });

    this.moodStats = this.moodOptions
      .map((option) => ({
        mood: option.value,
        emoji: option.emoji,
        label: option.label,
        count: moodCounts[option.value] || 0,
        percentage:
          total > 0 ? ((moodCounts[option.value] || 0) / total) * 100 : 0,
        gradient: option.gradient,
      }))
      .filter((stat) => stat.count > 0);
  }

  checkTodayMood() {
    const today = new Date();
    this.todayMoodLogged = this.moodEntries.some((entry) =>
      this.isSameDay(entry.date, today)
    );
  }

  calculateStreak() {
    let streak = 0;
    const today = new Date();

    for (let i = 0; i < 365; i++) {
      const checkDate = new Date(today);
      checkDate.setDate(checkDate.getDate() - i);

      const hasEntry = this.moodEntries.some((entry) =>
        this.isSameDay(entry.date, checkDate)
      );

      if (hasEntry) {
        streak++;
      } else {
        break;
      }
    }

    this.currentStreak = streak;
  }

  selectMood(mood: string) {
    if (!this.todayMoodLogged) {
      this.selectedMood = mood;
    }
  }

  saveMood() {
    if (!this.selectedMood) return;

    const newEntry: MoodEntry = {
      date: new Date(),
      mood: this.selectedMood,
      note: this.moodNote || undefined,
    };

    
    

    
    this.moodEntries.push(newEntry);
    this.todayMoodLogged = true;
    this.selectedMood = null;
    this.moodNote = '';

    
    this.generateWeekDays();
    this.calculateStats();
    this.calculateStreak();

    
    console.log('Mood saved! 🎉');
  }

  cancelMood() {
    this.selectedMood = null;
    this.moodNote = '';
  }

  getTodayMood() {
    const today = new Date();
    const entry = this.moodEntries.find((e) => this.isSameDay(e.date, today));
    return (
      this.moodOptions.find((o) => o.value === entry?.mood) ||
      this.moodOptions[2]
    );
  }

  getMoodGradient(mood: string): string {
    return this.moodOptions.find((o) => o.value === mood)?.gradient || '';
  }

  getMoodEmoji(mood: string): string {
    return this.moodOptions.find((o) => o.value === mood)?.emoji || '😐';
  }

  getMoodLabel(mood: string): string {
    return this.moodOptions.find((o) => o.value === mood)?.label || 'MOOD.Okay';
  }

  isToday(day: WeekDay): boolean {
    return this.isSameDay(day.date, new Date());
  }

  isFuture(day: WeekDay): boolean {
    return day.date > new Date();
  }

  isSameDay(date1: Date, date2: Date): boolean {
    return (
      date1.getFullYear() === date2.getFullYear() &&
      date1.getMonth() === date2.getMonth() &&
      date1.getDate() === date2.getDate()
    );
  }

  viewFullHistory() {
    
    console.log('View full history');
  }

  
  getMoodValue(mood: string): number {
    const values: { [key: string]: number } = {
      awful: 1,
      bad: 2,
      okay: 3,
      good: 4,
      great: 5,
    };
    return values[mood] || 3;
  }

  getX(index: number): number {
    const width = 400;
    const padding = 40;
    const step = (width - padding * 2) / 6;
    return padding + index * step;
  }

  getY(mood: string): number {
    const height = 150;
    const padding = 20;
    const value = this.getMoodValue(mood);
    return height - ((value - 1) / 4) * (height - padding * 2) - padding;
  }

  calculateLinePath(): string {
    if (this.weekDays.length === 0) return '';

    let path = '';
    let firstPoint = true;

    this.weekDays.forEach((day, i) => {
      if (!day.mood) return;
      const x = this.getX(i);
      const y = this.getY(day.mood);
      path += firstPoint ? `M ${x} ${y}` : ` L ${x} ${y}`;
      firstPoint = false;
    });

    return path;
  }

  
  previousWeek() {
    this.currentWeekOffset--;
    this.generateWeekDays();
    this.moodLinePath = this.calculateLinePath();
  }

  nextWeek() {
    if (this.currentWeekOffset < 0) {
      this.currentWeekOffset++;
      this.generateWeekDays();
      this.moodLinePath = this.calculateLinePath();
    }
  }

  
  toggleInfluence(factor: any) {
    factor.active = !factor.active;
    
  }

  getMoodColor(mood: string): string {
    return (
      this.moodOptions.find((o) => o.value === mood)?.gradient || '#A8D08D'
    );
  }
}
