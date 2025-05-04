package dev.psyconnect.profile_service.grpc;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;

import dev.psyconnect.profile_service.dto.request.FriendRRequest;
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
                    new FriendRRequest(request.getClientProfileId()), request.getTherapistProfileId());

        } catch (Exception e) {
            log.info("Error while {}", e.getMessage());
            MatchingCreatation.MatchingResponse response = new MatchingCreatation.MatchingResponse.Builder()
                    .setStatus(e.getMessage())
                    .setSuccess(false)
                    .build();
            responseObserver.onNext(response);
            responseObserver.onCompleted();
            return;
        }
        MatchingCreatation.MatchingResponse response = new MatchingCreatation.MatchingResponse.Builder()
                .setStatus("Create friend request success")
                .setSuccess(true)
                .build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
        return;
    }
}
