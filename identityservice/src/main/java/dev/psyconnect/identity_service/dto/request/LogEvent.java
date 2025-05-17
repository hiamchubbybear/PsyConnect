package dev.psyconnect.identity_service.dto.request;

import java.time.Instant;
import java.util.Map;

import lombok.*;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class LogEvent {
    private String eventId;
    private String service;
    private String action;
    private String level;
    private String userId;
    private String traceId;
    private Instant timestamp = Instant.now();
    private Map<String, Object> metadata;
    private String message;
}
