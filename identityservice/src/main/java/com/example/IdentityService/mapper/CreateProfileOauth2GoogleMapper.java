package com.example.IdentityService.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

import com.example.IdentityService.dto.request.CreateProfileOauth2GoogleRequest;
import com.example.IdentityService.dto.response.CreateProfileOauth2GoogleResponse;

@Mapper(componentModel = "spring")
public interface CreateProfileOauth2GoogleMapper {
    @Mapping(ignore = true, target = "roles")
    @Mapping(ignore = true, target = "status")
    CreateProfileOauth2GoogleResponse toCreateProfileOauth2Google(CreateProfileOauth2GoogleRequest request);
}
