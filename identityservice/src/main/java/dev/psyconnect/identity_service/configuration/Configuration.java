package dev.psyconnect.identity_service.configuration;

import jakarta.servlet.http.HttpServletResponse;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
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
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.client.authentication.OAuth2AuthenticationToken;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.UsernamePasswordAuthenticationFilter;

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
    private final String oauth2RedirectBase;

    @Autowired
    public Configuration(
            OAuth2Service oAuth2Service,
            JwtAuthFilter authFilter,
            UserAccountService userAccountService,
            @Value("${app.oauth2.redirect-base}") String oauth2RedirectBase) {
        this.oAuth2Service = oAuth2Service;
        this.authFilter = authFilter;
        this.userAccountService = userAccountService;
        this.oauth2RedirectBase = oauth2RedirectBase;
    }

    @Bean
    public UserDetailsService userDetailsService(UserAccountService userAccountService) {
        return userAccountService;
    }

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        return http.csrf(AbstractHttpConfigurer::disable)
                .sessionManagement(session -> session.sessionCreationPolicy(SessionCreationPolicy.IF_REQUIRED))
                .authenticationProvider(authenticationProvider(userAccountService))
                .addFilterBefore(authFilter, UsernamePasswordAuthenticationFilter.class)
                .authorizeRequests(requests -> requests.requestMatchers(
                        "/login",
                        "/oauth2/authorization/google",
                        "/oauth2/callback/google",
                        "/login/oauth2/code/google",
                        "/identity/**",
                        "/identity/create",
                        "/auth/login/**",
                        "/oauth2/userInfo/google",
                        "/oauth2/authorization/facebook",
                        "/auth/oauth2/callback/exchange-code",
                        "/oauth2/callback/facebook",
                        "/login/oauth2/code/facebook",
                        "/auth/internal/valid",
                        "/oauth2/userInfo/facebook",
                        "/favicon.ico")
                        .permitAll()
                        .requestMatchers("/auth/therapist/**")
                        .hasAuthority("ROLE_THERAPIST")
                        .requestMatchers("/auth/admin/**")
                        .hasAuthority("ROLE_ADMIN")
                        .anyRequest()
                        .authenticated())
                .exceptionHandling(
                        exception -> exception.authenticationEntryPoint(((request, response, authException) -> {
                            response.sendRedirect("/oauth2/authorization/google");
                        })))
                .oauth2Login(oauth2 -> oauth2.loginPage("/oauth2/authorization")
                        .authorizationEndpoint(config -> config.baseUri("/oauth2/authorization"))
                        .redirectionEndpoint(config -> config.baseUri("/oauth2/callback/*"))
                        .successHandler((request, response, authentication) -> {
                            try {
                                OAuth2AuthenticationToken token = (OAuth2AuthenticationToken) authentication;
                                String registrationId = token.getAuthorizedClientRegistrationId();
                                String email = "";
                                String avatarUri = "";
                                String redirectUrl = "";
                                if ("google".equals(registrationId)) {
                                    DefaultOidcUser user = (DefaultOidcUser) authentication.getPrincipal();
                                    log.info(((OAuth2User) authentication.getPrincipal())
                                            .getAttributes()
                                            .toString());
                                    email = user.getAttribute("email");
                                    avatarUri = user.getAttribute("picture");
                                    String platform = request.getParameter("platform");
                                    redirectUrl = String.format(
                                            "%s/oauth2/userInfo?provider=%s&email=%s&avatar=%s&platform=%s",
                                            oauth2RedirectBase, registrationId, email, avatarUri, platform);
                                } else if ("facebook".equals(registrationId)) {
                                    OAuth2User user = (OAuth2User) authentication.getPrincipal();
                                    log.info(((OAuth2User) authentication.getPrincipal())
                                            .getAttributes()
                                            .toString());
                                    avatarUri = (String) user.getAttribute("picture");
                                    email = (String) user.getAttribute("email");
                                    String platform = request.getParameter("platform");
                                    redirectUrl = String.format(
                                            "%s/oauth2/userInfo?provider=%s&email=%s&avatar=%s&platform=%s",
                                            oauth2RedirectBase, registrationId, email, avatarUri, platform);
                                }
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
                        }))
                .formLogin(formLogin -> formLogin.loginPage("/login").defaultSuccessUrl("/home").failureUrl("/login"))
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
