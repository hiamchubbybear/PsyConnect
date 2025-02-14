package dev.psyconnect.identity_service.dto.request;

import lombok.*;

@Getter
@AllArgsConstructor
@Setter
@NoArgsConstructor
@Builder
public class DeleteNotificationRequest {
    String token;
    String secret;
    String email;
}
