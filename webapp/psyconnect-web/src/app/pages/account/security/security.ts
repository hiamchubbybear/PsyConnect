import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import {
    FormBuilder,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { Router, RouterModule } from '@angular/router';

export interface SecurityItem {
  label: string;
  route: string;
}

export const securityItems: SecurityItem[] = [
  { label: 'Đổi mật khẩu', route: '/account/change-password' },
  { label: 'Xác thực hai yếu tố', route: '/account/2fa' },
  { label: 'Đăng nhập gần đây', route: '/account/recent-logins' },
  { label: 'Phương thức bảo mật', route: '/account/security-methods' },
];

export const securityCheckItems: SecurityItem[] = [
  { label: 'Thiết bị đang đăng nhập', route: '/account/devices' },
  {
    label: 'Hoạt động đăng nhập đáng ngờ',
    route: '/account/suspicious-activity',
  },
  { label: 'Email liên kết', route: '/account/email-check' },
  { label: 'Ứng dụng và quyền', route: '/account/app-permissions' },
];

@Component({
  selector: 'app-security',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, RouterModule],
  templateUrl: './security.html',
  styleUrls: ['./security.scss'],
})
export class SecuritySessionComponent {
  securityItems = securityItems;
  securityCheckItems = securityCheckItems;

  form: FormGroup;

  constructor(private fb: FormBuilder, private router: Router) {
    this.form = this.fb.group({
      currentPassword: ['', Validators.required],
      newPassword: ['', Validators.required],
      confirmPassword: ['', Validators.required],
    });
  }

  navigate(item: SecurityItem) {
    this.router.navigate([item.route]);
  }

  save() {
    if (this.form.valid) {
      console.log('Security saved:', this.form.value);
    } else {
      console.log('Form invalid');
      this.form.markAllAsTouched();
    }
  }
}
