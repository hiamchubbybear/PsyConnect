package dev.psyconnect.profile_service.service;

import java.time.OffsetDateTime;

import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.dto.request.*;
import dev.psyconnect.profile_service.dto.response.*;
import dev.psyconnect.profile_service.enums.FriendShipStatus;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.model.FriendRelationship;
import dev.psyconnect.profile_service.model.Profile;
import dev.psyconnect.profile_service.repository.FriendRepository;
import dev.psyconnect.profile_service.repository.ProfileRepository;
import dev.psyconnect.profile_service.utils.FriendShipUtils;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import lombok.extern.slf4j.Slf4j;

@Service
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
@Slf4j
public class FriendsService {
    private final ProfileRepository profileRepository;
    private final FriendShipUtils friendShipUtils;
    private final FriendRepository friendRepository;

    public FriendRequestResponse createFriend(FriendRRequest request, String senderId) {
        // Validate not sending to self
        if (request.getTarget().equals(senderId)) throw new CustomExceptionHandler(ErrorCode.CAN_NOT_REQUEST);
        Profile receiver = profileRepository
                .findById(request.getTarget())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        Profile sender = profileRepository
                .findById(senderId)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        if (friendShipUtils.hasAnyRelationship(sender, receiver.getProfileId())
                || friendShipUtils.hasAnyRelationship(receiver, request.getTarget()))
            throw new CustomExceptionHandler(ErrorCode.ALREADY_FRIEND_OR_REQUEST);
        FriendRelationship relationship = FriendRelationship.builder()
                .target(receiver)
                .createdAt(OffsetDateTime.now())
                .message(request.getMessage())
                .status(FriendShipStatus.PENDING)
                .build();
        sender.setFriendShip(relationship);
        Profile saved = profileRepository.save(sender);
        return FriendRequestResponse.builder()
                .friendRequestStatus("Success")
                .status(saved.toString())
                .build();
    }

    public FriendAcceptResponse acceptFriend(FriendAcceptationRequest request, String accepterId) {
        String reqId = request.getTarget();
        String tarId = accepterId;
        if (request.getTarget().equals(accepterId)) throw new CustomExceptionHandler(ErrorCode.CAN_NOT_REQUEST);
        if (!friendShipUtils.existByUser(reqId, tarId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        if ((friendRepository.exists(tarId, reqId, FriendShipStatus.ACCEPTED)
                || friendRepository.exists(reqId, tarId, FriendShipStatus.ACCEPTED)))
            throw new CustomExceptionHandler(ErrorCode.ALREADY_FRIEND_OR_REQUEST);
        // Xor if none or both of these user  have(not) pending request throw exception
        // Requester - Pending -> Accepter ( OK )
        // Requester <- Pending -> Accepter ( NOT OK )
        // Requester <-  -> Accepter ( NOT OK )
        if (!(friendRepository.exists(tarId, reqId, FriendShipStatus.PENDING)
                ^ friendRepository.exists(reqId, tarId, FriendShipStatus.PENDING)))
            throw new CustomExceptionHandler(ErrorCode.USER_DONT_HAVE_FRIEND_REQUEST);
        friendRepository.accept(reqId, tarId, FriendShipStatus.ACCEPTED);
        friendRepository.create(tarId, reqId, FriendShipStatus.ACCEPTED);
        return new FriendAcceptResponse("Success", FriendShipStatus.ACCEPTED);
    }

    public UnFriendResponse unFriend(UnfriendRequest request, String unfriendUserId) {
        String reqId = request.getTarget();
        String tarId = unfriendUserId;
        if (request.getTarget().equals(unfriendUserId)) throw new CustomExceptionHandler(ErrorCode.CAN_NOT_REQUEST);
        if (!friendShipUtils.existByUser(reqId, tarId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        // If requester or target doesn't have friend with each other
        // Requester - Accepted -> Target ( NOT OK )
        // Requester <- Accepted -> Target ( OK )
        if (!(friendRepository.exists(tarId, reqId, FriendShipStatus.ACCEPTED)
                && friendRepository.exists(reqId, tarId, FriendShipStatus.ACCEPTED)))
            throw new CustomExceptionHandler(ErrorCode.DO_NOT_FRIEND);
        friendRepository.delete(reqId, tarId, FriendShipStatus.ACCEPTED);
        return new UnFriendResponse("Success", FriendShipStatus.ACCEPTED);
    }

    public UndoFriendRequestResponse undoRequest(UndoFriendRRequest request, String undoUserId) {
        String reqId = request.getTarget();
        String tarId = undoUserId;
        if (request.getTarget().equals(undoUserId)) throw new CustomExceptionHandler(ErrorCode.CAN_NOT_REQUEST);
        if (!friendShipUtils.existByUser(reqId, tarId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        if (!(friendRepository.exists(tarId, reqId, FriendShipStatus.PENDING)
                ^ friendRepository.exists(reqId, tarId, FriendShipStatus.PENDING)))
            throw new CustomExceptionHandler(ErrorCode.USER_DONT_HAVE_FRIEND_REQUEST);
        friendRepository.delete(reqId, tarId, FriendShipStatus.PENDING);
        return new UndoFriendRequestResponse("Success", FriendShipStatus.UNFRIEND);
    }

    public DeclineFriendRequestResponse declineRequest(DeclineFriendRRequest request, String declineUserId) {
        String reqId = request.getTarget();
        String tarId = declineUserId;
        if (request.getTarget().equals(declineUserId)) throw new CustomExceptionHandler(ErrorCode.CAN_NOT_REQUEST);
        if (!friendShipUtils.existByUser(reqId, tarId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        if (!(friendRepository.exists(reqId, tarId, FriendShipStatus.PENDING)
                ^ friendRepository.exists(tarId, reqId, FriendShipStatus.PENDING)))
            throw new CustomExceptionHandler(ErrorCode.ALREADY_FRIEND_OR_REQUEST);
        friendRepository.delete(reqId, tarId, FriendShipStatus.PENDING);
        return new DeclineFriendRequestResponse("Success", FriendShipStatus.UNFRIEND);
    }
}
