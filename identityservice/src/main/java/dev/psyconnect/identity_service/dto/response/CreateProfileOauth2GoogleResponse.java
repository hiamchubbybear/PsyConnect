package dev.psyconnect.identity_service.dto.response;

import java.util.Set;

import dev.psyconnect.identity_service.configuration.ValidateLoginType;
import dev.psyconnect.identity_service.model.RoleEntity;
import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class CreateProfileOauth2GoogleResponse {
    String email;
    Boolean status;
    Set<RoleEntity> roles;

    @ValidateLoginType
    String creationType;

    String avatarUri;
}
