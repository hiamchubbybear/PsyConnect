package dev.psyconnect.profile_service.grpc;

import dev.psyconnect.grpc.ProfileCreationServiceGrpc;
import dev.psyconnect.grpc.consultation_profile.CheckExistedProfile;
import dev.psyconnect.grpc.consultation_profile.CheckProfileServiceGrpc;
import dev.psyconnect.profile_service.service.UserProfileService;
import io.grpc.stub.StreamObserver;
import lombok.AccessLevel;
import lombok.experimental.FieldDefaults;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;

@GrpcService
@Slf4j
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
class ProfileConsultationServer extends CheckProfileServiceGrpc.CheckProfileServiceImplBase {
    private  UserProfileService  userProfileService;

    ProfileConsultationServer(UserProfileService userProfileService) {
        this.userProfileService = userProfileService;
    }

    // Check whether profile exists in profile service before create a new therapist profile ______
    @Override
    public void checkProfileExists(CheckExistedProfile.ProfileRequest request, StreamObserver<CheckExistedProfile.ProfileResponse> responseObserver) {
        boolean result = userProfileService.checkProfileExisted(request.getProfileId());
        log.info("Check user request: {}", request.getProfileId());
        var response  = CheckExistedProfile.ProfileResponse.newBuilder().setExists(result).build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
    }
}
