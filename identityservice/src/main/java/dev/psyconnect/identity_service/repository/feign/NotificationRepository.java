package dev.psyconnect.identity_service.repository.feign;

import org.springframework.cloud.openfeign.FeignClient;
import org.springframework.http.MediaType;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

import dev.psyconnect.identity_service.dto.request.ActivateAccountNotificationRequest;
import dev.psyconnect.identity_service.dto.request.CreateAccountNotificationRequest;
import dev.psyconnect.identity_service.dto.request.DeleteNotificationRequest;
import feign.Response;

@FeignClient(name = "mail", url = "http://${baseUriNotificationService}:8082/noti")
public interface NotificationRepository {
    @PostMapping(
            path = "/internal/user",
            produces = MediaType.APPLICATION_JSON_VALUE,
            consumes = MediaType.APPLICATION_JSON_VALUE)
    Response sendCreateEmail(@RequestBody CreateAccountNotificationRequest createAccountEmailRequest);

    @PostMapping(
            path = "/internal/account/delete",
            produces = MediaType.APPLICATION_JSON_VALUE,
            consumes = MediaType.APPLICATION_JSON_VALUE)
    Response sendDeleteEmail(@RequestBody DeleteNotificationRequest deleteNotificationRequest);

    @PostMapping(
            path = "/internal/account/req/activate",
            produces = MediaType.APPLICATION_JSON_VALUE,
            consumes = MediaType.APPLICATION_JSON_VALUE)
    Response sendActivateToken(@RequestBody ActivateAccountNotificationRequest deleteNotificationRequest);
}
