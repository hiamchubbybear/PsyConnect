package dev.psyconnect.profile_service.controller.v1;

import java.io.IOException;
import java.util.List;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.configuration.filter.AllowedRoles;
import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.request.UserProfileUpdateRequest;
import dev.psyconnect.profile_service.dto.response.ProfileWithMoodSummaryDto;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileUpdateResponse;
import dev.psyconnect.profile_service.service.UserProfileService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController("profileControllerV1")
@RequestMapping(path = "/v1/profile")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserProfileControllerV1 {
    UserProfileService userProfileService;

    @PostMapping("/internal/user")
    ApiResponse<UserProfileCreationResponse> createUserProfile(@RequestBody UserProfileCreationRequest body)
            throws IOException {
        return new ApiResponse<>(userProfileService.create(body));
    }

    @PutMapping("/me")
    ApiResponse<UserProfileUpdateResponse> updateUserProfile(
            @RequestBody UserProfileUpdateRequest body, @RequestHeader(name = "X-Profile-Id") String userId) {
        return new ApiResponse<>(userProfileService.update(body, userId));
    }

    @GetMapping("/me")
    ApiResponse<?> getUserProfile(@RequestHeader(value = "X-Profile-Id") String profileId) {
        return new ApiResponse<>(userProfileService.get(profileId));
    }

    @GetMapping("/all")
    @AllowedRoles({"ADMIN"})
    ApiResponse<List<?>> getAllUserProfiles(@RequestParam int page, @RequestParam int size) {
        return new ApiResponse<>(userProfileService.getAll(page, size));
    }

    @GetMapping("/friends")
    ApiResponse<List<ProfileWithMoodSummaryDto>> getFriends(@RequestHeader(value = "X-Profile-Id") String profileId) {
        return new ApiResponse<>(userProfileService.getProfileWithMood(profileId));
    }
}
