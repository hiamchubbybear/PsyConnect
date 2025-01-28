package com.example.IdentityService.dto.request;

import com.example.IdentityService.configuration.ValidateLoginType;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class CreateProfileOauth2GoogleRequest {
    String email;

    @ValidateLoginType
    String creationType;

    String avatarUri;
}
