package com.example.profileservice.controller;

import com.example.profileservice.apiresponse.ApiResponse;
import com.example.profileservice.apiresponse.CustomResponseWrapper;
import com.example.profileservice.dto.request.UserProfileUpdateRequest;
import com.example.profileservice.dto.response.UserProfileUpdateResponse;
import com.google.protobuf.Api;
import org.springframework.context.annotation.Lazy;
import org.springframework.web.bind.annotation.*;

import com.example.profileservice.dto.request.UserProfileCreationRequest;
import com.example.profileservice.dto.response.UserProfileCreationResponse;
import com.example.profileservice.service.UserProfileService;

import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequestMapping(path = "/profile")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserProfileController {
    UserProfileService userProfileService;

    @PostMapping("/user")
    ApiResponse<UserProfileCreationResponse> createUserProfile(@RequestBody UserProfileCreationRequest body) {
        return new ApiResponse<>(userProfileService.createProfile(body));
    }
    @PostMapping("/update")
    ApiResponse<UserProfileUpdateResponse> updateUserProfile(@RequestBody UserProfileUpdateRequest body
            , @RequestParam String email) {
        return new ApiResponse<>(userProfileService.profileUpdateResponse(body, email));
    }
}
