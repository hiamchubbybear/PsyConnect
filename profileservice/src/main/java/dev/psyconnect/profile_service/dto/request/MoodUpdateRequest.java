package dev.psyconnect.profile_service.dto.request;

import lombok.*;

@Setter
@Getter
@NoArgsConstructor
@Builder
@AllArgsConstructor
public class MoodUpdateRequest {
    String mood;
    String moodDescription;
    String visibility;
}
