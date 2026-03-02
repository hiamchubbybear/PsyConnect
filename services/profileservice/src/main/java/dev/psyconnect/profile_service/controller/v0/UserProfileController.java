package dev.psyconnect.profile_service.controller.v0;

import java.io.IOException;
import java.util.List;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.configuration.filter.AllowedRoles;
import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.request.UserProfileUpdateRequest;
import dev.psyconnect.profile_service.dto.response.ProfileWithMoodSummaryDto;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileResponse;
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

    @Deprecated
    @PostMapping("/internal/user")
    ApiResponse<UserProfileCreationResponse> createUserProfile(@RequestBody UserProfileCreationRequest body)
            throws IOException {
        return new ApiResponse<>(userProfileService.create(body));
    }

    @Deprecated
    @PutMapping()
    ApiResponse<UserProfileUpdateResponse> updateUserProfile(
            @RequestBody UserProfileUpdateRequest body, @RequestHeader(name = "X-Profile-Id") String userId) {
        return new ApiResponse<>(userProfileService.update(body, userId));
    }

    @Deprecated
    @GetMapping()
    ApiResponse<?> getUserProfile(@RequestHeader(value = "X-Profile-Id") String profileId) {
        if (profileId == null || profileId.isEmpty()) {
            return new ApiResponse<>("Empty ProfileId");
        }
        return new ApiResponse<>(userProfileService.get(profileId));
    }

    @Deprecated
    @GetMapping("/all")
    @AllowedRoles({"ADMIN"})
    ApiResponse<List<?>> getAllUserProfiles(@RequestParam int page, @RequestParam int size) {
        return new ApiResponse<>(userProfileService.getAll(page, size));
    }

    @Deprecated
    @GetMapping("/friends")
    ApiResponse<List<ProfileWithMoodSummaryDto>> getFriends(@RequestHeader(value = "X-Profile-Id") String profileId) {
        return new ApiResponse<>(userProfileService.getProfileWithMood(profileId));
    }

    @GetMapping("/search")
    ApiResponse<List<UserProfileResponse>> searchProfiles(
            @RequestParam String query,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "10") int size) {
        return new ApiResponse<>(userProfileService.search(query, page, size));
    }
}
