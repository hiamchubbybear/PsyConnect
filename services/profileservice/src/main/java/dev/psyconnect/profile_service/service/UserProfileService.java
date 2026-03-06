package dev.psyconnect.profile_service.service;

import java.time.Instant;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.context.ApplicationEventPublisher;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.messaging.handler.annotation.Payload;
import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.dto.request.*;
import dev.psyconnect.profile_service.dto.response.*;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.grpc.client.IdentityGRPCClient;
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
    private final IdentityGRPCClient identityGRPCClient;
    ApplicationEventPublisher eventPublisher;

    @Autowired
    public UserProfileService(
            ProfileRepository userProfileRepository,
            UserProfileMapper userProfileMapper,
            KafkaService kafkaService,
            UserSettingService userSettingService,
            ApplicationEventPublisher eventPublisher,
            IdentityGRPCClient identityGRPCClient) {
        this.userProfileRepository = userProfileRepository;
        this.userProfileMapper = userProfileMapper;
        this.kafkaService = kafkaService;
        this.userSettingService = userSettingService;
        this.eventPublisher = eventPublisher;
        this.identityGRPCClient = identityGRPCClient;
    }

    public UserProfileCreationResponse create(UserProfileCreationRequest request) {
        Profile profile = userProfileMapper.toUserProfileMapper(request);
        String dobStr = request.getDob();
        profile.setDob((dobStr != null && !dobStr.isEmpty()) ? Time.parseFromString(dobStr) : null);
        var temp = userProfileRepository.save(profile);
        eventPublisher.publishEvent(new OnProfileCreatedEvent(this, temp.getProfileId()));

        userSettingService.resetSettings(temp.getProfileId());

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
            try {
                var userInfo = identityGRPCClient.getUserInfo(profileId);
                if (userInfo != null && userInfo.getSuccess()) {
                    response.setRole(userInfo.getRole());
                }
            } catch (Exception grpcException) {
                log.warn("Failed to fetch role from identity service for profile: {}", profileId, grpcException);
            }
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "profile-service",
                    profileId,
                    "Get profile",
                    "Failed",
                    Map.of("error", e.getMessage()),
                    LogLevel.ERROR));
            throw e;
        }
        return response;
    }

    public List<UserProfileResponse> getBulk(List<String> profileIds) {
        return userProfileRepository.findAllByIds(profileIds).stream()
                .map(userProfileMapper::toUserProfileRequest)
                .collect(Collectors.toList());
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
            updatedUser.setDob(userProfileUpdateRequest.getDob());
            log.info("Updated user {}", updatedUser.getDob());
            Profile savedUser = userProfileRepository.save(updatedUser);
            response = userProfileMapper.toUserProfileUpdateResponse(savedUser);
            kafkaService.sendLog(buildLog(
                    "profile-service", profileId, "Update profile", "Success", Map.of("data", response), LogLevel.LOG));
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "profile-service",
                    profileId,
                    "Update profile",
                    "Failed",
                    Map.of("error", e.getMessage()),
                    LogLevel.ERROR));
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
                    "profile-service",
                    "system",
                    "Get all profiles",
                    "Failed",
                    Map.of("error", e.getMessage()),
                    LogLevel.LOG));
            throw e;
        }
        return result;
    }

    public List<UserProfileResponse> search(String query, int page, int size) {
        try {
            return userProfileRepository.searchProfiles(query, page * size, size).stream()
                    .map(userProfileMapper::toUserProfileRequest)
                    .collect(Collectors.toList());
        } catch (Exception e) {
            throw e;
        }
    }

    public List<UserProfileResponse> getBatch(List<String> profileIds) {
        try {
            return userProfileRepository.findAllByIds(profileIds).stream()
                    .map(userProfileMapper::toUserProfileRequest)
                    .collect(Collectors.toList());
        } catch (Exception e) {
            throw e;
        }
    }

    public List<ProfileWithMoodSummaryDto> getProfileWithMood(String profileId) {
        try {
            return userProfileRepository.getProfileWithAllRelations(profileId).stream()
                    .map(dev.psyconnect.profile_service.dto.response.ProfileWithRelationShipResponse::toSummaryWithMood)
                    .collect(Collectors.toList());
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "profile-service",
                    profileId,
                    "Get profile with mood",
                    "Failed",
                    Map.of("error", e.getMessage()),
                    LogLevel.ERROR));
            throw e;
        }
    }

    @KafkaListener(topics = "profile.user-create-setting")
    public void handleOnCreateProfile(@Payload String raw) {
        String profileId = raw;
        userSettingService.resetSettings(profileId);
    }

    public Boolean checkProfileExisted(String profileId) {
        return userProfileRepository.existsById(profileId);
    }

    private LogEvent buildLog(
            String service,
            String userId,
            String action,
            String message,
            Map<String, Object> metadata,
            LogLevel level) {
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
