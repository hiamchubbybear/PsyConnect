package dev.psyconnect.profile_service.controller.v1;

import java.util.List;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.dto.request.MoodCreateRequest;
import dev.psyconnect.profile_service.dto.request.MoodUpdateRequest;
import dev.psyconnect.profile_service.dto.response.DeleteMoodResponse;
import dev.psyconnect.profile_service.dto.response.FriendMoodDTO;
import dev.psyconnect.profile_service.dto.response.GetMoodResponse;
import dev.psyconnect.profile_service.dto.response.MoodCreateResponse;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.service.MoodService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController("moodControllerV1")
@RequestMapping("/v1/mood")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class MoodControllerV1 {
    MoodService moodService;

    @PostMapping("/me/add")
    public ApiResponse<MoodCreateResponse> addMood(
            @RequestHeader(value = "X-Profile-Id") String profileId, @RequestBody MoodCreateRequest request) {
        if (profileId == null) throw new CustomExceptionHandler(ErrorCode.MISSING_TOKEN);
        return new ApiResponse<>(moodService.createMoodByProfileId(profileId, request));
    }

    @GetMapping("/me")
    public ApiResponse<GetMoodResponse> getMood(@RequestHeader(value = "X-Profile-Id") String profileId) {
        return new ApiResponse<>(moodService.getMoodById(profileId));
    }

    @PutMapping("/me")
    public ApiResponse<MoodCreateResponse> updateMoodById(
            @RequestHeader(value = "X-Profile-Id") String profileId, @RequestBody MoodUpdateRequest request) {
        return new ApiResponse<>(moodService.updateMoodByProfileId(profileId, request));
    }

    @DeleteMapping("/me")
    public ApiResponse<DeleteMoodResponse> deleteMood(@RequestHeader(value = "X-Profile-Id") String profileId) {
        return new ApiResponse<>(moodService.deleteMood(profileId));
    }

    @GetMapping("/me/friends")
    public ApiResponse<List<FriendMoodDTO>> getMoodById(@RequestHeader(value = "X-Profile-Id") String profileId) {
        return new ApiResponse<>(moodService.getFriendsMood(profileId));
    }
}
