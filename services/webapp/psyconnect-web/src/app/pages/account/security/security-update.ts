import { CommonModule } from '@angular/common';
import { Component, HostListener, OnInit } from '@angular/core';
import {
    FormBuilder,
    FormGroup,
    FormsModule,
    ReactiveFormsModule,
} from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { environment } from '../../../../environments/environment';
import { ProfileFieldDropdownComponent } from '../../../components/field-row/field-row';
import { SecureStorageService } from '../../../encrypt/secure';
import { PasswordService } from '../../../services/auth/password.service';
import { ToastType } from '../../../shared/toast/toast.model';
import { ToastService } from '../../../shared/toast/toast.service';
import { ProfileModel } from '../profile/profile-model';

interface PasswordData {
  current: string;
  new: string;
  confirm: string;
}

interface SecurityMethods {
  email: string;
  phone: string;
  key: string;
}

@Component({
  selector: 'app-security-section',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    ProfileFieldDropdownComponent,
    TranslateModule,
  ],
  templateUrl: './security-update.html',
  styleUrls: ['./security-update.scss'],
})
export class SecuritySectionComponent implements OnInit {
  PROFILE_KEY = environment.profileKey;
  saveName() {
    throw new Error('Method not implemented.');
  }
  confirmChanges() {
    throw new Error('Method not implemented.');
  }
  securityItems = [
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.ChangePassword',
      key: 'changePassword',
    },
    { label: 'ACCOUNT_MANAGEMENT.Security.Items.TwoFA', key: '2fa' },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.RecentLogins',
      key: 'recentLogins',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.SecurityMethods',
      key: 'securityMethods',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.LoggedDevices',
      key: 'loggedDevices',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.SuspiciousActivity',
      key: 'suspiciousActivity',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.LinkedEmails',
      key: 'linkedEmails',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Items.ThirdPartyApps',
      key: 'thirdPartyApps',
    },
  ];

  securityCheckItems = [
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Check.LoggedDevices',
      key: 'loggedDevices',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Check.LinkedEmails',
      key: 'linkedEmails',
    },
    {
      label: 'ACCOUNT_MANAGEMENT.Security.Check.ThirdPartyApps',
      key: 'thirdPartyApps',
    },
  ];

  editing: string | null = null;
  passwordData: PasswordData = {
    current: '',
    new: '',
    confirm: '',
  };

  securityMethods: SecurityMethods = {
    email: '',
    phone: '',
    key: '',
  };

  newEmail = '';
  newPhone = '';
  newKey = '';

  recentLogins: Array<{ device: string; time: Date; location: string }> = [];
  loggedDevices: Array<{ name: string; lastActive: Date }> = [];
  suspiciousActivities: Array<{ message: string; time: Date }> = [];
  linkedEmails: Array<{ address: string }> = [];
  thirdPartyApps: Array<{ name: string; permissions: string[] }> = [];

  form!: FormGroup;

  constructor(
    private fb: FormBuilder,
    private toastService: ToastService,
    private passwordService: PasswordService,
    private secureStorage: SecureStorageService
  ) {}

  ngOnInit(): void {
    this.form = this.fb.group({
      passwordCurrent: [''],
      passwordNew: [''],
      passwordConfirm: [''],
      email: [''],
      phone: [''],
      key: [''],
    });

    this.recentLogins = [
      { device: 'Chrome - Windows', time: new Date(), location: 'Hà Nội' },
      { device: 'Safari - iPhone', time: new Date(), location: 'HCM' },
    ];

    this.loggedDevices = [
      { name: 'Laptop Huy', lastActive: new Date() },
      { name: 'iPhone Huy', lastActive: new Date() },
    ];

    this.suspiciousActivities = [
      { message: 'Đăng nhập từ IP lạ', time: new Date() },
    ];

    this.linkedEmails = [{ address: 'huy@example.com' }];
    this.thirdPartyApps = [{ name: 'App XYZ', permissions: ['Đọc', 'Ghi'] }];
  }

  openEdit(field: string) {
    this.editing = field;
    switch (field) {
      case 'changePassword':
        this.passwordResetEmail = this.userEmail;
        break;
      case 'securityMethods':
        this.securityMethods = { ...this.securityMethods };
        break;
    }
  }

  passwordResetEmail = '';
  submittingReset = false;

  get userEmail(): string {
    return (
      this.secureStorage.getItem<ProfileModel>(this.PROFILE_KEY)?.firstName ||
      'user@example.com'
    );
  }

  sendPasswordReset() {
    if (!this.passwordResetEmail) {
      this.toastService.show(
        'TOAST.no_email_to_reset',
        'TOAST.error',
        ToastType.Error
      );
      return;
    }

    this.submittingReset = true;

    this.passwordService
      .requestReset({ email: this.passwordResetEmail })
      .subscribe({
        next: (res) => {
          this.submittingReset = false;
          if (res.code === 200) {
            this.toastService.show(
              'TOAST.reset_email_sent',
              'TOAST.success',
              ToastType.Success
            );
            this.editing = null;
          } else {
            this.toastService.show(
              'TOAST.error_generic',
              'TOAST.error',
              ToastType.Error
            );
          }
        },
        error: (err) => {
          this.submittingReset = false;

          this.toastService.show(
            'TOAST.error_generic',
            'TOAST.error',
            ToastType.Error
          );
        },
      });
  }

  saveSecurityMethods() {
    this.toastService.show(
      'TOAST.security_method_updated',
      'TOAST.success',
      ToastType.Success
    );
    this.editing = null;
  }

  enable2FA(method: 'authApp' | 'sms' | 'securityKey') {
    this.toastService.show(
      'TOAST.2fa_enabled',
      'TOAST.success',
      ToastType.Success
    );
    this.editing = null;
  }

  logoutDevice(device: { name: string; lastActive: Date }) {
    this.toastService.show(
      'TOAST.device_logged_out',
      'TOAST.success',
      ToastType.Success
    );
    this.loggedDevices = this.loggedDevices.filter((d) => d !== device);
  }

  resolveAlert(alert: { message: string; time: Date }) {
    this.toastService.show(
      'TOAST.alert_verified',
      'TOAST.success',
      ToastType.Success
    );
    this.suspiciousActivities = this.suspiciousActivities.filter(
      (a) => a !== alert
    );
  }

  addEmail() {
    if (this.newEmail.trim() !== '') {
      this.linkedEmails.push({ address: this.newEmail });
      this.newEmail = '';

      this.toastService.show(
        'TOAST.email_added',
        'TOAST.success',
        ToastType.Success
      );
    }
  }

  removeEmail(email: { address: string }) {
    this.linkedEmails = this.linkedEmails.filter((e) => e !== email);
    this.toastService.show(
      'TOAST.email_deleted',
      'TOAST.success',
      ToastType.Success
    );
  }

  revokeApp(app: { name: string; permissions: string[] }) {
    this.thirdPartyApps = this.thirdPartyApps.filter((a) => a !== app);

    this.toastService.show(
      'TOAST.app_permission_revoked',
      'TOAST.success',
      ToastType.Success
    );
  }

  cancel() {
    this.editing = null;
  }

  @HostListener('document:keydown', ['$event'])
  onEsc(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      this.cancel();
    }
  }
  getFieldDisplay(key: string): string {
    switch (key) {
      case 'changePassword':
        return '••••••••';
      case '2fa':
        return this.editing === '2fa' ? 'Enabled' : 'Disabled';
      case 'recentLogins':
        return `${this.recentLogins.length} logins`;
      case 'loggedDevices':
        return `${this.loggedDevices.length} devices`;
      case 'suspiciousActivity':
        return `${this.suspiciousActivities.length} alerts`;
      case 'linkedEmails':
        return this.linkedEmails.map((e) => e.address).join(', ');
      case 'thirdPartyApps':
        return this.thirdPartyApps.map((a) => a.name).join(', ');
      case 'securityMethods':
        return '••••';
      default:
        return '-';
    }
  }
  closeOverlay() {
    this.editing = null;
  }
  savePassword() {
    this.toastService.show(
      'TOAST.password_updated',
      'TOAST.success',
      ToastType.Success
    );
    this.editing = null;
  }
  saveOverlay() {
    switch (this.editing) {
      case 'changePassword':
        this.savePassword();
        break;
      case 'securityMethods':
        this.saveSecurityMethods();
        break;
      case '2fa':
        this.toastService.show(
          'TOAST.2fa_updated',
          'TOAST.success',
          ToastType.Success
        );
        break;
      default:
        this.toastService.show(
          'TOAST.update_success',
          'TOAST.success',
          ToastType.Success
        );
        break;
    }
    this.closeOverlay();
  }
}
