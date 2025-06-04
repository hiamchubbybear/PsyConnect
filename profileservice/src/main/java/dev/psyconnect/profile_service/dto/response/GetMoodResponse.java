package dev.psyconnect.profile_service.dto.response;

import java.io.Serializable;

import org.springframework.data.neo4j.core.schema.Id;

import lombok.*;

@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
@Builder
@ToString
public class GetMoodResponse implements Serializable {
    private String profileId;
    private String fullName;
    private String avatarUrl;
    private String moodId;
    private String mood;
    private String description;
    private String visibility;
    private long createdAt;
    private long expiresAt;
}
