package dev.psyconnect.identity_service.dto.request;

import lombok.*;

@Setter
@Getter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class RequestActivationAccount {
    String email;
    String username;
    String fullname;
}
