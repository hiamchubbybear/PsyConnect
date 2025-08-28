import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  HostListener,
  OnDestroy,
  OnInit,
} from '@angular/core';
import { MatMenuModule } from '@angular/material/menu';
import { Router } from '@angular/router';
import { Observable, Subscription } from 'rxjs';
import { Auth } from '../../services/auth/auth';
import { LoaderService } from '../../services/loader/loader';
import {
  UserContextService,
  UserProfile,
} from '../../services/profile/profile-service';
import { ThemeService } from '../../services/theme/theme-service';
import { AvatarMenuComponent } from '../avatar-menu/avatar-menu';
import { HeaderStateService } from './header-state';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [CommonModule, MatMenuModule, AvatarMenuComponent],
  templateUrl: './header.html',
  styleUrl: './header.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class Header implements OnInit, OnDestroy {
  isMini = false;

  userAvatarUrl = 'assets/images/avatar.jpeg';
  userName = 'Anonymous';
  isMenuOpen = false;
  isHidden = false;
  lastScrollTop = 30;
  loading$!: Observable<boolean>;
  isDark = false;

  private userSubscription?: Subscription;

  constructor(
    public headerState: HeaderStateService,
    public auth: Auth,
    private userContext: UserContextService,
    public loaderService: LoaderService,
    private cdr: ChangeDetectorRef,
    private router: Router,
    private themeService: ThemeService
  ) {}

  ngOnInit() {
    this.userSubscription = this.headerState.isMini$.subscribe((value) => {
      this.isMini = value;
      this.cdr.markForCheck();
    });
    this.loading$ = this.loaderService.loading$;
    this.userSubscription = this.userContext.user$.subscribe(
      (user: UserProfile | null) => {
        this.updateUserInfo(user);
        this.cdr.detectChanges();
      }
    );

    setTimeout(() => {
      const currentUser = this.userContext.getUser();

      if (currentUser) {
        this.updateUserInfo(currentUser);
        this.cdr.detectChanges();
      } else {
        this.userContext.reloadFromStorage();
      }
    }, 0);
  }

  ngOnDestroy() {
    this.userSubscription?.unsubscribe();
  }

  private updateUserInfo(user: UserProfile | null) {
    if (user) {
      const newAvatarUrl = user.avatarUri?.trim()
        ? user.avatarUri
        : 'assets/images/avatar.jpeg';
      const newUserName =
        `${user.firstName} ${user.lastName}`.trim() || 'Anonymous';

      this.userAvatarUrl = newAvatarUrl;
      this.userName = newUserName;
    } else {
      this.userAvatarUrl = 'assets/images/avatar.jpeg';
      this.userName = 'Anonymous';
    }
  }

  async logout() {
    this.auth.logout();
    await this.loaderService.showWithDelay(300);

    this.userContext.clear();
    window.location.href = '/';
  }

  toggleMenu() {
    this.isMenuOpen = !this.isMenuOpen;
    this.cdr.markForCheck();
  }
  toggleTheme() {
    this.themeService.toggleTheme();
    this.isDark = !this.isDark;
  }
  @HostListener('window:scroll', [])
  onScroll() {
    const scrollTop = window.scrollY || document.documentElement.scrollTop;
    if (scrollTop > this.lastScrollTop + 1) {
      this.headerState.setMini(true);
    } else if (scrollTop < this.lastScrollTop - 10) {
      this.headerState.setMini(false);
    }
    this.lastScrollTop = scrollTop <= 0 ? 0 : scrollTop;
    this.cdr.markForCheck();
  }
}
