import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { IncomingCall } from '../../services/webrtc/webrtc-signaling.service';

@Component({
  selector: 'app-incoming-call',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './incoming-call.component.html',
  styleUrls: ['./incoming-call.component.scss'],
})
export class IncomingCallComponent {
  @Input() call!: IncomingCall;
  @Output() accept = new EventEmitter<void>();
  @Output() reject = new EventEmitter<void>();
  @Output() remind = new EventEmitter<void>();

  onAccept() {
    this.accept.emit();
  }

  onReject() {
    this.reject.emit();
  }

  onRemind() {
    
    this.remind.emit();
  }
}
