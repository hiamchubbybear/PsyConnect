package dev.psyconnect.profile_service.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import dev.psyconnect.profile_service.dto.UserProfileResponse;
import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.request.UserProfileUpdateRequest;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileUpdateResponse;
import dev.psyconnect.profile_service.model.Profile;

@Mapper(componentModel = "spring")
public interface UserProfileMapper {
    @Mapping(ignore = true, target = "dob")
    Profile toUserProfileMapper(UserProfileCreationRequest userProfile);

    @Mapping(ignore = true, target = "dob")
    UserProfileCreationResponse toUserProfile(Profile profile);

    @Mapping(ignore = true, target = "dob")
    Profile toUserProfile(UserProfileUpdateRequest userProfile);

    @Mapping(ignore = true, target = "dob")
    UserProfileUpdateResponse toUserProfileUpdateResponse(Profile updatedUser);

    @Mapping(ignore = true, target = "dob")
    UserProfileUpdateResponse toUserProfileUpdateResponse(UserProfileUpdateRequest updatedUser);

    UserProfileUpdateResponse toUserProfileUpdateResponse(UserProfileUpdateResponse updatedUser);

    @Mapping(ignore = true, target = "dob")
    UserProfileResponse toUserProfileResponse(Profile profile);
}
