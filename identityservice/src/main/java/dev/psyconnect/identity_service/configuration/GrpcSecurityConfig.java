package dev.psyconnect.identity_service.configuration;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;

import io.grpc.Metadata;
import io.grpc.ServerCall;
import net.devh.boot.grpc.server.security.authentication.GrpcAuthenticationReader;

@Configuration
public class GrpcSecurityConfig {

    @Bean
    public GrpcAuthenticationReader grpcAuthenticationReader() {
        return new CustomGrpcAuthenticationReader();
    }

    private static class CustomGrpcAuthenticationReader implements GrpcAuthenticationReader {
        @Override
        public Authentication readAuthentication(ServerCall<?, ?> call, Metadata headers)
                throws AuthenticationException {
            return null;
        }
    }
}
