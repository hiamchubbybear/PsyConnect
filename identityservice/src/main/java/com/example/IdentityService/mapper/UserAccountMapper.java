package com.example.IdentityService.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import com.example.IdentityService.dto.request.UserAccountCreationRequest;
import com.example.IdentityService.dto.request.UserProfileCreationRequest;
import com.example.IdentityService.model.UserAccount;

@Mapper(componentModel = "spring")
public interface UserAccountMapper {
    @Mapping(target = "role", ignore = true)
    @Mapping(target = "userId", ignore = true)
    UserProfileCreationRequest toAccountResponse(UserAccountCreationRequest userAccount);

    @Mapping(target = "role", ignore = true)
    UserProfileCreationRequest toUserProfileCreationRequest(UserAccount userAccount);
}
