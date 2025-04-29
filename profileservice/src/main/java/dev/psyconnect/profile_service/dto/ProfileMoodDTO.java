package dev.psyconnect.profile_service.dto;

import dev.psyconnect.profile_service.model.Mood;
import dev.psyconnect.profile_service.model.Profile;
import lombok.*;

@Getter
@AllArgsConstructor
@NoArgsConstructor
@Setter
@Builder
public class ProfileMoodDTO {
    Profile profile;
    Mood mood;
}
