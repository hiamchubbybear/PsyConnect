package dev.psyconnect.profile_service.dto.request;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class ProfileSummaryDto implements Serializable {
    private static final long serialVersionUID = 1L;
    private String profileId;
    private String username;
    private String avatarUri;
}

