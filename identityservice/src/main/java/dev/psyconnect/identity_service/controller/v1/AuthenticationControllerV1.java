package dev.psyconnect.identity_service.controller.v1;

import java.text.ParseException;

import org.springframework.web.bind.annotation.*;

import com.nimbusds.jose.JOSEException;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.configuration.ValidateProviderType;
import dev.psyconnect.identity_service.dto.request.AuthenticationRequest;
import dev.psyconnect.identity_service.dto.response.v1.AuthenticationResponseV1;
import dev.psyconnect.identity_service.service.AuthenticationService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController("authenticationControllerV1")
@RequestMapping("/v1/auth")
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
@RequiredArgsConstructor
public class AuthenticationControllerV1 {
    AuthenticationService authenticationService;

    @PostMapping("/login")
    public ApiResponse<AuthenticationResponseV1> loginRequest(
            @RequestBody AuthenticationRequest authenticationRequest,
            @RequestParam("platform") String platform,
            @ValidateProviderType @RequestParam String provider) {
        return new ApiResponse<>(authenticationService.authenticateV1(authenticationRequest, provider, platform));
    }

    @PostMapping("/refresh")
    public ApiResponse<AuthenticationResponseV1> refreshToken(
            @RequestParam("token") String refreshToken,
            @RequestParam("username") String username,
            @RequestParam("platform") String platform,
            @ValidateProviderType @RequestParam String provider)
            throws ParseException, JOSEException {
        return new ApiResponse<>(authenticationService.refreshTokenV1(username, refreshToken, provider, platform));
    }
}
