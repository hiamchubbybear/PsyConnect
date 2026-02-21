import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { LoaderService } from '../../services/loader/loader';
import { ThemeService } from '../../services/theme/theme-service';
import { TranslationService } from '../../shared/translate/translate-service';
import {
  DropdownComponent,
  DropdownOption,
} from '../../shared/ui-atoms/dropdown/dropdown';

@Component({
  selector: 'app-auth-header',
  standalone: true,
  imports: [CommonModule, TranslateModule, DropdownComponent],
  templateUrl: './auth-header.html',
  styleUrls: ['./auth-header.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AuthHeaderComponent {
  isDark = false;
  selectedLanguage: any = null;
  displayLang: string = '';

  languageOptions: DropdownOption[] = [
    { value: 'en', label: 'English' },
    { value: 'vi', label: 'Vietnamese' },
  ];

  constructor(
    private themeService: ThemeService,
    private translateService: TranslationService,
    public loaderService: LoaderService,
  ) {
    this.translateService.currentLanguage$.subscribe((lang) => {
      this.displayLang = lang === 'en' ? 'English' : 'Vietnamese';
    });
  }

  toggleTheme() {
    this.themeService.toggleTheme();
    this.isDark = !this.isDark;
  }

  onLanguageChange(option: DropdownOption | null) {
    if (option) {
      this.translateService.setLanguage(option.value);
    }
  }
}
