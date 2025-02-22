package dev.psyconnect.profile_service.service;

import java.util.List;
import java.util.Optional;

import dev.psyconnect.profile_service.dto.response.ProfileWithRelationShipResponse;
import dev.psyconnect.profile_service.model.Profile;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.cache.annotation.EnableCaching;
import org.springframework.context.event.EventListener;
import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.dto.UserProfileResponse;
import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.request.UserProfileUpdateRequest;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileUpdateResponse;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.mapper.UserProfileMapper;
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
    private final UserSettingService userSettingService;

    @Autowired
    public UserProfileService(ProfileRepository userProfileRepository, UserProfileMapper userProfileMapper, UserSettingService userSettingService) {
        this.userProfileRepository = userProfileRepository;
        this.userProfileMapper = userProfileMapper;
        this.userSettingService = userSettingService;
    }

    // Create user profile
    public UserProfileCreationResponse create(UserProfileCreationRequest request) {
        log.debug("Create user with userId {}", request.getUserId());
        // Mapping user profile from request
        Profile profile = userProfileMapper.toUserProfileMapper(request);
        // Save user profile from request
        profile.setDob(Time.parseFromString(request.getDob()));
        var temp = userProfileRepository.save(profile);
        log.debug("Created profile: {}", profile.getProfileId());
        // Event listen to create profile to create default setting !!
        handleOnCreateProfile(temp.getProfileId());
        // Mapping to response
        return userProfileMapper.toUserProfile(temp);
    }

    @Cacheable(key = "#profileId", value = "profile")
    public UserProfileResponse get(String profileId) {
        // Find user by profileId
        Optional<Profile> profile = userProfileRepository.findById(profileId);
        // Mapping to response user profile
        return userProfileMapper.toUserProfileResponse(
                profile.orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND)));
    }

    // Update profile
    @CacheEvict(value = "profile", key = "#profileId")
    public UserProfileUpdateResponse update(UserProfileUpdateRequest userProfileUpdateRequest, String profileId) {
        // Find exists user from db
        Profile existingUser = userProfileRepository
                .findById(profileId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        log.debug("Updating profile for profileId: {}", existingUser.getProfileId());
        // Mapping to user profile for existed profile
        Profile updatedUser = userProfileMapper.toUserProfile(userProfileUpdateRequest);
        // Keep user profile and user id
        updatedUser.setProfileId(existingUser.getProfileId());
        updatedUser.setAccountId(existingUser.getAccountId());
        // Update to database
        Profile savedUser = userProfileRepository.save(updatedUser);
        //        log.info("Updated profile successfully: {}", savedUser);
        // Mapping to response user profile
        return userProfileMapper.toUserProfileUpdateResponse(savedUser);
    }

    // Get user profile with page and size of page
    public List<?> getAll(int page, int size) {
        return userProfileRepository.findAllProfilesPaged(page, size);
    }

    @Cacheable(key = "#profileId", value = "profile_mood")
    public ProfileWithRelationShipResponse getProfileWithMood(String profileId) {
        return userProfileRepository.getProfileWithAllRelations(profileId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
    }

    @EventListener
    // Handle Create Profile  -> Create def setting
    public void handleOnCreateProfile(String profileId) throws CustomExceptionHandler {
        log.info("Create default setting for account {}" , profileId);
        userSettingService.resetSettings(profileId);
    }

}
