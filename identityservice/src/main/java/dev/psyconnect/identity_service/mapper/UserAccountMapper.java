package dev.psyconnect.identity_service.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import dev.psyconnect.identity_service.dto.request.UserAccountCreationRequest;
import dev.psyconnect.identity_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.identity_service.model.Account;

@Mapper(componentModel = "spring")
public interface UserAccountMapper {
    @Mapping(target = "role", ignore = true)
    @Mapping(target = "accountId", ignore = true)
    UserProfileCreationRequest toAccountResponse(UserAccountCreationRequest userAccount);

    @Mapping(target = "role", ignore = true)
    UserProfileCreationRequest toUserProfileCreationRequest(Account userAccount);
}
