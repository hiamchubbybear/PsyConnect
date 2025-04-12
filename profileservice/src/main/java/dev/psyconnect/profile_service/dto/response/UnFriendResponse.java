package dev.psyconnect.profile_service.dto.response;

import dev.psyconnect.profile_service.enums.FriendShipStatus;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.NoArgsConstructor;

@AllArgsConstructor@NoArgsConstructor
@Builder
public class UnFriendResponse {
    public String status;
    public FriendShipStatus unfriendStatus;
}
