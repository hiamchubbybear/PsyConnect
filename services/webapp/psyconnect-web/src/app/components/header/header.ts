import { CommonModule } from '@angular/common';
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  HostBinding,
  HostListener,
  OnDestroy,
  OnInit,
} from '@angular/core';
import { MatMenuModule } from '@angular/material/menu';
import { NavigationEnd, Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { filter, Subscription } from 'rxjs';
import { Auth } from '../../services/auth/auth';
import { LoaderService } from '../../services/loader/loader';
import {
  UserContextService,
  UserProfile,
} from '../../services/profile/profile-service';
import { ThemeService } from '../../services/theme/theme-service';
import { AvatarFallbackPipe } from '../../shared/pipes/avatar-fallback.pipe';
import { TranslationService } from '../../shared/translate/translate-service';
import {
  DropdownComponent,
  DropdownOption,
} from '../../shared/ui-atoms/dropdown/dropdown';
import { AvatarMenuComponent } from '../avatar-menu/avatar-menu';
import { HeaderStateService } from './header-state';
import { NotificationDropdownComponent } from './notification-dropdown/notification-dropdown';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    CommonModule,
    RouterModule,
    TranslateModule,
    AvatarMenuComponent,
    MatMenuModule,
    DropdownComponent,
    NotificationDropdownComponent,
    AvatarFallbackPipe,
  ],
  templateUrl: './header.html',
  styleUrl: './header.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class Header implements OnInit, OnDestroy {
  currentLang: string = '';
  displayLang: string = '';

  selectedLanguage: any = null;
  private subscriptions: Subscription = new Subscription();

  languageOptions: DropdownOption[] = [
    { value: 'en', label: 'English' },
    { value: 'vi', label: 'Vietnamese' },
  ];

  isMini = false;
  isFeedRoute = false;

  currentUser: any = null;
  userAvatarUrl: string | null = null;
  userProfileId: string | null = null; 
  userName = 'Anonymous';
  isMenuOpen = false;
  isHidden = false;
  lastScrollTop = 30;
  isDark = false;

  private userSubscription?: Subscription;

  @HostBinding('class.host-hidden') get isHostHidden() {
    return this.isHidden;
  }

  constructor(
    public headerState: HeaderStateService,
    public auth: Auth,
    private userContext: UserContextService,
    public loaderService: LoaderService,
    private cdr: ChangeDetectorRef,
    private router: Router,
    private themeService: ThemeService,
    private translateService: TranslationService,
  ) {
    this.translateService.currentLanguage$.subscribe((lang) => {
      this.currentLang = lang;
      this.displayLang = lang === 'en' ? 'English' : 'Vietnamese';
    });

    this.router.events
      .pipe(filter((event) => event instanceof NavigationEnd))
      .subscribe((event: any) => {
        this.isFeedRoute = this.router.url.includes('/feature/feed');
        this.cdr.markForCheck();
      });
  }
  ngOnInit() {
    this.isFeedRoute = this.router.url.includes('/feature/feed');
    this.userSubscription = this.headerState.isMini$.subscribe((value) => {
      this.isMini = value;
      this.cdr.markForCheck();
    });
    this.userSubscription = this.userContext.user$.subscribe(
      (user: UserProfile | null) => {
        this.updateUserInfo(user);
        this.cdr.detectChanges();
      },
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
      this.userProfileId = user.profileId || user.accountId || 'guest';
      this.userAvatarUrl = user.avatarUri?.trim() ? user.avatarUri : null;
      this.userName =
        `${user.firstName} ${user.lastName}`.trim() || 'Anonymous';
    } else {
      this.userProfileId = 'guest';
      this.userAvatarUrl = null;
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
  toggleSearch() {
    const currentState = this.headerState.isSearchOpen();
    this.headerState.setSearchOpen(!currentState);
  }
  onLanguageChange(option: DropdownOption | null) {
    if (option) {
      this.setLanguage(option.value);
    }
  }

  setLanguage(lang: string) {
    this.translateService.setLanguage(lang);
    this.selectedLanguage = lang;
  }
  @HostListener('window:scroll', [])
  onScroll() {
    const scrollTop = window.scrollY || document.documentElement.scrollTop;
    const isSearchOpen = this.headerState.isSearchOpen();

    if (
      scrollTop > this.lastScrollTop + 10 &&
      scrollTop > 60 &&
      !isSearchOpen
    ) {
      this.isHidden = true;
      this.headerState.setMini(true);
    } else if (scrollTop < this.lastScrollTop - 5) {
      this.isHidden = false;
      this.headerState.setMini(false);
    }
    this.lastScrollTop = scrollTop <= 0 ? 0 : scrollTop;
    this.cdr.markForCheck();
  }
}
