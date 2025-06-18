package dev.psyconnect.identity_service.configuration;

import dev.psyconnect.identity_service.model.TokenRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.security.authentication.AuthenticationManager;
import org.springframework.security.authentication.AuthenticationProvider;
import org.springframework.security.authentication.dao.DaoAuthenticationProvider;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.authentication.configuration.AuthenticationConfiguration;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.config.http.SessionCreationPolicy;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.client.OAuth2AuthorizedClient;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

import dev.psyconnect.identity_service.dto.request.CreateProfileOauth2GoogleRequest;
import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.service.OAuth2Service;
import dev.psyconnect.identity_service.service.UserAccountService;
import lombok.AccessLevel;
import lombok.experimental.FieldDefaults;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;

@org.springframework.context.annotation.Configuration
@EnableWebSecurity
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class Configuration {
    private static final Logger log = LoggerFactory.getLogger(Configuration.class);
    OAuth2Service oAuth2Service;
    JwtAuthFilter authFilter;
    UserAccountService userAccountService;

    @Autowired
    public Configuration(OAuth2Service oAuth2Service, JwtAuthFilter authFilter, UserAccountService userAccountService) {
        this.oAuth2Service = oAuth2Service;
        this.authFilter = authFilter;
        this.userAccountService = userAccountService;
    }

    @Bean
    public UserDetailsService userDetailsService(UserAccountService userAccountService) {
        return userAccountService;
    }

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http, TokenRepository tokenRepository, UserAccountRepository userAccountRepository) throws Exception {
        return http
                .csrf(AbstractHttpConfigurer::disable)
                .sessionManagement(session -> session.sessionCreationPolicy(SessionCreationPolicy.IF_REQUIRED))
                .authenticationProvider(authenticationProvider(userAccountService))
                .addFilterBefore(authFilter, UsernamePasswordAuthenticationFilter.class)
                .authorizeRequests(requests -> requests
                        .requestMatchers(
                                "/login",
                                "/oauth2/authorization/google",
                                "/oauth2/callback/google",
                                "/login/oauth2/code/google",
                                "/identity/**",
                                "/identity/create",
                                "/auth/login/**",
                                "/oauth2/userInfo/google",
                                "/favicon.ico"
                        ).permitAll()
                        .requestMatchers("/auth/therapist/**").hasAuthority("ROLE_THERAPIST")
                        .requestMatchers("/auth/admin/**").hasAuthority("ROLE_ADMIN")
                        .anyRequest().authenticated()
                )
                .exceptionHandling(exception -> exception.authenticationEntryPoint(((request, response, authException) -> {
                    response.sendRedirect("/oauth2/authorization/google");
                })))
                .oauth2Login(oauth2 -> oauth2
                        .loginPage("/oauth2/authorization")
                        .authorizationEndpoint(config -> config.baseUri("/oauth2/authorization"))
                        .redirectionEndpoint(config -> config.baseUri("/oauth2/callback/*"))
                        .successHandler((request, response, authentication) -> {
                            try {
                                DefaultOidcUser user = (DefaultOidcUser) authentication.getPrincipal();
                                String email = user.getAttribute("email");
                                String avatarUri = user.getAttribute("picture");
                                if (email == null) {
                                    log.error("Email attribute not found");
                                    response.sendError(HttpServletResponse.SC_BAD_REQUEST, "Email not found");
                                    return;
                                }
                                oAuth2Service.processOAuthPostLoginGoogle(email, avatarUri,authentication);
                                String redirectUrl = String.format("/oauth2/userInfo/google?email=%s&avatar=%s",
                                        URLEncoder.encode(email, StandardCharsets.UTF_8),
                                        URLEncoder.encode(avatarUri, StandardCharsets.UTF_8));
                                response.sendRedirect(redirectUrl);


                            } catch (Exception e) {
                                log.error("OAuth2 success handler error", e);
                                response.sendError(HttpServletResponse.SC_INTERNAL_SERVER_ERROR);
                            }
                        })
                        .failureHandler((request, response, exception) -> {
                            log.error("OAuth2 login failed: {}", exception.getMessage());
                            response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
                            response.getWriter().write("{\"error\": \"OAuth2 login failed\"}");
                        })
                )
                .formLogin(formLogin -> formLogin
                        .loginPage("/login")
                        .defaultSuccessUrl("/home")
                        .failureUrl("/login")
                )
                .httpBasic(Customizer.withDefaults())
                .build();
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Bean
    public AuthenticationProvider authenticationProvider(UserAccountService userAccountService) {
        DaoAuthenticationProvider authenticationProvider = new DaoAuthenticationProvider();
        authenticationProvider.setUserDetailsService(userDetailsService(userAccountService));
        authenticationProvider.setPasswordEncoder(passwordEncoder());
        return authenticationProvider;
    }

    @Bean
    public AuthenticationManager authenticationManager(AuthenticationConfiguration config) throws Exception {
        return config.getAuthenticationManager();
    }
}
