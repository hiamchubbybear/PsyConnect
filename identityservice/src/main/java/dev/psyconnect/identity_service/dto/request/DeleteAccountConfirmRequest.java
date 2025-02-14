package dev.psyconnect.identity_service.dto.request;

import java.util.UUID;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class DeleteAccountConfirmRequest {
    String email;
    UUID session;
    String token;
}
