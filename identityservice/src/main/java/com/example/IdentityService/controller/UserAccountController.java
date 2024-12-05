package com.example.IdentityService.controller;

import com.example.IdentityService.dto.request.UserAccountCreationRequest;
import com.example.IdentityService.dto.respone.UserAccountCreationRespone;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.service.UserAccountService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.server.ResponseStatusException;

@RestController
@RequestMapping("/account")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserAccountController {
    UserAccountService userAccountService;

    @PostMapping("/create")
    public ResponseEntity<UserAccountCreationRespone> register(@RequestBody UserAccountCreationRequest accountRequest) {
        var respone = userAccountService.createAccount(accountRequest);

        if (respone != null) {
            return ResponseEntity.ok(respone);
        } else {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid request data");
        }

    }
}
