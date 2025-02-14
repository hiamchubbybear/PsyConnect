package dev.psyconnect.identity_service.controller;

import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.apiresponse.CustomResponseWrapper;
import dev.psyconnect.identity_service.dto.request.GoogleAuthenticationRequest;
import dev.psyconnect.identity_service.service.AuthenticationService;
import dev.psyconnect.identity_service.service.OAuth2Service;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@Controller
@RestController(value = "/oauth2")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class Oauth2Controller {
    OAuth2Service oAuth2Service;
    private final AuthenticationService authenticationService;

    @GetMapping("/userInfo/google")
    @CustomResponseWrapper
    public ApiResponse<String> oAuth2Google(Authentication authentication) {
        DefaultOidcUser user = (DefaultOidcUser) authentication.getPrincipal();
        String email = user.getAttribute("email");
        return new ApiResponse<>(
                authenticationService.generateGoogleAuthToken(new GoogleAuthenticationRequest(email), "GOOGLE"));
    }

    @GetMapping("/callback/google")
    @CustomResponseWrapper
    public ApiResponse<Boolean> oAuth2GoogleCallBack() {
        return new ApiResponse<>(400, "Failed", false);
    }
}
