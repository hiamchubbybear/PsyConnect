package dev.psyconnect.profile_service.service;

import java.time.Instant;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import dev.psyconnect.profile_service.dto.request.*;
import dev.psyconnect.profile_service.dto.response.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.kafka.service.KafkaService;
import dev.psyconnect.profile_service.mapper.UserProfileMapper;
import dev.psyconnect.profile_service.model.Profile;
import dev.psyconnect.profile_service.repository.ProfileRepository;
import lombok.AccessLevel;
import lombok.experimental.FieldDefaults;

@Service
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
@EnableCaching
public class UserProfileService {
    private static final Logger log = LoggerFactory.getLogger(UserProfileService.class);
    ProfileRepository userProfileRepository;
    UserProfileMapper userProfileMapper;
    KafkaService kafkaService;
    private final UserSettingService userSettingService;
    ApplicationEventPublisher eventPublisher;

    @Autowired
    public UserProfileService(
            ProfileRepository userProfileRepository,
            UserProfileMapper userProfileMapper,
            KafkaService kafkaService,
            UserSettingService userSettingService,
            ApplicationEventPublisher eventPublisher) {
        this.userProfileRepository = userProfileRepository;
        this.userProfileMapper = userProfileMapper;
        this.kafkaService = kafkaService;
        this.userSettingService = userSettingService;
        this.eventPublisher = eventPublisher;
    }

    public UserProfileCreationResponse create(UserProfileCreationRequest request) {
        Profile profile = userProfileMapper.toUserProfileMapper(request);
        profile.setDob(Time.parseFromString(request.getDob()));
        var temp = userProfileRepository.save(profile);
        eventPublisher.publishEvent(new OnProfileCreatedEvent(this, temp.getProfileId()));
        kafkaService.send("profile.user-create-setting", request.getProfileId());
        kafkaService.sendLog(
                buildLog("profile-service", request.getProfileId(), "Create profile", "Success", Map.of("metadata", temp.toString()), LogLevel.LOG));
        var response = userProfileMapper.toUserProfile(temp);
        response.setDob(request.getDob());
        return response;
    }

    public UserProfileResponse get(String profileId) {
        UserProfileResponse response;
        try {
            Profile profile = userProfileRepository
                    .findById(profileId)
                    .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
            response = userProfileMapper.toUserProfileRequest(profile);
        } catch (Exception e) {
            kafkaService.sendLog(
                    buildLog("profile-service", profileId, "Get profile", "Failed", Map.of("error", e.getMessage()), LogLevel.ERROR));
            throw e;
        }
        return response;
    }

    public UserProfileUpdateResponse update(UserProfileUpdateRequest userProfileUpdateRequest, String profileId) {
        UserProfileUpdateResponse response;
        try {
            Profile existingUser = userProfileRepository
                    .findById(profileId)
                    .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
            Profile updatedUser = userProfileMapper.toUserProfile(userProfileUpdateRequest);
            updatedUser.setProfileId(existingUser.getProfileId());
            updatedUser.setAccountId(existingUser.getAccountId());
            Profile savedUser = userProfileRepository.save(updatedUser);
            response = userProfileMapper.toUserProfileUpdateResponse(savedUser);
            kafkaService.sendLog(buildLog("profile-service", profileId, "Update profile", "Success", Map.of("data", response), LogLevel.LOG));
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "profile-service", profileId, "Update profile", "Failed", Map.of("error", e.getMessage()), LogLevel.ERROR));
            throw e;
        }
        return response;
    }

    public List<?> getAll(int page, int size) {
        List<?> result;
        try {
            result = userProfileRepository.findAllProfilesPaged(page, size);
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "profile-service", "system", "Get all profiles", "Failed", Map.of("error", e.getMessage()), LogLevel.LOG));
            throw e;
        }
        return result;
    }

    public List<ProfileWithMoodSummaryDto> getProfileWithMood(String profileId) {
        try {
            return userProfileRepository.getProfileWithAllRelations(profileId).stream()
                    .map(p -> ProfileWithMoodSummaryDto.builder()
                            .profile(ProfileSummaryDto.builder()
                                    .profileId(p.getProfile().getProfileId())
                                    .username(p.getProfile().getUsername())
                                    .avatarUri(p.getProfile().getAvatarUri())
                                    .build())
                            .mood(p.getMood())
                            .build())
                    .collect(Collectors.toList());
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "profile-service", profileId, "Get profile with mood", "Failed", Map.of("error", e.getMessage()), LogLevel.ERROR));
            throw e;
        }
    }


    @KafkaListener(topics = "profile.user-create-setting")
    public void handleOnCreateProfile(@Payload String raw) {
        String profileId = KafkaService.objectMapping(raw, String.class);
        userSettingService.resetSettings(profileId);
    }

    public Boolean checkProfileExisted(String profileId) {
        return userProfileRepository.existsById(profileId);
    }

    private LogEvent buildLog(
            String service, String userId, String action, String message, Map<String, Object> metadata, LogLevel level) {
        return LogEvent.builder()
                .service(service)
                .level(level)
                .timestamp(Instant.now().toString())
                .userId(userId)
                .action(action)
                .message(message)
                .metadata(metadata)
                .build();
    }
}
