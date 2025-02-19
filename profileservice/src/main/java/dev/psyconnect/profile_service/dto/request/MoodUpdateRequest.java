package dev.psyconnect.profile_service.dto.request;

import lombok.*;

@Setter
@Getter
@NoArgsConstructor
@Builder
@AllArgsConstructor
public class MoodUpdateRequest {
    private String id;
    private String title;
    private String description;
}
