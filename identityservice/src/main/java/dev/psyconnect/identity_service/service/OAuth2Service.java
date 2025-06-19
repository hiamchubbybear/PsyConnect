package dev.psyconnect.identity_service.service;

import java.util.Set;

import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import dev.psyconnect.identity_service.configuration.CallRestApi;
import dev.psyconnect.identity_service.dto.request.*;
import dev.psyconnect.identity_service.dto.response.AuthenticationResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.Authentication;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClientService;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;

import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
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
    OAuth2AuthorizedClientService authorizedClientService;
    CallRestApi callRestApi;
    RoleRepository roleRepository;
    private final AuthenticationService authenticationService;
    private final UserAccountService userAccountService;

    public AuthenticationResponse processOAuth2PostLogin(String email, String avatarUri, Authentication authentication, String loginProvider) {
        if (!userAccountRepository.existsByEmail(email)) {
            CreateProfileOauth2Request createProfileOauth2Request = null;
            String firstName = "", lastName = "", dob = "", gender = "";
            Provider provider = null;
            if (loginProvider.toUpperCase().equals(Provider.GOOGLE.toString())) {
                createProfileOauth2Request = extractDataFromJson(authentication);
                firstName = createProfileOauth2Request.getFirstName();
                lastName = createProfileOauth2Request.getLastName();
                dob = createProfileOauth2Request.getDob();
                gender = createProfileOauth2Request.getGender();
                provider = Provider.GOOGLE;
            } else if (loginProvider.toUpperCase().equals(Provider.FACEBOOK.toString())) {
                OAuth2User user = (OAuth2User) authentication.getPrincipal();
                firstName = user.getAttribute("firstName");
                lastName = user.getAttribute("lastName");
                provider = Provider.FACEBOOK;

            }
            try {
                Account existingUser = null;
                Set<RoleEntity> clientRoles = roleRepository.findAllByRoleId("Client");

                if (clientRoles.isEmpty()) {
                    log.error("Client role not found in database");
                    throw new CustomExceptionHandler(ErrorCode.ROLE_NOT_FOUND);
                }
                if (existingUser == null) {
                    userAccountService.createAccount(UserAccountCreationRequest.builder()
                            .email(email)
                            .dob(dob)
                            .gender(gender)
                            .firstName(firstName)
                            .lastName(lastName)
                            .role("Client")

                            .address("")
                            .avatarUri(avatarUri)
                            .build(), provider);
                }
                return authenticationService.generateOAuth2LoginToken(GoogleAuthenticationRequest.builder().email(email).build(), Provider.GOOGLE.toString());
            } catch (Exception ex) {
                log.error("Error during OAuth2 Google login: {}", ex.getMessage(), ex);
                throw new CustomExceptionHandler(ErrorCode.UNCATEGORIZED_EXCEPTION);
            }
        }
        return authenticationService.generateOAuth2LoginToken(GoogleAuthenticationRequest.builder().email(email).build(), Provider.GOOGLE.toString());
    }

    public CreateProfileOauth2Request extractDataFromJson(Authentication authentication) {
        DefaultOidcUser user = (DefaultOidcUser) authentication.getPrincipal();
        String email = user.getAttribute("email");
        OAuth2AuthenticationToken oauthToken = (OAuth2AuthenticationToken) authentication;

        OAuth2AuthorizedClient authorizedClient = authorizedClientService.loadAuthorizedClient(oauthToken.getAuthorizedClientRegistrationId(), oauthToken.getName());

        String accessToken = authorizedClient.getAccessToken().getTokenValue();
        String googleApiResponse = callRestApi.callGooglePeopleApi(accessToken);
        JsonObject jsonObject = new JsonParser().parse(googleApiResponse).getAsJsonObject();
        log.info("json parsed : {}", jsonObject.asMap());
        JsonArray genders = jsonObject.getAsJsonArray("genders");
        String gender = genders.size() > 0 ? genders.get(0).getAsJsonObject().get("value").getAsString() : null;
        String firstName = null;
        String lastName = null;

        JsonArray names = jsonObject.getAsJsonArray("names");
        if (names != null && names.size() > 0) {
            JsonObject nameObj = names.get(0).getAsJsonObject();
            if (nameObj.has("givenName")) {
                firstName = nameObj.get("givenName").getAsString();
            }
            if (nameObj.has("familyName")) {
                lastName = nameObj.get("familyName").getAsString();
            }
        }
        if (firstName == null || lastName == null) {
            log.warn("Missing name from Google People API, using email as username fallback.");
            firstName = "Unknown";
            lastName = "";
        }
        JsonArray birthdays = jsonObject.getAsJsonArray("birthdays");
        JsonObject birthdayObj = birthdays.size() > 0 ? birthdays.get(0).getAsJsonObject().getAsJsonObject("date") : null;

        int year = birthdayObj != null && birthdayObj.has("year") ? birthdayObj.get("year").getAsInt() : 0;
        int month = birthdayObj != null && birthdayObj.has("month") ? birthdayObj.get("month").getAsInt() : 0;
        int day = birthdayObj != null && birthdayObj.has("day") ? birthdayObj.get("day").getAsInt() : 0;
        String dob = String.format("%04d-%02d-%02d", year, month, day);

        return CreateProfileOauth2Request.builder().dob(dob).gender(gender).firstName((firstName)).lastName(lastName).build();

    }
}
