import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-payment-settings',
  standalone: true,
  imports: [CommonModule, TranslateModule, FormsModule],
  templateUrl: './settings.html',
  styleUrl: './settings.scss'
})
export class Settings {
  settings = {
    currency: 'USD',
    taxRate: 10,
    autoRefund: true,
    emailNotifications: true,
    invoicePrefix: 'INV',
    paymentGateway: 'stripe'
  };

  saveSettings() {
    console.log('Saving settings:', this.settings);
  }
}
