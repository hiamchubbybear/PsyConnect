package dev.psyconnect.profile_service.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import dev.psyconnect.grpc.ProfileCreationRequest;
import dev.psyconnect.grpc.ProfileCreationResponse;
import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.request.UserProfileUpdateRequest;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileResponse;
import dev.psyconnect.profile_service.dto.response.UserProfileUpdateResponse;
import dev.psyconnect.profile_service.model.Profile;
import org.mapstruct.Mappings;

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

    @Mappings({
            @Mapping(source = "accountId", target = "accountId"),
            @Mapping(source = "profileId", target = "profileId"),
            @Mapping(source = "firstName", target = "firstName"),
            @Mapping(source = "lastName", target = "lastName"),
            @Mapping(source = "dob", target = "dob", dateFormat = "yyyy-MM-dd"),
            @Mapping(source = "address", target = "address"),
            @Mapping(source = "gender", target = "gender"),
            @Mapping(source = "avatarUri", target = "avatarUri"),
    })
    UserProfileResponse toUserProfileRequest(Profile profile);

    UserProfileCreationRequest toUserProfileRequest(ProfileCreationRequest profile);

    ProfileCreationResponse toUserProfileResponse(UserProfileCreationResponse profileCreationResponse);

    @Mapping(source = "profileId", target = "profileId")
    @Mapping(source = "firstName", target = "firstName")
    @Mapping(source = "lastName", target = "lastName")
    @Mapping(source = "dob", target = "dob")
    @Mapping(source = "address", target = "address")
    @Mapping(source = "gender", target = "gender")
    @Mapping(source = "avatarUri", target = "avatarUri")
    UserProfileCreationResponse toDTO(ProfileCreationResponse proto);

    @Mapping(target = "profileId", source = "profileId")
    @Mapping(target = "firstName", source = "firstName")
    @Mapping(target = "lastName", source = "lastName")
    @Mapping(target = "dob", source = "dob")
    @Mapping(target = "address", source = "address")
    @Mapping(target = "gender", source = "gender")
    @Mapping(target = "avatarUri", source = "avatarUri")
    default ProfileCreationResponse toProto(UserProfileCreationResponse dto) {
        return ProfileCreationResponse.newBuilder()
                .setProfileId(dto.getProfileId())
                .setFirstName(dto.getFirstName())
                .setLastName(dto.getLastName())
                .setDob(dto.getDob())
                .setAddress(dto.getAddress())
                .setGender(dto.getGender())
                .setAvatarUri(dto.getAvatarUri())
                .build();
    }
}
