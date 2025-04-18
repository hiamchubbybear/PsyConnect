package dev.psyconnect.identity_service.mapper;

import org.mapstruct.Mapper;

import dev.psyconnect.identity_service.dto.request.CreateProfileOauth2GoogleRequest;
import dev.psyconnect.identity_service.dto.response.CreateProfileOauth2GoogleResponse;

@Mapper(componentModel = "spring")
public interface CreateProfileOauth2GoogleMapper {
    //    @Mapping(ignore = true, target = "roles")
    //    @Mapping(ignore = true, target = "status")
    CreateProfileOauth2GoogleResponse toCreateProfileOauth2Google(CreateProfileOauth2GoogleRequest request);
}
