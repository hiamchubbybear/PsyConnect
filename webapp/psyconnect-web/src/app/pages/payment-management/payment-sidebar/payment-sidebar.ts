import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'payment-sidebar',
  standalone: true,
  imports: [CommonModule, TranslateModule, RouterModule],
  templateUrl: './payment-sidebar.html',
  styleUrl: '../../consultation/sidebar/sidebar.scss',
})
export class ConsultationSidebar {}
