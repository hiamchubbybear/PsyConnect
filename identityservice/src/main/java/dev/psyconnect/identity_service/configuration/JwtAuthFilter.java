package dev.psyconnect.identity_service.configuration;

import java.io.IOException;
import java.text.ParseException;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import com.nimbusds.jose.JOSEException;

import dev.psyconnect.identity_service.dto.request.AuthenticationFilterRequest;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.interfaces.IUserAccountService;
import dev.psyconnect.identity_service.service.AuthenticationService;
import lombok.extern.slf4j.Slf4j;

@Component
@Slf4j
public class JwtAuthFilter extends OncePerRequestFilter {
    @Autowired
    private AuthenticationService authenticationService;

    @Autowired
    private IUserAccountService userAccountService;

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
            throws ServletException, IOException {
        String authHeader = request.getHeader("Authorization");
        String token = null;
        String username = null;

        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            token = authHeader.substring(7);
            log.info("JWT token: {}", token);
//            if (authenticationService.isTokenValid(token))
//                throw new CustomExceptionHandler(ErrorCode.USER_UNAUTHENTICATED);
            try {
                username = authenticationService.extractUsername(token);
            } catch (ParseException e) {
                throw new RuntimeException(e);
            } catch (JOSEException e) {
                throw new RuntimeException(e);
            }
        }

        if (username != null && SecurityContextHolder.getContext().getAuthentication() == null) {
            AuthenticationFilterRequest userAccount = userAccountService.loadUserByUsername(username);
            try {
                if (authenticationService.validateToken(token, userAccount)) {
                    UsernamePasswordAuthenticationToken authToken =
                            new UsernamePasswordAuthenticationToken(userAccount, null, userAccount.getAuthorities());
                    authToken.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));
                    SecurityContextHolder.getContext().setAuthentication(authToken);
                }
            } catch (ParseException e) {
                throw new RuntimeException(e);
            } catch (JOSEException e) {
                throw new RuntimeException(e);
            }
        }
        filterChain.doFilter(request, response);
    }
}
