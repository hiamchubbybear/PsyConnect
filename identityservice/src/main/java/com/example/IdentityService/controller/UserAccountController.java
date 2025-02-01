package com.example.IdentityService.controller;

import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import com.example.IdentityService.apiresponse.ApiResponse;
import com.example.IdentityService.apiresponse.CustomResponseWrapper;
import com.example.IdentityService.dto.request.UserAccountCreationRequest;
import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import com.example.IdentityService.service.UserAccountService;

import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequestMapping("/identity")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserAccountController {
    UserAccountService userAccountService;

    @PostMapping(value = "/create", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    @CustomResponseWrapper
    public ApiResponse<UserAccountCreationResponse> register(
            @RequestPart("json") UserAccountCreationRequest accountRequest, @RequestPart("file") MultipartFile file) {
        return new ApiResponse<>(userAccountService.createAccount(accountRequest, file));
    }

    @GetMapping("/hello")
    @CustomResponseWrapper
    public ApiResponse<String> register() {
        return new ApiResponse<>("Hello World");
    }
}
