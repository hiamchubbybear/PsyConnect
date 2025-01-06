package com.example.IdentityService.controller;

import com.example.IdentityService.apiresponse.ApiResponse;
import com.example.IdentityService.apiresponse.CustomResponseWrapper;
import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import com.example.IdentityService.service.OAuth2Service;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import java.security.Principal;

@Controller
@RestController
@RequiredArgsConstructor
public class Oauth2Controller {
    OAuth2Service oAuth2Service;

    @GetMapping("/oauth2/userInfo/google")
    @CustomResponseWrapper
    public ApiResponse<Boolean> oAuth2Google() {
        return new ApiResponse<>(true);
    }

    @GetMapping("/oauth2/callback/google")
    @CustomResponseWrapper
    public ApiResponse<Boolean> oAuth2GoogleCallBack() {
        return new ApiResponse<>(400, "Failed", false);
    }
}
