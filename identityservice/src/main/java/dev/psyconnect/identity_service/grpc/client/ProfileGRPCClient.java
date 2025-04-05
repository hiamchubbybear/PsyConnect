package dev.psyconnect.identity_service.grpc.client;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;


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

    public ProfileGRPCClient(UserAccountMapper userAccountMapper) {
        ManagedChannel channel = ManagedChannelBuilder.forAddress("127.0.0.1", 9091)
                .usePlaintext()
                .build();
        stub = ProfileCreationServiceGrpc.newBlockingStub(channel);
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
            log.error("Error creating profile", e); // Use error level for clarity
            throw e; // Re-throw to see the full stack trace
        }
    }
}
