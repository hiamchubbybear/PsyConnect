package dev.psyconnect.identity_service.controller;

import java.text.ParseException;
import javax.naming.AuthenticationException;

import org.springframework.web.bind.annotation.*;

import com.nimbusds.jose.JOSEException;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.configuration.ValidateLoginType;
import dev.psyconnect.identity_service.dto.request.AuthenticationRequest;
import dev.psyconnect.identity_service.dto.response.AuthenticationResponse;
import dev.psyconnect.identity_service.dto.response.LogoutRequest;
import dev.psyconnect.identity_service.dto.response.LogoutResponse;
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

    @PostMapping("/login")
    public ApiResponse<AuthenticationResponse> loginRequest(
            @RequestBody AuthenticationRequest authenticationRequest, @RequestParam @ValidateLoginType String loginType)
            throws AuthenticationException {
        return new ApiResponse<>(authenticationService.authenticate(authenticationRequest, loginType));
    }

    @PostMapping("/introspect")
    public ApiResponse<LogoutResponse> introspectRequest(@RequestBody LogoutRequest introspectRequest)
            throws ParseException, JOSEException {
        return new ApiResponse<>(authenticationService.logout(introspectRequest));
    }

    // For access token validate
    @PostMapping("internal/valid/{token}")
    public Boolean introspectRequest(@PathVariable String token) throws ParseException, JOSEException {
        return (authenticationService.isTokenValid(token));
    }
}
