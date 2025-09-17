import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-profile-overlay',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: './profile-overlay.html',
  styleUrls: ['./profile-overlay.scss'],
})
export class ProfileOverlayComponent implements OnInit {
  @Input() initialValue: string = '';

  @Output() saveField = new EventEmitter<{ field: string; value: string }>();
  @Output() cancel = new EventEmitter<void>();
  @Input() editValue: string = '';
  @Output() editValueChange = new EventEmitter<string>();
  @Input() editing: string | null = null;
  @Input() draftData: {
    firstName: string;
    middleName: string;
    lastName: string;
  } = {
    firstName: '',
    middleName: '',
    lastName: '',
  };
  @Output() saveName = new EventEmitter<void>();
  addressSuggestions: string[] = [];

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.editValue = this.initialValue;
  }

  onAddressInput(query: string) {
    if (query.length > 2) {
      const url = `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(
        query
      )}&format=json&addressdetails=1&limit=10`;

      this.http.get<any[]>(url).subscribe((data) => {
        this.addressSuggestions = data.map((item) => item.display_name);
      });
    } else {
      this.addressSuggestions = [];
    }
  }

  selectAddress(suggestion: string) {
    this.editValue = suggestion;
    this.addressSuggestions = [];
  }

  onSaveField(field: string) {
    this.saveField.emit({ field, value: this.editValue });
  }
}
