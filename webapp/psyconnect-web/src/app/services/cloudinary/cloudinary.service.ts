import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment_secret } from '../../../environments/environment.secret';

@Injectable({ providedIn: 'root' })
export class CloudinaryService {
  private cloudName = `${environment_secret.cloudName}`;
  private uploadPreset = `${environment_secret.uploadPreset}`;

  constructor(private http: HttpClient) {}

  uploadImage(imageFile: File, username: string): Promise<string | null> {
    const formData = new FormData();
    formData.append('file', imageFile);
    formData.append('upload_preset', this.uploadPreset);
    formData.append('public_id', username.toLowerCase());
    formData.append('filename_override', username.toLowerCase());

    const url = `https://api.cloudinary.com/v1_1/${this.cloudName}/image/upload`;

    return this.http
      .post<any>(url, formData)
      .toPromise()
      .then((res) => res.secure_url as string)
      .catch((err) => {
        console.error('Upload error', err);
        return null;
      });
  }
}
