package dev.psyconnect.identity_service.dto.request;

import java.time.Instant;
import java.util.Map;

import dev.psyconnect.identity_service.dto.LogLevel;
import lombok.*;

import javax.print.attribute.standard.JobKOctets;

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