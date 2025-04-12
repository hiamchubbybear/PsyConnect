package dev.psyconnect.profile_service.dto.response;

import java.io.Serializable;
import java.util.List;

import dev.psyconnect.profile_service.model.ActivityLog;
import dev.psyconnect.profile_service.model.FriendRelationship;
import dev.psyconnect.profile_service.model.Mood;
import dev.psyconnect.profile_service.model.Setting;
import org.springframework.data.neo4j.core.schema.Relationship;
import org.springframework.format.annotation.DateTimeFormat;

import com.fasterxml.jackson.annotation.JsonFormat;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class UserProfileResponse implements Serializable {
    String accountId;
    String profileId;
    String firstName;
    String lastName;
    @JsonFormat(pattern = "yyyy-MM-dd", shape = JsonFormat.Shape.STRING)
    @DateTimeFormat(pattern = "yyyy-MM-dd")
    String dob;
    String address;
    String gender;
    String role;
    String avatarUri;
}
