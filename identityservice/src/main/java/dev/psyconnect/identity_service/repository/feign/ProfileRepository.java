package dev.psyconnect.identity_service.repository.feign;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

import dev.psyconnect.identity_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.identity_service.dto.response.UserAccountCreationResponse;

@FeignClient(name = "profileservice", url = "http://${baseUriProfileService}:8081/profile")
public interface ProfileRepository {
    @PostMapping(
            path = "/internal/user",
            produces = MediaType.APPLICATION_JSON_VALUE,
            consumes = MediaType.APPLICATION_JSON_VALUE)
    UserAccountCreationResponse createProfile(@RequestBody UserProfileCreationRequest request);
}
