package com.example.IdentityService.controller;

import com.example.IdentityService.apiresponse.ApiResponse;
import com.example.IdentityService.apiresponse.CustomResponseWrapper;
import com.example.IdentityService.dto.request.GoogleAuthenticationRequest;
import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import com.example.IdentityService.service.AuthenticationService;
import com.example.IdentityService.service.OAuth2Service;
import jakarta.servlet.http.HttpServletRequest;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.parameters.P;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import java.security.Principal;

@Controller
@RestController
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class Oauth2Controller {
    OAuth2Service oAuth2Service;
    private final AuthenticationService authenticationService;

    //    @GetMapping("/oauth2/userInfo/google")
//    @CustomResponseWrapper
//    public ApiResponse<String> oAuth2Google(Authentication authentication) {
//        DefaultOidcUser user = (DefaultOidcUser) authentication.getPrincipal();
//        String email = user.getAttribute("email");
//        return new ApiResponse<>(authenticationService.generateGoogleAuthToken(new GoogleAuthenticationRequest(email),"GOOGLE"));
//    }
    @GetMapping("/oauth2/userInfo/google")
    @CustomResponseWrapper
    public ApiResponse<Principal> oAuth2Google() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        Principal principal = (Principal) authentication.getPrincipal();
        return new ApiResponse<>(principal);
    }

    @GetMapping("/oauth2/callback/google")
    @CustomResponseWrapper
    public ApiResponse<Boolean> oAuth2GoogleCallBack() {
        return new ApiResponse<>(400, "Failed", false);
    }
}
