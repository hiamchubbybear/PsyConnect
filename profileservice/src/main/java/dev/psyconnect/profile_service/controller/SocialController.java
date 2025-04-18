package dev.psyconnect.profile_service.controller;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import dev.psyconnect.profile_service.apiresponse.ApiResponse;
import dev.psyconnect.profile_service.dto.request.*;
import dev.psyconnect.profile_service.dto.response.*;
import dev.psyconnect.profile_service.service.FriendsService;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@RestController
@RequestMapping("/profile")
public class SocialController {
    private final FriendsService friendsService;

    @Autowired
    public SocialController(FriendsService friendsService) {
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
    public ApiResponse<UnFriendResponse> friendAccept(
            @RequestBody UnfriendRequest receiver, @RequestHeader("X-Profile-Id") String sender) {
        return new ApiResponse<>(friendsService.unFriend(receiver, sender));
    }

    @PostMapping("/friend/undo")
    public ApiResponse<UndoFriendRequestResponse> friendAccept(
            @RequestBody UndoFriendRRequest receiver, @RequestHeader("X-Profile-Id") String sender) {
        return new ApiResponse<>(friendsService.undoRequest(receiver, sender));
    }

    @DeleteMapping("/friend/request")
    public ApiResponse<DeclineFriendRequestResponse> friendAccept(
            @RequestBody DeclineFriendRRequest receiver, @RequestHeader("X-Profile-Id") String sender) {
        return new ApiResponse<>(friendsService.declineRequest(receiver, sender));
    }
}
