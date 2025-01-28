package com.example.IdentityService.dto.response;

import java.util.Set;

import com.example.IdentityService.configuration.ValidateLoginType;
import com.example.IdentityService.model.RoleEntity;

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
