package dev.psyconnect.profile_service.grpc.client;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import dev.psyconnect.grpc.api_gateway.IdentityServiceGrpc;
import dev.psyconnect.grpc.api_gateway.TokenCheckIdentity;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

@Service
public class IdentityGRPCClient {
    private static final Logger log = LoggerFactory.getLogger(IdentityGRPCClient.class);
    private final IdentityServiceGrpc.IdentityServiceBlockingStub stub;

    public IdentityGRPCClient(@Value("${IDENTITY_SERVICE_URL:http://identityservice}") String identityServiceUrl) {
        // Extract host from URL (e.g. http://identity-service:8080 -> identity-service)
        String host = identityServiceUrl.replace("http://", "").replace("https://", "");
        if (host.contains(":")) {
            host = host.split(":")[0];
        }

        // Use standard gRPC port 50051 for internal communication
        ManagedChannel channel =
                ManagedChannelBuilder.forAddress(host, 50051).usePlaintext().build();
        stub = IdentityServiceGrpc.newBlockingStub(channel);
        log.info("gRPC Identity Service connected to: {}:50051 (Original URL: {})", host, identityServiceUrl);
    }

    public TokenCheckIdentity.UserInfoResponse getUserInfo(String profileId) {
        try {
            TokenCheckIdentity.UserInfoRequest request = TokenCheckIdentity.UserInfoRequest.newBuilder()
                    .setProfileId(profileId)
                    .build();
            return stub.getUserInfoByProfileId(request);
        } catch (Exception e) {
            log.error("Error fetching user info from identity service", e);
            return TokenCheckIdentity.UserInfoResponse.newBuilder()
                    .setSuccess(false)
                    .build();
        }
    }
}
