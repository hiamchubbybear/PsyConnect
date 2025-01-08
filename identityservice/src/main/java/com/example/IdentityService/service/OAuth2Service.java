package com.example.IdentityService.service;

import com.example.IdentityService.dto.request.CreateProfileOauth2GoogleRequest;
import com.example.IdentityService.globalexceptionhandle.CustomExceptionHandler;
import com.example.IdentityService.globalexceptionhandle.ErrorCode;
import com.example.IdentityService.mapper.CreateProfileOauth2GoogleMapper;
import com.example.IdentityService.mapper.UserAccountMapper;
import com.example.IdentityService.model.RoleEntity;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.num.Provider;
import com.example.IdentityService.repository.RoleRepository;
import com.example.IdentityService.repository.UserAccountRepository;
import com.example.IdentityService.repository.https.ProfileRepository;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import java.util.Set;

@Service
@RequiredArgsConstructor
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class OAuth2Service {
    private static final Logger log = LoggerFactory.getLogger(OAuth2Service.class);
    UserAccountRepository userAccountRepository;
    RoleRepository roleRepository;
    UserAccountMapper userAccountMapper;
    ProfileRepository profileRepository;
    private final CreateProfileOauth2GoogleMapper createProfileOauth2GoogleMapper;

    public void processOAuthPostLoginGoogle(CreateProfileOauth2GoogleRequest createProfileOauth2GoogleRequest) {
        String email = createProfileOauth2GoogleRequest.getEmail();
        String avatarUri = createProfileOauth2GoogleRequest.getAvatarUri();
        UserAccount existUser = userAccountRepository.findByEmail(email).orElse(null);
        Set<RoleEntity> roleEntities = roleRepository.findAllByRoleId("Client");
        log.info("Client Role {} :", roleEntities.stream().findFirst().toString());
        if (existUser == null) {
            UserAccount newUser = UserAccount.builder()
                    .email(email).provider(com.example.IdentityService.num.Provider.GOOGLE)
                    .role(roleEntities).build();
            log.info("Creating new user account with provider GOOGLE {} ", email);
            var savedUser = userAccountRepository.save(newUser);
            log.info("Create new user account with accountId {} ", savedUser.getUserId());
            var createProfileRequest = userAccountMapper.toUserProfileCreationRequest(savedUser);
            createProfileRequest.setAvatarUri(avatarUri);
            log.info("Trigger to profile service {} ", savedUser.getUserId());
            profileRepository.createProfile(createProfileRequest);
            var response = createProfileOauth2GoogleMapper.toCreateProfileOauth2Google(createProfileOauth2GoogleRequest);
            response.setCreationType(Provider.GOOGLE.toString());
            response.setRoles(roleEntities);
            response.setAvatarUri(avatarUri);
            return;
        }
        throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
    }
}