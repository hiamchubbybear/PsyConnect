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
      return 'LOGIN.REQUEST.Form.EmailRequired';
    }
    if (field.errors['email']) {
      return 'LOGIN.REQUEST.Form.EmailInvalid';
    }
    return 'LOGIN.REQUEST.Form.Invalid';
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
    setTimeout(() => {
      this.submitting = false;
      this.router.navigate(['/auth/request-reset-success'], {
        state: { email },
      });
    }, 1500);
  }

  onCancel() {
    this.router.navigate(['/auth/login']);
  }
}
