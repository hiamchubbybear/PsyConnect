package dev.psyconnect.identity_service.configuration;

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
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

import dev.psyconnect.identity_service.dto.request.CreateProfileOauth2GoogleRequest;
import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.service.OAuth2Service;
import dev.psyconnect.identity_service.service.UserAccountService;
import lombok.AccessLevel;
import lombok.experimental.FieldDefaults;

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
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        return http.authorizeRequests(requests -> {
                    requests.requestMatchers(
                                    "/",
                                    "/oauth2/register/google",
                                    "/oauth2/userInfo/google",
                                    "/login",
                                    "/oauth2/authorization/google",
                                    "/identity/**",
                                    "/identity/create",
                                    "auth/login/**")
                            .permitAll()
                            .requestMatchers("/auth/therapist/**")
                            .hasAuthority("ROLE_THERAPIST")
                            .requestMatchers("/auth/admin/**")
                            .hasAuthority("ROLE_ADMIN")
                            .requestMatchers("/account/**")
                            .authenticated();
                })
                .sessionManagement(session -> session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
                .authenticationProvider(authenticationProvider(userAccountService))
                .addFilterBefore(authFilter, UsernamePasswordAuthenticationFilter.class)
                .oauth2Login(oAuth2Login -> {
                    oAuth2Login.loginPage("/custom-login").successHandler((request, response, authentication) -> {
                        try {
                            DefaultOidcUser user = (DefaultOidcUser) authentication.getPrincipal();
                            String email = user.getAttribute("email");
                            String avatarUri = user.getAttribute("picture");
                            if (email == null) {
                                log.error("Email attribute not found in OIDC user information");
                                response.sendError(
                                        HttpServletResponse.SC_BAD_REQUEST, "Email not found in user information");
                                return;
                            }
                            log.info("OAuth2 Login successful. Email: {}", email);
                            var trans = oAuth2Service.processOAuthPostLoginGoogle(
                                    new CreateProfileOauth2GoogleRequest(email, Provider.GOOGLE.toString(), avatarUri));
                            log.info("OAuth2 Create account successful. Email: {}", trans.getEmail());
                            String redirectUrl =
                                    (String) request.getSession().getAttribute("SPRING_SECURITY_SAVED_REQUEST");
                            if (redirectUrl == null || redirectUrl.isEmpty()) {
                                redirectUrl = "/oauth2/userInfo/google";
                            }
                            response.sendRedirect(redirectUrl);
                        } catch (Exception ex) {
                            log.error("Error handling successful authentication", ex);
                            response.sendError(HttpServletResponse.SC_INTERNAL_SERVER_ERROR, "Internal Server Error");
                        }
                    });
                    oAuth2Login.failureHandler(
                            (HttpServletRequest request,
                                    HttpServletResponse response,
                                    AuthenticationException exception) -> {
                                log.error("OAuth2 login failed: {}", exception.getMessage());
                                response.setContentType("application/json");
                                response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
                                response.getWriter()
                                        .write("{\"error\": \"Authentication failed: " + exception.getMessage()
                                                + "\"}");
                            });
                    oAuth2Login
                            .authorizationEndpoint(authorizationEndpointConfig ->
                                    authorizationEndpointConfig.baseUri("/oauth2/authorization/google"))
                            .redirectionEndpoint(redirectionEndpointConfig ->
                                    redirectionEndpointConfig.baseUri("/oauth2/callback/google"));
                })
                .formLogin(formLogin -> {
                    formLogin.loginPage("/login").defaultSuccessUrl("/home").failureUrl("/login");
                })
                .httpBasic(Customizer.withDefaults())
                .csrf(AbstractHttpConfigurer::disable)
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
