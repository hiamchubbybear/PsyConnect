import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { ThemeService } from '../../services/theme/theme-service';
import { ToastService } from '../../shared/toast/toast.service';
import { ToastType } from '../../shared/toast/toast.model';

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
    this.toastService.show({
      message: 'TOAST.theme_change_success',
      title: 'TOAST.success',
      type: ToastType.Success
    });
  }
}
