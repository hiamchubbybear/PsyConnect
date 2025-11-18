import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, finalize, map } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Friend } from '../../models/chat.models';
import { LoaderService } from '../loader/loader';

@Injectable({ providedIn: 'root' })
export class FriendService {
  private http = inject(HttpClient);
  private baseUrl = environment.apiUrl;
  private version = environment.apiVersion;

  constructor(private loader: LoaderService) {}

  getMyFriends(): Observable<Friend[]> {
    const url = `${this.baseUrl}/${this.version}/profile/friends/me`;
    this.loader.show();

    return this.http.get<any>(url).pipe(
      map((res) => {
        console.log('[FriendService] raw response:', res);
        if (res.code !== 200 || !Array.isArray(res.data)) return [];

        const mapped = res.data.map((friend: any) => ({
          profileId: friend.profileId,
          firstName: friend.firstName,
          lastName: friend.lastName,
          avatarUri: friend.avatarUri,
        })) as Friend[];

        console.log('[FriendService] mapped friends:', mapped);
        return mapped;
      }),
      finalize(() => this.loader.hide())
    );
  }
}
