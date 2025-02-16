package dev.psyconnect.identity_service.dto.request;

import java.sql.Timestamp;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@Builder
public class ActivateAccountRequest {
    private String token;
    private String email;
    private Timestamp verifedTime;
}
