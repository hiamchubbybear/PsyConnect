package com.example.IdentityService.service;

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
import com.example.IdentityService.repository.https.ProfileRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

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
        if (userAccountRepository.existsByUsername(request.getUsername()))
            throw new CustomExceptionHandler(ErrorCode.USER_NOTFOUND);
        if (userAccountRepository.existsByEmail(request.getEmail()))
            throw new CustomExceptionHandler(ErrorCode.USER_NOTFOUND);
        Set<RoleEntity> roles = new HashSet<>();
        roleRepository.findById(request.getRole()).ifPresentOrElse(roles::add, () -> {
            throw new CustomExceptionHandler(ErrorCode.ROLE_NOT_FOUND);
        });
        roles.forEach(roleEntity -> log.info("Role {}", roleEntity));
        UserAccount account = UserAccount.builder()
                .username(request.getUsername())
                .email(request.getEmail())
                .password(PasswordEncodingService.encoder(request.getPassword()))
                .id(UUID.randomUUID().toString())
                .role(roles)
                .build();
        var savedAccount = accountRepository.save(account);
        log.info("Created account: {}", savedAccount);
        UserProfileCreationRequest temp = mapper.toAccountResponse(request);
        temp.setUserId(savedAccount.getId());
        ObjectMapper objectMapper = new ObjectMapper();
        objectMapper.registerModule(new JavaTimeModule());
        var profileResponseRaw = profileRepository.createProfile(temp);
        UserProfileCreationRequest profileResponse = objectMapper.convertValue(profileResponseRaw, UserProfileCreationRequest.class);
        return UserAccountCreationResponse.builder().username(request.getUsername())
                .password(account.getPassword()).email(account.getEmail())
                .role(roles.stream().findFirst().get().getName()).build();
    }
}
