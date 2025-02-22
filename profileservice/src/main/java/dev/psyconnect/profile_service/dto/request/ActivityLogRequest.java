package dev.psyconnect.profile_service.dto.request;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class ActivityLogRequest {
    private String action;
    private String targetId;
    private String targetType;
    private String extraData;
}

