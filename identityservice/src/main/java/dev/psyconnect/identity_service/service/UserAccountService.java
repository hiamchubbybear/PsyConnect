package dev.psyconnect.identity_service.service;

import java.sql.Timestamp;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.*;

import jakarta.transaction.Transactional;

import lombok.AllArgsConstructor;
import lombok.NoArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.oauth2.server.resource.authentication.JwtAuthenticationToken;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;

import dev.psyconnect.identity_service.dto.request.*;
import dev.psyconnect.identity_service.dto.response.ActivateAccountResponse;
import dev.psyconnect.identity_service.dto.response.DeleteAccountResponse;
import dev.psyconnect.identity_service.dto.response.UserAccountCreationResponse;
import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.mapper.UserAccountMapper;
import dev.psyconnect.identity_service.model.RoleEntity;
import dev.psyconnect.identity_service.model.Token;
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
public class UserAccountService implements UserDetailsService {
    private static final Logger log = LoggerFactory.getLogger(UserAccountService.class);
    RoleRepository roleRepository;
    UserAccountRepository userAccountRepository;
    ProfileRepository profileRepository;
    UserAccountMapper mapper;
    UserAccountRepository accountRepository;
    NotificationRepository notificationRepository;
    ActivateRepository activateRepository;

    @Transactional
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
        // Trigger to Profile Service Request -> POST -x http://localhost:8081/profile/internal/user
        try {
            var profileResponse = profileRepository.createProfile(temp);
            if (profileResponse == null) {
                // If there are no response from Profile Service -> throw run time error exception
                throw new CustomExceptionHandler(ErrorCode.RUNTIME_ERROR);
            }
        } catch (Exception e) {
            // If there are no response from Profile Service -> throw run time error exception
            log.error("Error creating profile: {}", e.getMessage());
            throw new RuntimeException("Profile creation failed, rolling back...", e);
        }


