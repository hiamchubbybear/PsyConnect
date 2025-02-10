package dev.psyconnect.identity_service.dto.request;

import dev.psyconnect.identity_service.configuration.ValidateLoginType;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class CreateProfileOauth2GoogleRequest {
    String email;

    @ValidateLoginType
    String creationType;

    String avatarUri;
}
