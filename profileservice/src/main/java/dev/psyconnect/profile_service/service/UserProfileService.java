package dev.psyconnect.profile_service.service;

import java.util.List;

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

import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.request.UserProfileUpdateRequest;
import dev.psyconnect.profile_service.dto.response.ProfileWithRelationShipResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileUpdateResponse;
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
        kafkaService.send("profile.user-create-setting", request.getProfileId());
        var response = userProfileMapper.toUserProfile(temp);
        response.setDob(request.getDob());
        return response;
    }

    @Cacheable(key = "#profileId", value = "profile")
    public UserProfileResponse get(String profileId) {
        Profile profile = userProfileRepository
                .findById(profileId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        return userProfileMapper.toUserProfileRequest(profile);
    }

    @CacheEvict(value = "profile", key = "#profileId")
    public UserProfileUpdateResponse update(UserProfileUpdateRequest userProfileUpdateRequest, String profileId) {
        Profile existingUser = userProfileRepository
                .findById(profileId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        Profile updatedUser = userProfileMapper.toUserProfile(userProfileUpdateRequest);
        updatedUser.setProfileId(existingUser.getProfileId());
        updatedUser.setAccountId(existingUser.getAccountId());
        Profile savedUser = userProfileRepository.save(updatedUser);
        return userProfileMapper.toUserProfileUpdateResponse(savedUser);
    }

    public List<?> getAll(int page, int size) {
        return userProfileRepository.findAllProfilesPaged(page, size);
    }

    @Cacheable(key = "#profileId", value = "profile_mood")
    public ProfileWithRelationShipResponse getProfileWithMood(String profileId) {
        return userProfileRepository
                .getProfileWithAllRelations(profileId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
    }

    @KafkaListener(topics = "profile.user-create-setting")
    public void handleOnCreateProfile(@Payload String raw) {
        String profileId = KafkaService.objectMapping(raw, String.class);
        userSettingService.resetSettings(profileId);
    }

    public Boolean checkProfileExisted(String profileId) {
        return userProfileRepository.existsById(profileId);
    }
}