        // Trigger to Send Verified Code Request -> POST -x http://localhost:8082/noti/internal/create and get the
        // response.
        var response = notificationRepository.sendCreateEmail(CreateAccountNotificationRequest.builder()
                .fullname(request.getFirstName() + request.getLastName())
                .username(request.getUsername())
                .email(request.getEmail())
                .code(activateCode)
                .build());
        // If response is 200 create a new Activate Model
        if (response.status() == 200) {
            // Create a new ActivateModel and set the expired time and issue time
            activateRepository.save(Token.builder()
                    .username(savedAccount.getUsername())
                    .activateId(savedAccount.getUserId())
                    .expires(Timestamp.from(
                            savedAccount.getCreatedAt().toInstant().plus(10, ChronoUnit.MINUTES)))
                    .issuedAt(Timestamp.from(savedAccount.getCreatedAt().toInstant()))
                    .token(activateCode)
                    .revoked(false)
                    .build());
        }
        //  If status code is 500 log an error
        else if (response.status() == 500) {
            log.error("There are error while sending token {}", response.body());
        }
        // Mapping and logging and return to client response
        return UserAccountCreationResponse.builder()
                .username(request.getUsername())
                .email(account.getEmail())
                .imageUri(request.getAvatarUri())
                .role(roles)
                .build();
    }

    @Transactional
    public ActivateAccountResponse activateAccount(ActivateAccountRequest activateAccountRequest) {
        // Declared and assign email
        String email = activateAccountRequest.getEmail();
        // Log an email for audit
        log.debug("Activation Model account for user email {}", email);
        // Check if email not exist throw an error 404 NOT FOUND
        if (!userAccountRepository.existsByEmail(email)) {
            throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        }
        // Declared and assign foundObject for Token Object if found assign it or else assign for null
        var foundObject = activateRepository
                .findByToken(activateAccountRequest.getToken())
                .orElse(null);
        // If foundObject is null return a response facilitate failed
        if (foundObject == null) {
            return ActivateAccountResponse.builder()
                    .isSuccess(false)
                    .message("Invalid or expired activation token")
                    .build();
        }
        // If expire time before current time return activation token has expired
        if (foundObject.getExpires().before(new Date())) {
            return ActivateAccountResponse.builder()
                    .isSuccess(false)
                    .message("Activation token has expired")
                    .build();
        }
        // Activate account if there are no validation
        userAccountRepository.activateUser(email);
        return ActivateAccountResponse.builder()
                .isSuccess(true)
                .message("Account activation successful")
                .build();
    }

    // Delete account internal use
    public boolean deleteAccountWithoutCheck(UUID userId) {
        try {
            // Delete without check call only
            userAccountRepository.deleteById(userId);
            // Builder a response for account
            return true;
        } catch (Exception e) {
            throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
        }
    }

    // Delete account external use
    public DeleteAccountResponse deleteAccount(DeleteAccountRequest deleteAccountRequest, UUID uuid) {
        // Find the user account object
        UserAccount userObject = userAccountRepository
                .findById(uuid)
                .orElseThrow(() ->
                        // If not found throw an exception 404
                        new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        log.debug("Deleting account {}", userObject.getUsername());
        // Check if username not match or password doesn't match
        if (!userObject.getUsername().equals(deleteAccountRequest.getUsername())
                || PasswordEncodingService.getBCryptPasswordEncoder()
                .matches(userObject.getPassword(), deleteAccountRequest.getPassword()))
            // If not found throw an exception 418
            throw new CustomExceptionHandler(ErrorCode.DELETE_ACCOUNT_FAILED);
        // Check if token doesn't match
        if (!userObject.getToken().getToken().equals(deleteAccountRequest.getToken()))
            // If not match throw an exception 407
            throw new CustomExceptionHandler(ErrorCode.TOKEN_INVALID);
        // Build an object to trigger to notification service
        var sendObject = DeleteNotificationRequest.builder()
                .email(deleteAccountRequest.getEmail())
                .secret(SecretEnum.getRandomMessage())
                .token(userObject.getToken().getToken())
                .build();
        // Trigger to Notification Service Request -> POST -x http://localhost:8081/noti/internal/account/delete
        if (notificationRepository.sendDeleteEmail(sendObject).status() != 200) {
            // If send error throw an exception 429
            throw new CustomExceptionHandler(ErrorCode.RUNTIME_ERROR);
        } else
            return DeleteAccountResponse.builder()
                    .isSuccess(true)
                    .secret(sendObject.getSecret())
                    .session(UUID.randomUUID())
                    .session(userObject.getSession())
                    .timestamp(deleteAccountRequest.getTimestamp())
                    .build();
    }

    public boolean deleteAccountRequest(DeleteAccountConfirmRequest removeAccountRequest, UUID uuid) {
        // Find the user account object
        UserAccount userObject = userAccountRepository
                .findById(uuid)
                .orElseThrow(() ->
                        // If not found throw an exception 404
                        new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        // Check if session match
        if (!removeAccountRequest.getSession().equals(userObject.getSession())
                || !userObject.getToken().getToken().equals(removeAccountRequest.getToken()))
            // If not throw a 418 status code
            throw new CustomExceptionHandler(ErrorCode.DELETE_ACCOUNT_FAILED);
        try {
            // Delete account
            userAccountRepository.deleteById(uuid);
            log.warn("User account {} deleted", userObject.getUsername());
        } catch (Exception e) {
            e.printStackTrace();
            // If not throw a 418 status code
            throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
        }
        // return a true
        return true;
    }

    public UserAccount getUserAccount() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        String username = authentication.getName();
        return userAccountRepository.findByUsername(username).orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
    }

    @Override
    public AuthenticationFilterRequest loadUserByUsername(String username) throws CustomExceptionHandler {
        var foundObject = userAccountRepository.findByUsername(username).orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        return AuthenticationFilterRequest.builder()
                .username(foundObject.getUsername())
                .password(foundObject.getPassword())
                .build();
    }
}
