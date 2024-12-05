package com.example.profileservice.mapper;

import com.example.profileservice.dto.request.UserProfileCreationRequest;
import com.example.profileservice.dto.respone.UserProfileCreationRespone;
import com.example.profileservice.model.UserProfile;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")
public interface UserProfileMapper {
    UserProfile toUserProfileMapper(UserProfileCreationRequest userProfile);
    UserProfileCreationRespone toUserProfile(UserProfile userProfile);
}

