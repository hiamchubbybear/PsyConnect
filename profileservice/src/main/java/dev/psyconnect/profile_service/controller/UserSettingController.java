package dev.psyconnect.profile_service.controller;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.dto.request.UserSettingRequest;
import dev.psyconnect.profile_service.dto.response.UserSettingResponse;
import dev.psyconnect.profile_service.model.Setting;
import dev.psyconnect.profile_service.repository.DeleteResponse;
import dev.psyconnect.profile_service.service.UserSettingService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/user-setting")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserSettingController {
    UserSettingService userSettingService;

    @PostMapping("/add/{id}")
    public ApiResponse<UserSettingResponse> addUserSetting(@PathVariable String id, @RequestBody UserSettingRequest request) {
        return new ApiResponse<>(userSettingService.createUserSetting(id, request));
    }

    @GetMapping("/{id}")
    public ApiResponse<Setting> getUserSetting(@PathVariable String id) {
        return new ApiResponse<>(userSettingService.getUserSettingById(id));
    }

    @PutMapping("/{id}")
    public ApiResponse<UserSettingResponse> updateUserSetting(@PathVariable String id, @RequestBody UserSettingRequest request) {
        return new ApiResponse<>(userSettingService.updateUserSetting(id, request));
    }

    @DeleteMapping("/{id}")
    public ApiResponse<DeleteResponse> deleteUserSetting(@PathVariable String id) {
        return new ApiResponse<>(userSettingService.deleteUserSetting(id));
    }
    @PostMapping("/{id}/default")
    public ApiResponse<Setting> setDefaultSetting(@PathVariable String id) {
        return new ApiResponse<>(userSettingService.resetSettings(id));
    }
}

