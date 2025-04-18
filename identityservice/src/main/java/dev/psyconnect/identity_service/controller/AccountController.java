package dev.psyconnect.identity_service.controller;

import java.lang.reflect.GenericArrayType;
import java.util.List;
import java.util.UUID;

import dev.psyconnect.identity_service.configuration.AllowedRoles;
import dev.psyconnect.identity_service.model.Account;
import org.springframework.data.domain.Page;
import org.springframework.data.repository.query.Param;
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
@RequiredArgsConstructor()
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class AccountController {
    UserAccountService userAccountService;

    @PostMapping("/delete")
    public ApiResponse<DeleteAccountResponse> deleteAccount(
            @RequestBody DeleteAccountRequest request, @RequestHeader(value = "X-User-Id") UUID userId) {
        return new ApiResponse<>(userAccountService.deleteAccount(request, userId));
    }

    @DeleteMapping("/delete")
    public ApiResponse<Boolean> deleteAccountConfirm(
            @RequestBody DeleteAccountConfirmRequest request, @RequestHeader(value = "X-User-Id") UUID userId) {
        return new ApiResponse<>(userAccountService.deleteAccountRequest(request, userId));
    }

    @GetMapping("/info")
    public ApiResponse<UserInfoResponse> getAccount(@RequestHeader(value = "X-User-Id") UUID username) {
        return new ApiResponse<>(userAccountService.getUserAccount(username));
    }

    @GetMapping("/all/{page}")
    @AllowedRoles({"ADMIN"})
    public ApiResponse<Page<Account>> getAll(@PathVariable int page) {
        return new ApiResponse<>(userAccountService.getAllAccount(page));
    }

    @PutMapping("/update")
    public ApiResponse<UpdateAccountResponse> updateAccount(
            @RequestBody UpdateAccountRequest request, @RequestHeader(value = "X-User-Id") UUID uuid) {
        return new ApiResponse<>(userAccountService.updateAccount(request, uuid));
    }
}
