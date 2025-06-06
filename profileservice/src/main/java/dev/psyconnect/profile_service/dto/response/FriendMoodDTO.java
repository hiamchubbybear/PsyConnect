package dev.psyconnect.profile_service.dto.response;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class FriendMoodDTO implements Serializable {
    private static final long serialVersionUID = 1L;
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
