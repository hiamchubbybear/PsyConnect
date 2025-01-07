package com.example.IdentityService.service;

import com.example.IdentityService.globalexceptionhandle.CustomExceptionHandler;
import com.example.IdentityService.globalexceptionhandle.ErrorCode;
import com.example.IdentityService.mapper.UserAccountMapper;
import com.example.IdentityService.model.RoleEntity;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.repository.RoleRepository;
import com.example.IdentityService.repository.UserAccountRepository;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.util.Set;

@Service
@RequiredArgsConstructor
public class OAuth2Service {
    private static final Logger log = LoggerFactory.getLogger(OAuth2Service.class);
    private final UserAccountRepository userAccountRepository;
    private final RoleRepository roleRepository;
    private final UserAccountMapper userAccountMapper;

    public void processOAuthPostLoginGoogle(String email)  {
        UserAccount existUser = userAccountRepository.findByEmail(email).orElse(null);
        Set<RoleEntity> roleEntities = roleRepository.findAllByRoleId("Client");
        log.info("Client Role {} :",roleEntities.stream().findFirst().toString());
        if (existUser == null) {
            UserAccount newUser = UserAccount.builder()
                    .email(email).provider(com.example.IdentityService.num.Provider.GOOGLE)
                    .role(roleEntities).build();
            log.info("Creating new user account with provider GOOGLE {} " , email);
            userAccountMapper.toUserProfileCreationResponse(userAccountRepository.save(newUser));
            return;
        }
        throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
    }
}