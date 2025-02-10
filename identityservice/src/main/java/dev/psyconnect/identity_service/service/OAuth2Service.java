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
import dev.psyconnect.identity_service.model.RoleEntity;
import dev.psyconnect.identity_service.model.UserAccount;
import dev.psyconnect.identity_service.repository.RoleRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import dev.psyconnect.identity_service.repository.feign.ProfileRepository;
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
    ProfileRepository profileRepository;
    private final CreateProfileOauth2GoogleMapper createProfileOauth2GoogleMapper;

    public CreateProfileOauth2GoogleResponse processOAuthPostLoginGoogle(
            CreateProfileOauth2GoogleRequest createProfileOauth2GoogleRequest) {
        // Get email from request
        String email = createProfileOauth2GoogleRequest.getEmail();
        // Get avatarUri from request
        String avatarUri = createProfileOauth2GoogleRequest.getAvatarUri();
        // Find the existed user by email
        UserAccount existUser = userAccountRepository.findByEmail(email).orElse(null);
        // Find the role of user by Client id
        Set<RoleEntity> roleEntities = roleRepository.findAllByRoleId("Client");
        // Logging
        log.info("Client Role {} :", roleEntities.stream().findFirst().toString());
        // If null create a new user only have email  , provider , role
        if (existUser == null) {
            UserAccount newUser = UserAccount.builder()
                    .email(email)
                    .provider(Provider.GOOGLE)
                    .role(roleEntities)
                    .build();
            log.info("Creating new user account with provider GOOGLE {} ", email);
            // Write new user to db
            var savedUser = userAccountRepository.save(newUser);
            // Logging
            log.info("Create new user account with accountId {} ", savedUser.getUserId());
            // Mapping request to trigger to profile service with http:localhost:8081/profile/create
            var createProfileRequest = userAccountMapper.toUserProfileCreationRequest(savedUser);
            // Set avatarUri for profile
            createProfileRequest.setAvatarUri(avatarUri);
            // Logging
            log.info("Trigger to profile service {} ", savedUser.getUserId());
            // Trigger to profile service with http:localhost:8081/profile/create
            profileRepository.createProfile(createProfileRequest);
            // Mapping
            var response =
                    createProfileOauth2GoogleMapper.toCreateProfileOauth2Google(createProfileOauth2GoogleRequest);
            // Set attribute for response
            response.setCreationType(Provider.GOOGLE.toString());
            response.setRoles(roleEntities);
            response.setAvatarUri(avatarUri);
            return response;
        }
        // /Throw new Uncategorized exception
        throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
    }
}
