package dev.psyconnect.profile_service.grpc;

import dev.psyconnect.grpc.*;
import dev.psyconnect.profile_service.dto.request.UserProfileCreationRequest;
import dev.psyconnect.profile_service.dto.response.UserProfileCreationResponse;
import dev.psyconnect.profile_service.mapper.UserProfileMapper;
import dev.psyconnect.profile_service.service.UserProfileService;
import io.grpc.stub.StreamObserver;
import lombok.AccessLevel;
import lombok.experimental.FieldDefaults;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;


@GrpcService
@Slf4j
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class ProfileServer extends ProfileCreationServiceGrpc.ProfileCreationServiceImplBase {
    private final UserProfileService userProfileService;
    private final UserProfileMapper profileMapper;

    public ProfileServer(UserProfileService userProfileService, UserProfileMapper profileMapper) {
        this.userProfileService = userProfileService;
        this.profileMapper = profileMapper;
        System.out.println("ProfileServer gRPC initialized");
    }

    @Override
    public void helloword(Hello request, StreamObserver<HelloResponse> responseObserver) {
        HelloResponse response = HelloResponse.newBuilder()
                .setMessage(request.getMessage())
                .build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }

    @Override
    public void createUser(ProfileCreationRequest request, StreamObserver<ProfileCreationResponse> responseObserver) {
        log.info("Received createUser request: {}", request.getProfileId());
        UserProfileCreationRequest creationRequest = profileMapper.toUserProfileRequest(request);
        log.info("Created user datetime1: {}", creationRequest.getDob());
        UserProfileCreationResponse response = userProfileService.create(creationRequest);
        log.info("Created user datetime1: {}", response.getDob());
        ProfileCreationResponse responseMapped = ProfileCreationResponse.newBuilder()
                .setProfileId(response.getProfileId())
                .setFirstName(response.getFirstName())
                .setLastName(response.getLastName())
                .setDob(response.getDob())
                .setAddress(response.getAddress())
                .setGender(response.getGender())
                .setAvatarUri(response.getAvatarUri())
                .build();
        ;
        responseObserver.onNext(responseMapped);
        responseObserver.onCompleted();
    }
}
