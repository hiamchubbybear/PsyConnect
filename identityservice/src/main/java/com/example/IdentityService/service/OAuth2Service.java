package com.example.IdentityService.service;

import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import com.example.IdentityService.mapper.UserAccountMapper;
import com.example.IdentityService.model.RoleEntity;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.repository.RoleRepository;
import com.example.IdentityService.repository.UserAccountRepository;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.web.authentication.AuthenticationFailureHandler;
import org.springframework.stereotype.Service;

import java.util.Set;

@Service
@RequiredArgsConstructor
public class OAuth2Service {
    private static final Logger log = LoggerFactory.getLogger(OAuth2Service.class);
    private final UserAccountRepository userAccountRepository;
    private final RoleRepository roleRepository;
    private final UserAccountMapper userAccountMapper;

    public UserAccountCreationResponse processOAuthPostLoginGoogle(String email) {
        UserAccount existUser = userAccountRepository.findByEmail(email);
        Set<RoleEntity> roleEntities = roleRepository.findAllByName("Client");
        if (existUser == null) {
            UserAccount newUser = UserAccount.builder().email(email).provider(com.example.IdentityService.num.Provider.GOOGLE).role(roleEntities).build();
            return userAccountMapper.toUserProfileCreationResponse(userAccountRepository.save(newUser));
        }
        throw new IllegalArgumentException("User already exists");
    }
}