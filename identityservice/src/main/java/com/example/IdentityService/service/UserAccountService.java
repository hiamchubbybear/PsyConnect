package com.example.IdentityService.service;

import java.util.HashSet;
import java.util.Set;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import com.example.IdentityService.dto.request.UserAccountCreationRequest;
import com.example.IdentityService.dto.request.UserProfileCreationRequest;
import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import com.example.IdentityService.globalexceptionhandle.CustomExceptionHandler;
import com.example.IdentityService.globalexceptionhandle.ErrorCode;
import com.example.IdentityService.mapper.UserAccountMapper;
import com.example.IdentityService.model.RoleEntity;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.repository.RoleRepository;
import com.example.IdentityService.repository.UserAccountRepository;
import com.example.IdentityService.repository.feign.ProfileRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;

import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@Service
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserAccountService {
    private static final Logger log = LoggerFactory.getLogger(UserAccountService.class);
    RoleRepository roleRepository;
    UserAccountRepository userAccountRepository;
    ProfileRepository profileRepository;
    UserAccountMapper mapper;
    UserAccountRepository accountRepository;

    public UserAccountCreationResponse createAccount(UserAccountCreationRequest request) {
        // Validate if the user not found
        if (userAccountRepository.existsByUsername(request.getUsername()))
            throw new CustomExceptionHandler(ErrorCode.USERNAME_ALREADY_EXISTS);
        // Validate if the email not found
        if (userAccountRepository.existsByEmail(request.getEmail()))
            throw new CustomExceptionHandler(ErrorCode.EMAIL_ALREADY_EXISTS);
        // Find all role and save it to our response data
        Set<RoleEntity> roles = new HashSet<>(); // Get the first role Admin ->  Therapist -> Client
        roleRepository.findById(request.getRole()).ifPresentOrElse(roles::add, () -> {
            throw new CustomExceptionHandler(ErrorCode.ROLE_NOT_FOUND);
        });
        // Logging
        roles.forEach(roleEntity -> log.info("Role {}", roleEntity));
        UserAccount account = UserAccount.builder()
                .username(request.getUsername())
                .email(request.getEmail())
                .password(PasswordEncodingService.encoder(request.getPassword()))
                .role(roles)
                .build();
        // Get the data saved to "savedAccount"
        var savedAccount = accountRepository.save(account);
        // Logging "savedAccount" to audit action and activity
        log.info("Created account: {}", savedAccount.getUserId());

        UserProfileCreationRequest temp = mapper.toAccountResponse(request);
        // Set the userId to identity the account id
        temp.setUserId(account.getUserId());
        ObjectMapper objectMapper = new ObjectMapper();
        objectMapper.registerModule(new JavaTimeModule());
        // Trigger to Profile Service Request -> POST -x http://localhost:8081/profile/user
        profileRepository.createProfile(temp);
        // Mapping and logging and return to client response
        return UserAccountCreationResponse.builder()
                .username(request.getUsername())
                .password(account.getPassword())
                .email(account.getEmail())
                .role(roles.stream().findFirst().get().getName())
                .build();
    }
}
