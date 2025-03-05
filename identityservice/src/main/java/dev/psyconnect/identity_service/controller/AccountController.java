package dev.psyconnect.identity_service.controller;

import java.util.UUID;

import org.springframework.web.bind.annotation.*;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.dto.request.DeleteAccountConfirmRequest;
import dev.psyconnect.identity_service.dto.request.DeleteAccountRequest;
import dev.psyconnect.identity_service.dto.request.UpdateAccountRequest;
import dev.psyconnect.identity_service.dto.response.DeleteAccountResponse;
import dev.psyconnect.identity_service.dto.response.UpdateAccountResponse;
import dev.psyconnect.identity_service.dto.response.UserInfoResponse;
import dev.psyconnect.identity_service.service.UserAccountService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RequestMapping("/account")
@RestController
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
@RequiredArgsConstructor
public class AccountController {
    UserAccountService userAccountService;

    @PostMapping("/delete")
    public ApiResponse<DeleteAccountResponse> deleteAccount(
            @RequestBody DeleteAccountRequest request,
            @RequestHeader(value = "X-User-Id", required = true) UUID userId) {
        return new ApiResponse<>(userAccountService.deleteAccount(request, userId));
    }

    //                        .header("X-User-Id", accountId)
    //                        .header("X-Profile-Id", profileId)
    @DeleteMapping("/delete")
    public ApiResponse<Boolean> deleteAccountConfirm(
            @RequestBody DeleteAccountConfirmRequest request,
            @RequestHeader(value = "X-User-Id", required = true) UUID userId) {
        return new ApiResponse<>(userAccountService.deleteAccountRequest(request, userId));
    }

    @GetMapping("/info")
    public ApiResponse<UserInfoResponse> getAccount() {
        return new ApiResponse<>(userAccountService.getUserAccount());
    }

    @PutMapping("/update/{uuid}")
    public ApiResponse<UpdateAccountResponse> updateAccount(
            @RequestBody UpdateAccountRequest request, @PathVariable String uuid) {
        return new ApiResponse<>(userAccountService.updateAccount(request, uuid));
    }
}
