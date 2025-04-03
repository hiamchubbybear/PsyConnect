package dev.psyconnect.identity_service.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import dev.psyconnect.grpc.ProfileCreationRequest;
import dev.psyconnect.grpc.ProfileCreationResponse;
import dev.psyconnect.identity_service.dto.request.UserAccountCreationRequest;
import dev.psyconnect.identity_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.identity_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.identity_service.model.Account;

@Mapper(componentModel = "spring")
public interface UserAccountMapper {
    @Mapping(target = "role", ignore = true)
    UserProfileCreationRequest toAccountResponse(UserAccountCreationRequest userAccount);

    @Mapping(target = "role", ignore = true)
    UserProfileCreationRequest toUserProfileCreationRequest(Account userAccount);

    @Mapping(target = "role", ignore = true)
    ProfileCreationRequest toUserProfileCreateRequest(UserProfileCreationRequest userAccount);

    @Mapping(target = "dob", ignore = true)
    UserProfileCreationResponse toUserProfileCreateResponse(ProfileCreationResponse userAccount);
}
