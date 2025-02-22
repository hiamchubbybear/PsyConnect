package dev.psyconnect.profile_service.dto.response;

import lombok.*;
import org.springframework.data.neo4j.core.schema.Id;

import java.io.Serializable;

@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class GetMoodResponse implements Serializable {
    @Id
    private String moodId;
    private String mood;
    private String description;
    private String visibility;
    private long createdAt;
    private long expiresAt;
}
