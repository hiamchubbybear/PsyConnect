package dev.psyconnect.identity_service.dto.request;

import lombok.*;

@Getter@Setter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class AuthenticationRequest {
    private String username;
    private String password;
}
