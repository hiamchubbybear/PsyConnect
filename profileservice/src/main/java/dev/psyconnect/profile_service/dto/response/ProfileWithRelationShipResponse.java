package dev.psyconnect.profile_service.dto.response;


import dev.psyconnect.profile_service.model.Mood;
import dev.psyconnect.profile_service.model.Profile;
import dev.psyconnect.profile_service.model.Setting;
import lombok.*;

import java.io.Serializable;
import java.util.List;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class ProfileWithRelationShipResponse implements Serializable {
    private Profile profile;
    private Mood mood;
    private Setting setting;
//    private List<Profile> friends;
}

