package dev.psyconnect.identity_service.configuration;

import com.nimbusds.jose.JOSEException;
import dev.psyconnect.identity_service.dto.request.AuthenticationFilterRequest;
import dev.psyconnect.identity_service.dto.request.AuthenticationRequest;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.interfaces.IUserAccountService;
import dev.psyconnect.identity_service.model.UserAccount;
import dev.psyconnect.identity_service.service.AuthenticationService;
import dev.psyconnect.identity_service.service.UserAccountService;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.web.authentication.WebAuthenticationDetailsSource;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.text.ParseException;

@Component
@Slf4j
public class JwtAuthFilter extends OncePerRequestFilter {
    @Autowired
    private AuthenticationService authenticationService;

    @Autowired
    private IUserAccountService userAccountService;

    @Override
    protected void doFilterInternal(HttpServletRequest request, HttpServletResponse response, FilterChain filterChain) throws ServletException, IOException {
        // Retrieve the Authorization header
        String authHeader = request.getHeader("Authorization");
        String token = null;
        String username = null;

        // Check if the header starts with "Bearer "
        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            token = authHeader.substring(7); // Extract token
            log.info("JWT token: {}", token);
            // Check if token is not in black list
            if(authenticationService.isTokenInvalid(token)) throw new CustomExceptionHandler(ErrorCode.USER_UNAUTHENTICATED);
            try {
                username = authenticationService.extractUsername(token); // Extract username from token
            } catch (ParseException e) {
                throw new RuntimeException(e);
            } catch (JOSEException e) {
                throw new RuntimeException(e);
            }
        }


        // If the token is valid and no authentication is set in the context
        if (username != null && SecurityContextHolder.getContext().getAuthentication() == null) {
            AuthenticationFilterRequest userAccount = userAccountService.loadUserByUsername(username);
            // Validate token and set authentication
            try {
                if (authenticationService.validateToken(token, userAccount)) {
                    UsernamePasswordAuthenticationToken authToken = new UsernamePasswordAuthenticationToken(
                            userAccount,
                            null,
                            userAccount.getAuthorities()
                    );
                    authToken.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));
                    SecurityContextHolder.getContext().setAuthentication(authToken);
                }
            } catch (ParseException e) {
                throw new RuntimeException(e);
            } catch (JOSEException e) {
                throw new RuntimeException(e);
            }
        }
        // Continue the filter chain
        filterChain.doFilter(request, response);
    }
}
