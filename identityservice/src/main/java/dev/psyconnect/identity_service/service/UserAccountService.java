package dev.psyconnect.identity_service.service;

import java.sql.Timestamp;
import java.time.Instant;
import java.time.ZoneOffset;
import java.time.temporal.ChronoUnit;
import java.util.*;

import jakarta.transaction.Transactional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;

import dev.psyconnect.identity_service.dto.request.AcitvateAccountRequest;
import dev.psyconnect.identity_service.dto.request.CreateAccountNotificationRequest;
import dev.psyconnect.identity_service.dto.request.UserAccountCreationRequest;
import dev.psyconnect.identity_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.identity_service.dto.response.ActivateAccountResponse;
import dev.psyconnect.identity_service.dto.response.UserAccountCreationResponse;
import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.mapper.UserAccountMapper;
import dev.psyconnect.identity_service.model.ActivationModel;
import dev.psyconnect.identity_service.model.RoleEntity;
import dev.psyconnect.identity_service.model.UserAccount;
import dev.psyconnect.identity_service.repository.ActivateRepository;
import dev.psyconnect.identity_service.repository.RoleRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import dev.psyconnect.identity_service.repository.feign.NotificationRepository;
import dev.psyconnect.identity_service.repository.feign.ProfileRepository;
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
    NotificationRepository notificationRepository;
    ActivateRepository activateRepository;
    ZoneOffset offset = ZoneOffset.ofHours(7);

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
                .isActivated(false)
                .createdAt(Timestamp.from(Instant.now()))
                .provider(Provider.ORDINARY)
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
        String activateCode = String.valueOf(new Random().nextInt(90000) + 10000);
        // Trigger to Send Verified Code Request -> POST -x http://localhost:8082/noti/internal/create
        // If response is 200 create a new Activate Model
        if (notificationRepository
                        .sendCreateEmail(CreateAccountNotificationRequest.builder()
                                .fullname(request.getFirstName() + request.getLastName())
                                .username(request.getUsername())
                                .email(request.getEmail())
                                .code(activateCode)
                                .build())
                        .status()
                == 200) {
            // Create a new ActivateModel and set the expired time and issue time
            activateRepository.save(ActivationModel.builder()
                    .username(savedAccount.getUsername())
                    .activateId(savedAccount.getUserId())
                    .expires(Timestamp.from(
                            savedAccount.getCreatedAt().toInstant().plus(10, ChronoUnit.MINUTES)))
                    .issuedAt(Timestamp.from(savedAccount.getCreatedAt().toInstant()))
                    .token(activateCode)
                    .revoked(false)
                    .build());
        }
        // Trigger to Profile Service Request -> POST -x http://localhost:8081/profile/internal/user
        if (profileRepository.createProfile(temp) == null)
            // If there are no response from Profile Service -> throw a uncategorized exception
            throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
        // Mapping and logging and return to client response
        return UserAccountCreationResponse.builder()
                .username(request.getUsername())
                .email(account.getEmail())
                .imageUri(request.getAvatarUri())
                .role(roles)
                .build();
    }

    @Transactional
    public ActivateAccountResponse activateAccount(AcitvateAccountRequest activateAccountRequest) {
        String email = activateAccountRequest.getEmail();
        log.debug("Activation Model account for user email {}", email);
        if (!userAccountRepository.existsByEmail(email)) {
            throw new CustomExceptionHandler(ErrorCode.USER_NOTFOUND);
        }
        var foundObject = activateRepository
                .findByToken(activateAccountRequest.getToken())
                .orElse(null);
        if (foundObject == null) {
            return ActivateAccountResponse.builder()
                    .isSuccess(false)
                    .message("Invalid or expired activation token")
                    .build();
        }
        if (foundObject.getExpires().before(new Date())) {
            return ActivateAccountResponse.builder()
                    .isSuccess(false)
                    .message("Activation token has expired")
                    .build();
        }
        userAccountRepository.activateUser(email);
        return ActivateAccountResponse.builder()
                .isSuccess(true)
                .message("Account activation successful")
                .build();
    }
}
