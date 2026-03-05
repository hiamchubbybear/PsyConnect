package dev.psyconnect.profile_service.grpc;

import java.util.List;
import java.util.stream.Collectors;

import dev.psyconnect.grpc.api_gateway.TokenCheckIdentity;
import dev.psyconnect.grpc.profile.*;
import dev.psyconnect.profile_service.dto.response.UserProfileResponse;
import dev.psyconnect.profile_service.grpc.client.IdentityGRPCClient;
import dev.psyconnect.profile_service.service.UserProfileService;
import io.grpc.stub.StreamObserver;
import lombok.AccessLevel;
import lombok.experimental.FieldDefaults;
import lombok.extern.slf4j.Slf4j;
import net.devh.boot.grpc.server.service.GrpcService;

@GrpcService
@Slf4j
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class ProfileDataServer extends ProfileServiceGrpc.ProfileServiceImplBase {

    private final UserProfileService userProfileService;
    private final IdentityGRPCClient identityGRPCClient;

    public ProfileDataServer(UserProfileService userProfileService, IdentityGRPCClient identityGRPCClient) {
        this.userProfileService = userProfileService;
        this.identityGRPCClient = identityGRPCClient;
    }

    @Override
    public void getProfile(GetProfileRequest request, StreamObserver<GetProfileResponse> responseObserver) {
        try {
            UserProfileResponse profile = userProfileService.get(request.getProfileId());
            TokenCheckIdentity.UserInfoResponse userInfo = identityGRPCClient.getUserInfo(request.getProfileId());

            ProfileData data = ProfileData.newBuilder()
                    .setProfileId(profile.getProfileId())
                    .setFirstName(profile.getFirstName())
                    .setLastName(profile.getLastName())
                    .setAvatar(profile.getAvatarUri())
                    .setGender(profile.getGender())
                    .setEmail(userInfo.getSuccess() ? userInfo.getEmail() : "")
                    .build();

            responseObserver.onNext(GetProfileResponse.newBuilder()
                    .setSuccess(true)
                    .setProfile(data)
                    .build());
        } catch (Exception e) {
            log.error("Error in GetProfile", e);
            responseObserver.onNext(GetProfileResponse.newBuilder()
                    .setSuccess(false)
                    .setErrorMessage(e.getMessage())
                    .build());
        } finally {
            responseObserver.onCompleted();
        }
    }

    @Override
    public void getProfileBatch(
            GetProfileBatchRequest request, StreamObserver<GetProfileBatchResponse> responseObserver) {
        try {
            List<UserProfileResponse> profiles = userProfileService.getBatch(request.getProfileIdsList());

            List<ProfileData> dataList = profiles.stream()
                    .map(p -> {
                        TokenCheckIdentity.UserInfoResponse userInfo = identityGRPCClient.getUserInfo(p.getProfileId());
                        return ProfileData.newBuilder()
                                .setProfileId(p.getProfileId())
                                .setFirstName(p.getFirstName())
                                .setLastName(p.getLastName())
                                .setAvatar(p.getAvatarUri())
                                .setGender(p.getGender())
                                .setEmail(userInfo.getSuccess() ? userInfo.getEmail() : "")
                                .build();
                    })
                    .collect(Collectors.toList());

            responseObserver.onNext(GetProfileBatchResponse.newBuilder()
                    .setSuccess(true)
                    .addAllProfiles(dataList)
                    .build());
        } catch (Exception e) {
            log.error("Error in GetProfileBatch", e);
            responseObserver.onNext(GetProfileBatchResponse.newBuilder()
                    .setSuccess(false)
                    .setErrorMessage(e.getMessage())
                    .build());
        } finally {
            responseObserver.onCompleted();
        }
    }
}
