package dev.psyconnect.profile_service.service;

import dev.psyconnect.profile_service.dto.request.UserSettingRequest;
import dev.psyconnect.profile_service.dto.response.UserSettingResponse;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.model.Setting;
import dev.psyconnect.profile_service.repository.DeleteResponse;
import dev.psyconnect.profile_service.repository.ProfileRepository;
import dev.psyconnect.profile_service.repository.SettingRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Service;

import java.util.Optional;

@Service
@RequiredArgsConstructor
@Slf4j
public class UserSettingService {
    private final SettingRepository userSettingRepository;
    private final ProfileRepository profileRepository;

    @CacheEvict(key = "#profileId", value = "setting")
    public UserSettingResponse createUserSetting(String profileId, UserSettingRequest request) {
        if (!profileRepository.existsById(profileId) || userSettingRepository.getUserSetting(profileId).isPresent()) {
            throw new CustomExceptionHandler(ErrorCode.RUNTIME_ERROR);
        }
        var response = userSettingRepository.createUserSetting(
                profileId, request.getPrivacyLevel(), request.isShowLastSeen(),
                request.isShowProfilePicture(), request.isShowMood(), request.isNotificationsEnabled(),
                request.isEmailNotifications(), request.isPushNotifications(), request.isSmsNotifications(),
                request.isTwoFactorAuth(), request.isAllowLoginAlerts(), request.getTrustedDevices(),
                request.getLanguage(), request.getTheme(), request.isAutoDeleteOldMoods()
        ).orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));

        return UserSettingResponse.builder().profileId(profileId).isSuccess(true).build();
    }

    @CacheEvict(key = "#profileId", value = "setting")
    public UserSettingResponse updateUserSetting(String profileId, UserSettingRequest request) {
        if (!profileRepository.existsById(profileId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        var response = userSettingRepository.updateUserSetting(
                profileId, request.getPrivacyLevel(), request.isShowLastSeen(),
                request.isShowProfilePicture(), request.isShowMood(), request.isNotificationsEnabled(),
                request.isEmailNotifications(), request.isPushNotifications(), request.isSmsNotifications(),
                request.isTwoFactorAuth(), request.isAllowLoginAlerts(), request.getTrustedDevices(),
                request.getLanguage(), request.getTheme(), request.isAutoDeleteOldMoods()
        ).orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));

        return UserSettingResponse.builder().profileId(profileId).isSuccess(true).build();
    }

    @Cacheable(key = "#profileId", value = "setting")
    public Setting getUserSettingById(String profileId) {
        return userSettingRepository.getUserSetting(profileId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
    }

    @CacheEvict(key = "#profileId", value = "setting")
    public DeleteResponse deleteUserSetting(String profileId) {
        if (!profileRepository.existsById(profileId)) throw new CustomExceptionHandler(ErrorCode.QUERY_FAILED);
        userSettingRepository.deleteUserSetting(profileId);
        return new DeleteResponse(true, "Deleted successfully");
    }

    @CacheEvict(key = "#profileId", value = "setting")
    public Setting resetSettings(String profileId) {
        return userSettingRepository.createUserSetting(
                profileId, "PRIVATE",
                true, true,
                false, false,
                false, false,
                false, false,
                false, "",
                "en", "light",
                true
        ).orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
    }
}
