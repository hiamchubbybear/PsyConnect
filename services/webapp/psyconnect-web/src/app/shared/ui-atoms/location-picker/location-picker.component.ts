import { CommonModule } from '@angular/common';
import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  EventEmitter,
  OnDestroy,
  Output,
  ViewChild,
  inject,
} from '@angular/core';
import { FormControl, ReactiveFormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import maplibregl, { Map, Marker } from 'maplibre-gl';
import { catchError, firstValueFrom, of } from 'rxjs';
import {
  MeetingLocationSearchResult,
  MeetingLocationService,
} from '../../../services/map/meeting-location.service';

export interface PickedLocation {
  address: string;
  latitude: number;
  longitude: number;
}

@Component({
  selector: 'app-location-picker',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule],
  templateUrl: './location-picker.component.html',
  styleUrl: './location-picker.component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LocationPickerComponent implements AfterViewInit, OnDestroy {
  @ViewChild('mapContainer') mapContainer?: ElementRef<HTMLDivElement>;
  @Output() locationSelected = new EventEmitter<PickedLocation>();

  private readonly meetingLocationService = inject(MeetingLocationService);
  private map?: Map;
  private marker?: Marker;

  readonly searchControl = new FormControl('', { nonNullable: true });

  isResolving = false;
  isSearching = false;
  errorKey = '';
  selectedAddress = '';
  selectedLatitude?: number;
  selectedLongitude?: number;
  searchResults: MeetingLocationSearchResult[] = [];
  hasSubmittedSearch = false;

  ngAfterViewInit(): void {
    queueMicrotask(() => this.initMap());
  }

  ngOnDestroy(): void {
    this.map?.remove();
  }

  useCurrentLocation(): void {
    if (!navigator.geolocation) {
      this.errorKey = 'CONSULTATION.Booking.LocationPicker.Errors.Unsupported';
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords;
        this.setLocation(latitude, longitude);
      },
      () => {
        this.errorKey =
          'CONSULTATION.Booking.LocationPicker.Errors.CurrentLocationFailed';
      },
      { enableHighAccuracy: true, timeout: 10000 },
    );
  }

  selectSearchResult(result: MeetingLocationSearchResult): void {
    this.searchResults = [];
    this.hasSubmittedSearch = false;
    this.searchControl.setValue(result.address, { emitEvent: false });
    void this.setLocation(result.latitude, result.longitude, result.address);
  }

  onSearchInput(): void {
    this.errorKey = '';
    this.hasSubmittedSearch = false;

    if (this.searchControl.value.trim().length < 3) {
      this.searchResults = [];
    }
  }

  onSearchEnter(event: Event): void {
    event.preventDefault();
    void this.searchPlaces();
  }

  async searchPlaces(): Promise<void> {
    const trimmedQuery = this.searchControl.value.trim();
    this.errorKey = '';
    this.hasSubmittedSearch = true;

    if (trimmedQuery.length < 3 || this.isSearching) {
      return;
    }

    this.isSearching = true;
    try {
      this.searchResults = await firstValueFrom(
        this.meetingLocationService.searchPlaces(trimmedQuery).pipe(
          catchError(() => {
            this.errorKey =
              'CONSULTATION.Booking.LocationPicker.Errors.SearchFailed';
            return of<MeetingLocationSearchResult[]>([]);
          }),
        ),
      );
    } finally {
      this.isSearching = false;
    }
  }

  private initMap(): void {
    const container = this.mapContainer?.nativeElement;
    if (!container || this.map) return;

    this.map = new maplibregl.Map({
      container,
      style: this.meetingLocationService.getMapStyleUrl(),
      center: [106.7009, 10.7769],
      zoom: 13,
    });
    this.map.addControl(new maplibregl.NavigationControl(), 'top-left');
    this.map.on('click', (event) => {
      this.setLocation(event.lngLat.lat, event.lngLat.lng);
    });

    setTimeout(() => this.map?.resize(), 0);
  }

  private async setLocation(
    latitude: number,
    longitude: number,
    preferredAddress?: string,
  ): Promise<void> {
    this.errorKey = '';
    this.selectedLatitude = latitude;
    this.selectedLongitude = longitude;

    if (this.map) {
      this.map.flyTo({ center: [longitude, latitude], zoom: 16 });
      if (!this.marker) {
        const markerElement = this.createMarkerElement();
        this.marker = new maplibregl.Marker({
          element: markerElement,
          anchor: 'bottom',
        })
          .setLngLat([longitude, latitude])
          .addTo(this.map);
      } else {
        this.marker.setLngLat([longitude, latitude]);
      }
    }

    if (preferredAddress) {
      this.selectedAddress = preferredAddress;
      this.hasSubmittedSearch = false;
      this.locationSelected.emit({
        address: preferredAddress,
        latitude,
        longitude,
      });
      return;
    }

    this.isResolving = true;
    try {
      const data = await firstValueFrom(
        this.meetingLocationService
        .reverseGeocode(latitude, longitude)
        .pipe(
          catchError(() => {
            this.errorKey =
              'CONSULTATION.Booking.LocationPicker.Errors.ReverseGeocodeFailed';
            return of({
              address: `${latitude.toFixed(6)}, ${longitude.toFixed(6)}`,
              latitude,
              longitude,
            });
          }),
        ),
      );

      this.selectedAddress = data.address;
      this.hasSubmittedSearch = false;
      this.searchControl.setValue(data.address, { emitEvent: false });
      this.locationSelected.emit({
        address: data.address,
        latitude,
        longitude,
      });
    } finally {
      this.isResolving = false;
    }
  }

  private createMarkerElement(): HTMLDivElement {
    const markerElement = document.createElement('div');
    markerElement.style.width = '28px';
    markerElement.style.height = '28px';
    markerElement.style.borderRadius = '50% 50% 50% 0';
    markerElement.style.transform = 'rotate(-45deg)';
    markerElement.style.background =
      'linear-gradient(180deg, #4f8cff 0%, #2563eb 100%)';
    markerElement.style.border = '3px solid #ffffff';
    markerElement.style.boxShadow = '0 10px 24px rgba(37, 99, 235, 0.35)';

    const centerDot = document.createElement('div');
    centerDot.style.width = '8px';
    centerDot.style.height = '8px';
    centerDot.style.borderRadius = '999px';
    centerDot.style.background = '#ffffff';
    centerDot.style.position = 'absolute';
    centerDot.style.top = '8px';
    centerDot.style.left = '8px';
    centerDot.style.transform = 'rotate(45deg)';

    markerElement.appendChild(centerDot);
    return markerElement;
  }
}
