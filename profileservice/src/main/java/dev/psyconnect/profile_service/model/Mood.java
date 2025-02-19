package dev.psyconnect.profile_service.model;

import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.Node;

import lombok.*;

@Node("mood")
@Getter
@Builder
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class Mood {
    @Id
    private String moodId;

    private String mood;
    private String description;
    private String visibility;
    private long createdAt;
    private long expiresAt;
}
