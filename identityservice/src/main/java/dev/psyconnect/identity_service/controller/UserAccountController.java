package dev.psyconnect.identity_service.controller;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.dto.request.*;
import dev.psyconnect.identity_service.dto.response.ActivateAccountResponse;
import dev.psyconnect.identity_service.dto.response.UserAccountCreationResponse;
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

    @PostMapping(value = "/create")
    public ApiResponse<UserAccountCreationResponse> register(@RequestBody UserAccountCreationRequest accountRequest) {
        return new ApiResponse<>(userAccountService.createAccount(accountRequest));
    }

    @PostMapping(value = "/activate")
    public ApiResponse<ActivateAccountResponse> activateAccount(
            @RequestBody ActivateAccountRequest acitvateAccountRequest) {
        return new ApiResponse<>(userAccountService.activateAccount(acitvateAccountRequest));
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
}
