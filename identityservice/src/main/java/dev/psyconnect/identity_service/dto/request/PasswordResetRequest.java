package dev.psyconnect.identity_service.dto.request;

import com.google.type.DateTime;
import io.grpc.netty.shaded.io.netty.util.internal.SuppressJava6Requirement;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.web.bind.annotation.GetMapping;

import java.lang.annotation.Target;

@Getter
@Setter
@NoArgsConstructor
public class PasswordResetRequest {
    String newPassword;
    String username;
    String email;
    String resetToken;
}
