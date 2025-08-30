import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ThemeService } from '../../services/theme/theme-service';
import { ToastType } from '../../shared/toast/toast-type';
import { ToastService } from '../../shared/toast/toast.service';
import { SharedTranslateModule } from '../../shared/translate/translate.module';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-notfound',
  standalone: true,
  imports: [RouterModule, TranslateModule],
  templateUrl: './notfound.html',
  styleUrl: './notfound.scss',
})
export class Notfound {
  constructor(
    private themeService: ThemeService,
    private toastService: ToastService
  ) {}
  toggleTheme() {
    this.themeService.toggleTheme();
    console.log('Show toast');
    this.toastService.show('Login success', 'Success', ToastType.Success);
  }
}
