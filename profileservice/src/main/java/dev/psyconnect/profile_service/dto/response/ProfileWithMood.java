package dev.psyconnect.profile_service.dto.response;

import dev.psyconnect.profile_service.model.Mood;
import dev.psyconnect.profile_service.model.Profile;
import lombok.*;
import org.springframework.context.annotation.Bean;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class ProfileWithMood {
    Profile profile;
    Mood mood;
}
