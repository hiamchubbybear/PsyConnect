package dev.psyconnect.identity_service.service;

import java.util.Set;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import dev.psyconnect.identity_service.dto.request.CreateProfileOauth2GoogleRequest;
import dev.psyconnect.identity_service.dto.response.CreateProfileOauth2GoogleResponse;
import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.mapper.CreateProfileOauth2GoogleMapper;
import dev.psyconnect.identity_service.mapper.UserAccountMapper;
import dev.psyconnect.identity_service.model.Account;
import dev.psyconnect.identity_service.model.RoleEntity;
import dev.psyconnect.identity_service.repository.RoleRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@Service
@RequiredArgsConstructor
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class OAuth2Service {
    private static final Logger log = LoggerFactory.getLogger(OAuth2Service.class);
    UserAccountRepository userAccountRepository;
    RoleRepository roleRepository;
    UserAccountMapper userAccountMapper;
    private final CreateProfileOauth2GoogleMapper createProfileOauth2GoogleMapper;

    public CreateProfileOauth2GoogleResponse processOAuthPostLoginGoogle(
            CreateProfileOauth2GoogleRequest createProfileOauth2GoogleRequest) {
        String email = createProfileOauth2GoogleRequest.getEmail();
        String avatarUri = createProfileOauth2GoogleRequest.getAvatarUri();
        Account existUser = userAccountRepository.findByEmail(email).orElse(null);
        Set<RoleEntity> roleEntities = roleRepository.findAllByRoleId("Client");
        log.info("Client Role {} :", roleEntities.stream().findFirst().toString());
        if (existUser == null) {
            Account newUser = Account.builder()
                    .email(email)
                    .provider(Provider.GOOGLE)
                    .role(roleEntities)
                    .isActivated(true)
                    .build();
            log.info("Creating new user account with provider GOOGLE {} ", email);
            var savedUser = userAccountRepository.save(newUser);
            log.info("Create new user account with accountId {} ", savedUser.getAccountId());
            var createProfileRequest = userAccountMapper.toUserProfileCreationRequest(savedUser);
            createProfileRequest.setAvatarUri(avatarUri);
            log.info("Trigger to profile service {} ", savedUser.getAccountId());
            var response =
                    createProfileOauth2GoogleMapper.toCreateProfileOauth2Google(createProfileOauth2GoogleRequest);
            response.setCreationType(Provider.GOOGLE.toString());
            response.setRoles(roleEntities);
            response.setAvatarUri(avatarUri);
            return response;
        }
        throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
    }
}
