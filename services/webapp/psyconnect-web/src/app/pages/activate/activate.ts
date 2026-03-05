import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { ActivateService } from '../../services/activate/activate.service';
import { LoaderService } from '../../services/loader/loader';
import { ToastType } from '../../shared/toast/toast.model';
import { ToastService } from '../../shared/toast/toast.service';

@Component({
  standalone: true,
  selector: 'app-activate',
  templateUrl: './activate.html',
  styleUrls: ['./activate.scss'],
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    TranslateModule,
    RouterModule,
  ],
})
export class ActivateComponent implements OnInit {
  activateForm!: FormGroup;
  isLoading = false;
  email = '';

  constructor(
    private fb: FormBuilder,
    private route: ActivatedRoute,
    private router: Router,
    private activateService: ActivateService,
    private loaderService: LoaderService,
    private toastService: ToastService
  ) {}

  ngOnInit() {
    
    this.route.queryParams.subscribe((params) => {
      this.email = params['email'] || '';
    });

    this.activateForm = this.fb.group({
      email: [this.email, [Validators.required, Validators.email]],
      token: [
        '',
        [Validators.required, Validators.minLength(5), Validators.maxLength(5)],
      ],
    });
  }

  onCodeInput(event: any, index: number) {
    const input = event.target;
    const value = input.value;

    
    if (value && !/^\d$/.test(value)) {
      input.value = '';
      return;
    }

    
    const inputs = document.querySelectorAll('.code-input');
    let code = '';
    inputs.forEach((inp: any) => {
      code += inp.value || '';
    });
    this.activateForm.patchValue({ token: code });

    
    if (value && index < 4) {
      const nextInput = inputs[index + 1] as HTMLInputElement;
      nextInput?.focus();
    }
  }

  onCodeKeyDown(event: KeyboardEvent, index: number) {
    const input = event.target as HTMLInputElement;

    
    if (event.key === 'Backspace' && !input.value && index > 0) {
      const inputs = document.querySelectorAll('.code-input');
      const prevInput = inputs[index - 1] as HTMLInputElement;
      prevInput?.focus();
    }
  }

  onCodePaste(event: ClipboardEvent) {
    event.preventDefault();
    const pastedData = event.clipboardData?.getData('text') || '';
    const digits = pastedData.replace(/\D/g, '').slice(0, 5);

    const inputs = document.querySelectorAll('.code-input');
    digits.split('').forEach((digit, index) => {
      if (inputs[index]) {
        (inputs[index] as HTMLInputElement).value = digit;
      }
    });

    this.activateForm.patchValue({ token: digits });

    
    const lastIndex = Math.min(digits.length, 4);
    (inputs[lastIndex] as HTMLInputElement)?.focus();
  }

  onActivate() {
    if (this.activateForm.invalid) {
      this.activateForm.markAllAsTouched();
      return;
    }

    this.loaderService.show();
    this.isLoading = true;

    const request = {
      email: this.activateForm.value.email,
      token: this.activateForm.value.token,
      verifedTime: new Date().toISOString(),
    };

    this.activateService.activate(request).subscribe({
      next: (res) => {
        if (res.code === 200) {
          this.toastService.show(
            'activate.success.message',
            'activate.success.title',
            ToastType.Success
          );
          setTimeout(() => {
            this.loaderService.hide();
            this.router.navigate(['/auth/login']);
          }, 1500);
        }
      },
      error: (err) => {
        const errorCode = err.error?.code;

        if (errorCode === 103) {
          
          this.toastService.show(
            'activate.error.tokenExpired',
            'activate.error.title',
            ToastType.Error
          );
        } else if (errorCode === 604) {
          
          this.toastService.show(
            'activate.error.alreadyActivated',
            'activate.error.title',
            ToastType.Warning
          );
        } else {
          this.toastService.show(
            'activate.error.unexpected',
            'activate.error.title',
            ToastType.Error
          );
        }
        this.loaderService.hide();
        this.isLoading = false;
      },
      complete: () => {
        this.loaderService.hide();
        this.isLoading = false;
      },
    });
  }

  onResendCode() {
    if (!this.activateForm.value.email) {
      this.toastService.show(
        'activate.error.emailRequired',
        'activate.error.title',
        ToastType.Error
      );
      return;
    }

    this.loaderService.show();

    this.activateService
      .resendActivation(this.activateForm.value.email)
      .subscribe({
        next: (res) => {
          if (res.code === 200) {
            this.toastService.show(
              'activate.resend.success',
              'activate.success.title',
              ToastType.Success
            );
          }
          this.loaderService.hide();
        },
        error: () => {
          this.toastService.show(
            'activate.resend.error',
            'activate.error.title',
            ToastType.Error
          );
          this.loaderService.hide();
        },
      });
  }
}
