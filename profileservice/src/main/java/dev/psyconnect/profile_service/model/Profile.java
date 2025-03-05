package dev.psyconnect.profile_service.model;

import java.io.Serializable;
import java.time.LocalDate;
import java.util.List;

import jakarta.validation.constraints.NotNull;

import org.springframework.data.annotation.Id;
import org.springframework.data.neo4j.core.schema.Node;
import org.springframework.data.neo4j.core.schema.Property;
import org.springframework.data.neo4j.core.schema.Relationship;

import com.fasterxml.jackson.annotation.JsonFormat;

import lombok.*;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@Setter
@Getter
@Node("user_profile")
public class Profile implements Serializable {
    private static final long serialVersionUID = 1L;
    private @Id String profileId;
    private @Property("accountId") @NotNull String accountId;
    private @Property("username") String username;
    private @Property("firstName") String firstName;
    private @Property("lastName") String lastName;
    private @Property("dob") @JsonFormat(pattern = "yyyy-MM-dd") LocalDate dob;
    private @Property("address") String address;
    private @Property("gender") String gender;
    private @Property("avatarUri") String avatarUri;

    @Builder.Default
    private @Property(value = "description") String description = "Welcome to my wall!!";

    private @Relationship(direction = Relationship.Direction.INCOMING, value = "HAS_MOOD") Mood moodList;
    private @Relationship(type = "FRIEND", direction = Relationship.Direction.OUTGOING) List<Profile> friends;
    private @Relationship(type = "HAS_ACTIVITY", direction = Relationship.Direction.OUTGOING) List<ActivityLog>
            activityLogs;
    private @Relationship(type = "HAS_SETTING", direction = Relationship.Direction.OUTGOING) Setting settings;
}
