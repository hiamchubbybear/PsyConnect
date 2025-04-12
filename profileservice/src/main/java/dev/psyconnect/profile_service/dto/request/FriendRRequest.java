package dev.psyconnect.profile_service.dto.request;

import lombok.*;

@Builder
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class FriendRRequest {
    private  String target;

}
