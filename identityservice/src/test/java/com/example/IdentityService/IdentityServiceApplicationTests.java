package com.example.IdentityService;

import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

import dev.psyconnect.identity_service.model.UserAccount;
import dev.psyconnect.identity_service.repository.RoleRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import dev.psyconnect.identity_service.service.AuthenticationService;

@ExtendWith(MockitoExtension.class)
class AuthenticationServiceTest {

    @InjectMocks
    private AuthenticationService authenticationService;

    @Mock
    private UserAccountRepository userAccountRepository;

    @Mock
    private RoleRepository roleRepository;

    @Value("${SIGNER_KEY}")
    private String SIGNER_KEY = "my-signing-key"; // Mock giá trị

    private final PasswordEncoder passwordEncoder = new BCryptPasswordEncoder(10);

    private static final String TEST_USERNAME = "testUser";
    private static final String TEST_PASSWORD = "password123";
    private static final String TEST_ROLE = "USER";
    private static final String TEST_EMAIL = "test@gmail.com";
    private static final String TEST_LOGIN_TYPE = "default";

    private UserAccount mockUserAccount() {
        var user = new UserAccount();
        user.setUsername(TEST_USERNAME);
        user.setPassword(passwordEncoder.encode(TEST_PASSWORD));
        user.setEmail(TEST_EMAIL);
        return user;
    }
}
