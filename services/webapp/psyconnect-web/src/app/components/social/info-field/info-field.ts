import { CommonModule } from '@angular/common';
import { Component, Input, TemplateRef } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-info-field',
  standalone: true,
  imports: [CommonModule ,TranslateModule],
  templateUrl: './info-field.html',
  styleUrls: ['./info-field.scss'],
})
export class InfoFieldComponent {
  @Input() label: string = '';
  @Input() value: string | TemplateRef<any> = '';

  isTemplateRef(value: any): value is TemplateRef<any> {
    return value instanceof TemplateRef;
  }
}
