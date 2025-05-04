package dev.psyconnect.profile_service.dto.response;

import dev.psyconnect.profile_service.model.Mood;
import dev.psyconnect.profile_service.model.Profile;
import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class ProfileWithMood {
    Profile profile;
    Mood mood;
}
