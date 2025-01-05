package com.example.IdentityService.model;

import com.example.IdentityService.apiresponse.ApiResponse;
import com.example.IdentityService.apiresponse.CustomResponseWrapper;
import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import com.example.IdentityService.service.OAuth2Service;
import lombok.RequiredArgsConstructor;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserService;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.security.Principal;
import java.util.Map;

@Controller
@RestController
@RequiredArgsConstructor
public class Oauth2Controller {
    OAuth2Service oAuth2UserService;

    @GetMapping("/oauth2/register/google")
    public String redirectToGoogle() {
        return "redirect:/oauth2/authorization/google";
    }

    @GetMapping("/oauth2/userInfo/google")
    @CustomResponseWrapper
    public ApiResponse<UserAccountCreationResponse> oAuth2Google(Principal principal) {
        return new ApiResponse<>(oAuth2UserService.createAccountWithGoogle(principal));
    }
}
