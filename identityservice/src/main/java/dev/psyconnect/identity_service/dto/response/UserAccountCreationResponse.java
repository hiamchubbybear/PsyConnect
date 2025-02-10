package dev.psyconnect.identity_service.dto.response;

import java.util.Set;

import dev.psyconnect.identity_service.model.RoleEntity;
import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class UserAccountCreationResponse {
    String username;
    String email;
    Set<RoleEntity> role;
    String imageUri;
}
