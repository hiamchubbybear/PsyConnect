package dev.psyconnect.api_service.configuration;

import dev.psyconnect.api_service.configuration.producer.KafkaService;
import dev.psyconnect.api_service.dto.LogEvent;
import dev.psyconnect.api_service.dto.LogLevel;
import dev.psyconnect.api_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.api_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.api_service.service.TokenValidateService;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.JwtException;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.SignatureException;
import lombok.extern.slf4j.Slf4j;
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

    private final String SIGNER_KEY = "5f2a6e8ffb0fd3cbac079162c3168a2c2dfb5b115a7119ef38dd5fb1553295d52b29ed0f70ff57320edf7de1ceeee22ed574988301bdd9eebe1ff9d3abb05ed8f8471dbb397b0375583a3b3887d108cac2694da14335c3a7379a30755672fa22a77615fb5bba04ee6d8d58faac42b704512b4174dff2a8bd8f678c0e3332c404dec6c2dc94dc50901d1ae58a033c3f8b50aedcde149a07c0f25dbfdf567ec156664f9d89bffd7741dc0dce0d62515000d6bac88f91a7793fd93d474e8717e51e1108897893b8374b526f51a29d1890412d5d3b7481cd025ed9d28861ffdfdf2c37fcc44eaf2477ba7be4118d5325b2a93f1f21d04d915d6fd9bd6b8cba111d41a7f2c1f8889212156a137ed21a84477a099ef5565d174a5d34d08cc7b616316fe32d3ce1f988012321458f5251719757abf43453ba16edb26faf15d70453210a3210ae7c65d877d20fb5073dce9e1630bb5a0ce051e4b82bb263c4b456b8a4f50270419234315b9e41e341171a8bbebb6ec4d974d71649b92b3d8a3c8d2caf04b5d96805e53ba79458bf5161e1a14db97447cd982bf6e2f72375d3598cca1bb2da2ea27acd6a5184276598d6516a7b6c034dd7f8d73127f3ce99a1140ddd295fc4e53daaa4b77cf16c9db1c6ff09bc0640024b43554f71df6899995dfb1d95e2bc0f73487079024aedd7cbae122bc831a312d2f0b7f58528d75c2bdbe2378b1d";
    private final TokenValidateService service;
    private final KafkaService kafkaService;

    public JwtAuthFilter(TokenValidateService service, KafkaService kafkaService) {
        this.service = service;
        this.kafkaService = kafkaService;
    }

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        String userId = null;
        String profileId = null;
        String scopes = null;

        String authHeader = exchange.getRequest().getHeaders().getFirst(HttpHeaders.AUTHORIZATION);
        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            String token = authHeader.substring(7);
            if (!service.isTokenValid(token)) {
                throw new CustomExceptionHandler(ErrorCode.UNAUTHORIZED);
            }
            try {
                Claims claims = getClaimsFromToken(token);
                userId = claims.get("accountId", String.class);
                profileId = claims.get("profileId", String.class);
                scopes = claims.get("scope", String.class);
                String iss = claims.get("iss", String.class);

                if (!iss.equals("PsyConnect Authentication Service")) {
                    throw new JwtException("Token authentication failed");
                }

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
                throw new CustomExceptionHandler(ErrorCode.UNAUTHORIZED);
            } catch (JwtException e) {
                log.error("JWT validation failed: {}", e.getMessage());
                throw new CustomExceptionHandler(ErrorCode.UNAUTHORIZED);
            }
        }

        // Audit logging
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
            String service, String userId, String action, String message,
            Map<String, Object> metadata, LogLevel level) {
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
