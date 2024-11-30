package com.example.IdentityService.Service;

import com.example.IdentityService.DTO.Request.AccountRequest;
import com.example.IdentityService.DTO.Respone.AccountRespone;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
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
public class UserController {
    @Autowired
    public  UserAccountService userAccountService;
    @PostMapping("/create")
    public ResponseEntity<AccountRespone> register(@RequestBody AccountRequest accountRequest) {
        var  respone  = userAccountService.createAccount(accountRequest);
        if(respone!=null) {
            return ResponseEntity.ok(respone);
        }else {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid request data");
        }

    }
}
