package dev.psyconnect.profile_service.service;

import java.net.InetAddress;
import java.time.Instant;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.CompletableFuture;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.kafka.support.SendResult;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;

import dev.psyconnect.profile_service.dto.LogEvent;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Service
public class KafkaLoggerService {

    private final KafkaTemplate<String, String> kafkaTemplate;
    private final ObjectMapper objectMapper;
    private final String serviceName;
    private final String environment;
    private final String version;
    private final String hostname;
    private final String topic;

    public KafkaLoggerService(
            KafkaTemplate<String, String> kafkaTemplate,
            ObjectMapper objectMapper,
            @Value("${spring.application.name:profile-service}") String serviceName,
            @Value("${app.environment:development}") String environment,
            @Value("${app.version:1.0.0}") String version,
            @Value("${kafka.topic.logging:logging-service}") String topic) {
        this.kafkaTemplate = kafkaTemplate;
        this.objectMapper = objectMapper;
        this.serviceName = serviceName;
        this.environment = environment;
        this.version = version;
        this.topic = topic;
        this.hostname = getHostname();
    }

    private String getHostname() {
        try {
            return InetAddress.getLocalHost().getHostName();
        } catch (Exception e) {
            return "unknown";
        }
    }

    public void log(String level, String message, Map<String, Object> fields) {
        try {
            LogEvent event = LogEvent.builder()
                    .timestamp(Instant.now().toString())
                    .level(level)
                    .service(serviceName)
                    .message(message)
                    .environment(environment)
                    .version(version)
                    .hostname(hostname)
                    .build();

            if (fields != null) {
                // Extract known fields
                event.setTraceId((String) fields.get("traceId"));
                event.setUserId((String) fields.get("userId"));
                event.setAction((String) fields.get("action"));
                event.setError((String) fields.get("error"));
                event.setStackTrace((String) fields.get("stackTrace"));
                event.setMethod((String) fields.get("method"));
                event.setPath((String) fields.get("path"));
                event.setStatusCode((Integer) fields.get("statusCode"));
                event.setDuration((Long) fields.get("duration"));
                event.setIp((String) fields.get("ip"));
                event.setUserAgent((String) fields.get("userAgent"));

                // Remaining fields go to metadata
                Map<String, Object> metadata = new HashMap<>(fields);
                metadata.remove("traceId");
                metadata.remove("userId");
                metadata.remove("action");
                metadata.remove("error");
                metadata.remove("stackTrace");
                metadata.remove("method");
                metadata.remove("path");
                metadata.remove("statusCode");
                metadata.remove("duration");
                metadata.remove("ip");
                metadata.remove("userAgent");

                if (!metadata.isEmpty()) {
                    event.setMetadata(metadata);
                }
            }

            String json = objectMapper.writeValueAsString(event);

            // Send to Kafka asynchronously
            CompletableFuture<SendResult<String, String>> future = kafkaTemplate.send(topic, json);

            future.whenComplete((result, ex) -> {
                if (ex != null) {
                    log.error("Failed to send log to Kafka: {}", ex.getMessage());
                }
            });

        } catch (Exception e) {
            log.error("Failed to create log event: {}", e.getMessage());
        }
    }

    public void debug(String message, Map<String, Object> fields) {
        log("DEBUG", message, fields);
    }

    public void info(String message, Map<String, Object> fields) {
        log("INFO", message, fields);
    }

    public void warn(String message, Map<String, Object> fields) {
        log("WARN", message, fields);
    }

    public void error(String message, Map<String, Object> fields) {
        log("ERROR", message, fields);
    }

    public void fatal(String message, Map<String, Object> fields) {
        log("FATAL", message, fields);
    }

    public void audit(String message, Map<String, Object> fields) {
        log("AUDIT", message, fields);
    }

    public void logHttpRequest(String method, String path, int statusCode, long duration, Map<String, Object> fields) {
        if (fields == null) {
            fields = new HashMap<>();
        }
        fields.put("method", method);
        fields.put("path", path);
        fields.put("statusCode", statusCode);
        fields.put("duration", duration);

        String level = "INFO";
        if (statusCode >= 500) {
            level = "ERROR";
        } else if (statusCode >= 400) {
            level = "WARN";
        }

        String message = String.format("%s %s - %d (%dms)", method, path, statusCode, duration);
        log(level, message, fields);
    }

    public void logAction(String action, String userId, String message, Map<String, Object> fields) {
        if (fields == null) {
            fields = new HashMap<>();
        }
        fields.put("action", action);
        fields.put("userId", userId);

        log("AUDIT", message, fields);
    }
}
