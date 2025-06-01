package dev.psyconnect.profile_service.dto.response;

import dev.psyconnect.profile_service.dto.request.ProfileSummaryDto;
import dev.psyconnect.profile_service.model.Mood;
import lombok.*;

import java.io.Serializable;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ProfileWithMoodSummaryDto implements Serializable {
    private ProfileSummaryDto profile;
    private Mood mood;
}
