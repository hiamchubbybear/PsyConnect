package com.example.IdentityService.service;

import com.example.IdentityService.dto.request.UserAccountCreationRequest;
import com.example.IdentityService.dto.request.UserProfileCreationRequest;
import com.example.IdentityService.dto.respone.UserAccountCreationRespone;
import com.example.IdentityService.model.RoleEntity;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.mapper.UserAccountMapper;
import com.example.IdentityService.repository.RoleRepository;
import com.example.IdentityService.repository.UserAccountRepository;
import com.example.IdentityService.repository.https.ProfileRepository;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.Collection;
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
    PasswordEncodingService passwordEncodingService;

    public UserAccount createAccount(UserAccountCreationRequest request) {
        if (userAccountRepository.existsByUsername(request.getUsername())
                || userAccountRepository.existsByEmail(request.getEmail()))
            throw new IllegalArgumentException("User existed");
        Set<RoleEntity> roles = new HashSet<>();
        roleRepository.findById(request.getRole()).ifPresentOrElse(roles::add, () -> {
            throw new IllegalArgumentException("Role not exist");
        });
        roles.stream().forEach(roleEntity -> log.info("Role {}", roleEntity));
        UserAccount account = UserAccount.builder()
                .username(request.getUsername())
                .email(request.getEmail())
                .password(new BCryptPasswordEncoder().encode(request.getPassword()))
                .id(UUID.randomUUID().toString())
                .role(roles)
                .build();
        var savedAccount = accountRepository.save(account);
        log.info("Created account: {}", savedAccount);
        UserProfileCreationRequest temp = mapper.toAccountRespone(request);
        temp.setUserId(savedAccount.getId());
        ObjectMapper objectMapper = new ObjectMapper();
        var profileResponseRaw = profileRepository.createProfile(temp);
        UserProfileCreationRequest profileResponse = objectMapper.convertValue(profileResponseRaw, UserProfileCreationRequest.class);
        var res = mapper.toUserAccountRespone(profileResponse);
        return savedAccount;
    }


}
