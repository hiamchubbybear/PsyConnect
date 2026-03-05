import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  OnDestroy,
  OnInit,
} from '@angular/core';
import { Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import {
  InAppNotification,
  NotificationService,
} from '../../../services/notification/notification.service';
import { UserContextService } from '../../../services/profile/profile-service';

@Component({
  selector: 'app-notification-dropdown',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './notification-dropdown.html',
  styleUrl: './notification-dropdown.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class NotificationDropdownComponent implements OnInit, OnDestroy {
  notifications: InAppNotification[] = [];
  unreadCount = 0;
  isOpen = false;
  private subscriptions = new Subscription();

  constructor(
    private notificationService: NotificationService,
    private userContext: UserContextService,
    private cdr: ChangeDetectorRef,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.subscriptions.add(
      this.notificationService.notifications$.subscribe(
        (list: InAppNotification[]) => {
          this.notifications = list;
          this.cdr.markForCheck();
        }
      )
    );

    this.subscriptions.add(
      this.notificationService.unreadCount$.subscribe((count: number) => {
        this.unreadCount = count;
        this.cdr.markForCheck();
      })
    );

    
    const user = this.userContext.getUser();
    if (user) {
      this.notificationService.fetchNotifications(user.accountId);
    }
  }

  ngOnDestroy(): void {
    this.subscriptions.unsubscribe();
  }

  toggleDropdown(event: Event): void {
    event.stopPropagation();
    this.isOpen = !this.isOpen;
  }

  markAsRead(event: Event, notification: InAppNotification): void {
    event.stopPropagation();
    if (notification.isRead) return;

    const user = this.userContext.getUser();
    if (user) {
      this.notificationService.markAsRead(user.accountId, notification.id);
    }
  }

  onNotificationClick(notification: InAppNotification): void {
    this.isOpen = false;

    
    if (!notification.isRead) {
      const user = this.userContext.getUser();
      if (user) {
        this.notificationService.markAsRead(user.accountId, notification.id);
      }
    }

    
    if (notification.type === 'like' || notification.type === 'comment') {
      const postId = notification.metadata?.postId;
      if (postId) {
        
        
        console.log('Navigating to post:', postId);
      }
    } else if (notification.type === 'follow') {
      const followerId = notification.metadata?.followerId;
      if (followerId) {
        this.router.navigate(['/profile', followerId]);
      }
    }
  }
}
