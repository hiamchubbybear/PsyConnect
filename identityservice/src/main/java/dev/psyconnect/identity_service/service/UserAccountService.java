package dev.psyconnect.identity_service.service;

import java.sql.Timestamp;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.*;

import jakarta.transaction.Transactional;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.stereotype.Service;

import dev.psyconnect.identity_service.dto.LogEvent;
import dev.psyconnect.identity_service.dto.LogLevel;
import dev.psyconnect.identity_service.dto.request.*;
import dev.psyconnect.identity_service.dto.response.*;
import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.grpc.client.ProfileGRPCClient;
import dev.psyconnect.identity_service.interfaces.IUserAccountService;
import dev.psyconnect.identity_service.kafka.producer.KafkaService;
import dev.psyconnect.identity_service.mapper.UserAccountMapper;
import dev.psyconnect.identity_service.model.*;
import dev.psyconnect.identity_service.repository.ActivateRepository;
import dev.psyconnect.identity_service.repository.RoleRepository;
import dev.psyconnect.identity_service.repository.TokenRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import lombok.extern.slf4j.Slf4j;

@Service
@Slf4j
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE)
public class UserAccountService implements UserDetailsService, IUserAccountService {
    public int MINUTE_EXPIRED = 10;

    @Value("${NOTIFICATION_CREATE_TOPIC}")
    String NOTIFICATION_CREATE_TOPIC;

    static String NOTIFICATION_RESET_TOPIC = "notification.user-reset";
    final KafkaService kafkaService;
    final RoleRepository roleRepository;
    final UserAccountRepository userAccountRepository;
    final UserAccountMapper mapper;
    final UserAccountRepository accountRepository;
    final ActivateRepository activateRepository;
    final TokenRepository tokenRepository;
    final ProfileGRPCClient profileGRPCClient;

