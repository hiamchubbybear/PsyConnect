package dev.psyconnect.identity_service.dto.request;

import java.sql.Timestamp;

import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class DeleteAccountRequest {
    String email;
    String username;
    String password;
    String token;
    Timestamp timestamp;
}
