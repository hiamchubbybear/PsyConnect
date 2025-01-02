package com.example.IdentityService.controller;

import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import com.example.IdentityService.apiresponse.ApiResponse;
import com.example.IdentityService.apiresponse.CustomResponseWrapper;
import com.example.IdentityService.dto.request.UserAccountCreationRequest;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.service.UserAccountService;

import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequestMapping("/account")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserAccountController {
    UserAccountService userAccountService;

    @PostMapping("/create")
    @CustomResponseWrapper
    public ApiResponse<UserAccountCreationResponse> register(@RequestBody UserAccountCreationRequest accountRequest) {
        return new ApiResponse<>(userAccountService.createAccount(accountRequest));
    }
}
