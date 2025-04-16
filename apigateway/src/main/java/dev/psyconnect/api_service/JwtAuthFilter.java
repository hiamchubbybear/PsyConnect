package dev.psyconnect.api_service;

import com.sun.jdi.InternalException;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.SignatureException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.cloud.gateway.filter.GlobalFilter;
import org.springframework.core.Ordered;
import org.springframework.http.HttpHeaders;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.context.ReactiveSecurityContextHolder;
import org.springframework.security.core.context.SecurityContext;
import org.springframework.security.core.context.SecurityContextImpl;
import org.springframework.security.core.userdetails.User;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

import java.nio.charset.StandardCharsets;
import java.util.Collections;

@Component
@Slf4j
public class JwtAuthFilter implements GlobalFilter, Ordered {

    @Value("${SIGNER_KEY}")
    private String SIGNER_KEY;

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {

        String authHeader = exchange.getRequest().getHeaders().getFirst(HttpHeaders.AUTHORIZATION);
        log.info("Security filters: {}", authHeader);

        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            String token = authHeader.substring(7);
            try {
                Claims claims = getClaimsFromToken(token);
                String username = claims.getSubject();
                String userId = claims.get("accountId", String.class);
                String profileId = claims.get("profileId", String.class);
                String scopes = claims.get("scope", String.class);
                String iss = claims.get("iss", String.class);
                // If token issuer is not PsyConnect Authentication Service throw exception
                if (!iss.equals("PsyConnect Authentication Service"))
                    throw new InternalException("Token authentication failed");
                // Add accountId and profileId into header request
                ServerHttpRequest mutatedRequest = exchange.getRequest().mutate()
                        .header("X-User-Id", userId)
                        .header("X-Profile-Id", profileId)
                        .header("X-Roles", scopes)
                        .build();
                return chain.filter(exchange.mutate().request(mutatedRequest).build());
            } catch (SignatureException e) {
                log.error("Invalid JWT signature: {}", e.getMessage());
            }
        }

        return chain.filter(exchange);
    }

    private Claims getClaimsFromToken(String token) {
        return Jwts.parserBuilder()
                .setSigningKey(SIGNER_KEY.getBytes(StandardCharsets.UTF_8))
                .build()
                .parseClaimsJws(token)
                .getBody();
    }

    @Override
    public int getOrder() {
        return -1;
    }
}
