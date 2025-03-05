package dev.psyconnect.profile_service.controller;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.dto.request.MoodCreateRequest;
import dev.psyconnect.profile_service.dto.request.MoodUpdateRequest;
import dev.psyconnect.profile_service.dto.response.DeleteMoodResponse;
import dev.psyconnect.profile_service.dto.response.GetMoodResponse;
import dev.psyconnect.profile_service.dto.response.MoodCreateResponse;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.service.MoodService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequestMapping("/mood")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class MoodController {
    MoodService moodService;

    @PostMapping("/add")
    public ApiResponse<MoodCreateResponse> addMood(
            @RequestHeader(value = "X-Profile-Id", required = true) String profileId,
            @RequestBody MoodCreateRequest request) {
        if (profileId == null) {
            throw new CustomExceptionHandler(ErrorCode.MISSING_TOKEN);
        }
        return new ApiResponse<>(moodService.createMoodByProfileId(profileId, request));
    }

    @GetMapping
    public ApiResponse<GetMoodResponse> getMood(
            @RequestHeader(value = "X-Profile-Id", required = true) String profileId) {

        return new ApiResponse<>(moodService.getMoodById(profileId));
    }

    @PutMapping
    public ApiResponse<MoodCreateResponse> updateMood(
            @RequestHeader(value = "X-Profile-Id", required = true) String profileId,
            @RequestBody MoodUpdateRequest request) {
        return new ApiResponse<>(moodService.updateMoodByProfileId(profileId, request));
    }

    @DeleteMapping
    public ApiResponse<DeleteMoodResponse> updateMood(
            @RequestHeader(value = "X-Profile-Id", required = true) String profileId) {
        return new ApiResponse<>(moodService.deleteMood(profileId));
    }
}
