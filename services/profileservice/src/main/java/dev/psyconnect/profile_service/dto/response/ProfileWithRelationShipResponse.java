package dev.psyconnect.profile_service.dto.response;

import java.io.Serializable;
import java.util.List;

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
    private List<Mood> moods;

    public ProfileWithMoodSummaryDto toSummaryWithMood() {
        return ProfileWithMoodSummaryDto.builder()
                .profile(ProfileSummaryDto.builder()
                        .profileId(profile.getProfileId())
                        .avatarUri(profile.getAvatarUri())
                        .build())
                .mood(moods != null && !moods.isEmpty() ? moods.get(0) : null)
                .build();
    }
}
