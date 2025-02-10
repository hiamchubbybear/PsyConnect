package dev.psyconnect.identity_service.dto.request;

import lombok.*;
import lombok.experimental.FieldDefaults;

@Getter
@Setter
@Builder
@FieldDefaults(level = AccessLevel.PRIVATE)
@AllArgsConstructor
@NoArgsConstructor
public class CreateAccountNotificationRequest {
    String code;
    String email;
    String username;
    String fullname;
}
