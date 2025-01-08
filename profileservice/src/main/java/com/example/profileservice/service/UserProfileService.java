package com.example.profileservice.service;

import java.util.Optional;

import com.example.profileservice.controller.UserProfileController;
import com.example.profileservice.dto.request.UserProfileUpdateRequest;
import com.example.profileservice.dto.response.UserProfileUpdateResponse;
import com.example.profileservice.globalexceptionhandle.CustomExceptionHandler;
import com.example.profileservice.globalexceptionhandle.ErrorCode;
import jakarta.validation.constraints.Null;
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

    public UserProfileCreationResponse createProfile(UserProfileCreationRequest request) {
        log.info("Create user with userId {}", request.getUserId());
        UserProfile profile = userProfileMapper.toUserProfileMapper(request);
        var temp = userProfileRepository.save(profile);
        log.info("Created profile: {}", profile.getUserId());
        return userProfileMapper.toUserProfile(temp);
    }

    public UserProfile getProfile(String profileId) {
        Optional<UserProfile> profile = userProfileRepository.findById(profileId);
        return profile.orElse(null);
    }

    public UserProfileUpdateResponse profileUpdateResponse(UserProfileUpdateRequest userProfileUpdateRequest, String email) {
        var updateUser = userProfileRepository.findByEmail(email);
        UserProfile updatedUser = userProfileRepository.save(updateUser.orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOTFOUND)));
        return userProfileMapper.toUserProfileUpdateResponse(updatedUser);
    }
}