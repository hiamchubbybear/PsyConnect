package dev.psyconnect.identity_service.dto.response;

import java.sql.Timestamp;
import java.util.UUID;

import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class DeleteAccountResponse {
    boolean isSuccess;
    UUID session;
    String secret;
    Timestamp timestamp;
}
