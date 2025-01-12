package com.example.profileservice.controller;

import com.example.profileservice.apiresponse.ApiResponse;
import com.example.profileservice.apiresponse.CustomResponseWrapper;
import com.example.profileservice.dto.UserProfileResponse;
import com.example.profileservice.dto.request.UserProfileUpdateRequest;
import com.example.profileservice.dto.response.UserProfileUpdateResponse;
import com.example.profileservice.model.UserProfile;

import org.springframework.data.domain.Page;
import org.springframework.web.bind.annotation.*;

import com.example.profileservice.dto.request.UserProfileCreationRequest;
import com.example.profileservice.dto.response.UserProfileCreationResponse;
import com.example.profileservice.service.UserProfileService;

import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

import java.util.List;


@RestController
@RequestMapping(path = "/profile")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserProfileController {
    UserProfileService userProfileService;

    //Create user profile
    /*
     Only trigger with openfeign to listen http request
     */
    @PostMapping("/user")
    ApiResponse<UserProfileCreationResponse> createUserProfile(@RequestBody UserProfileCreationRequest body) {
        return new ApiResponse<>(userProfileService.create(body));
    }

    //Update user profile by user id
    @PostMapping()
    ApiResponse<UserProfileUpdateResponse> updateUserProfile(@RequestBody UserProfileUpdateRequest body
            , @RequestParam String userId) {
        return new ApiResponse<>(userProfileService.update(body, userId));
    }

    //Get user profile by profile id
    @GetMapping()
    ApiResponse<UserProfileResponse> getUserProfile(@RequestParam String profileId) {
        return new ApiResponse<>(userProfileService.get(profileId));
    }

    //Get all user profiles with page and size of page
    @GetMapping("/all")
    ApiResponse<List<?>> getAllUserProfiles(int page, int size) {
        return new ApiResponse<>(userProfileService.getAll(page, size));

    }
}
