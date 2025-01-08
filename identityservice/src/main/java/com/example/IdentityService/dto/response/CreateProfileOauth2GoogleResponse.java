package com.example.IdentityService.dto.response;

import com.example.IdentityService.configuration.ValidateLoginType;
import com.example.IdentityService.model.RoleEntity;
import lombok.*;
import org.springframework.security.core.parameters.P;
import org.springframework.stereotype.Service;

import java.util.Set;

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
