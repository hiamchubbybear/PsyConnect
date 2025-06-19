package dev.psyconnect.identity_service.controller;

import dev.psyconnect.identity_service.configuration.CallRestApi;
import dev.psyconnect.identity_service.dto.response.AuthenticationResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.service.OAuth2Service;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class Oauth2Controller {
    private static final Logger log = LoggerFactory.getLogger(Oauth2Controller.class);
    CallRestApi callRestApi;
    private final OAuth2Service oAuth2Service;
    private final OAuth2AuthorizedClientService authorizedClientService;

    @GetMapping("/oauth2/userInfo")
    public ApiResponse<AuthenticationResponse> oAuth2Google(
            @RequestParam String provider,
            @RequestParam String email,
            @RequestParam String avatar,
            Authentication authentication) {
        return new ApiResponse<>(
                oAuth2Service.processOAuth2PostLogin(email, avatar, authentication, provider)
        );
    }

    @GetMapping("/oauth2/callback/google")
    public ApiResponse<Boolean> oAuth2GoogleCallBack() {
        return new ApiResponse<>(400, "Failed", false);
    }
}
