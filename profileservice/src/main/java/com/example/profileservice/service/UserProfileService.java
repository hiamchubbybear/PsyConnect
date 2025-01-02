package com.example.profileservice.service;

import java.util.Optional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import com.example.profileservice.dto.request.UserProfileCreationRequest;
import com.example.profileservice.dto.respone.UserProfileCreationResponse;
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
        UserProfile profile = userProfileMapper.toUserProfileMapper(request);
        var temp = userProfileRepository.save(profile);
        log.info("Created profile: {}", profile);
        return userProfileMapper.toUserProfile(temp);
    }

    public UserProfile getProfile(String profileId) {
        Optional<UserProfile> profile = userProfileRepository.findById(profileId);
        return profile.orElse(null);
    }
}
