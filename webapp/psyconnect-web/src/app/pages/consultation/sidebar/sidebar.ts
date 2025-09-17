import { CommonModule } from '@angular/common';
import {
    Component
} from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'consultation-sidebar',
  imports: [CommonModule, TranslateModule, RouterModule],
  templateUrl: './sidebar.html',
  styleUrl: './sidebar.scss',
})
export class ConsultationSidebar {}
