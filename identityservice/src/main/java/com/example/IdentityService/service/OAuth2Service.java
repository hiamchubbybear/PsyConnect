package com.example.IdentityService.service;

import com.example.IdentityService.dto.request.UserProfileCreationRequest;
import com.example.IdentityService.dto.response.UserAccountCreationResponse;
import com.example.IdentityService.mapper.UserAccountMapper;
import com.example.IdentityService.model.Provider;
import com.example.IdentityService.model.RoleEntity;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.repository.RoleRepository;
import com.example.IdentityService.repository.UserAccountRepository;
import com.example.IdentityService.repository.https.ProfileRepository;
import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.format.annotation.DateTimeFormat;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.stereotype.Service;

import java.security.Principal;
import java.time.LocalDate;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class OAuth2Service {
    private static final Logger log = LoggerFactory.getLogger(OAuth2Service.class);
    private final UserAccountRepository userAccountRepository;
    private final RoleRepository roleRepository;
    private final ProfileRepository profileRepository;
    private final UserAccountMapper userAccountMapper;

    public UserAccountCreationResponse createAccountWithGoogle(Principal principal) {
        Authentication authentication = (Authentication) principal;
        Map<String, Object> attributes = ((OAuth2AuthenticationToken) authentication).getPrincipal().getAttributes();
        if (!Boolean.TRUE.equals(attributes.get("email_verified"))) {
            throw new BadCredentialsException("Email verification failed");
        }
        String email = attributes.get("email").toString();
        String username = email.substring(0, email.indexOf("@")).toLowerCase();
        String password = null;
        Set<RoleEntity> roles = roleRepository.findAllByName("Client");
        UserAccount savedUser = userAccountRepository.save(
                new UserAccount(
                        UUID.randomUUID().toString(),
                        username,
                        password,
                        email,
                        Provider.GOOGLE,
                        roles
                )
        );
        log.info("Created user account with email: " + savedUser.getEmail());
        String firstName = attributes.getOrDefault("given_name", "").toString();
        String lastName = attributes.getOrDefault("family_name", "").toString();
        String address = attributes.getOrDefault("address", "").toString();
        String gender = attributes.getOrDefault("gender", "").toString();
        LocalDate dob = null;
        if (attributes.containsKey("dob")) {
            try {
                dob = LocalDate.parse(attributes.get("dob").toString());
            } catch (Exception e) {
                log.warn("Invalid date of birth format from Google OAuth2: " + attributes.get("dob"));
            }
        }
        UserProfileCreationRequest userProfileCreationRequest = UserProfileCreationRequest.builder()
                .username(username)
                .address(address)
                .dob(dob)
                .gender(gender)
                .firstName(firstName)
                .lastName(lastName)
                .build();
        log.info("Created user profile creation request: " + userProfileCreationRequest);
        profileRepository.createProfile(userProfileCreationRequest);
        UserAccountCreationResponse response = userAccountMapper.toUserProfileCreationResponse(savedUser);
        response.setRole("Client");
        return response;
    }
}
