import { CommonModule } from '@angular/common';
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
import { CloudinaryService } from '../../services/cloudinary/cloudinary.service';
import {
    RegisterRequest,
    RegisterService,
} from '../../services/signup/register';
import { ToastType } from '../../shared/toast/toast.model';
import { ToastService } from '../../shared/toast/toast.service';
@Component({
  standalone: true,
  selector: 'signup',
  templateUrl: './signup.html',
  styleUrls: ['./signup.scss'],
  imports: [CommonModule, ReactiveFormsModule, FormsModule],
})
export class MultiStepRegisterComponent implements OnInit, AfterViewInit {
  @ViewChild('emailInput') emailInputRef!: ElementRef;
  @ViewChild('usernameInput') usernameInputRef!: ElementRef;
  @ViewChild('imageInput') imageInputRef!: ElementRef;
  @ViewChild('dobInput') dobInputRef!: ElementRef;
  @ViewChild('addressInput') addressInputRef!: ElementRef;
  @ViewChild('firstnameInput') firstnameInput!: ElementRef;
  ngAfterViewInit() {
    this.focusCurrentStepInput();
  }
  focusCurrentStepInput() {
    console.log(this.currentStep);
    setTimeout(() => {
      switch (this.currentStep) {
        case 1:
          this.imageInputRef?.nativeElement.focus();
          break;
        case 2:
          this.firstnameInput.nativeElement.focus();
          break;
        case 3:
          this.dobInputRef.nativeElement.focus();
          break;
        case 4:
          this.addressInputRef.nativeElement.focus();
          break;
        case 5:
          this.emailInputRef.nativeElement.focus();
          break;
        case 6:
          this.usernameInputRef.nativeElement.focus();
          break;
      }
    });
  }
  currentStep = 0;
  totalSteps = 7;
  isLoading = false;
  selectedImage: File | null = null;
  imagePreview: string | null = null;
  addressSuggestions: string[] = [];

  avatarForm!: FormGroup;
  nameForm!: FormGroup;
  personalForm!: FormGroup;
  addressForm!: FormGroup;
  emailForm!: FormGroup;
  roleForm!: FormGroup;
  credentialsForm!: FormGroup;

  constructor(
    private fb: FormBuilder,
    private cloudinaryService: CloudinaryService,
    private registerService: RegisterService,
    private toastService: ToastService,
    private cdr: ChangeDetectorRef
  ) {}

  ngOnInit() {
    this.initializeForms();
  }
  shakeErrors = false;

  private triggerShake() {
    this.shakeErrors = false;
    requestAnimationFrame(() => {
      this.shakeErrors = true;
      setTimeout(() => (this.shakeErrors = false), 320);
    });
  }
  initializeForms() {
    this.avatarForm = this.fb.group({
      avatar: [''],
    });

    this.nameForm = this.fb.group({
      firstName: ['', [Validators.required, Validators.minLength(2)]],
      lastName: ['', [Validators.required, Validators.minLength(2)]],
    });

    this.personalForm = this.fb.group({
      dateOfBirth: ['', Validators.required],
      gender: ['', Validators.required],
    });

    this.addressForm = this.fb.group({
      address: ['', Validators.required],
    });

    this.emailForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
    });

    this.roleForm = this.fb.group({
      role: ['', Validators.required],
    });

    this.credentialsForm = this.fb.group(
      {
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
      { validators: this.passwordMatchValidator }
    );
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

  getCurrentForm(): FormGroup {
    const forms = [
      this.avatarForm,
      this.nameForm,
      this.personalForm,
      this.addressForm,
      this.emailForm,
      this.roleForm,
      this.credentialsForm,
    ];
    return forms[this.currentStep];
  }

  nextStep() {
    if (this.currentStep < this.totalSteps - 1) {
      if (this.getCurrentForm().valid) {
        this.currentStep++;
      } else {
        this.getCurrentForm().markAllAsTouched();
        this.triggerShake();
        this.scrollToFirstInvalid();
      }
    }
  }

  previousStep() {
    if (this.currentStep > 0) {
      this.currentStep--;
    }
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
      this.personalForm.get('dateOfBirth')?.setErrors({ underAge: true });
    }
  }

  onAddressInput(event: any) {
    const query = event.target.value;
    if (query.length > 2) {
      this.addressSuggestions = [
        `${query}, Ho Chi Minh City, Vietnam`,
        `${query}, Hanoi, Vietnam`,
        `${query}, Da Nang, Vietnam`,
      ];
    } else {
      this.addressSuggestions = [];
    }
  }

  selectAddress(address: string) {
    this.addressForm.get('address')?.setValue(address);
    this.addressSuggestions = [];
  }

  async submitFinalRegister() {
    if (this.getCurrentForm().invalid) {
      this.getCurrentForm().markAllAsTouched();
      this.triggerShake();
      return;
    }

    this.isLoading = true;

    let avatarUri = '';
    if (this.selectedImage) {
      const uploadedUrl = await this.cloudinaryService.uploadImage(
        this.selectedImage,
        this.credentialsForm.value.username
      );
      if (!uploadedUrl) {
        this.toastService.show(
          'Image upload failed!',
          'Error',
          ToastType.Error
        );
        this.isLoading = false;
        return;
      }
      avatarUri = uploadedUrl;
    }

    const request: RegisterRequest = {
      username: this.credentialsForm.value.username.toLowerCase(),
      password: this.credentialsForm.value.password,
      firstName: this.nameForm.value.firstName,
      lastName: this.nameForm.value.lastName,
      address: this.addressForm.value.address,
      gender: this.personalForm.value.gender,
      email: this.emailForm.value.email,
      role: this.roleForm.value.role,
      avatarUri,
      dob: this.personalForm.value.dateOfBirth,
    };

    this.registerService.register(request).subscribe({
      next: (res) => {
        if (res.code === 200) {
          this.toastService.show(
            'Register success',
            'Success',
            ToastType.Success
          );
        }
      },
      error: (err) => {
        const status = err.status;
        if (status === 409) {
          const errorMessage = err.error.message || 'Register failed';
          const errorCode = err.error.code;
          console.log('Error code ', errorCode);
          this.toastService.show(errorMessage, 'Error', ToastType.Error);
          if (errorCode === 202) {
            const control = this.credentialsForm.get('username');
            const currentErrors = control?.errors || {};
            control?.setErrors({
              ...currentErrors,
              exists: errorMessage || 'Username already exists',
            });
            this.emailForm.markAllAsTouched();
            this.cdr.detectChanges();
          } else if (errorCode === 201) {
            this.currentStep -= 2;
            const control = this.emailForm.get('email');
            const currentErrors = control?.errors || {};
            control?.setErrors({
              ...currentErrors,
              exists: errorMessage || 'Email already exists',
            });
            this.emailForm.markAllAsTouched();
            this.cdr.detectChanges();
          }
          this.isLoading = false;
        } else {
          this.toastService.show(
            'Unexpected error occurred',
            'Error',
            ToastType.Error
          );
          this.isLoading = false;
        }
      },
      complete: () => {
        this.isLoading = false;
      },
    });
  }
  get progressPercentage(): number {
    return ((this.currentStep + 1) / this.totalSteps) * 100;
  }
  private scrollToFirstInvalid() {
    setTimeout(() => {
      const el = document.querySelector('.form-input.error');
      el?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }, 0);
  }
}
