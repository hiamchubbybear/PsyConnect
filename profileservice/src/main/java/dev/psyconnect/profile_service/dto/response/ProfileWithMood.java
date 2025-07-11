package dev.psyconnect.profile_service.dto.response;

import java.io.Serializable;

import dev.psyconnect.profile_service.model.Mood;
import dev.psyconnect.profile_service.model.Profile;
import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class ProfileWithMood implements Serializable {
    private static final long serialVersionUID = 1L;
    Profile profile;
    Mood mood;
}
