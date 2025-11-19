package dev.psyconnect.identity_service.dto.response;

import lombok.*;

@Getter
@Builder
public class UpdateAccountResponse {
    private String username;
    private String password;
    private String email;
}
