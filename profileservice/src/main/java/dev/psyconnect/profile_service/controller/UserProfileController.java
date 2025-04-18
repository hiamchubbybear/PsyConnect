package dev.psyconnect.profile_service.controller;

import java.io.IOException;
import java.util.List;

import dev.psyconnect.profile_service.configuration.filter.AllowedRoles;
import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.request.UserProfileUpdateRequest;
import dev.psyconnect.profile_service.dto.response.ProfileWithRelationShipResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileUpdateResponse;
import dev.psyconnect.profile_service.service.UserProfileService;
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
    @PostMapping("/internal/user")
    ApiResponse<UserProfileCreationResponse> createUserProfile(@RequestBody UserProfileCreationRequest body)
            throws IOException {
        return new ApiResponse<>(userProfileService.create(body));
    }

    // Update user profile by user id
    @PostMapping()
    ApiResponse<UserProfileUpdateResponse> updateUserProfile(
            @RequestBody UserProfileUpdateRequest body, @RequestHeader(name = "X-User-Id") String userId) {
        return new ApiResponse<>(userProfileService.update(body, userId));
    }

    @GetMapping()
    ApiResponse<?> getUserProfile(@RequestHeader(value = "X-Profile-Id") String profileId) {
        if (profileId == null || profileId.isEmpty()) {
            return new ApiResponse<>("Empty ProfileId");
        }
        return new ApiResponse<>(userProfileService.get(profileId));
    }

    // Get all user profiles with page and size of page
    @GetMapping("/all")
    @AllowedRoles({"ADMIN"})
    ApiResponse<List<?>> getAllUserProfiles(@RequestParam int page, @RequestParam int size) {
        return new ApiResponse<>(userProfileService.getAll(page, size));
    }

    @GetMapping("/friends")
    ApiResponse<ProfileWithRelationShipResponse> getFriends(@RequestHeader(value = "X-Profile-Id") String profileId) {
        return new ApiResponse<>(userProfileService.getProfileWithMood(profileId));
    }
}
