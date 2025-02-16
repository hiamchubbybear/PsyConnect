package dev.psyconnect.identity_service.controller;

import java.util.UUID;

import dev.psyconnect.identity_service.model.UserAccount;
import org.springframework.web.bind.annotation.*;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.dto.request.DeleteAccountConfirmRequest;
import dev.psyconnect.identity_service.dto.request.DeleteAccountRequest;
import dev.psyconnect.identity_service.dto.response.DeleteAccountResponse;
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
            @RequestBody DeleteAccountRequest request, @PathVariable UUID uuid) {
        return new ApiResponse<>(userAccountService.deleteAccount(request, uuid));
    }

    @DeleteMapping("/delete")
    public ApiResponse<Boolean> deleteAccountConfirm(
            @RequestBody DeleteAccountConfirmRequest request, @PathVariable UUID uuid) {
        return new ApiResponse<>(userAccountService.deleteAccountRequest(request, uuid));
    }
    @GetMapping("/info")
    public ApiResponse<UserAccount> getAccount(){
        return new ApiResponse<>(userAccountService.getUserAccount());
    }
}
