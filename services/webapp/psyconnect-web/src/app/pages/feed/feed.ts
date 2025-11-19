import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

interface Post {
  id: string;
  author: string;
  role?: string;
  time: string;
  titleKey: string;
  excerptKey: string;
  likes: number;
  comments: number;
  avatar?: string;
  media?: string;
}
interface Therapist {
  name: string;
  avatar: string;
  specialtyKey: string;
  rating: number;
  statusKey: string;
}

interface SupportGroup {
  nameKey: string;
  members: number;
}
@Component({
  selector: 'app-feed',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule],
  templateUrl: './feed.html',
  styleUrls: ['./feed.scss'],
})
export class FeedComponent {
  constructor(private router: Router) {}
  posts: Post[] = [
    {
      id: 'p1',
      author: 'Dr. Sarah Johnson',
      role: 'Therapist',
      time: '2 hours ago',
      titleKey: 'FEED.Posts.P1.Title',
      excerptKey: 'FEED.Posts.P1.Excerpt',
      likes: 42,
      comments: 8,
      avatar: 'https://randomuser.me/api/portraits/women/44.jpg',
      media: 'https://images.pexels.com/photos/414886/pexels-photo-414886.jpeg',
    },
    {
      id: 'p2',
      author: 'Anonymous User',
      role: 'Community Member',
      time: '4 hours ago',
      titleKey: 'FEED.Posts.P2.Title',
      excerptKey: 'FEED.Posts.P2.Excerpt',
      likes: 28,
      comments: 12,
      avatar: 'https://randomuser.me/api/portraits/men/35.jpg',
      media:
        'https://images.pexels.com/photos/1051838/pexels-photo-1051838.jpeg',
    },
  ];

  moodWeek = ['😊', '😐', '🙂', '😕', '😁', '😢', '🙂'];

  therapists: Therapist[] = [
    {
      name: 'Dr. Emily Chen',
      avatar: 'https://randomuser.me/api/portraits/women/65.jpg',
      specialtyKey: 'FEED.Therapists.P1.Specialty',
      rating: 4.9,
      statusKey: 'FEED.Therapists.P1.Status',
    },
    {
      name: 'Dr. Michael Torres',
      avatar: 'https://randomuser.me/api/portraits/men/43.jpg',
      specialtyKey: 'FEED.Therapists.P2.Specialty',
      rating: 4.7,
      statusKey: 'FEED.Therapists.P2.Status',
    },
  ];
  supportGroups: SupportGroup[] = [
    { nameKey: 'FEED.SupportGroups.P1', members: 127 },
    { nameKey: 'FEED.SupportGroups.P2', members: 89 },
  ];
  toggleMoodLog() {
    alert('Open mood log (implement modal)');
  }

  bookTherapist(therapist: any) {
    alert(`Book ${therapist.name} (implement booking)`);
  }

  newPost() {
    this.router.navigate(['/feature/article/create']);
  }

  likePost(p: Post) {
    p.likes++;
  }

  openPost(p: Post) {}
}
