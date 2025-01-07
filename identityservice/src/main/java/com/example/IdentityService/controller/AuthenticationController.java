package com.example.IdentityService.controller;

import com.example.IdentityService.apiresponse.ApiResponse;
import com.example.IdentityService.apiresponse.CustomResponseWrapper;
import com.example.IdentityService.configuration.ValidateLoginType;
import com.example.IdentityService.dto.request.AuthenticationRequest;
import com.example.IdentityService.dto.response.AuthenticationResponse;
import com.example.IdentityService.num.Provider;
import com.example.IdentityService.service.AuthenticationService;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.naming.AuthenticationException;


@RestController
@RequestMapping("/auth")
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
@RequiredArgsConstructor
public class AuthenticationController {
    AuthenticationService authenticationService;

    @CustomResponseWrapper
    @PostMapping("/login")
    public ApiResponse<AuthenticationResponse> loginRequest
            (@RequestBody AuthenticationRequest authenticationRequest
                    , @RequestParam @ValidateLoginType String loginType) throws AuthenticationException {
        return new ApiResponse<>(authenticationService.authenticate(authenticationRequest, loginType));
    }

}
