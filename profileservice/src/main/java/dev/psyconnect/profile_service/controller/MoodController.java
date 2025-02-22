package dev.psyconnect.profile_service.controller;

import dev.psyconnect.profile_service.dto.request.MoodCreateRequest;
import dev.psyconnect.profile_service.dto.request.MoodUpdateRequest;
import dev.psyconnect.profile_service.dto.response.DeleteMoodResponse;
import dev.psyconnect.profile_service.dto.response.GetMoodResponse;
import dev.psyconnect.profile_service.service.MoodService;
import org.springframework.web.bind.annotation.*;
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

    @GetMapping("/{id}")
    public ApiResponse<GetMoodResponse> getMood(@PathVariable String id) {
        return new ApiResponse<>(moodService.getMoodById(id));
    }

    @PutMapping("/{id}")
    public ApiResponse<MoodCreateResponse> updateMood(@PathVariable String id, @RequestBody MoodUpdateRequest request) {
        return new ApiResponse<>(moodService.updateMoodByProfileId(id, request));
    }

    @DeleteMapping("/{id}")
    public ApiResponse<DeleteMoodResponse> updateMood(@PathVariable String id) {
        return new ApiResponse<>(moodService.deleteMood(id));
    }
}