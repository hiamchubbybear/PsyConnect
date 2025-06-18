package dev.psyconnect.identity_service.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import dev.psyconnect.identity_service.configuration.CallRestApi;
import dev.psyconnect.identity_service.dto.response.AuthenticationResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.apiresponse.CustomResponseWrapper;
import dev.psyconnect.identity_service.dto.request.GoogleAuthenticationRequest;
import dev.psyconnect.identity_service.service.AuthenticationService;
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

    @GetMapping("/oauth2/userInfo/google")
    public ApiResponse<AuthenticationResponse> oAuth2Google(@RequestParam String email,
                                                            @RequestParam String avatar,
                                                            Authentication authentication) {
        return new ApiResponse<>(
                oAuth2Service.processOAuthPostLoginGoogle(email, avatar, authentication)
        );
    }

    @GetMapping("/oauth2/callback/google")
    public ApiResponse<Boolean> oAuth2GoogleCallBack() {
        return new ApiResponse<>(400, "Failed", false);
    }
}
