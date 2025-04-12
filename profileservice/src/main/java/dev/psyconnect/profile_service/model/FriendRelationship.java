package dev.psyconnect.profile_service.model;

import dev.psyconnect.profile_service.enums.FriendShipStatus;
import dev.psyconnect.profile_service.model.Profile;
import dev.psyconnect.profile_service.service.UserProfileService;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
import lombok.Builder;
import lombok.Getter;
import lombok.Setter;
import org.springframework.data.neo4j.core.schema.GeneratedValue;
import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.RelationshipProperties;
import org.springframework.data.neo4j.core.schema.TargetNode;

import java.io.Serializable;
import java.time.OffsetDateTime;

@RelationshipProperties
@Builder
@Setter
@Getter
public class FriendRelationship  implements Serializable {
    @Id
    @GeneratedValue
    private String id;

    @TargetNode
    private Profile target;

    @Enumerated(EnumType.STRING)
    private FriendShipStatus status;

    private OffsetDateTime createdAt;
}
