package dev.psyconnect.identity_service.dto.response;

import lombok.*;
import lombok.experimental.FieldDefaults;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE)
public class CloudinarySignResponse {
    String signature;
    long timestamp;
    String apiKey;
    String cloudName;
}
