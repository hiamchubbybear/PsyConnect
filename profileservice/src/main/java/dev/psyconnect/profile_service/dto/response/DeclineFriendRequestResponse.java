package dev.psyconnect.profile_service.dto.response;

import dev.psyconnect.profile_service.enums.FriendShipStatus;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
public class DeclineFriendRequestResponse {
    public String status;
    public FriendShipStatus acceptStatus;
}