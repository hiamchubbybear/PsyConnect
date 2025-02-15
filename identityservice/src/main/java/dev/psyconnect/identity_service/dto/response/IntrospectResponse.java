package dev.psyconnect.identity_service.dto.response;

import lombok.*;

@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class IntrospectResponse {
    boolean isSuccessful;
}
