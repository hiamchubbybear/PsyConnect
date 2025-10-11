import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import {
    Component,
    EventEmitter,
    Input,
    OnDestroy,
    OnInit,
    Output,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import {
    Subject,
    Subscription,
    debounceTime,
    distinctUntilChanged,
    of,
    switchMap,
} from 'rxjs';

@Component({
  selector: 'app-profile-overlay',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: './profile-overlay.html',
  styleUrls: ['./profile-overlay.scss'],
})
export class ProfileOverlayComponent implements OnInit, OnDestroy {
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

  private searchSubject = new Subject<string>();
  private subscription!: Subscription;

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.editValue = this.initialValue;

    this.subscription = this.searchSubject
      .pipe(
        debounceTime(400),
        distinctUntilChanged(),
        switchMap((query) => {
          if (query.length > 2) {
            const url = `https://nominatim.openstreetmap.org/search?q=${encodeURIComponent(
              query
            )}&format=json&addressdetails=1&limit=10`;
            return this.http.get<any[]>(url);
          } else {
            return of([]);
          }
        })
      )
      .subscribe((data) => {
        this.addressSuggestions = data.map((item: any) => item.display_name);
      });
  }

  onAddressInput(query: string) {
    this.searchSubject.next(query);
  }

  selectAddress(suggestion: string) {
    this.editValue = suggestion;
    this.addressSuggestions = [];
  }

  onSaveField(field: string) {
    this.saveField.emit({ field, value: this.editValue });
  }

  ngOnDestroy() {
    if (this.subscription) {
      this.subscription.unsubscribe();
    }
  }
}
