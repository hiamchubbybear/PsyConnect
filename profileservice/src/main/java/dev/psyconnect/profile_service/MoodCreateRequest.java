package dev.psyconnect.profile_service;

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
