package dev.psyconnect.profile_service.dto.response;

import java.io.Serializable;

import dev.psyconnect.profile_service.dto.request.ProfileSummaryDto;
import dev.psyconnect.profile_service.model.Mood;
import lombok.*;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ProfileWithMoodSummaryDto implements Serializable {
    private ProfileSummaryDto profile;
    private Mood mood;
}
