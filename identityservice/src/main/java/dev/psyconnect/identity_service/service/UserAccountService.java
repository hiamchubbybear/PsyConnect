package dev.psyconnect.identity_service.service;

import java.sql.Timestamp;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.*;

import jakarta.transaction.Transactional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.stereotype.Service;

import dev.psyconnect.identity_service.dto.request.*;
import dev.psyconnect.identity_service.dto.response.*;
import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.grpc.client.ProfileGRPCClient;
import dev.psyconnect.identity_service.interfaces.IUserAccountService;
import dev.psyconnect.identity_service.mapper.UserAccountMapper;
import dev.psyconnect.identity_service.model.*;
import dev.psyconnect.identity_service.repository.ActivateRepository;
import dev.psyconnect.identity_service.repository.NotificationRepository;
import dev.psyconnect.identity_service.repository.RoleRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@Service
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class UserAccountService implements UserDetailsService, IUserAccountService {
    private NotificationRepository notificationRepository;
    private static final Logger log = LoggerFactory.getLogger(UserAccountService.class);
    private static final int MINUTE_EXPIRED = 10;
    RoleRepository roleRepository;
    UserAccountRepository userAccountRepository;
    UserAccountMapper mapper;
    UserAccountRepository accountRepository;
    ActivateRepository activateRepository;
    private final PasswordEncodingService passwordEncodingService;
    private final TokenRepository tokenRepository;
    private final ProfileGRPCClient profileGRPCClient;

    @Transactional
    public UserAccountCreationResponse createAccount(UserAccountCreationRequest request) {
        // Validate if the user not found
        if (userAccountRepository.existsByUsername(request.getUsername()))
            throw new CustomExceptionHandler(ErrorCode.USERNAME_ALREADY_EXISTS);

        // Validate if the email not found
        if (userAccountRepository.existsByEmail(request.getEmail()))
            throw new CustomExceptionHandler(ErrorCode.EMAIL_ALREADY_EXISTS);

        // Generate and send activation code
        String token = generateActivationCode();
        log.info("Generated activation code: {}", token);
        var activateTokenEntity = sendNotification(ActivateAccountNotificationRequest.builder()
                .code(token)
                .email(request.getEmail())
                .fullname(request.getFirstName() + " " + request.getLastName())
                .username(request.getUsername())
                .build());
        var savedToken = tokenRepository.save(activateTokenEntity);
        log.info("Activate account entity for: {} ", activateTokenEntity.getUsername());
        // Find role by id and save it to our response data
        Set<RoleEntity> roles = roleRepository.findAllByRoleId(request.getRole().toUpperCase());
        UUID profileId = UUID.randomUUID();
        // Create Account and assign roles
        Account account = Account.builder()
                .username(request.getUsername())
                .isActivated(false)
                .createdAt(Timestamp.from(Instant.now()))
                .provider(Provider.ORDINARY)
                .email(request.getEmail())
                .password(PasswordEncodingService.encoder(request.getPassword()))
                .role(roles)
                .token(savedToken)
                .profileId(profileId)
                .build();

        // Save the Account object
        var savedAccount = accountRepository.save(account);
        log.info("Created account: {}", savedAccount.getAccountId());

        // Prepare response
        UserProfileCreationRequest temp = mapper.toAccountResponse(request);
        temp.setAccountId(account.getAccountId().toString());
        temp.setProfileId(profileId.toString());
        temp.setRole(request.getRole());

        // Send request to Profile Service
        try {
            var profileResponse = profileGRPCClient.createProfile(temp);
            // Call to gRPC client
            if (UUID.fromString(profileResponse.getProfileId())
                    .equals(savedAccount.getProfileId().toString())) {
                throw new CustomExceptionHandler(ErrorCode.RUNTIME_ERROR);
            }
        } catch (Exception e) {
            log.error("Error creating profile: {}", e);
            throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
        }
        String fullName = request.getFirstName() + " " + request.getLastName();
        String username = request.getUsername();
        String email = savedAccount.getEmail();
        notificationRepository.notificationSend(CreateAccountNotificationRequest.builder()
                .code(token)
                .username(username)
                .email(email)
                .fullname(fullName)
                .build());
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
        String activateToken = activateAccountRequest.getToken();
        log.info("Activate account entity for: {} ", activateToken);
        if (userAccountRepository.isActiveByEmail(email)) {
            throw new CustomExceptionHandler(ErrorCode.ACTIVATED);
        }
        // Log an email for audit
        log.info("Activation Model account for user email {}", email);
        // Check if email not exist throw an error 404 NOT FOUND
        if (!userAccountRepository.existsByEmail(email)) {
            throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        }
        // Declared and assign foundObject for Token Object if found assign it or else assign for null
        var foundObject = activateRepository
                .findByToken(activateToken)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.ACTIVATION_FAILED));
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
    public boolean deleteAccountWithoutCheck(UUID accountId) {
        try {
            // Delete without check call only
            userAccountRepository.deleteById(accountId);
            // Builder a response for account
            return true;
        } catch (Exception e) {
            throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
        }
    }

    // Delete account external use
    public DeleteAccountResponse deleteAccount(DeleteAccountRequest deleteAccountRequest, UUID uuid) {
        // Find the user account object
        Account userObject = userAccountRepository
                .findById(uuid)
                .orElseThrow(() ->
                        // If not found throw an exception 404
                        new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        log.info("Deleting account {}", userObject.getUsername());
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
        //        if (notificationRepository.sendDeleteEmail(sendObject).status() != 200) {
        //            // If send error throw an exception 429
        //            throw new CustomExceptionHandler(ErrorCode.RUNTIME_ERROR);
        //        } else
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
        Account userObject = userAccountRepository
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

    public UserInfoResponse getUserAccount() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        String username = authentication.getName();
        var userResponse = userAccountRepository
                .findByUsername(username)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        if (userAccountRepository.isActive(username))
            return UserInfoResponse.builder()
                    .username(username)
                    .accountId(userResponse.getAccountId())
                    .createdAt(userResponse.getCreatedAt().toLocalDateTime().toLocalDate())
                    .role(userResponse.getRole())
                    .email(userResponse.getEmail())
                    .isActivated(userResponse.isActivated())
                    .provider(userResponse.getProvider())
                    .build();
        else throw new CustomExceptionHandler(ErrorCode.ACCOUNT_INACTIVE);
    }

    @Override
    public AuthenticationFilterRequest loadUserByUsername(String username) throws CustomExceptionHandler {
        // Load user with username and parse it to security filter
        var foundObject = userAccountRepository
                .findByUsername(username)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        // Build object
        return AuthenticationFilterRequest.builder()
                .username(foundObject.getUsername())
                .password(foundObject.getPassword())
                .build();
    }

    @Transactional
    public UpdateAccountResponse updateAccount(UpdateAccountRequest updateAccountRequest, String accountId) {
        // Find existing account
        UUID uuid = UUID.fromString(accountId);
        log.info("Updating account {}", uuid);
        UUID userAccountId = UUID.fromString(AuthenticationService.extractUUIDClaim(
                SecurityContextHolder.getContext().getAuthentication()));
        log.info("UUID parsed: {}", uuid);
        // If account id from JWT Token doesn't match with account id provide >> throw Account Not Match exception
        if (userAccountId != uuid) {
            throw new CustomExceptionHandler(ErrorCode.ACCOUNT_NOT_MATCH);
        }
        var foundObject = userAccountRepository
                .findById(uuid)
                // If not found throw exception 401 NOT FOUND
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        if (foundObject.getUsername() != updateAccountRequest.getUsername())
            foundObject.setUsername(foundObject.getUsername());
        if (!passwordEncodingService
                .getBCryptPasswordEncoder()
                .matches(updateAccountRequest.getPassword(), foundObject.getPassword())) {
            foundObject.setPassword(updateAccountRequest.getPassword());
        }
        if (foundObject.getEmail() != updateAccountRequest.getEmail()) {
            foundObject.setEmail(updateAccountRequest.getEmail());
        }
        userAccountRepository.save(foundObject);
        return UpdateAccountResponse.builder()
                .email(foundObject.getEmail())
                .password(foundObject.getPassword())
                .username(foundObject.getUsername())
                .build();
    }

    @Transactional
    // Send email for activate request
    public Boolean requestActivateAccount(RequestActivationAccount requestActivationAccount) {
        // Check if account is activated throw an exception
        if (accountRepository.isActive(requestActivationAccount.getUsername()))
            throw new CustomExceptionHandler(ErrorCode.ACTIVATED);
        sendNotification(ActivateAccountNotificationRequest.builder()
                .email(requestActivationAccount.getEmail())
                .fullname(requestActivationAccount.getFullname())
                .code(generateActivationCode())
                .username(requestActivationAccount.getUsername())
                .build());
        return true;
    }

    public String generateActivationCode() {
        // Generate Activation Code
        return String.valueOf(new Random().nextInt(90000) + 10000);
    }

    public Token sendNotification(ActivateAccountNotificationRequest request) {
        log.info("Sending activation notification {}", request.getCode());
        // Trigger to Send Verified Code Request -> POST -x http://localhost:8082/noti/internal/create and get there
        // sponse.
        //        var response = notificationRepository.sendCreateEmail(CreateAccountNotificationRequest.builder()
        //                .fullname(request.getFullname())
        //                .username(request.getUsername())
        //                .email(request.getEmail())
        //                .code(request.getCode())
        //                .build());
        //        //  If status code is 500 log an error
        //        if (response.status() != 200) {
        //            log.error("Error while sending notification {}", response.body());
        //            throw new CustomExceptionHandler(ErrorCode.SEND_FAILED);
        //        }
        // Create a new ActivateModel and set the expired time and issue time;
        return saveActivationModel(request.getUsername(), request.getCode());
    }

    public Token saveActivationModel(String username, String activateCode) {
        log.info("Activate account {}", activateCode);
        var foundObject = userAccountRepository.findByUsername(username).orElse(null);
        if (foundObject == null) {
            return Token.builder()
                    .username(username)
                    .expires(Timestamp.from(Instant.now().plus(MINUTE_EXPIRED, ChronoUnit.MINUTES)))
                    .issuedAt(Timestamp.from(Instant.now()))
                    .token(activateCode)
                    .revoked(false)
                    .build();
        } else {
            Token token = foundObject.getToken();
            token.setToken(activateCode);
            token.setExpires(Timestamp.from(Instant.now().plus(MINUTE_EXPIRED, ChronoUnit.MINUTES)));
            return tokenRepository.save(token);
        }
    }
}
