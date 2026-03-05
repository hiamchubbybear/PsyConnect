package dev.psyconnect.identity_service.grpc.client;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import dev.psyconnect.grpc.ProfileCreationRequest;
import dev.psyconnect.grpc.ProfileCreationResponse;
import dev.psyconnect.grpc.ProfileCreationServiceGrpc;
import dev.psyconnect.identity_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.identity_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.identity_service.mapper.UserAccountMapper;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;

@Service
public class ProfileGRPCClient {
    private static final Logger log = LoggerFactory.getLogger(ProfileGRPCClient.class);
    private final ProfileCreationServiceGrpc.ProfileCreationServiceBlockingStub stub;
    private final UserAccountMapper userAccountMapper;

    public ProfileGRPCClient(
            UserAccountMapper userAccountMapper,
            @Value("${baseUriProfileService:localhost}") String profileServiceUrl) {

        String host = profileServiceUrl.replace("http://", "").replace("https://", "");
        if (host.contains(":")) {
            host = host.split(":")[0];
        }

        ManagedChannel channel =
                ManagedChannelBuilder.forAddress(host, 50051).usePlaintext().build();
        stub = ProfileCreationServiceGrpc.newBlockingStub(channel);
        log.info("gRPC Profile Service connected to: {}:50051 (Original URL: {})", host, profileServiceUrl);
        this.userAccountMapper = userAccountMapper;
    }

    public UserProfileCreationResponse createProfile(UserProfileCreationRequest userProfileCreationRequest) {
        try {
            log.info("Request: {}", userProfileCreationRequest.getProfileId());
            ProfileCreationRequest request = userAccountMapper.toUserProfileCreateRequest(userProfileCreationRequest);
            log.info("Request: {}", request.toString());
            ProfileCreationResponse response = stub.createUser(request);
            return userAccountMapper.toUserProfileCreateResponse(response);
        } catch (Exception e) {
            log.error("Error creating profile", e);
            throw e;
        }
    }
}
