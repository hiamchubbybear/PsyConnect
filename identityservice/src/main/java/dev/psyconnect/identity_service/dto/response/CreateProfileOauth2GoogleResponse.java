package dev.psyconnect.identity_service.dto.response;

import java.util.Set;

import dev.psyconnect.identity_service.configuration.ValidateProviderType;
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

    @ValidateProviderType
    String creationType;

    String avatarUri;
}
