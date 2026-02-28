import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import * as crypto from 'crypto-js';
import { environment } from '../../../environments/environment';
import { SKIP_AUTH } from '../auth/auth.interceptor';

@Injectable({
  providedIn: 'root',
})
export class CloudinaryService {
  private cloudName = `${environment.cloudName}`;
  private apiKey = `${environment.cloudinaryApiKey}`;
  private apiSecret = `${environment.cloudinaryApiSecret}`;

  constructor(private http: HttpClient) {}

  async uploadImage(imageFile: File, username: string): Promise<string> {
    try {
      const timestamp = Math.round(new Date().getTime() / 1000);
      // Sanitize username to avoid spaces or special chars in public_id
      const sanitizedUsername = username.replace(/[^a-zA-Z0-9]/g, '_').toLowerCase();
      const basePublicId = `avatar_${sanitizedUsername}`;

      const uploadParams: any = {
        timestamp: timestamp,
        public_id: basePublicId,
        folder: 'avatars',
        overwrite: true,
        invalidate: true,
      };

      // If uploadPreset is defined, it MUST be part of the signature if we send it
      if (environment.uploadPreset) {
        uploadParams.upload_preset = environment.uploadPreset;
      }

      const signature = this.generateSignature(uploadParams);
      
      console.log('[DEBUG CLOUDINARY UPLOAD]', {
        cloudName: this.cloudName,
        apiKey: this.apiKey ? `${this.apiKey.substring(0, 4)}...` : 'MISSING',
        apiSecretSet: !!this.apiSecret && this.apiSecret !== 'PLACEHOLDER_API_SECRET',
        timestamp,
        basePublicId,
        uploadParams,
        signature
      });

      const formData = new FormData();
      formData.append('file', imageFile);
      formData.append('api_key', this.apiKey);
      formData.append('timestamp', timestamp.toString());
      formData.append('signature', signature);
      formData.append('public_id', basePublicId);
      formData.append('folder', 'avatars');
      formData.append('overwrite', 'true');
      formData.append('invalidate', 'true');
      
      if (environment.uploadPreset) {
        formData.append('upload_preset', environment.uploadPreset);
      }

      const url = `https://api.cloudinary.com/v1_1/${this.cloudName}/image/upload`;

      const response = await this.http
        .post<any>(url, formData, {
          context: new HttpContext().set(SKIP_AUTH, true),
        })
        .toPromise();

      if (!response || !response.secure_url) {
        throw new Error(
          'Invalid response from Cloudinary - missing secure_url',
        );
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

  private generateSignature(params: any): string {
    const sortedParams = Object.keys(params)
      .sort()
      .map((key) => `${key}=${params[key]}`)
      .join('&');

    const stringToSign = sortedParams + this.apiSecret;
    // Cloudinary expects SHA-1 or SHA-256 hex string
    return crypto.SHA1(stringToSign).toString();
  }
}
