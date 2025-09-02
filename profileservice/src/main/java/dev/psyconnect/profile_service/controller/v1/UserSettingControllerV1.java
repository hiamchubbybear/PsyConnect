package dev.psyconnect.profile_service.controller.v1;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.dto.request.UserSettingRequest;
import dev.psyconnect.profile_service.dto.response.UserSettingResponse;
import dev.psyconnect.profile_service.kafka.service.KafkaService;
import dev.psyconnect.profile_service.model.Setting;
import dev.psyconnect.profile_service.service.UserSettingService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController("settingControllerV1")
@RequestMapping("/v1/user-setting")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserSettingControllerV1 {
    private final UserSettingService userSettingService;
    private final KafkaService kafkaService;

    @GetMapping("/me")
    public ApiResponse<Setting> getUserSetting(@RequestHeader(value = "X-Profile-Id") String id) {
        return new ApiResponse<>(userSettingService.getUserSettingById(id));
    }

    @PutMapping("/me")
    public ApiResponse<UserSettingResponse> updateUserSetting(
            @RequestHeader(value = "X-Profile-Id") String id, @RequestBody UserSettingRequest request) {
        request.setProfileId(id);
        return new ApiResponse<>(userSettingService.updateUserSetting(id, request));
    }

    @PostMapping("/me/default")
    public ApiResponse<Setting> setDefaultSetting(@RequestHeader(value = "X-Profile-Id") String id) {
        return new ApiResponse<>(userSettingService.resetSettings(id));
    }
}
