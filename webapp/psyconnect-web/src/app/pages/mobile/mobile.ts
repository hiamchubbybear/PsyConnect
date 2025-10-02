import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-mobile-required',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './mobile.html',
  styleUrls: ['./mobile.scss'],
})
export class MobileRequiredComponent {}