    @Transactional
    @Cacheable(key = "#request.accountId", value = "account")
    public UserAccountCreationResponse createAccount(UserAccountCreationRequest request, Provider provider) {

        String encodedPassword = null, session = "";
        if (userAccountRepository.existsByUsername(request.getUsername())) {
            throw new CustomExceptionHandler(ErrorCode.USERNAME_ALREADY_EXISTS);
        }

        if (userAccountRepository.existsByEmail(request.getEmail())) {
            throw new CustomExceptionHandler(ErrorCode.EMAIL_ALREADY_EXISTS);
        }
        if (provider != Provider.ORDINARY) {
            request.setUsername(request.getEmail());
            session = request.getOauth2Session();
        } else {
            encodedPassword = PasswordEncodingService.encoder(request.getPassword());
        }

        String activationCode = generateActivationCode();
        Token activateToken = saveActivationModel(request.getUsername(), activationCode, provider);

        Set<RoleEntity> roles = roleRepository.findAllByRoleId(request.getRole().toUpperCase());
        UUID profileId = UUID.randomUUID();
        Account account = Account.builder()
                .username(request.getUsername())
                .isActivated(false)
                .createdAt(Timestamp.from(Instant.now()))
                .provider(provider)
                .email(request.getEmail())
                .password(encodedPassword)
                .isActivated((provider != Provider.ORDINARY))
                .role(roles)
                .session(session)
                .token(activateToken)
                .profileId(profileId)
                .build();

        var savedAccount = accountRepository.save(account);

        UserProfileCreationRequest profileRequest = mapper.toAccountResponse(request);
        profileRequest.setAccountId(savedAccount.getAccountId().toString());
        profileRequest.setProfileId(profileId.toString());
        profileRequest.setRole(request.getRole());
        try {
            var profileResponse = profileGRPCClient.createProfile(profileRequest);
            if (!profileResponse
                    .getProfileId()
                    .equals(savedAccount.getProfileId().toString())) {
                throw new CustomExceptionHandler(ErrorCode.RUNTIME_ERROR);
            }
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "identity-service",
                    request.getUsername(),
                    "Create account",
                    e.getMessage(),
                    Map.of("error", e.getMessage()),
                    LogLevel.ERROR));
            log.info(e.getMessage());
            throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
        }
        String fullName = request.getFirstName() + " " + request.getLastName();
        if (provider == Provider.ORDINARY) {
            kafkaService.send(
                    NOTIFICATION_CREATE_TOPIC,
                    CreateAccountNotificationRequest.builder()
                            .code(activationCode)
                            .username(request.getUsername())
                            .email(savedAccount.getEmail())
                            .fullname(fullName)
                            .build());
        }
        kafkaService.sendLog(buildLog(
                "identity-service",
                request.getUsername(),
                "Create account",
                "Success",
                Map.of("fullname", fullName, "email", savedAccount.getEmail()),
                LogLevel.LOG));

        return UserAccountCreationResponse.builder()
                .username(request.getUsername())
                .email(account.getEmail())
                .imageUri(request.getAvatarUri())
                .role(roles)
                .build();
    }

    @Transactional
    @CachePut(key = "#request.accountId", value = "account")
    public boolean deleteAccountWithoutCheck(UUID accountId) {
        try {
            userAccountRepository.deleteById(accountId);
            kafkaService.sendLog(buildLog(
                    "identity-service",
                    accountId.toString(),
                    "Delete account (force)",
                    "Success",
                    Map.of("status", true),
                    LogLevel.AUDIT));
            return true;
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "identity-service",
                    accountId.toString(),
                    "Delete account (force)",
                    "Failed",
                    Map.of("error", e.getMessage()),
                    LogLevel.ERROR));
            throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
        }
    }

    @CachePut(key = "uuid", value = "account")
    public DeleteAccountResponse deleteAccount(DeleteAccountRequest deleteAccountRequest, UUID uuid) {
        Account userObject = userAccountRepository
                .findById(uuid)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));

        if (!userObject.getUsername().equals(deleteAccountRequest.getUsername())
                || PasswordEncodingService.getBCryptPasswordEncoder()
                        .matches(userObject.getPassword(), deleteAccountRequest.getPassword()))
            throw new CustomExceptionHandler(ErrorCode.DELETE_ACCOUNT_FAILED);
        if (!userObject.getToken().getToken().equals(deleteAccountRequest.getToken()))
            throw new CustomExceptionHandler(ErrorCode.TOKEN_INVALID);

        kafkaService.sendLog(buildLog(
                "identity-service",
                userObject.getUsername(),
                "Delete account",
                "Requested",
                Map.of("email", deleteAccountRequest.getEmail()),
                LogLevel.LOG));

        return DeleteAccountResponse.builder()
                .isSuccess(true)
                .secret(SecretEnum.getRandomMessage())
                .session(userObject.getSession())
                .timestamp(deleteAccountRequest.getTimestamp())
                .build();
    }

    public boolean deleteAccountRequest(DeleteAccountConfirmRequest removeAccountRequest, UUID uuid) {
        Account userObject = userAccountRepository
                .findById(uuid)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        if (!removeAccountRequest.getSession().equals(userObject.getSession())
                || !userObject.getToken().getToken().equals(removeAccountRequest.getToken()))
            throw new CustomExceptionHandler(ErrorCode.DELETE_ACCOUNT_FAILED);
        try {
            userAccountRepository.deleteById(uuid);
            kafkaService.sendLog(buildLog(
                    "identity-service",
                    userObject.getUsername(),
                    "Delete account",
                    "Success",
                    Map.of("metadata", true),
                    LogLevel.LOG));
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "identity-service",
                    userObject.getUsername(),
                    "Delete account",
                    "Failed",
                    Map.of("error", e.getMessage()),
                    LogLevel.ERROR));
            throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
        }
        return true;
    }

    @CachePut(key = "#result.accountId", value = "account")
    public UserInfoResponse getUserAccount(UUID accountId) {
        var userResponse = userAccountRepository
                .findById((accountId))
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        if (userAccountRepository.existsByUsernameAndIsActivatedTrue(userResponse.getUsername())) {
            kafkaService.sendLog(buildLog(
                    "identity-service",
                    userResponse.getUsername(),
                    "Get user account",
                    "Success",
                    Map.of("status", true),
                    LogLevel.LOG));
            return UserInfoResponse.builder()
                    .username(userResponse.getUsername())
                    .accountId(userResponse.getAccountId())
                    .createdAt(userResponse.getCreatedAt().toLocalDateTime().toLocalDate())
                    .role(userResponse.getRole())
                    .email(userResponse.getEmail())
                    .isActivated(userResponse.isActivated())
                    .provider(userResponse.getProvider())
                    .build();
        } else throw new CustomExceptionHandler(ErrorCode.ACCOUNT_INACTIVE);
    }

    @Override
    public AuthenticationFilterRequest loadUserByUsername(String username) throws CustomExceptionHandler {
        var foundObject = userAccountRepository
                .findByUsername(username)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        return AuthenticationFilterRequest.builder()
                .username(foundObject.getUsername())
                .password(foundObject.getPassword())
                .build();
    }

    @Override
    public Page<Account> getAllAccount(int page) {
        Pageable pageable = PageRequest.of(page, 10);
        return userAccountRepository.findAll(pageable);
    }

    @Transactional
    @CacheEvict(key = "#accountId", value = "account")
    public UpdateAccountResponse updateAccount(UpdateAccountRequest updateAccountRequest, UUID accountId) {
        Account foundObject = userAccountRepository
                .findById(accountId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        boolean updated = false;
        if (!Objects.equals(foundObject.getEmail(), updateAccountRequest.getEmail())) {
            boolean emailExists = userAccountRepository.existsByEmail(updateAccountRequest.getEmail());
            if (emailExists) throw new CustomExceptionHandler(ErrorCode.EMAIL_ALREADY_EXISTS);
            foundObject.setEmail(updateAccountRequest.getEmail());
            updated = true;
        }
        if (updated) {
            userAccountRepository.save(foundObject);
            kafkaService.send(
                    "notification.account-change",
                    UpdateAccountNotificatorRequest.builder()
                            .email(foundObject.getEmail())
                            .username(foundObject.getUsername())
                            .build());
            kafkaService.sendLog(buildLog(
                    "identity-service",
                    foundObject.getUsername(),
                    "Update account",
                    "Success",
                    Map.of("newEmail", foundObject.getEmail()),
                    LogLevel.LOG));
        }
        return UpdateAccountResponse.builder()
                .email(foundObject.getEmail())
                .username(foundObject.getUsername())
                .build();
    }

    @Transactional
    @CachePut(key = "#request.accountId", value = "account")
    public ActivateAccountResponse activateAccount(ActivateAccountRequest request) {
        String email = request.getEmail();
        String activateToken = request.getToken();

        if (userAccountRepository.existsByEmailAndIsActivatedTrue(email)) {
            throw new CustomExceptionHandler(ErrorCode.ACTIVATED);
        }
        if (!userAccountRepository.existsByEmail(email)) {
            throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        }
        var foundObject = activateRepository
                .findByToken(activateToken)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.ACTIVATION_FAILED));

        if (foundObject.getExpires().before(new Date())) {
            kafkaService.sendLog(buildLog(
                    "identity-service",
                    request.getEmail(),
                    "Activate account",
                    "Failed",
                    Map.of("reason", "Token expired"),
                    LogLevel.WARN));

            return ActivateAccountResponse.builder()
                    .isSuccess(false)
                    .message("Activation token has expired")
                    .build();
        }
        userAccountRepository.activateUser(email);
        foundObject.setToken("");
        foundObject.setExpires(null);
        activateRepository.save(foundObject);
        kafkaService.sendLog(buildLog(
                "identity-service",
                request.getEmail(),
                "Activate account",
                "Success",
                Map.of("status", true),
                LogLevel.AUDIT));
        return ActivateAccountResponse.builder()
                .isSuccess(true)
                .message("Account activation successful")
                .build();
    }

    public Boolean requestActivateAccount(RequestActivationAccount requestActivationAccount) {
        if (accountRepository.existsByUsernameAndIsActivatedTrue(requestActivationAccount.getUsername()))
            throw new CustomExceptionHandler(ErrorCode.ACTIVATED);
        ActivateAccountNotificationRequest req = ActivateAccountNotificationRequest.builder()
                .email(requestActivationAccount.getEmail())
                .fullname(requestActivationAccount.getFullname())
                .code(generateActivationCode())
                .username(requestActivationAccount.getUsername())
                .build();
        kafkaService.send((NOTIFICATION_CREATE_TOPIC), req);
        sendNotification(
                ActivateAccountNotificationRequest.builder()
                        .username(requestActivationAccount.getUsername())
                        .email(requestActivationAccount.getEmail())
                        .build(),
                Provider.ORDINARY);
        kafkaService.sendLog(buildLog(
                "identity-service",
                requestActivationAccount.getUsername(),
                "Request account activation",
                "Success",
                Map.of("metadata", req),
                LogLevel.LOG));
        return true;
    }

    public Boolean requestPasswordReset(RequestPasswordReset requestActivationAccount) {
        var userFound = accountRepository
                .findByEmail(requestActivationAccount.getEmail())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        ActivateAccountNotificationRequest req = ActivateAccountNotificationRequest.builder()
                .email(userFound.getEmail())
                .code(generateActivationCode())
                .username(userFound.getUsername())
                .build();
        log.info("Request found {}", req);
        kafkaService.send((NOTIFICATION_RESET_TOPIC), req);
        sendNotification(
                ActivateAccountNotificationRequest.builder()
                        .username(userFound.getUsername())
                        .email(requestActivationAccount.getEmail())
                        .code(req.getCode())
                        .build(),
                Provider.ORDINARY);
        kafkaService.sendLog(buildLog(
                "identity-service",
                userFound.getUsername(),
                "Request account password reset",
                "Success",
                Map.of("metadata", req),
                LogLevel.LOG));
        return true;
    }

    private LogEvent buildLog(
            String service, String userId, String action, String message, Map<String, ?> metadata, LogLevel level) {
        return LogEvent.builder()
                .service(service)
                .level(level)
                .timestamp(Instant.now().toString())
                .userId(userId)
                .action(action)
                .message(message)
                .metadata(metadata)
                .build();
    }

    public String generateActivationCode() {
        return String.valueOf(new Random().nextInt(90000) + 10000);
    }

    public Token sendNotification(ActivateAccountNotificationRequest request, Provider provider) {
        return saveActivationModel(request.getUsername(), request.getCode(), provider);
    }

    public Token saveActivationModel(String username, String activationCode, Provider provider) {
        if (username == null || username.isBlank()) {
            throw new CustomExceptionHandler(ErrorCode.NULL_EXCEPTION);
        }
        Token token = tokenRepository
                .findById(username)
                .orElse(Token.builder()
                        .token(activationCode)
                        .username(username)
                        .issuedAt(Timestamp.from(Instant.now()))
                        .provider(provider.toString())
                        .revoked(false)
                        .build());
        log.info("Token saved {}", token);
        token.setToken(activationCode);
        token.setRevoked(false);
        token.setExpires(Timestamp.from(Instant.now().plus(MINUTE_EXPIRED, ChronoUnit.MINUTES)));
        return tokenRepository.save(token);
    }

    public Boolean resetPassword(PasswordResetRequest passwordResetRequest) {
        var userFound = userAccountRepository
                .findByEmail(passwordResetRequest.getEmail())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        var tokenFound = tokenRepository.findById(userFound.getUsername()).get().getToken();
        log.info("Token found!!!!!!! {}", tokenFound);
        if (userFound.getToken() == null || userFound.getToken().getToken() == null) {
            throw new CustomExceptionHandler(ErrorCode.TOKEN_NOT_FOUND);
        }
        if (userFound.getToken().getToken().equals(passwordResetRequest.getResetToken())) {
            userFound.setPassword(PasswordEncodingService.encoder(passwordResetRequest.getNewPassword()));
            Token userToken = userFound.getToken();
            userToken.setToken("");
            tokenRepository.save(userToken);
            userAccountRepository.save(userFound);
            return true;
        } else {
            throw new CustomExceptionHandler(ErrorCode.TOKEN_INVALID);
        }
    }
}
