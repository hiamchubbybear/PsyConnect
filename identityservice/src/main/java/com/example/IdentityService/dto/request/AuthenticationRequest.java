package com.example.IdentityService.dto.request;

import com.example.IdentityService.configuration.ValidateLoginType;

import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class AuthenticationRequest {
    private String username;
    private String password;

    @ValidateLoginType
    private String loginType;
}
