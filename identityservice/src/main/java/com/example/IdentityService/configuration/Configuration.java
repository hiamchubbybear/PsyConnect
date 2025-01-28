package com.example.IdentityService.configuration;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.web.SecurityFilterChain;

import com.example.IdentityService.dto.request.CreateProfileOauth2GoogleRequest;
import com.example.IdentityService.enumeration.Provider;
import com.example.IdentityService.service.OAuth2Service;

import lombok.AccessLevel;
import lombok.experimental.FieldDefaults;

@org.springframework.context.annotation.Configuration
@EnableWebSecurity
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class Configuration {
    private static final Logger log = LoggerFactory.getLogger(Configuration.class);
    OAuth2Service oAuth2Service;

    @Autowired
    public Configuration(OAuth2Service oAuth2Service) {
        this.oAuth2Service = oAuth2Service;
    }

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        return http.authorizeRequests(requests -> {
                    // Config Allowed Endpoint
                    requests.requestMatchers(
                                    "/",
                                    "/oauth2/register/google",
                                    "/oauth2/userInfo/google",
                                    "/login",
                                    "/oauth2/authorization/google",
                                    "/identity/create",
                                    "auth/login")
                            .permitAll();
                    requests.anyRequest().authenticated();
                })
                .oauth2Login(oAuth2Login -> {
                    // Process When Process Successful
                    oAuth2Login.successHandler(
                            (HttpServletRequest request,
                                    HttpServletResponse response,
                                    Authentication authentication) -> {
                                try {
                                    DefaultOidcUser user = (DefaultOidcUser) authentication.getPrincipal();
                                    String email = user.getAttribute("email");
                                    String avatarUri = user.getAttribute("picture");
                                    if (email == null) {
                                        // Logging
                                        log.error("Email attribute not found in OIDC user information");
                                        response.sendError(
                                                HttpServletResponse.SC_BAD_REQUEST,
                                                "Email not found in user information");
                                        return;
                                    }
                                    log.info("OAuth2 Login successful. Email: {}", email);
                                    // Logic after Authentication success
                                    var trans = oAuth2Service.processOAuthPostLoginGoogle(
                                            new CreateProfileOauth2GoogleRequest(
                                                    email, Provider.GOOGLE.toString(), avatarUri));
                                    // Logging
                                    log.info("OAuth2 Create account successful. Email: {}", trans.getEmail());
                                    // Get url from session
                                    String redirectUrl =
                                            (String) request.getSession().getAttribute("SPRING_SECURITY_SAVED_REQUEST");
                                    // If session not found, get the default
                                    if (redirectUrl == null || redirectUrl.isEmpty()) {
                                        redirectUrl = "/oauth2/userInfo/google";
                                    }
                                    response.sendRedirect(redirectUrl);
                                } catch (Exception ex) {
                                    log.error("Error handling successful authentication", ex);
                                    response.sendError(
                                            HttpServletResponse.SC_INTERNAL_SERVER_ERROR, "Internal Server Error");
                                }
                            });
                    oAuth2Login.failureHandler(
                            (HttpServletRequest request,
                                    HttpServletResponse response,
                                    AuthenticationException exception) -> {
                                // Logging
                                log.error("OAuth2 login failed: {}", exception.getMessage());
                                // If failed response a error
                                response.setContentType("application/json");
                                response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
                                response.getWriter()
                                        .write("{\"error\": \"Authentication failed: " + exception.getMessage()
                                                + "\"}");
                            });
                    // Authentication uri
                    oAuth2Login
                            .authorizationEndpoint(authorizationEndpointConfig ->
                                    authorizationEndpointConfig.baseUri("/oauth2/authorization/google"))
                            // Redirect endpoint after failed
                            .redirectionEndpoint(redirectionEndpointConfig ->
                                    redirectionEndpointConfig.baseUri("/oauth2/callback/google"));
                })
                .formLogin(formLogin -> {
                    // Config Login Page
                    formLogin
                            .loginPage("/login")
                            .defaultSuccessUrl("/home") // Navigate after create success
                            .failureUrl("/login"); // Error page after create failed
                })
                .httpBasic(Customizer.withDefaults())
                .csrf(AbstractHttpConfigurer::disable) // Disable CSRF (Only for JWT)
                .build();
    }
}
