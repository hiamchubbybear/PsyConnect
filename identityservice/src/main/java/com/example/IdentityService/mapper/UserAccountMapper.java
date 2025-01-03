package com.example.IdentityService.mapper;

import com.example.IdentityService.dto.request.UserAccountCreationRequest;
import com.example.IdentityService.dto.request.UserProfileCreationRequest;
import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import com.example.IdentityService.model.UserAccount;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface UserAccountMapper {
    @Mapping(ignore = true, target = "userId")
    @Mapping(target = "role" ,ignore = true)
    UserProfileCreationRequest toAccountResponse(UserAccountCreationRequest userAccount);

    @Mapping(ignore = true, target = "userId")
    @Mapping(target = "role" ,ignore = true)
    UserProfileCreationRequest toUserProfileCreationRequest(UserAccount userAccount);
    @Mapping(target = "role" ,ignore = true)
    UserAccountCreationResponse toUserAccountResponse(UserProfileCreationRequest userProfileCreationRequest);
}