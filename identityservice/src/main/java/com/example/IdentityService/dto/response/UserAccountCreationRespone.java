package com.example.IdentityService.dto.response;

import com.example.IdentityService.model.RoleEntity;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class UserAccountCreationRespone {
    String username;
    String password;
    String email;
    RoleEntity role;
}
