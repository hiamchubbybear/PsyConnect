
import { CommonModule } from '@angular/common';
import {
  Component,
  EventEmitter,
  Input,
  Output,
  TemplateRef,
} from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { InfoFieldComponent } from '../info-field/info-field';

export interface InfoField {
  label: string;
  value: string | TemplateRef<any>;
}

@Component({
  selector: 'app-info-block',
  standalone: true,
  imports: [CommonModule, InfoFieldComponent, TranslateModule],
  templateUrl: './info-block.html',
  styleUrls: ['./info-block.scss'],
})
export class InfoBlockComponent {
  @Input() title: string = '';
  @Input() fields: InfoField[] = [];
  @Input() showEdit: boolean = true;
  @Output() edit = new EventEmitter<void>();

  onEdit(): void {
    this.edit.emit();
  }
}
