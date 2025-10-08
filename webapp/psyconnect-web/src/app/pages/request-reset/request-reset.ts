import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
    FormBuilder,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { PasswordService } from '../../services/auth/password.service';
import { ToastService } from '../../shared/toast/toast.service';

@Component({
  selector: 'app-request-reset',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule, RouterModule],
  templateUrl: './request-reset.html',
  styleUrls: ['./request-reset.scss'],
})
export class RequestResetComponent {
  requestForm: FormGroup;
  submitting = false;
  shakeErrors = false;

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private passwordService: PasswordService,
    private toastService: ToastService,
    private translate: TranslateService
  ) {
    this.requestForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
    });
  }

  isFieldInvalid(fieldName: string): boolean {
    const field = this.requestForm.get(fieldName);
    return !!(field && field.invalid && (field.dirty || field.touched));
  }

  getFieldError(fieldName: string): string {
    const field = this.requestForm.get(fieldName);
    if (!field || !field.errors) return '';

    if (field.errors['required']) {
      return this.translate.instant('LOGIN.REQUEST.Form.EmailRequired');
    }
    if (field.errors['email']) {
      return this.translate.instant('LOGIN.REQUEST.Form.EmailInvalid');
    }
    return this.translate.instant('LOGIN.REQUEST.Form.Invalid');
  }

  private triggerShake() {
    this.shakeErrors = false;
    requestAnimationFrame(() => {
      this.shakeErrors = true;
      setTimeout(() => (this.shakeErrors = false), 320);
    });
  }

  onSubmit() {
    if (this.requestForm.invalid || this.submitting) {
      this.triggerShake();
      return;
    }

    this.submitting = true;
    const email = this.requestForm.value.email;
    /**
     * @deprecated The field should be remove next api update
     */
    const username = '';
    this.passwordService
      .requestReset({
        email: email,
      })
      .subscribe({
        next: (res) => {
          if (res.code == 200 && res.data) {
            setTimeout(() => {
              this.submitting = false;
              this.router.navigate(['/auth/request-reset-success'], {
                state: { email },
              });
            }, 1500);
          } else {
            this.toastService.error('Error', res.message, 1500);
            this.triggerShake();
          }
        },
        error: (err) => {
          this.submitting = false;
          this.triggerShake();
          if (err.error && err.error.message) {
            this.toastService.error('Error', err.error.message, 2000);
          } else {
            this.toastService.error(
              'Error',
              err.statusText || 'Unknown error',
              2000
            );
            console.log('HTTP error:', err);
          }
        },
      });
  }

  onCancel() {
    this.router.navigate(['/auth/login']);
  }
}
