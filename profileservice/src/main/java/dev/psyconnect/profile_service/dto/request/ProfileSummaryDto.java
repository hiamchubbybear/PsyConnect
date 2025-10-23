package dev.psyconnect.profile_service.dto.request;

import java.io.Serializable;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ProfileSummaryDto implements Serializable {
    private static final long serialVersionUID = 1L;
    private String profileId;
    private String avatarUri;
}
