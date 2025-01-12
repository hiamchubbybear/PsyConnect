package com.example.profileservice.mapper;

import com.example.profileservice.controller.UserProfileController;
import com.example.profileservice.dto.UserProfileResponse;
import com.example.profileservice.dto.request.UserProfileUpdateRequest;
import com.example.profileservice.dto.response.UserProfileUpdateResponse;
import org.mapstruct.Mapper;

import com.example.profileservice.dto.request.UserProfileCreationRequest;
import com.example.profileservice.dto.response.UserProfileCreationResponse;
import com.example.profileservice.model.UserProfile;

@Mapper(componentModel = "spring")
public interface UserProfileMapper {
    UserProfile toUserProfileMapper(UserProfileCreationRequest userProfile);

    UserProfileCreationResponse toUserProfile(UserProfile userProfile);

    UserProfile toUserProfile(UserProfileUpdateRequest userProfile);

    UserProfileUpdateResponse toUserProfileUpdateResponse(UserProfile updatedUser);

    UserProfileUpdateResponse toUserProfileUpdateResponse(UserProfileUpdateRequest updatedUser);

    UserProfileUpdateResponse toUserProfileUpdateResponse(UserProfileUpdateResponse updatedUser);

    UserProfileResponse toUserProfileResponse(UserProfile userProfile);
}
