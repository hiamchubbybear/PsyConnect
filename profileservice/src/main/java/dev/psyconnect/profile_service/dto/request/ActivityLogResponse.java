package dev.psyconnect.profile_service.dto.request;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class ActivityLogResponse {
    private String id;
    private String profileId;
    private String action;
    private long timestamp;
    private String targetId;
    private String targetType;
    private String extraData;
}
