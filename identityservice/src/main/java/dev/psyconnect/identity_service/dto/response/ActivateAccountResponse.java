package dev.psyconnect.identity_service.dto.response;

import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

@Getter
@Setter
@Builder
public class ActivateAccountResponse {
    boolean isSuccess;
    String message;
}
