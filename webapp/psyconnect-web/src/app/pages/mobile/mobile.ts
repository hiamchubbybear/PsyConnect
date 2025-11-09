import { CommonModule } from '@angular/common';
import { Component, computed } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { ThemeService } from '../../services/theme/theme-service';

@Component({
  selector: 'app-mobile-required',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './mobile.html',
  styleUrls: ['./mobile.scss'],
})
export class MobileRequiredComponent {
  logoPath = computed(() =>
    this.themeService.getTheme() === 'dark'
      ? './assets/sig-logo/iconic-logo-dark.svg'
      : './assets/sig-logo/iconic-logo-light.svg'
  );

  constructor(private themeService: ThemeService) {}
}
