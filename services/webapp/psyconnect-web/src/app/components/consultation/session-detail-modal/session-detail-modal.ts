import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CalendarEvent } from 'angular-calendar';
import { IconComponent } from '../../../shared/ui-atoms/icon/icon.component';

@Component({
  selector: 'app-session-detail-modal',
  standalone: true,
  imports: [CommonModule, IconComponent],
  templateUrl: './session-detail-modal.html',
  styleUrls: ['./session-detail-modal.scss'],
})
export class SessionDetailModalComponent {
  @Input() event!: CalendarEvent;
  @Output() close = new EventEmitter<void>();

  get meta() {
    return this.event.meta;
  }

  closeModal() {
    this.close.emit();
  }

  get addressQuery(): string {
    return this.meta.address || `${this.meta.partnerName || 'Phòng khám PsyConnect'}, Việt Nam`;
  }

  openGoogleMaps() {
    if (this.meta.locationInfo?.coordinates?.length >= 2) {
      // Backend provides [latitude, longitude]
      const lat = this.meta.locationInfo.coordinates[0];
      const lng = this.meta.locationInfo.coordinates[1];
      window.open(`https://www.google.com/maps/search/?api=1&query=${lat},${lng}`, '_blank');
      return;
    }
    const address = encodeURIComponent(this.addressQuery);
    window.open(`https://www.google.com/maps/search/?api=1&query=${address}`, '_blank');
  }

  openAppleMaps() {
    if (this.meta.locationInfo?.coordinates?.length >= 2) {
      const lat = this.meta.locationInfo.coordinates[0];
      const lng = this.meta.locationInfo.coordinates[1];
      window.open(`http://maps.apple.com/?ll=${lat},${lng}&q=Phòng khám`, '_blank');
      return;
    }
    const address = encodeURIComponent(this.addressQuery);
    window.open(`http://maps.apple.com/?q=${address}`, '_blank');
  }

  copyAddress() {
    if (!this.meta.address) return;
    navigator.clipboard.writeText(this.meta.address);
    // Optional: Add a toast notification here if available
  }
}
