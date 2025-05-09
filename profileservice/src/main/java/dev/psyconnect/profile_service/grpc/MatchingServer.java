package dev.psyconnect.profile_service.grpc;

import io.grpc.Status;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;

import dev.psyconnect.profile_service.dto.request.DeclineFriendRRequest;
import dev.psyconnect.profile_service.dto.request.FriendAcceptationRequest;
import dev.psyconnect.profile_service.dto.request.FriendRRequest;
import dev.psyconnect.profile_service.dto.response.DeclineFriendRequestResponse;
import dev.psyconnect.profile_service.dto.response.FriendAcceptResponse;
import dev.psyconnect.profile_service.service.FriendsService;
import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import profile.CreateMatchingRequestGrpc;
import profile.MatchingCreatation;

@GrpcService
public class MatchingServer extends CreateMatchingRequestGrpc.CreateMatchingRequestImplBase {
    private static final Logger log = LoggerFactory.getLogger(MatchingServer.class);
    private FriendsService friendsService;

    @Autowired
    public MatchingServer(FriendsService friendsService) {
        this.friendsService = friendsService;
    }

    @Override
    public void createMatchingRequest(
            MatchingCreatation.MatchingRequest request,
            StreamObserver<MatchingCreatation.MatchingResponse> responseObserver) {
        try {
            friendsService.createFriend(
                    new FriendRRequest(request.getClientProfileId(), request.getMessage()),
                    request.getTherapistProfileId());

        } catch (Exception e) {
            log.info("Error while {}", e.getMessage());
            responseObserver.onError(Status.INTERNAL.withDescription(e.getMessage()).withCause(e).asException());
        }
        MatchingCreatation.MatchingResponse response = MatchingCreatation.MatchingResponse.newBuilder()
                .setStatus("Create friend request success")
                .setSuccess(true)
                .build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
        return;
    }

    @Override
    public void responseRequestMatching(
            MatchingCreatation.ResponseMatchingRequest request,
            StreamObserver<MatchingCreatation.ResponseMatchingResponse> responseObserver) {
        try {
            if (request.getOption().trim().equalsIgnoreCase("accept")) {
                FriendAcceptResponse response = friendsService.acceptFriend(
                        FriendAcceptationRequest.builder()
                                .target(request.getClientProfileId())
                                .build(),
                        request.getTherapistProfileId());
                responseObserver.onNext(MatchingCreatation.ResponseMatchingResponse.newBuilder()
                        .setStatus(response.status)
                        .setSuccess(true)
                        .build());
            } else if (request.getOption().trim().equalsIgnoreCase("decline")) {
                DeclineFriendRequestResponse response = friendsService.declineRequest(
                        DeclineFriendRRequest.builder()
                                .target(request.getClientProfileId())
                                .build(),
                        request.getTherapistProfileId());
                responseObserver.onNext(MatchingCreatation.ResponseMatchingResponse.newBuilder()
                        .setStatus(response.status)
                        .setSuccess(true)
                        .build());
            }
            responseObserver.onCompleted();
        } catch (Exception e) {
            log.error("Failed to process the request: {}", e.getMessage(), e);
            responseObserver.onError(Status.ALREADY_EXISTS.withDescription(e.getMessage()).withCause(e).asRuntimeException());
        }
    }
}
