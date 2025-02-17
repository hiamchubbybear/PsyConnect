package dev.psyconnect.identity_service.dto.request;

import lombok.*;
import org.hibernate.annotations.GeneratedColumn;

@Getter
@Builder
public class UpdateAccountRequest {
    private String username;
    private String password;
    private String email;
}
