package dev.psyconnect.profile_service.dto.response;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class MoodProfileResponse {
    String avatarUrl;
    String fullName;
}
