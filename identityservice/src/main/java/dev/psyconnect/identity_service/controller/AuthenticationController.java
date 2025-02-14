package dev.psyconnect.identity_service.controller;

import javax.naming.AuthenticationException;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.apiresponse.CustomResponseWrapper;
import dev.psyconnect.identity_service.configuration.ValidateLoginType;
import dev.psyconnect.identity_service.dto.request.AuthenticationRequest;
import dev.psyconnect.identity_service.dto.response.AuthenticationResponse;
import dev.psyconnect.identity_service.service.AuthenticationService;
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
