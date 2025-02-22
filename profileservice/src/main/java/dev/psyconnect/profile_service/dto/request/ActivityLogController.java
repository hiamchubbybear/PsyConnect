package dev.psyconnect.profile_service.dto.request;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/activity-log")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class ActivityLogController {
    ActivityLogService activityLogService;

    @PostMapping("/{profileId}")
    public ApiResponse<ActivityLogResponse> addLog(
            @PathVariable String profileId,
            @RequestBody ActivityLogRequest request
    ) {
        return new ApiResponse<>(activityLogService.createLog(profileId, request));
    }

    @GetMapping("/{profileId}")
    public ApiResponse<List<ActivityLogResponse>> getLogs(@PathVariable String profileId) {
        return new ApiResponse<>(activityLogService.getLogsByProfileId(profileId));
    }

    @DeleteMapping("/{logId}")
    public ApiResponse<String> deleteLog(@PathVariable String logId) {
        activityLogService.deleteLog(logId);
        return new ApiResponse<>("Delete success");
    }
}
