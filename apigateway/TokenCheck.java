package dev.psyconnect.api_service.service;

import dev.psyconnect.grpc.api_gateway.TokenCheckIdentity;
import dev.psyconnect.grpc.api_gateway.TokenServiceGrpc;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import jakarta.annotation.PreDestroy;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

@Slf4j
@Service
public class TokenCheck {

    private final ManagedChannel channel;
    private final TokenServiceGrpc.TokenServiceBlockingStub stub;

    public TokenCheck(@Value("${grpc.identity.host:localhost}") String host,
                      @Value("${grpc.identity.port:9090}") int port) {
        this.channel = ManagedChannelBuilder.forAddress(host, port)
                .usePlaintext()
                .build();
        this.stub = TokenServiceGrpc.newBlockingStub(channel);
        log.info("TokenCheck gRPC client initialized on {}:{}", host, port);
    }

    public boolean isValidToken(String token) {
        try {
            TokenCheckIdentity.TokenRequest request =
                    TokenCheckIdentity.TokenRequest.newBuilder().setToken(token).build();
            TokenCheckIdentity.TokenResponse response = stub.tokenCheckValid(request);
            return response.getValid();
        } catch (Exception e) {
            log.error("Error while checking token via gRPC: {}", e.getMessage(), e);
            throw e;
        }
    }

    @PreDestroy
    public void shutdown() {
        log.info("Shutting down TokenCheck gRPC channel...");
        if (channel != null && !channel.isShutdown()) {
            channel.shutdownNow();
        }
    }
}
