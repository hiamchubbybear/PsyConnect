package com.example.profileservice.service;

import java.util.List;
import java.util.Optional;

import com.example.profileservice.dto.UserProfileResponse;
import com.example.profileservice.dto.request.UserProfileUpdateRequest;
import com.example.profileservice.dto.response.UserProfileUpdateResponse;
import com.example.profileservice.globalexceptionhandle.CustomExceptionHandler;
import com.example.profileservice.globalexceptionhandle.ErrorCode;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import com.example.profileservice.dto.request.UserProfileCreationRequest;
import com.example.profileservice.dto.response.UserProfileCreationResponse;
import com.example.profileservice.mapper.UserProfileMapper;
import com.example.profileservice.model.UserProfile;
import com.example.profileservice.repository.ProfileRepository;

import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@Service
@RequiredArgsConstructor
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class UserProfileService {
    private static final Logger log = LoggerFactory.getLogger(UserProfileService.class);
    ProfileRepository userProfileRepository;
    UserProfileMapper userProfileMapper;

    // Create user profile
    public UserProfileCreationResponse create(UserProfileCreationRequest request) {
        log.debug("Create user with userId {}", request.getUserId());
        // Mapping user profile from request
        UserProfile profile = userProfileMapper.toUserProfileMapper(request);
        // Save user profile from request
        var temp = userProfileRepository.save(profile);
        log.debug("Created profile: {}", profile.getUserId());
        // Mapping to response
        return userProfileMapper.toUserProfile(temp);
    }

    public UserProfileResponse get(String profileId) {
        // Find user by profileId
        Optional<UserProfile> profile = userProfileRepository.findById(profileId);
        // Mapping to response user profile
        return userProfileMapper.toUserProfileResponse(profile.orElse(null));
    }

    // Update profile
    public UserProfileUpdateResponse update(UserProfileUpdateRequest userProfileUpdateRequest, String userId) {
        // Find exists user from db
        UserProfile existingUser = userProfileRepository.findByUserId(userId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOTFOUND));
        log.debug("Updating profile for userId: {}", existingUser.getUserId());
        // Mapping to user profile for existed profile
        UserProfile updatedUser = userProfileMapper.toUserProfile(userProfileUpdateRequest);
        // Keep user profile and user id
        updatedUser.setProfileId(existingUser.getProfileId());
        updatedUser.setUserId(existingUser.getUserId());
        // Update to database
        UserProfile savedUser = userProfileRepository.save(updatedUser);
        log.info("Updated profile successfully: {}", savedUser);
        // Mapping to response user profile
        return userProfileMapper.toUserProfileUpdateResponse(savedUser);
    }

    // Get user profile with page and size of page
    public List<?> getAll(int page, int size) {
        return userProfileRepository.findAllProfilesPaged(page, size);
    }
}