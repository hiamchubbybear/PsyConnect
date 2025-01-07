package com.example.IdentityService.configuration;

import com.example.IdentityService.dto.request.AuthenticationRequest;
import com.example.IdentityService.dto.request.GoogleAuthenticationRequest;
import com.example.IdentityService.globalexceptionhandle.CustomExceptionHandler;
import com.example.IdentityService.globalexceptionhandle.ErrorCode;
import com.example.IdentityService.mapper.UserAccountMapper;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.repository.UserAccountRepository;
import com.example.IdentityService.service.AuthenticationService;
import com.example.IdentityService.service.OAuth2Service;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.http.ResponseEntity;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configurers.AbstractHttpConfigurer;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.oauth2.core.oidc.user.DefaultOidcUser;
import org.springframework.security.web.SecurityFilterChain;

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
    public SecurityFilterChain securityFilterChain(
            HttpSecurity http
    ) throws Exception {
        return http
                .authorizeRequests(requests -> {
                    // Cấu hình các endpoint được phép truy cập
                    requests.requestMatchers("/", "/oauth2/register/google", "/oauth2/userInfo/google", "/login", "/oauth2/authorization/google", "/identity/create").permitAll();
                    requests.anyRequest().authenticated();
                })
                .oauth2Login(oAuth2Login -> {
                    // Xử lý đăng nhập thành công
                    oAuth2Login.successHandler((HttpServletRequest request, HttpServletResponse response, Authentication authentication) -> {
                        try {
                            DefaultOidcUser user = (DefaultOidcUser) authentication.getPrincipal();
                            String email = user.getAttribute("email");

                            if (email == null) {
                                log.error("Email attribute not found in OIDC user information");
                                response.sendError(HttpServletResponse.SC_BAD_REQUEST, "Email not found in user information");
                                return;
                            }
                            log.info("OAuth2 Login successful. Email: {}", email);
                            // Logic after Authentication success
                            oAuth2Service.processOAuthPostLoginGoogle(email);
                            // Get url from session
                            String redirectUrl = (String) request.getSession()
                                    .getAttribute("SPRING_SECURITY_SAVED_REQUEST");
                            // If session not found, get the default
                            if (redirectUrl == null || redirectUrl.isEmpty()) {
                                redirectUrl = "/oauth2/userInfo/google";
                            }
                            response.sendRedirect(redirectUrl);
                        } catch (Exception ex) {
                            log.error("Error handling successful authentication", ex);
                            response.sendError(HttpServletResponse.SC_INTERNAL_SERVER_ERROR, "Internal Server Error");
                        }});
                    oAuth2Login.failureHandler((HttpServletRequest request, HttpServletResponse response, AuthenticationException exception) -> {
                        log.error("OAuth2 login failed: {}", exception.getMessage());
                        response.setContentType("application/json");
                        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
                        response.getWriter().write("{\"error\": \"Authentication failed: " + exception.getMessage() + "\"}");
                    });

                    oAuth2Login.authorizationEndpoint(authorizationEndpointConfig ->
                                    authorizationEndpointConfig.baseUri("/oauth2/authorization/google"))
                            .redirectionEndpoint(redirectionEndpointConfig ->
                                    redirectionEndpointConfig.baseUri("/oauth2/callback/google"));
                })
                .formLogin(formLogin -> {
                    // Cấu hình form login
                    formLogin.loginPage("/login")
                            .defaultSuccessUrl("/home") // Trang đích sau khi login thành công
                            .failureUrl("/login?error=true"); // Trang lỗi khi login thất bại
                })
                .httpBasic(Customizer.withDefaults())
                .csrf(AbstractHttpConfigurer::disable) // Disable CSRF (chỉ khi dùng JWT)
                .build();
    }

}