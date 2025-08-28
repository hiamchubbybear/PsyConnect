import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ThemeService } from '../../services/theme/theme-service';
import { ToastType } from '../../shared/toast/toast-type';
import { ToastService } from '../../shared/toast/toast.service';

@Component({
  selector: 'app-notfound',
  standalone: true,
  imports: [RouterModule],
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
