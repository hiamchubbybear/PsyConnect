package dev.psyconnect.identity_service.dto.request;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class Oauth2AuthenticationRequest {
    private String email;
    private String sessionToken;
}
