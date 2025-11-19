package dev.psyconnect.profile_service.dto.request;

import lombok.*;

@Builder
@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
public class MoodCreateRequest {
    String mood;
    String moodDescription;
    String visibility;
}
