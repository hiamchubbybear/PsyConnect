import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { SKIP_AUTH } from '../auth/auth.interceptor';
import { firstValueFrom } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class CloudinaryService {
  private cloudName = `${environment.cloudName}`;
  private apiUrl = `${environment.apiUrl}/v1/identity/cloudinary/sign`;

  constructor(private http: HttpClient) {}

  async uploadImage(imageFile: File, username: string): Promise<string> {
    try {
      const timestamp = Math.round(new Date().getTime() / 1000);
      const sanitizedUsername = username.replace(/[^a-zA-Z0-9]/g, '_').toLowerCase();
      const basePublicId = `avatar_${sanitizedUsername}`;

      // 1. Get signature from Backend
      const signParams: any = {
        timestamp: timestamp,
        public_id: basePublicId,
        folder: 'avatars',
        overwrite: true,
        invalidate: true,
      };

      if (environment.uploadPreset) {
        signParams.upload_preset = environment.uploadPreset;
      }

      console.log('[DEBUG CLOUDINARY] Requesting signature from backend...');
      const signResponse = await firstValueFrom(
        this.http.post<any>(this.apiUrl, { params: signParams })
      );

      if (!signResponse || !signResponse.data || !signResponse.data.signature) {
        throw new Error('Failed to get signature from backend');
      }

      const { signature, apiKey, cloudName } = signResponse.data;

      // 2. Upload to Cloudinary using the signature
      const formData = new FormData();
      formData.append('file', imageFile);
      formData.append('api_key', apiKey);
      formData.append('timestamp', timestamp.toString());
      formData.append('signature', signature);
      formData.append('public_id', basePublicId);
      formData.append('folder', 'avatars');
      formData.append('overwrite', 'true');
      formData.append('invalidate', 'true');
      
      if (environment.uploadPreset) {
        formData.append('upload_preset', environment.uploadPreset);
      }

      const uploadUrl = `https://api.cloudinary.com/v1_1/${cloudName}/image/upload`;

      console.log('[DEBUG CLOUDINARY] Uploading to Cloudinary...');
      const response = await firstValueFrom(
        this.http.post<any>(uploadUrl, formData, {
          context: new HttpContext().set(SKIP_AUTH, true),
        })
      );

      if (!response || !response.secure_url) {
        throw new Error('Invalid response from Cloudinary - missing secure_url');
      }

      return response.secure_url;
    } catch (error: any) {
      console.error('Detailed Cloudinary error:', error);
      
      let errorMessage = 'Cloudinary upload failed';
      if (error.error && error.error.error && error.error.error.message) {
        errorMessage += `: ${error.error.error.message}`;
      } else if (error.message) {
        errorMessage += `: ${error.message}`;
      }

      throw new Error(errorMessage);
    }
  }
}
