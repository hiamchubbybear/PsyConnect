package dev.psyconnect.identity_service.dto.response.v1;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class AuthenticationResponseV1 {
    private String token;
    private String refreshToken;
    private boolean isSuccessful;
}
