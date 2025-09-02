import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import * as crypto from 'crypto-js';
import { environment_secret } from '../../../environments/environment.secret';

@Injectable({
  providedIn: 'root',
})
export class CloudinaryService {
  private cloudName = `${environment_secret.cloudName}`;
  private apiKey = `${environment_secret.cloudinaryApiKey}`;
  private apiSecret = `${environment_secret.cloudinaryApiSecret}`;

  constructor(private http: HttpClient) {}

  async uploadImage(imageFile: File, username: string): Promise<string> {
    try {
      const timestamp = Math.round(new Date().getTime() / 1000);
      const basePublicId = `avatar_${username.toLowerCase()}`;

      // Parameters for signed upload
      const uploadParams = {
        timestamp: timestamp,
        public_id: basePublicId,
        folder: 'avatars',
        overwrite: true,
        invalidate: true,
      };

      // Generate signature
      const signature = this.generateSignature(uploadParams);

      // Prepare form data
      const formData = new FormData();
      formData.append('file', imageFile);
      formData.append('api_key', this.apiKey);
      formData.append('timestamp', timestamp.toString());
      formData.append('signature', signature);
      formData.append('public_id', basePublicId);
      formData.append('folder', 'avatars');
      formData.append('overwrite', 'true');
      formData.append('invalidate', 'true');

      const url = `https://api.cloudinary.com/v1_1/${this.cloudName}/image/upload`;

      const response = await this.http.post<any>(url, formData).toPromise();

      if (!response || !response.secure_url) {
        throw new Error(
          'Invalid response from Cloudinary - missing secure_url'
        );
      }

      console.log('✅ Upload success - Overwritten:', {
        public_id: response.public_id,
        version: response.version,
      });

      return response.secure_url;
    } catch (error) {
      console.error('❌ Upload error:', error);

      if (error instanceof Error) {
        throw new Error(`Cloudinary upload failed: ${error.message}`);
      } else {
        throw new Error('Cloudinary upload failed: Unknown error');
      }
    }
  }

  private generateSignature(params: any): string {
    // Sort parameters alphabetically and create query string
    const sortedParams = Object.keys(params)
      .sort()
      .map((key) => `${key}=${params[key]}`)
      .join('&');

    // Add API secret and generate SHA1 hash
    const stringToSign = sortedParams + this.apiSecret;
    return crypto.SHA1(stringToSign).toString();
  }
}
