package dev.psyconnect.identity_service.dto.request;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class GoogleAuthenticationRequest {
    private String email;
}
