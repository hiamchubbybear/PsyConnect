package dev.psyconnect.profile_service.service;

import java.time.OffsetDateTime;
import java.util.Map;

import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.dto.request.*;
import dev.psyconnect.profile_service.dto.response.*;
import dev.psyconnect.profile_service.enums.FriendShipStatus;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.kafka.service.KafkaService;
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
    private final KafkaService kafkaService;

    public FriendRequestResponse createFriend(FriendRRequest request, String senderId) {
        try {
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

            kafkaService.sendLog(buildLog(
                    "friends-service",
                    senderId,
                    "Create friend request",
                    "Success",
                    Map.of("target", request.getTarget())));

            return FriendRequestResponse.builder()
                    .friendRequestStatus("Success")
                    .status(saved.toString())
                    .build();
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "friends-service",
                    senderId,
                    "Create friend request",
                    "Failed",
                    Map.of("target", request.getTarget(), "error", e.getMessage())));
            throw e;
        }
    }

    public FriendAcceptResponse acceptFriend(FriendAcceptationRequest request, String accepterId) {
        try {
            String reqId = request.getTarget();
            String tarId = accepterId;
            if (request.getTarget().equals(accepterId)) throw new CustomExceptionHandler(ErrorCode.CAN_NOT_REQUEST);
            if (!friendShipUtils.existByUser(reqId, tarId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
            if ((friendRepository.exists(tarId, reqId, FriendShipStatus.ACCEPTED)
                    || friendRepository.exists(reqId, tarId, FriendShipStatus.ACCEPTED)))
                throw new CustomExceptionHandler(ErrorCode.ALREADY_FRIEND_OR_REQUEST);
            if (!(friendRepository.exists(tarId, reqId, FriendShipStatus.PENDING)
                    ^ friendRepository.exists(reqId, tarId, FriendShipStatus.PENDING)))
                throw new CustomExceptionHandler(ErrorCode.USER_DONT_HAVE_FRIEND_REQUEST);
            friendRepository.accept(reqId, tarId, FriendShipStatus.ACCEPTED);
            friendRepository.create(tarId, reqId, FriendShipStatus.ACCEPTED);

            kafkaService.sendLog(buildLog(
                    "friends-service", accepterId, "Accept friend request", "Success", Map.of("target", reqId)));

            return new FriendAcceptResponse("Success", FriendShipStatus.ACCEPTED);
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "friends-service",
                    accepterId,
                    "Accept friend request",
                    "Failed",
                    Map.of("target", request.getTarget(), "error", e.getMessage())));
            throw e;
        }
    }

    public UnFriendResponse unFriend(UnfriendRequest request, String unfriendUserId) {
        try {
            String reqId = request.getTarget();
            String tarId = unfriendUserId;
            if (request.getTarget().equals(unfriendUserId)) throw new CustomExceptionHandler(ErrorCode.CAN_NOT_REQUEST);
            if (!friendShipUtils.existByUser(reqId, tarId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
            if (!(friendRepository.exists(tarId, reqId, FriendShipStatus.ACCEPTED)
                    && friendRepository.exists(reqId, tarId, FriendShipStatus.ACCEPTED)))
                throw new CustomExceptionHandler(ErrorCode.DO_NOT_FRIEND);
            friendRepository.delete(reqId, tarId, FriendShipStatus.ACCEPTED);

            kafkaService.sendLog(
                    buildLog("friends-service", unfriendUserId, "Unfriend user", "Success", Map.of("target", reqId)));

            return new UnFriendResponse("Success", FriendShipStatus.ACCEPTED);
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "friends-service",
                    unfriendUserId,
                    "Unfriend user",
                    "Failed",
                    Map.of("target", request.getTarget(), "error", e.getMessage())));
            throw e;
        }
    }

    public UndoFriendRequestResponse undoRequest(UndoFriendRRequest request, String undoUserId) {
        try {
            String reqId = request.getTarget();
            String tarId = undoUserId;
            if (request.getTarget().equals(undoUserId)) throw new CustomExceptionHandler(ErrorCode.CAN_NOT_REQUEST);
            if (!friendShipUtils.existByUser(reqId, tarId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
            if (!(friendRepository.exists(tarId, reqId, FriendShipStatus.PENDING)
                    ^ friendRepository.exists(reqId, tarId, FriendShipStatus.PENDING)))
                throw new CustomExceptionHandler(ErrorCode.USER_DONT_HAVE_FRIEND_REQUEST);
            friendRepository.delete(reqId, tarId, FriendShipStatus.PENDING);

            kafkaService.sendLog(
                    buildLog("friends-service", undoUserId, "Undo friend request", "Success", Map.of("target", reqId)));

            return new UndoFriendRequestResponse("Success", FriendShipStatus.UNFRIEND);
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "friends-service",
                    undoUserId,
                    "Undo friend request",
                    "Failed",
                    Map.of("target", request.getTarget(), "error", e.getMessage())));
            throw e;
        }
    }

    public DeclineFriendRequestResponse declineRequest(DeclineFriendRRequest request, String declineUserId) {
        try {
            String reqId = request.getTarget();
            String tarId = declineUserId;
            if (request.getTarget().equals(declineUserId)) throw new CustomExceptionHandler(ErrorCode.CAN_NOT_REQUEST);
            if (!friendShipUtils.existByUser(reqId, tarId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
            if (!(friendRepository.exists(reqId, tarId, FriendShipStatus.PENDING)
                    ^ friendRepository.exists(tarId, reqId, FriendShipStatus.PENDING)))
                throw new CustomExceptionHandler(ErrorCode.ALREADY_FRIEND_OR_REQUEST);
            friendRepository.delete(reqId, tarId, FriendShipStatus.PENDING);

            kafkaService.sendLog(buildLog(
                    "friends-service", declineUserId, "Decline friend request", "Success", Map.of("target", reqId)));

            return new DeclineFriendRequestResponse("Success", FriendShipStatus.UNFRIEND);
        } catch (Exception e) {
            kafkaService.sendLog(buildLog(
                    "friends-service",
                    declineUserId,
                    "Decline friend request",
                    "Failed",
                    Map.of("target", request.getTarget(), "error", e.getMessage())));
            throw e;
        }
    }

    private LogEvent buildLog(
            String service, String userId, String action, String message, Map<String, Object> metadata) {
        return LogEvent.builder()
                .service(service)
                .eventId("LOG")
                .userId(userId)
                .action(action)
                .message(message)
                .metadata(metadata)
                .build();
    }
}
