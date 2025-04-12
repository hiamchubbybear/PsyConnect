package dev.psyconnect.profile_service.dto.response;

import dev.psyconnect.profile_service.enums.FriendShipStatus;
import lombok.*;

@Builder
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class UndoFriendRequestResponse {
    public String status;
    public FriendShipStatus undoStatus;
}
