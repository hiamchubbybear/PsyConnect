package com.example.profileservice.service;

import java.io.IOException;
import java.util.Base64;
import java.util.Map;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import com.cloudinary.Cloudinary;
import com.cloudinary.utils.ObjectUtils;

@Service
public class ImageService {

    private static final Logger log = LoggerFactory.getLogger(ImageService.class);
    private final Cloudinary cloudinary;

    public ImageService(@Value("${CLOUDINARY_URL}") String cloudinaryUrl) {
        this.cloudinary = new Cloudinary(cloudinaryUrl);
    }

    public Map<String, Object> uploadImageFromBase64(String base64Image, UUID imageId) throws IOException {
        try {
            if (base64Image.contains(",")) {
                base64Image = base64Image.split(",")[1];
            }
            byte[] decodedBytes = Base64.getDecoder().decode(base64Image);

            Map<String, Object> params = ObjectUtils.asMap(
                    "public_id", imageId.toString(),
                    "overwrite", true,
                    "resource_type", "image",
                    "folder", "uploads",
                    "format", "png",
                    "quality", "auto:low",
                    "width", 800,
                    "height", 800,
                    "crop", "limit");

            return cloudinary.uploader().upload("data:image/png;base64," + base64Image, params);
        } catch (IOException e) {
            log.error("Failed to upload image to Cloudinary: {}", e.getMessage());
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during image upload: {}", e.getMessage());
            throw new IOException("Failed to upload image", e);
        }
    }

    public Map<String, Object> uploadImageFromMultiPartedFile(MultipartFile file, UUID imageId) throws IOException {
        try {

            Map<String, Object> params = ObjectUtils.asMap(
                    "public_id", imageId.toString(),
                    "overwrite", true,
                    "resource_type", "image",
                    "folder", "uploads",
                    "format", "png",
                    "quality", "auto:low",
                    "width", 800,
                    "height", 800,
                    "crop", "limit");

            return cloudinary.uploader().upload(file.getBytes(), params);
        } catch (IOException e) {
            log.error("Failed to upload image to Cloudinary: {}", e.getMessage());
            throw e;
        } catch (Exception e) {
            log.error("Unexpected error during image upload: {}", e.getMessage());
            throw new IOException("Failed to upload image", e);
        }
    }
}
