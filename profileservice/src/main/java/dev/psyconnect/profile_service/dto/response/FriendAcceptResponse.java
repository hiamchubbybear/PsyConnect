package dev.psyconnect.profile_service.dto.response;

import dev.psyconnect.profile_service.enums.FriendShipStatus;
import lombok.*;

@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
public class FriendAcceptResponse {
    public String status;
    public  FriendShipStatus acceptStatus;
}
