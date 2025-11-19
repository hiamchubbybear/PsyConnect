package dev.psyconnect.identity_service.dto.request;

import lombok.*;

@Getter
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Setter
public class UpdateAccountNotificatorRequest {
    String username;
    String email;
}
