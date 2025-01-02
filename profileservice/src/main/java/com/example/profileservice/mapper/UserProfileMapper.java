package com.example.profileservice.mapper;

import org.mapstruct.Mapper;

import com.example.profileservice.dto.request.UserProfileCreationRequest;
import com.example.profileservice.dto.respone.UserProfileCreationResponse;
import com.example.profileservice.model.UserProfile;

@Mapper(componentModel = "spring")
public interface UserProfileMapper {
    UserProfile toUserProfileMapper(UserProfileCreationRequest userProfile);

    UserProfileCreationResponse toUserProfile(UserProfile userProfile);
}
