package com.example.IdentityService.service;

import com.example.IdentityService.dto.request.AccountCreationRequest;
import com.example.IdentityService.model.UserAccount;
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
@RequestMapping("/api/v1")
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserController {
    UserAccountService userAccountService;

    @PostMapping("/create")
    public ResponseEntity<UserAccount> register(@RequestBody AccountCreationRequest accountRequest) {
        var respone = userAccountService.createAccount(accountRequest);
        if (respone != null) {
            return ResponseEntity.ok(respone);
        } else {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid request data");
        }

    }
}
