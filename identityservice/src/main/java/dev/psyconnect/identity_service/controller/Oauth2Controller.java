package dev.psyconnect.identity_service.controller;

import java.io.IOException;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;

import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import jakarta.servlet.http.HttpServletResponse;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.configuration.ValidateProviderType;
import dev.psyconnect.identity_service.dto.request.Oauth2AuthenticationRequest;
import dev.psyconnect.identity_service.dto.response.AuthenticationResponse;
import dev.psyconnect.identity_service.service.OAuth2Service;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@RestController
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class Oauth2Controller {
    private static final Logger log = LoggerFactory.getLogger(Oauth2Controller.class);
    private final OAuth2Service oAuth2Service;

    @GetMapping("/oauth2/userInfo")
    public void oAuth2Google(
            @RequestParam String provider,
            @RequestParam String email,
            @RequestParam String avatar,
            Authentication authentication,
            HttpServletResponse response)
            throws IOException {

        var res = oAuth2Service.processOAuth2PreLogin(email, avatar, authentication, provider);

        if (res.isSuccessful()) {
            String accessToken = res.getToken();

            if (accessToken != null && !accessToken.isBlank()
                    && email != null && !email.isBlank()
                    && provider != null && !provider.isBlank()) {

                String deepLinkUrl = String.format(
                        "psyconnect://oauth2/callback/code?code=%s&email=%s&provider=%s",
                        URLEncoder.encode(accessToken, StandardCharsets.UTF_8),
                        URLEncoder.encode(email, StandardCharsets.UTF_8),
                        URLEncoder.encode(provider, StandardCharsets.UTF_8)
                );

                String redirectHtml = String.format(
                        """
                                <!DOCTYPE html>
                                <html>
                                <head>
                                    <title>Redirecting to PsyConnect...</title>
                                    <meta charset="UTF-8">
                                </head>
                                <body>
                                    <div style="text-align: center; padding: 50px; font-family: Arial, sans-serif;">
                                        <h2>Login Successful!</h2>
                                        <p>Redirecting you back to PsyConnect app...</p>
                                        <p>If you're not redirected automatically, <a href="%s">click here</a></p>
                                    </div>
                                    <script>
                                        window.location.href = '%s';
                                        setTimeout(function() {
                                            window.close();
                                        }, 3000);
                                    </script>
                                </body>
                                </html>
                                """, deepLinkUrl, deepLinkUrl
                );

                response.setContentType("text/html; charset=UTF-8");
                response.getWriter().write(redirectHtml);
            } else {
                throw new CustomExceptionHandler(ErrorCode.NULL_EXCEPTION);
            }
        } else {
            response.sendRedirect("/oauth2/callback/error");
        }
    }


    @GetMapping("/oauth2/callback/error")
    public ApiResponse<Boolean> oAuth2CallBackFailed() {
        return new ApiResponse<>(400, "Failed", false);
    }

    @GetMapping("/oauth2/callback/code")
    public ApiResponse<Boolean> oAuth2CallBackSuccess() {
        return new ApiResponse<>(200, "Success", true);
    }

    @GetMapping("/auth/oauth2/callback/exchange-code")
    public ApiResponse<AuthenticationResponse> oAuth2PostSuccessExchangeCode(
            @RequestParam String code,
            @RequestParam String email,
            @ValidateProviderType @RequestParam String provider,
            @RequestParam String platform) {
        return new ApiResponse<>(oAuth2Service.processOAuth2Login(
                Oauth2AuthenticationRequest.builder()
                        .email(email)
                        .sessionToken(code)
                        .provider(provider)
                        .build(),
                platform));
    }
}
