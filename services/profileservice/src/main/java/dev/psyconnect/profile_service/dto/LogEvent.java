package dev.psyconnect.profile_service.dto;

import java.util.Map;

import com.fasterxml.jackson.annotation.JsonInclude;

import lombok.*;

@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
@JsonInclude(JsonInclude.Include.NON_NULL)
public class LogEvent {
    private String timestamp;
    private String level;
    private String service;
    private String message;
    private String traceId;
    private String userId;
    private String action;
    private Map<String, Object> metadata;
    private String error;
    private String stackTrace;
    private String method;
    private String path;
    private Integer statusCode;
    private Long duration;
    private String ip;
    private String userAgent;
    private String environment;
    private String version;
    private String hostname;
}
