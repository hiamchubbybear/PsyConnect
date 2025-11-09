package dev.psyconnect.api_service.configuration;

import com.sun.jdi.InternalException;
import dev.psyconnect.api_service.configuration.producer.KafkaService;
import dev.psyconnect.api_service.dto.LogEvent;
import dev.psyconnect.api_service.dto.LogLevel;
import dev.psyconnect.api_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.api_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.api_service.service.TokenValidateService;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.SignatureException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.cloud.gateway.filter.GlobalFilter;
import org.springframework.cloud.gateway.route.Route;
import org.springframework.cloud.gateway.support.ServerWebExchangeUtils;
import org.springframework.core.Ordered;
import org.springframework.http.HttpHeaders;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.Map;

@Component
@Slf4j
public class JwtAuthFilter implements GlobalFilter, Ordered {
    @Value("${SIGNER_KEY}")
    private String SIGNER_KEY;
    @Autowired
    TokenValidateService service;
    @Autowired
    private KafkaService kafkaService;

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {

        String userId = null;
        String profileId = null;
        String scopes = null;

        String authHeader = exchange.getRequest().getHeaders().getFirst(HttpHeaders.AUTHORIZATION);
        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            String token = authHeader.substring(7);
            if (!service.isTokenValid(token)) throw new CustomExceptionHandler(ErrorCode.UNAUTHORIZED);
            try {
                Claims claims = getClaimsFromToken(token);
                userId = claims.get("accountId", String.class);
                profileId = claims.get("profileId", String.class);
                scopes = claims.get("scope", String.class);
                String iss = claims.get("iss", String.class);
                if (!iss.equals("PsyConnect Authentication Service"))
                    throw new InternalException("Token authentication failed");
                ServerHttpRequest mutatedRequest = exchange.getRequest().mutate()
                        .header("X-User-Id", userId)
                        .header("X-Profile-Id", profileId)
                        .header("X-Roles", scopes)
                        .build();
                
                
                return chain.filter(exchange.mutate().request(mutatedRequest).build());
            } catch (SignatureException e) {
                kafkaService.sendLog(buildLog(
                        "api-gateway", "",
                        "login", "Invalid JWT signature",
                        null, LogLevel.ERROR
                ));
                log.error("Invalid JWT signature: {}", e.getMessage());
            }
        }
        Route route = exchange.getAttribute(ServerWebExchangeUtils.GATEWAY_ROUTE_ATTR);
        if (route != null) {
            log.info("Audit: userId={} profileId={} route={}", userId, profileId, route.getId());
        } else {
            log.warn("Audit: userId={} profileId={} route not resolved yet", userId, profileId);
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

    private LogEvent buildLog(
            String service, String userId, String action, String message, Map<String, Object> metadata, LogLevel level) {
        return LogEvent.builder()
                .service(service)
                .level(level)
                .timestamp(Instant.now().toString())
                .userId(userId)
                .action(action)
                .message(message)
                .metadata(metadata)
                .build();
    }
}