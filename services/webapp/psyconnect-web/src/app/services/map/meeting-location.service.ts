import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { map, Observable, of } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface MeetingLocationSearchResult {
  address: string;
  latitude: number;
  longitude: number;
}

@Injectable({ providedIn: 'root' })
export class MeetingLocationService {
  private readonly http = inject(HttpClient);
  private readonly searchEndpoint = `${environment.apiUrl}/v1/consultation/locations/search`;
  private readonly reverseEndpoint = `${environment.apiUrl}/v1/consultation/locations/reverse`;
  private readonly styleUrl = 'https://tiles.openfreemap.org/styles/liberty';
  private readonly searchCache = new Map<string, MeetingLocationSearchResult[]>();

  getMapStyleUrl(): string {
    return this.styleUrl;
  }

  searchPlaces(query: string): Observable<MeetingLocationSearchResult[]> {
    const normalizedQuery = query.trim().toLowerCase();
    const cachedResults = this.searchCache.get(normalizedQuery);
    if (cachedResults) {
      return of(cachedResults);
    }

    const params = new HttpParams()
      .set('q', query.trim())
      .set('format', 'jsonv2')
      .set('limit', '6')
      .set('addressdetails', '1')
      .set('accept-language', 'vi,en');

    return this.http
      .get<{ data?: MeetingLocationSearchResult[] }>(this.searchEndpoint, {
        params,
      })
      .pipe(
        map((response) => {
          const results = response.data ?? [];
          this.searchCache.set(normalizedQuery, results);
          return results;
        }),
      );
  }

  reverseGeocode(
    latitude: number,
    longitude: number,
  ): Observable<MeetingLocationSearchResult> {
    const params = new HttpParams()
      .set('format', 'jsonv2')
      .set('lat', String(latitude))
      .set('lon', String(longitude));

    return this.http
      .get<{ data?: MeetingLocationSearchResult }>(this.reverseEndpoint, {
        params,
      })
      .pipe(
        map(
          (response) =>
            response.data || {
              address: `${latitude.toFixed(6)}, ${longitude.toFixed(6)}`,
              latitude,
              longitude,
            },
        ),
      );
  }
}
