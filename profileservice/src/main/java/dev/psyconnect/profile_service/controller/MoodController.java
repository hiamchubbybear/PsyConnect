package dev.psyconnect.profile_service.controller;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.MoodCreateRequest;
import dev.psyconnect.profile_service.MoodService;
import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.dto.response.MoodCreateResponse;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequestMapping("/mood")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class MoodController {
    MoodService moodService;

    @PostMapping("/add/{id}")
    public ApiResponse<MoodCreateResponse> addMood(@PathVariable String id, @RequestBody MoodCreateRequest request) {
        return new ApiResponse<>(moodService.createMoodByProfileId(id, request));
    }
}
