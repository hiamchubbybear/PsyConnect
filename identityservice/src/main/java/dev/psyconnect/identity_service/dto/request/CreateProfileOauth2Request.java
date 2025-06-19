package dev.psyconnect.identity_service.dto.request;

import dev.psyconnect.identity_service.configuration.ValidateLoginType;
import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class CreateProfileOauth2Request {
    String accountId;
    String profileId;
    String firstName;
    String lastName;
    String email;
    @ValidateLoginType
    String creationType;
    String avatarUri;
    String dob;
    String gender;
}
