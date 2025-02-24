package dev.psyconnect.profile_service.controller;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.dto.request.UserSettingRequest;
import dev.psyconnect.profile_service.dto.response.UserSettingResponse;
import dev.psyconnect.profile_service.model.Setting;
import dev.psyconnect.profile_service.repository.DeleteResponse;
import dev.psyconnect.profile_service.service.UserSettingService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequestMapping("/user-setting")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserSettingController {
    UserSettingService userSettingService;

    @PostMapping("/add")
    public ApiResponse<UserSettingResponse> addUserSetting(
            @RequestHeader(value = "X-Profile-Id" , required = true) String id, @RequestBody UserSettingRequest request) {
        return new ApiResponse<>(userSettingService.createUserSetting(id, request));
    }

    @GetMapping("")
    public ApiResponse<Setting> getUserSetting(@RequestHeader(value = "X-Profile-Id" , required = true) String id) {
        return new ApiResponse<>(userSettingService.getUserSettingById(id));
    }

    @PutMapping("")
    public ApiResponse<UserSettingResponse> updateUserSetting(
            @RequestHeader(value = "X-Profile-Id" , required = true) String id, @RequestBody UserSettingRequest request) {
        return new ApiResponse<>(userSettingService.updateUserSetting(id, request));
    }

    @DeleteMapping("")
    public ApiResponse<DeleteResponse> deleteUserSetting(@RequestHeader(value = "X-Profile-Id" , required = true) String id) {
        return new ApiResponse<>(userSettingService.deleteUserSetting(id));
    }

    @PostMapping("/default")
    public ApiResponse<Setting> setDefaultSetting(@RequestHeader(value = "X-Profile-Id" , required = true) String id) {
        return new ApiResponse<>(userSettingService.resetSettings(id));
    }
}
