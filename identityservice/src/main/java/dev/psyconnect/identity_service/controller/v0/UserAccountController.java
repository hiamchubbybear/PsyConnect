package dev.psyconnect.identity_service.controller.v0;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.dto.request.*;
import dev.psyconnect.identity_service.dto.response.ActivateAccountResponse;
import dev.psyconnect.identity_service.dto.response.UserAccountCreationResponse;
import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.kafka.producer.KafkaService;
import dev.psyconnect.identity_service.service.UserAccountService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequestMapping("/identity")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserAccountController {
    UserAccountService userAccountService;
    KafkaService kafkaService;

    @PostMapping(value = "/create")
    public ApiResponse<UserAccountCreationResponse> register(@RequestBody UserAccountCreationRequest accountRequest) {
        return new ApiResponse<>(userAccountService.createAccount(accountRequest, Provider.ORDINARY));
    }

    @PostMapping(value = "/activate")
    public ApiResponse<ActivateAccountResponse> activateAccount(
            @RequestBody ActivateAccountRequest activateAccountRequest) {
        return new ApiResponse<>(userAccountService.activateAccount(activateAccountRequest));
    }

    @PostMapping(value = "/req/activate")
    public ApiResponse<Boolean> activateAccount(
            @RequestBody RequestActivationAccount activateAccountNotificationRequest) {
        return new ApiResponse<>(userAccountService.requestActivateAccount(activateAccountNotificationRequest));
    }

    @GetMapping("/hello")
    public ApiResponse<String> register() {
        return new ApiResponse<>("Hello World");
    }

    @PutMapping("/password")
    public ApiResponse<Boolean> register(@RequestBody PasswordResetRequest passwordResetRequest) {
        return new ApiResponse<>(userAccountService.resetPassword(passwordResetRequest));
    }

    @PostMapping(value = "/req/reset-password")
    public ApiResponse<Boolean> requestPasswordRequest(@RequestBody RequestPasswordReset requestPasswordReset) {
        return new ApiResponse<>(userAccountService.requestPasswordReset(requestPasswordReset));
    }
}
