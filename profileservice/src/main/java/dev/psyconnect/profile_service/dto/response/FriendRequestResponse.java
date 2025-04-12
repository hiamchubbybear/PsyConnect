package dev.psyconnect.profile_service.dto.response;

import lombok.*;

@Builder
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class FriendRequestResponse {
    String status ;
    String friendRequestStatus;
}
