import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-request-reset-success',
  standalone: true,
  imports: [CommonModule, TranslateModule, RouterModule],
  templateUrl: './request-reset-success.html',
  styleUrls: ['./request-reset-success.scss'],
})
export class RequestResetSuccessComponent implements OnInit, OnDestroy {
  email = '';
  countdown = 30;
  private timer: any;

  constructor(private router: Router) {
    const nav = this.router.getCurrentNavigation();
    this.email = (nav?.extras.state as any)?.email || '';
  }

  ngOnInit() {
    this.startCountdown();
  }

  ngOnDestroy() {
    if (this.timer) clearInterval(this.timer);
  }

  private startCountdown() {
    this.countdown = 30;
    this.timer = setInterval(() => {
      if (this.countdown > 0) {
        this.countdown--;
      } else {
        clearInterval(this.timer);
      }
    }, 1000);
  }

  onResend() {
    console.log('Resend reset link to:', this.email);
    this.startCountdown();
  }
}
