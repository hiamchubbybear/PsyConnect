package dev.psyconnect.profile_service.model;

import java.io.Serializable;
import java.time.OffsetDateTime;

import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;

import org.springframework.data.neo4j.core.schema.GeneratedValue;
import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.RelationshipProperties;
import org.springframework.data.neo4j.core.schema.TargetNode;

import dev.psyconnect.profile_service.enums.FriendShipStatus;
import lombok.Builder;
import lombok.Getter;
import lombok.Setter;

@RelationshipProperties
@Builder
@Setter
@Getter
public class FriendRelationship implements Serializable {
    @Id
    @GeneratedValue
    private String id;

    @TargetNode
    private Profile target;

    @Enumerated(EnumType.STRING)
    private FriendShipStatus status;

    private OffsetDateTime createdAt;
}
