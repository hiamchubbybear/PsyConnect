package dev.psyconnect.profile_service.service;

import java.util.List;
import java.util.Optional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.dto.UserProfileResponse;
import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.request.UserProfileUpdateRequest;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileUpdateResponse;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.mapper.UserProfileMapper;
import dev.psyconnect.profile_service.model.UserProfile;
import dev.psyconnect.profile_service.repository.ProfileRepository;
import lombok.AccessLevel;
import lombok.experimental.FieldDefaults;

@Service
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class UserProfileService {
    private static final Logger log = LoggerFactory.getLogger(UserProfileService.class);
    ProfileRepository userProfileRepository;
    UserProfileMapper userProfileMapper;

    @Autowired
    public UserProfileService(ProfileRepository userProfileRepository, UserProfileMapper userProfileMapper) {
        this.userProfileRepository = userProfileRepository;
        this.userProfileMapper = userProfileMapper;
    }

    // Create user profile
    public UserProfileCreationResponse create(UserProfileCreationRequest request) {
        log.debug("Create user with userId {}", request.getUserId());
        // Mapping user profile from request
        UserProfile profile = userProfileMapper.toUserProfileMapper(request);
        // Save user profile from request
        profile.setDob(DateTimeFormatting.parseFromString(request.getDob()));
        var temp = userProfileRepository.save(profile);
        log.debug("Created profile: {}", profile.getUserId());
        // Mapping to response
        return userProfileMapper.toUserProfile(temp);
    }

    public UserProfileResponse get(String profileId) {
        // Find user by profileId
        Optional<UserProfile> profile = userProfileRepository.findById(profileId);
        // Mapping to response user profile
        return userProfileMapper.toUserProfileResponse(
                profile.orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND)));
    }

    // Update profile
    public UserProfileUpdateResponse update(UserProfileUpdateRequest userProfileUpdateRequest, String userId) {
        // Find exists user from db
        UserProfile existingUser = userProfileRepository
                .findByUserId(userId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        log.debug("Updating profile for userId: {}", existingUser.getUserId());
        // Mapping to user profile for existed profile
        UserProfile updatedUser = userProfileMapper.toUserProfile(userProfileUpdateRequest);
        // Keep user profile and user id
        updatedUser.setProfileId(existingUser.getProfileId());
        updatedUser.setUserId(existingUser.getUserId());
        // Update to database
        UserProfile savedUser = userProfileRepository.save(updatedUser);
        //        log.info("Updated profile successfully: {}", savedUser);
        // Mapping to response user profile
        return userProfileMapper.toUserProfileUpdateResponse(savedUser);
    }

    // Get user profile with page and size of page
    public List<?> getAll(int page, int size) {
        return userProfileRepository.findAllProfilesPaged(page, size);
    }
}
