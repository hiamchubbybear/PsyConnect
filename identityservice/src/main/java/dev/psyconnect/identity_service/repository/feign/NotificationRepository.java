package dev.psyconnect.identity_service.repository.feign;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

import dev.psyconnect.identity_service.dto.request.CreateAccountNotificationRequest;
import feign.Response;

@FeignClient(name = "mail", url = "http://localhost:8082/noti")
public interface NotificationRepository {
    @PostMapping(
            path = "/internal/user",
            produces = MediaType.APPLICATION_JSON_VALUE,
            consumes = MediaType.APPLICATION_JSON_VALUE)
    Response sendCreateEmail(@RequestBody CreateAccountNotificationRequest createAccountEmailRequest);
}
