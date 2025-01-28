package com.example.IdentityService.dto.response;

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
