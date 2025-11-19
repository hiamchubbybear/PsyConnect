package dev.psyconnect.profile_service.dto.request;

import lombok.*;

@Getter
@Setter
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class FriendAcceptationRequest {
    String target;
}
