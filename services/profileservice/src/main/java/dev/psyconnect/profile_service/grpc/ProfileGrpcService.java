package dev.psyconnect.profile_service.grpc;

import java.util.Optional;

import dev.psyconnect.grpc.profile.*;
import dev.psyconnect.profile_service.entity.UserProfile;
import dev.psyconnect.profile_service.repository.UserProfileRepository;
import io.grpc.stub.StreamObserver;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;

@Slf4j
@GrpcService
@RequiredArgsConstructor
public class ProfileGrpcService extends ProfileServiceGrpc.ProfileServiceImplBase {

    private final UserProfileRepository userProfileRepository;

    @Override
    public void getProfile(GetProfileRequest request, StreamObserver<GetProfileResponse> responseObserver) {
        log.info("gRPC GetProfile called for profileId: {}", request.getProfileId());

        try {
            Optional<UserProfile> profileOpt = userProfileRepository.findById(request.getProfileId());

            if (profileOpt.isEmpty()) {
                GetProfileResponse response = GetProfileResponse.newBuilder()
                        .setSuccess(false)
                        .setErrorMessage("Profile not found")
                        .build();
                responseObserver.onNext(response);
                responseObserver.onCompleted();
                return;
            }

            UserProfile profile = profileOpt.get();

            ProfileData profileData = ProfileData.newBuilder()
                    .setProfileId(profile.getId())
                    .setFirstName(profile.getFirstName() != null ? profile.getFirstName() : "")
                    .setLastName(profile.getLastName() != null ? profile.getLastName() : "")
                    .setAvatar(profile.getAvatar() != null ? profile.getAvatar() : "")
                    .setEmail(profile.getEmail() != null ? profile.getEmail() : "")
                    .setGender(profile.getGender() != null ? profile.getGender() : "")
                    .setDateOfBirth(
                            profile.getDateOfBirth() != null
                                    ? profile.getDateOfBirth().toString()
                                    : "")
                    .build();

            GetProfileResponse response = GetProfileResponse.newBuilder()
                    .setSuccess(true)
                    .setProfile(profileData)
                    .build();

            responseObserver.onNext(response);
            responseObserver.onCompleted();

            log.info("Successfully returned profile for: {}", request.getProfileId());

        } catch (Exception e) {
            log.error("Error fetching profile: {}", e.getMessage(), e);

            GetProfileResponse response = GetProfileResponse.newBuilder()
                    .setSuccess(false)
                    .setErrorMessage("Internal server error: " + e.getMessage())
                    .build();

            responseObserver.onNext(response);
            responseObserver.onCompleted();
        }
    }

    @Override
    public void getProfileBatch(
            GetProfileBatchRequest request, StreamObserver<GetProfileBatchResponse> responseObserver) {
        log.info("gRPC GetProfileBatch called for {} profiles", request.getProfileIdsCount());

        try {
            GetProfileBatchResponse.Builder responseBuilder =
                    GetProfileBatchResponse.newBuilder().setSuccess(true);

            for (String profileId : request.getProfileIdsList()) {
                Optional<UserProfile> profileOpt = userProfileRepository.findById(profileId);

                if (profileOpt.isPresent()) {
                    UserProfile profile = profileOpt.get();
                    ProfileData profileData = ProfileData.newBuilder()
                            .setProfileId(profile.getId())
                            .setFirstName(profile.getFirstName() != null ? profile.getFirstName() : "")
                            .setLastName(profile.getLastName() != null ? profile.getLastName() : "")
                            .setAvatar(profile.getAvatar() != null ? profile.getAvatar() : "")
                            .build();

                    responseBuilder.addProfiles(profileData);
                }
            }

            responseObserver.onNext(responseBuilder.build());
            responseObserver.onCompleted();

        } catch (Exception e) {
            log.error("Error in batch profile fetch: {}", e.getMessage(), e);

            GetProfileBatchResponse response = GetProfileBatchResponse.newBuilder()
                    .setSuccess(false)
                    .setErrorMessage("Internal server error: " + e.getMessage())
                    .build();

            responseObserver.onNext(response);
            responseObserver.onCompleted();
        }
    }
}
