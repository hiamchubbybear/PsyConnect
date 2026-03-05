import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, catchError, finalize, map, of } from 'rxjs';
import { environment } from '../../../environments/environment';
import { LoaderService } from '../loader/loader';

export interface Friend {
  profileId: string;
  firstName: string;
  lastName: string;
  avatarUri: string;
  role?: string;
}

@Injectable({ providedIn: 'root' })
export class FriendService {
  private http = inject(HttpClient);
  private baseUrl = environment.apiUrl;
  private version = environment.apiVersion;

  constructor(private loader: LoaderService) {}

  getMyFriends(): Observable<Friend[]> {
    const url = `${this.baseUrl}/${this.version}/profile/friends`;
    this.loader.show();

    return this.http.get<any>(url).pipe(
      map((res) => {
        console.log('[FriendService] getMyFriends response:', res);
        if (res.code !== 200 || !Array.isArray(res.data)) return [];

        return res.data.map((f: any) => ({
          profileId: f.profileId,
          firstName: f.firstName,
          lastName: f.lastName,
          avatarUri: f.avatarUri,
          role: f.role,
        })) as Friend[];
      }),
      finalize(() => this.loader.hide()),
    );
  }

  getReceivedRequests(): Observable<Friend[]> {
    const url = `${this.baseUrl}/${this.version}/profile/friends/received`;
    this.loader.show();

    return this.http.get<any>(url).pipe(
      map((res) => {
        console.log('[FriendService] getReceivedRequests response:', res);
        if (res.code !== 200 || !Array.isArray(res.data)) return [];

        return res.data.map((f: any) => ({
          profileId: f.profileId,
          firstName: f.firstName,
          lastName: f.lastName,
          avatarUri: f.avatarUri,
          role: f.role,
        })) as Friend[];
      }),
      finalize(() => this.loader.hide()),
    );
  }

  getFriendSuggestions(): Observable<Friend[]> {
    const url = `${this.baseUrl}/${this.version}/profile/friends/suggestions`;
    this.loader.show();

    return this.http.get<any>(url).pipe(
      map((res) => {
        console.log('[FriendService] getFriendSuggestions response:', res);
        if (res.code !== 200 || !Array.isArray(res.data)) return [];

        return res.data.map((f: any) => ({
          profileId: f.profileId,
          firstName: f.firstName,
          lastName: f.lastName,
          avatarUri: f.avatarUri,
          role: f.role,
        })) as Friend[];
      }),
      finalize(() => this.loader.hide()),
    );
  }
  sendFriendRequest(targetId: string): Observable<any> {
    const url = `${this.baseUrl}/${this.version}/profile/friend/request`;
    const body = { target: targetId };
    this.loader.show();

    return this.http.post<any>(url, body).pipe(
      map((res) => {
        console.log('[FriendService] sendFriendRequest response:', res);
        return res;
      }),
      finalize(() => this.loader.hide()),
    );
  }
  acceptFriendRequest(targetId: string): Observable<any> {
    const url = `${this.baseUrl}/${this.version}/profile/friend/accept`;
    const body = { target: targetId };
    this.loader.show();

    return this.http.post<any>(url, body).pipe(
      map((res) => {
        console.log('[FriendService] acceptFriendRequest response:', res);
        return res;
      }),
      finalize(() => this.loader.hide()),
    );
  }
  searchProfiles(query: string, page = 0, size = 10): Observable<Friend[]> {
    const url = `${this.baseUrl}/${this.version}/profile/search?query=${encodeURIComponent(query)}&page=${page}&size=${size}`;
    return this.http.get<any>(url).pipe(
      map((res) => {
        if (res.code !== 200 || !Array.isArray(res.data)) return [];
        return res.data.map((f: any) => ({
          profileId: f.profileId,
          firstName: f.firstName,
          lastName: f.lastName,
          avatarUri: f.avatarUri,
          role: f.role,
        })) as Friend[];
      }),
    );
  }

  getProfilesBatch(profileIds: string[]): Observable<Friend[]> {
    if (!profileIds || profileIds.length === 0) return of([]);
    const url = `${this.baseUrl}/${this.version}/profile/batch`;
    return this.http.post<any>(url, profileIds).pipe(
      map((res) => {
        if (res.code !== 200 || !Array.isArray(res.data)) return [];
        return res.data.map((f: any) => ({
          profileId: f.profileId,
          firstName: f.firstName,
          lastName: f.lastName,
          avatarUri: f.avatarUri,
          role: f.role,
        })) as Friend[];
      }),
      catchError(() => of([])),
    );
  }

  unfriend(targetId: string): Observable<any> {
    const url = `${this.baseUrl}/${this.version}/profile/friend/unfriend`;
    const body = { target: targetId };
    this.loader.show();

    return this.http.post<any>(url, body).pipe(
      map((res) => {
        console.log('[FriendService] unfriend response:', res);
        return res;
      }),
      finalize(() => this.loader.hide()),
    );
  }
}
