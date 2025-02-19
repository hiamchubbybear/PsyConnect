package dev.psyconnect.profile_service.model;

import org.springframework.data.annotation.Id;
import org.springframework.data.neo4j.core.schema.GeneratedValue;
import org.springframework.data.neo4j.core.schema.Node;
import org.springframework.data.neo4j.core.support.UUIDStringGenerator;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Node
public class ActivityLog {
    @Id
    @GeneratedValue(generatorClass = UUIDStringGenerator.class)
    private String id;

    private String profileId;
    private String action;
    private long timestamp;
}
