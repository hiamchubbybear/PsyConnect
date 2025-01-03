package com.example.IdentityService.dto.response;

import com.example.IdentityService.model.RoleEntity;

import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class UserAccountCreationResponse {
    String username;
    String password;
    String email;
    String role;
}
