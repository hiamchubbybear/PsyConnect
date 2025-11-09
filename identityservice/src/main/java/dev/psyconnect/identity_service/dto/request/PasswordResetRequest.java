package dev.psyconnect.identity_service.dto.request;

import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
public class PasswordResetRequest {
    String newPassword;
    String email;
    String resetToken;
}
