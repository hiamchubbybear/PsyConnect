package dev.psyconnect.profile_service.dto.response;

import java.io.Serializable;

import dev.psyconnect.profile_service.model.Mood;
import dev.psyconnect.profile_service.model.Profile;
import dev.psyconnect.profile_service.model.Setting;
import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class ProfileWithRelationShipResponse implements Serializable {
    private Profile profile;
    private Mood mood;
}
