package com.example.IdentityService.repository.https;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestPart;
import org.springframework.web.multipart.MultipartFile;

import com.example.IdentityService.configuration.FeignConfiguration;
import com.example.IdentityService.dto.request.UserProfileCreationRequest;

@FeignClient(name = "profile-service", url = "http://localhost:8081/profile", configuration = FeignConfiguration.class)
public interface ProfileRepository {

    @PostMapping(
            path = "/user",
            produces = MediaType.APPLICATION_JSON_VALUE,
            consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    Object createProfile(
            @RequestPart("json") UserProfileCreationRequest request, @RequestPart("file") MultipartFile file);
}
