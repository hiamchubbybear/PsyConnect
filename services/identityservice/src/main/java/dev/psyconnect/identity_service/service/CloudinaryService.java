package dev.psyconnect.identity_service.service;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Map;
import java.util.TreeMap;
import java.util.stream.Collectors;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.crypto.codec.Hex;
import org.springframework.stereotype.Service;

import dev.psyconnect.identity_service.dto.request.CloudinarySignRequest;
import dev.psyconnect.identity_service.dto.response.CloudinarySignResponse;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE)
@Slf4j
public class CloudinaryService {

    @Value("${cloudinary.cloud-name}")
    String cloudName;

    @Value("${cloudinary.api-key}")
    String apiKey;

    @Value("${cloudinary.api-secret}")
    String apiSecret;

    public CloudinarySignResponse generateSignature(CloudinarySignRequest request) {
        Map<String, Object> params = request.getParams();

        // Cloudinary requires parameters to be sorted alphabetically
        TreeMap<String, Object> sortedParams = new TreeMap<>(params);

        String stringToSign = sortedParams.entrySet().stream()
                .map(e -> e.getKey() + "=" + e.getValue())
                .collect(Collectors.joining("&"));

        // Append API Secret
        stringToSign += apiSecret;

        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-1");
            byte[] hash = digest.digest(stringToSign.getBytes(StandardCharsets.UTF_8));
            String signature = new String(Hex.encode(hash));

            long timestamp = Long.parseLong(params.get("timestamp").toString());

            return CloudinarySignResponse.builder()
                    .signature(signature)
                    .timestamp(timestamp)
                    .apiKey(apiKey)
                    .cloudName(cloudName)
                    .build();
        } catch (NoSuchAlgorithmException e) {
            log.error("Error generating Cloudinary signature", e);
            throw new RuntimeException("Could not generate signature");
        }
    }
}
