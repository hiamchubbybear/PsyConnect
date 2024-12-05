package com.example.IdentityService.repository.https;

import com.example.IdentityService.dto.request.UserProfileCreationRequest;
import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

@FeignClient(name = "profile-service" , url = "http://localhost:8081/profile")
public interface ProfileRepository {
    @PostMapping(path = "/user", produces = MediaType.APPLICATION_JSON_VALUE)
    Object createProfile(@RequestBody UserProfileCreationRequest request);
}