package dev.psyconnect.profile_service.dto.response;

import java.io.Serializable;

import dev.psyconnect.profile_service.dto.request.ProfileSummaryDto;
import dev.psyconnect.profile_service.model.Mood;
import dev.psyconnect.profile_service.model.Profile;
import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class ProfileWithRelationShipResponse implements Serializable {
    private Profile profile;
    private Mood mood;

    public ProfileWithMoodSummaryDto toSummaryWithMood() {
        return ProfileWithMoodSummaryDto.builder()
                .profile(ProfileSummaryDto.builder()
                        .profileId(profile.getProfileId())
                        .username(profile.getUsername())
                        .avatarUri(profile.getAvatarUri())
                        .build())
                .mood(mood)
                .build();
    }

}
