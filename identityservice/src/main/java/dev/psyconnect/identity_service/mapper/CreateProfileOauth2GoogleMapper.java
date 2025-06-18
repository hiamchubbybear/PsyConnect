package dev.psyconnect.identity_service.mapper;

import dev.psyconnect.identity_service.dto.request.UserAccountCreationRequest;
import dev.psyconnect.identity_service.dto.request.UserProfileCreationRequest;
import org.mapstruct.Mapper;

import dev.psyconnect.identity_service.dto.request.CreateProfileOauth2GoogleRequest;
import dev.psyconnect.identity_service.dto.response.CreateProfileOauth2GoogleResponse;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")
public interface CreateProfileOauth2GoogleMapper {
    @Mapping(target = "username", source = "email")
    @Mapping(target = "password", constant = "")
    @Mapping(target = "address", constant = "")
    @Mapping(target = "role", constant = "USER")
    UserAccountCreationRequest fromOauth2Google(CreateProfileOauth2GoogleRequest request);

}
