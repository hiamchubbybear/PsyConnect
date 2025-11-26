import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

interface PaymentMethod {
  id: string;
  type: 'card' | 'bank' | 'paypal';
  name: string;
  last4?: string;
  expiryDate?: string;
  isDefault: boolean;
}

@Component({
  selector: 'app-methods',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './methods.html',
  styleUrl: './methods.scss'
})
export class Methods implements OnInit {
  methods: PaymentMethod[] = [];

  ngOnInit() {
    this.loadMethods();
  }

  loadMethods() {
    this.methods = [
      {
        id: '1',
        type: 'card',
        name: 'Visa ending in 4242',
        last4: '4242',
        expiryDate: '12/2025',
        isDefault: true
      },
      {
        id: '2',
        type: 'card',
        name: 'Mastercard ending in 5555',
        last4: '5555',
        expiryDate: '06/2026',
        isDefault: false
      },
      {
        id: '3',
        type: 'paypal',
        name: 'PayPal (user@example.com)',
        isDefault: false
      },
      {
        id: '4',
        type: 'bank',
        name: 'Bank Account ****1234',
        last4: '1234',
        isDefault: false
      }
    ];
  }

  getMethodIcon(type: string): string {
    switch (type) {
      case 'card': return 'fa-credit-card';
      case 'bank': return 'fa-university';
      case 'paypal': return 'fa-paypal';
      default: return 'fa-wallet';
    }
  }

  setDefault(method: PaymentMethod) {
    this.methods.forEach(m => m.isDefault = false);
    method.isDefault = true;
  }

  deleteMethod(id: string) {
    if (confirm('Delete this payment method?')) {
      this.methods = this.methods.filter(m => m.id !== id);
    }
  }
}
