package dev.psyconnect.identity_service.dto.response;

import lombok.*;

@Builder
@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
public class LogoutRequest {
    String token;
}
