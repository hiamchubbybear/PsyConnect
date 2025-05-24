package dev.psyconnect.profile_service.dto.request;

import java.time.Instant;
import java.util.Map;


import lombok.*;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class LogEvent {
    private String service;
    private String action;
    private LogLevel level;
    private String userId;
    private String traceId;
    private String timestamp;
    private Map<String, ?> metadata;
    private String message;
}
