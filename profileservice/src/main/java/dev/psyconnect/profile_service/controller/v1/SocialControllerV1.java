package dev.psyconnect.profile_service.controller.v1;

import java.util.List;
import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.dto.request.*;
import dev.psyconnect.profile_service.dto.response.*;
import dev.psyconnect.profile_service.service.FriendsService;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@RestController("socialControllerV1")
@RequestMapping("/v1/profile")
public class SocialControllerV1 {

    private final FriendsService friendsService;

    @Autowired
    public SocialControllerV1(FriendsService friendsService) {
        this.friendsService = friendsService;
    }

    @PostMapping("/friend/request")
    public ApiResponse<FriendRequestResponse> friendRequest(
            @RequestBody FriendRRequest receiver, @RequestHeader("X-Profile-Id") String sender) {
        return new ApiResponse<>(friendsService.createFriend(receiver, sender));
    }

    @PostMapping("/friend/accept")
    public ApiResponse<FriendAcceptResponse> friendAccept(
            @RequestBody FriendAcceptationRequest receiver, @RequestHeader("X-Profile-Id") String sender) {
        return new ApiResponse<>(friendsService.acceptFriend(receiver, sender));
    }

    @PostMapping("/friend/unfriend")
    public ApiResponse<UnFriendResponse> unfriend(
            @RequestBody UnfriendRequest receiver, @RequestHeader("X-Profile-Id") String sender) {

        return new ApiResponse<>(friendsService.unFriend(receiver, sender));
    }

    @PostMapping("/friend/undo")
    public ApiResponse<UndoFriendRequestResponse> undoFriendRequest(
            @RequestBody UndoFriendRRequest receiver, @RequestHeader("X-Profile-Id") String sender) {
        return new ApiResponse<>(friendsService.undoRequest(receiver, sender));
    }

    @DeleteMapping("/friend/request")
    public ApiResponse<DeclineFriendRequestResponse> declineFriendRequest(
            @RequestBody DeclineFriendRRequest receiver, @RequestHeader("X-Profile-Id") String sender) {
        return new ApiResponse<>(friendsService.declineRequest(receiver, sender));
    }

    @GetMapping("/friends/me")
    public ApiResponse<List<Map<String, Object>>> getAllFriendsById(@RequestHeader("X-Profile-Id") String profileId) {
        return new ApiResponse<>(friendsService.getAllFriendsById(profileId));
    }

    @GetMapping("/friends/received")
    public ApiResponse<List<Map<String, Object>>> getReceivedFriendRequests(
            @RequestHeader("X-Profile-Id") String profileId) {
        return new ApiResponse<>(friendsService.getReceivedRequests(profileId));
    }

    @GetMapping("/friends/sent")
    public ApiResponse<List<Map<String, Object>>> getSentFriendRequests(
            @RequestHeader("X-Profile-Id") String profileId) {
        return new ApiResponse<>(friendsService.getSentRequests(profileId));
    }

    @GetMapping("/friends/mutual")
    public ApiResponse<List<Map<String, Object>>> getMutualFriends(
            @RequestHeader("X-Profile-Id") String profileId, @RequestParam("with") String targetId) {
        return new ApiResponse<>(friendsService.getMutualFriends(profileId, targetId));
    }

    @GetMapping("/friends/suggestions")
    public ApiResponse<List<Map<String, Object>>> getFriendSuggestions(
            @RequestHeader("X-Profile-Id") String profileId) {
        return new ApiResponse<>(friendsService.getFriendSuggestions(profileId));
    }
}
