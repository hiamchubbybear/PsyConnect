package dev.psyconnect.identity_service.dto.request;

import dev.psyconnect.identity_service.configuration.ValidateLoginType;
import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class AuthenticationRequest {
    private String username;
    private String password;

    @ValidateLoginType
    private String loginType;
}
