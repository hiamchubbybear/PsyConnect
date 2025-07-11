package dev.psyconnect.identity_service.dto.response;

import java.sql.Timestamp;

import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class DeleteAccountResponse {
    boolean isSuccess;
    String session;
    String secret;
    Timestamp timestamp;
}
