package com.example.profileservice.controller;

import java.io.IOException;
import java.util.List;

import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

import com.example.profileservice.apiresponse.ApiResponse;
import com.example.profileservice.dto.UserProfileResponse;
import com.example.profileservice.dto.request.UserProfileCreationRequest;
import com.example.profileservice.dto.request.UserProfileUpdateRequest;
import com.example.profileservice.dto.response.UserProfileCreationResponse;
import com.example.profileservice.dto.response.UserProfileUpdateResponse;
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

    // Create user profile
    /*
    Only trigger with openfeign to listen http request
     */
    @PostMapping("/user")
    ApiResponse<UserProfileCreationResponse> createUserProfile(
            @RequestBody UserProfileCreationRequest body, @RequestPart MultipartFile file) throws IOException {
        return new ApiResponse<>(userProfileService.create(body, file));
    }

    // Update user profile by user id
    @PostMapping()
    ApiResponse<UserProfileUpdateResponse> updateUserProfile(
            @RequestBody UserProfileUpdateRequest body, @RequestParam String userId) {
        return new ApiResponse<>(userProfileService.update(body, userId));
    }

    // Get user profile by profile id
    @GetMapping()
    ApiResponse<UserProfileResponse> getUserProfile(@RequestParam String profileId) {
        return new ApiResponse<>(userProfileService.get(profileId));
    }

    // Get all user profiles with page and size of page
    @GetMapping("/all")
    ApiResponse<List<?>> getAllUserProfiles(int page, int size) {
        return new ApiResponse<>(userProfileService.getAll(page, size));
    }
}
