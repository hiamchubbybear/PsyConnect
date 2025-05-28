package dev.psyconnect.api_service.dto;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.util.Map;

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