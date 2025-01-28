package com.example.IdentityService.controller;

import javax.naming.AuthenticationException;

import org.springframework.web.bind.annotation.*;

import com.example.IdentityService.apiresponse.ApiResponse;
import com.example.IdentityService.apiresponse.CustomResponseWrapper;
import com.example.IdentityService.configuration.ValidateLoginType;
import com.example.IdentityService.dto.request.AuthenticationRequest;
import com.example.IdentityService.dto.response.AuthenticationResponse;
import com.example.IdentityService.service.AuthenticationService;

import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequestMapping("/auth")
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
@RequiredArgsConstructor
public class AuthenticationController {
    AuthenticationService authenticationService;

    @CustomResponseWrapper
    @PostMapping("/login")
    public ApiResponse<AuthenticationResponse> loginRequest(
            @RequestBody AuthenticationRequest authenticationRequest, @RequestParam @ValidateLoginType String loginType)
            throws AuthenticationException {
        return new ApiResponse<>(authenticationService.authenticate(authenticationRequest, loginType));
    }
}
