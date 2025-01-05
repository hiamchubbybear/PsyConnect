package com.example.IdentityService.controller;

import com.example.IdentityService.apiresponse.ApiResponse;
import com.example.IdentityService.apiresponse.CustomResponseWrapper;
import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import com.example.IdentityService.service.OAuth2Service;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.security.Principal;

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
