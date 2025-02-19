package dev.psyconnect.identity_service.dto.request;

import lombok.*;

@Getter
@Builder
public class UpdateAccountRequest {
    private String username;
    private String password;
    private String email;
}
