package com.example.IdentityService.dto.request;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class AccountCreationRequest {
    String username;
    String password;
    String email;
    String role;
}
