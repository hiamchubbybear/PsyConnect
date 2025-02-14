package dev.psyconnect.profile_service.mapper;

import org.mapstruct.Mapper;

import dev.psyconnect.profile_service.dto.UserProfileResponse;
import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.request.UserProfileUpdateRequest;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileUpdateResponse;
import dev.psyconnect.profile_service.model.UserProfile;

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
