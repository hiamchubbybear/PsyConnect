import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import {
  AfterViewInit,
  ChangeDetectorRef,
  Component,
  ElementRef,
  OnInit,
  ViewChild,
} from '@angular/core';
import {
  FormBuilder,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { debounceTime, distinctUntilChanged } from 'rxjs';
import { CloudinaryService } from '../../services/cloudinary/cloudinary.service';
import { LoaderService } from '../../services/loader/loader';
import {
  RegisterRequest,
  RegisterService,
} from '../../services/signup/register';
import { ToastType } from '../../shared/toast/toast.model';
import { ToastService } from '../../shared/toast/toast.service';

import { PsyButtonComponent } from '../../shared/ui-atoms/button/psy-button.component';
import { CheckboxComponent } from '../../shared/ui-atoms/checkbox/checkbox.component';
import { InputComponent } from '../../shared/ui-atoms/input/input.component';

@Component({
  standalone: true,
  selector: 'signup',
  templateUrl: './signup.html',
  styleUrls: ['./signup.scss'],
  imports: [
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    TranslateModule,
    InputComponent,
    CheckboxComponent,
    PsyButtonComponent,
  ],
})
export class MultiStepRegisterComponent implements OnInit, AfterViewInit {
  @ViewChild('imageInput') imageInputRef!: ElementRef;

  ngAfterViewInit() {
    // Auto-focus first input
    setTimeout(() => {
      const firstInput = document.querySelector(
        'input:not([type="file"])',
      ) as HTMLElement;
      firstInput?.focus();
    }, 100);
  }

  isLoading = false;
  selectedImage: File | null = null;
  imagePreview: string | null = null;
  addressSuggestions: String[] = [];
  isLoadingAddress = false;
  shakeErrors = false;

  // Combined form
  registerForm!: FormGroup;

  constructor(
    private fb: FormBuilder,
    private cloudinaryService: CloudinaryService,
    private registerService: RegisterService,
    private toastService: ToastService,
    private cdr: ChangeDetectorRef,
    private http: HttpClient,
    private router: Router,
    private loaderService: LoaderService,
  ) {}

  ngOnInit() {
    this.initializeForm();
    this.setupAddressAutocomplete();
    this.setupProgressTracking();
  }

  private triggerShake() {
    this.shakeErrors = false;
    requestAnimationFrame(() => {
      this.shakeErrors = true;
      setTimeout(() => (this.shakeErrors = false), 320);
    });
  }

  initializeForm() {
    this.registerForm = this.fb.group(
      {
        // Avatar (optional)
        avatar: [''],

        // Name fields
        firstName: ['', [Validators.required, Validators.minLength(2)]],
        lastName: ['', [Validators.required, Validators.minLength(2)]],

        // Personal info
        dateOfBirth: ['', Validators.required],
        gender: ['', Validators.required],

        // Address
        address: ['', Validators.required],

        // Email
        email: ['', [Validators.required, Validators.email]],

        // Role
        role: ['', Validators.required],

        // Credentials
        username: ['', [Validators.required, Validators.minLength(4)]],
        password: [
          '',
          [
            Validators.required,
            Validators.minLength(8),
            this.passwordValidator,
          ],
        ],
        confirmPassword: ['', Validators.required],
        acceptTerms: [false, Validators.requiredTrue],
      },
      { validators: this.passwordMatchValidator },
    );
  }

  setupAddressAutocomplete() {
    this.registerForm
      .get('address')!
      .valueChanges.pipe(debounceTime(500), distinctUntilChanged())
      .subscribe((query: string) => {
        if (query && query.length > 2) {
          this.isLoadingAddress = true;
          const url = `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(
            query,
          )}&format=json&addressdetails=1&limit=10`;

          this.http.get<any[]>(url).subscribe({
            next: (data) => {
              this.addressSuggestions = data.map((item) => item.display_name);
              this.isLoadingAddress = false;
            },
            error: () => {
              this.isLoadingAddress = false;
              this.addressSuggestions = [];
            },
          });
        } else {
          this.addressSuggestions = [];
          this.isLoadingAddress = false;
        }
      });
  }

  setupProgressTracking() {
    // Track form changes to update progress
    this.registerForm.valueChanges.subscribe(() => {
      this.cdr.detectChanges();
    });
  }

  passwordValidator(control: any) {
    const value = control.value;
    if (!value) return null;

    const hasNumber = /[0-9]/.test(value);
    const hasLetter = /[a-zA-Z]/.test(value);
    const hasSpecial = /[!@#$%^&*(),.?":{}|<>]/.test(value);

    if (hasNumber && hasLetter && hasSpecial) {
      return null;
    }
    return { passwordStrength: true };
  }

  passwordMatchValidator(group: FormGroup) {
    const password = group.get('password');
    const confirmPassword = group.get('confirmPassword');

    if (
      password &&
      confirmPassword &&
      password.value !== confirmPassword.value
    ) {
      return { passwordMismatch: true };
    }
    return null;
  }

  onImageSelected(event: any) {
    const file = event.target.files[0];
    if (file) {
      this.selectedImage = file;
      const reader = new FileReader();
      reader.onload = (e: any) => {
        this.imagePreview = e.target.result;
      };
      reader.readAsDataURL(file);
    }
  }

  onDateChange(event: any) {
    const selectedDate = new Date(event.target.value);
    const today = new Date();
    const age = today.getFullYear() - selectedDate.getFullYear();

    if (age < 18) {
      this.registerForm.get('dateOfBirth')?.setErrors({ underAge: true });
    }
  }

  selectAddress(address: String) {
    this.registerForm.get('address')?.setValue(address);
    this.addressSuggestions = [];
  }

  async submitRegister() {
    if (this.registerForm.invalid) {
      this.registerForm.markAllAsTouched();
      this.triggerShake();
      this.scrollToFirstInvalid();
      return;
    }

    this.loaderService.show();
    this.isLoading = true;

    let avatarUri = '';
    if (this.selectedImage) {
      const uploadedUrl = await this.cloudinaryService.uploadImage(
        this.selectedImage,
        this.registerForm.value.username,
      );
      if (!uploadedUrl) {
        this.toastService.show(
          'toast.image_upload_failed',
          'toast.error',
          ToastType.Error,
        );
        this.isLoading = false;
        return;
      }
      avatarUri = uploadedUrl;
    }

    const request: RegisterRequest = {
      username: this.registerForm.value.username.toLowerCase(),
      password: this.registerForm.value.password,
      firstName: this.registerForm.value.firstName,
      lastName: this.registerForm.value.lastName,
      address: this.registerForm.value.address,
      gender: this.registerForm.value.gender,
      email: this.registerForm.value.email,
      role: this.registerForm.value.role,
      avatarUri,
      dob: this.registerForm.value.dateOfBirth,
    };

    this.registerService.register(request).subscribe({
      next: (res) => {
        if (res.code === 200) {
          this.toastService.show(
            'register.success.message',
            'register.success.title',
            ToastType.Success,
          );
          // Redirect to activate page after 1.5 seconds
          setTimeout(() => {
            this.loaderService.hide();
            this.router.navigate(['/activate'], {
              queryParams: { email: this.registerForm.value.email },
            });
          }, 1500);
        }
      },
      error: (err) => {
        console.log('=== Registration Error Details ===');
        console.log('Full error object:', err);
        console.log('Error type:', typeof err);
        console.log('Error constructor:', err.constructor.name);
        console.log('Error keys:', Object.keys(err));
        console.log('Error values:', Object.values(err));
        console.log('Error message property:', err.message);
        console.log('Error name property:', err.name);
        console.log('Status:', err.status);
        console.log('Status text:', err.statusText);
        console.log('Error body:', err.error);
        console.log('Error body type:', typeof err.error);
        console.log('Error code:', err.error?.code);
        console.log('Error message:', err.error?.message);
        console.log('Headers:', err.headers);
        console.log('URL:', err.url);

        // Log all enumerable properties
        console.log('All properties:');
        for (let key in err) {
          console.log(`  ${key}:`, err[key]);
        }

        // Try to parse if it's a string
        if (typeof err.error === 'string') {
          try {
            const parsed = JSON.parse(err.error);
            console.log('Parsed error:', parsed);
          } catch (e) {
            console.log('Could not parse error as JSON');
          }
        }
        console.log('==================================');
        const status = err.status;

        if (status === 409) {
          const errorCode = err.error?.code;
          const errorMessage = err.error?.message;

          console.log(
            'Handling 409 Conflict - Code:',
            errorCode,
            'Message:',
            errorMessage,
          );

          if (errorCode === 202) {
            const control = this.registerForm.get('username');
            control?.setErrors({
              ...control.errors,
              exists: true,
            });
            control?.markAsTouched();
            this.toastService.show(
              'register.error.usernameExists',
              'register.error.title',
              ToastType.Error,
            );
          } else if (errorCode === 201) {
            const control = this.registerForm.get('email');
            control?.setErrors({
              ...control.errors,
              exists: true,
            });
            control?.markAsTouched();
            this.toastService.show(
              'register.error.emailExists',
              'register.error.title',
              ToastType.Error,
            );
          } else {
            // Other conflict errors
            this.toastService.show(
              'register.error.conflict',
              'register.error.title',
              ToastType.Error,
            );
          }
          this.scrollToFirstInvalid();
          this.cdr.detectChanges();
        } else if (status === 400) {
          // Bad request
          this.toastService.show(
            'register.error.invalidData',
            'register.error.title',
            ToastType.Error,
          );
        } else {
          this.toastService.show(
            'register.error.unexpected',
            'register.error.title',
            ToastType.Error,
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

  get progressPercentage(): number {
    const fields = [
      'firstName',
      'lastName',
      'dateOfBirth',
      'gender',
      'address',
      'email',
      'role',
      'username',
      'password',
      'confirmPassword',
      'acceptTerms',
    ];

    const validFields = fields.filter((field) => {
      const control = this.registerForm.get(field);
      return control?.valid && control?.value;
    });

    return (validFields.length / fields.length) * 100;
  }

  get canSubmit(): boolean {
    return this.registerForm.valid && !this.isLoading;
  }

  private scrollToFirstInvalid() {
    setTimeout(() => {
      const el = document.querySelector(
        '.form-input.error, .form-select.error',
      );
      el?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }, 0);
  }
}
