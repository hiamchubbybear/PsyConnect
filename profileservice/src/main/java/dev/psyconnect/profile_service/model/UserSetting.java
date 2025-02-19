package dev.psyconnect.profile_service.model;

import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.Node;

@Node
public class UserSetting {
    @Id
    String profileId;

    private String privacyLevel;
    private boolean notificationsEnabled;
}
