package dev.psyconnect.identity_service.service;

import java.util.Collections;
import java.util.Map;

import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserService;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.user.DefaultOAuth2User;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;

@Service
public class CustomerOauth2UserService implements OAuth2UserService<OAuth2UserRequest, OAuth2User> {
    @Override
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        OAuth2UserService<OAuth2UserRequest, OAuth2User> delegate = new DefaultOAuth2UserService();
        OAuth2User oAuth2User = delegate.loadUser(userRequest);
        String registrationId = userRequest.getClientRegistration().getRegistrationId();
        Map<String, Object> attributes = oAuth2User.getAttributes();
        String email = (String) attributes.get("email");
        String picture = null;
        String firstName = null;
        String lastName = null;

        if ("google".equals(registrationId)) {
            picture = (String) attributes.get("picture");
            firstName = (String) attributes.get("given_name");
            lastName = (String) attributes.get("family_name");
        } else if ("facebook".equals(registrationId)) {
            picture = ((Map<String, Object>) ((Map<String, Object>) attributes.get("picture")).get("data"))
                    .get("url")
                    .toString();
            String name = (String) attributes.get("name");
            firstName = name.split(" ")[0];
            lastName = name.contains(" ") ? name.substring(name.indexOf(" ") + 1) : "";
        }
        return new DefaultOAuth2User(
                Collections.singleton(new SimpleGrantedAuthority("ROLE_USER")),
                Map.of(
                        "email", email,
                        "picture", picture,
                        "firstName", firstName,
                        "lastName", lastName,
                        "provider", registrationId),
                "email");
    }
}
