import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import {
    FormBuilder,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { PasswordService } from '../../services/auth/password.service';
import { ToastService } from '../../shared/toast/toast.service';

@Component({
  selector: 'app-password-reset',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule, RouterModule],
  templateUrl: './password-reset.html',
  styleUrls: ['./password-reset.scss'],
})
export class ResetPasswordComponent implements OnInit {
  resetForm!: FormGroup;
  submitting = false;

  username = '';
  email = '';
  token = '';
  provider = '';

  shakeErrors = false;

  passwordStrength = 0;
  passwordStrengthLabel = '';
  passwordStrengthClass = 'password-meter__bar';

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private route: ActivatedRoute,
    private translate: TranslateService,
    private toastService: ToastService,
    private passwordService: PasswordService
  ) {}

  ngOnInit() {
    this.resetForm = this.createForm();

    this.route.queryParams.subscribe((params) => {
      this.username = params['username'] || '';
      this.email = params['email'] || '';
      this.token = params['token'] || '';
      this.provider = params['provider'] || '';
    });

    this.resetForm
      .get('password')
      ?.valueChanges.subscribe((password) =>
        this.updatePasswordStrength(password || '')
      );
  }

  private createForm(): FormGroup {
    return this.fb.group(
      {
        password: ['', [Validators.required, Validators.minLength(8)]],
        confirmPassword: ['', [Validators.required]],
      },
      { validators: this.passwordMatchValidator }
    );
  }

  private passwordMatchValidator(form: FormGroup) {
    const password = form.get('password');
    const confirmPassword = form.get('confirmPassword');

    if (
      password &&
      confirmPassword &&
      password.value !== confirmPassword.value
    ) {
      confirmPassword.setErrors({ mismatch: true });
    } else if (confirmPassword?.hasError('mismatch')) {
      confirmPassword.setErrors(null);
    }

    return null;
  }

  private updatePasswordStrength(password: string) {
    if (!password) {
      this.passwordStrength = 0;
      this.passwordStrengthLabel = this.translate.instant(
        'LOGIN.RESET.Form.Empty'
      );
      this.passwordStrengthClass = 'password-meter__bar';
      return;
    }

    let score = 0;
    if (password.length >= 8) score += 25;
    if (password.length >= 12) score += 25;
    if (/[a-z]/.test(password)) score += 15;
    if (/[A-Z]/.test(password)) score += 15;
    if (/[0-9]/.test(password)) score += 10;
    if (/[^A-Za-z0-9]/.test(password)) score += 10;

    this.passwordStrength = score;

    if (score < 40) {
      this.passwordStrengthLabel = this.translate.instant(
        'LOGIN.RESET.Form.Weak'
      );
      this.passwordStrengthClass =
        'password-meter__bar password-meter__bar--weak';
    } else if (score < 70) {
      this.passwordStrengthLabel = this.translate.instant(
        'LOGIN.RESET.Form.Medium'
      );
      this.passwordStrengthClass =
        'password-meter__bar password-meter__bar--medium';
    } else {
      this.passwordStrengthLabel = this.translate.instant(
        'LOGIN.RESET.Form.Strong'
      );
      this.passwordStrengthClass =
        'password-meter__bar password-meter__bar--strong';
    }
  }

  isFieldInvalid(fieldName: string): boolean {
    const field = this.resetForm.get(fieldName);
    return !!(field && field.invalid && (field.dirty || field.touched));
  }

  getFieldError(fieldName: string): string {
    const field = this.resetForm.get(fieldName);
    if (!field || !field.errors) return '';

    if (field.errors['required']) {
      return this.translate.instant(
        fieldName === 'password'
          ? 'LOGIN.RESET.Form.PasswordRequired'
          : 'LOGIN.RESET.Form.ConfirmPasswordRequired'
      );
    }

    if (field.errors['minlength']) {
      return this.translate.instant('LOGIN.RESET.Form.PasswordMinLength');
    }

    if (field.errors['mismatch']) {
      return this.translate.instant('LOGIN.RESET.Form.PasswordMismatch');
    }

    return this.translate.instant('LOGIN.RESET.Form.Invalid');
  }

  onSubmit() {
    if (this.resetForm.invalid || this.submitting) {
      this.triggerShake();
      return;
    }

    this.submitting = true;

    const payload = {
      username: this.username,
      email: this.email,
      newPassword: this.resetForm.value.password,
      resetToken: this.token,
    };

    this.passwordService.confirmReset(payload).subscribe({
      next: (res) => {
        this.submitting = false;
        if (res.code === 200 && res.data) {
          this.toastService.success('Success', res.message, 2000);
          this.router.navigate(['/auth/login'], {
            queryParams: { message: res.message },
          });
        } else {
          this.toastService.error('Error', res.message, 2000);
          this.triggerShake();
        }
      },
      error: (err) => {
        this.submitting = false;
        this.triggerShake();
        const msg = err.error?.message || 'Something went wrong';
        this.toastService.error('Error', msg, 2000);
      },
    });
  }

  onCancel() {
    this.router.navigate(['/auth/login']);
  }

  private triggerShake() {
    this.shakeErrors = false;
    requestAnimationFrame(() => {
      this.shakeErrors = true;
      setTimeout(() => (this.shakeErrors = false), 320);
    });
  }
}
