package dev.psyconnect.identity_service.kafka.producer;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;

import dev.psyconnect.identity_service.dto.LogEvent;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;

@Service
public class KafkaService {
    private static final Logger log = LoggerFactory.getLogger(KafkaService.class);
    private final KafkaTemplate<String, Object> kafkaTemplate;
    private static final ObjectMapper objectMapper = new ObjectMapper();

    @Autowired
    public KafkaService(KafkaTemplate<String, Object> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    @Async
    public void send(String topic, Object payload) {
        try {
            log.info("Send to " + topic + " values : " + payload);
            String objectMapper = new ObjectMapper().writeValueAsString(payload);
            kafkaTemplate.send(topic, objectMapper);
        } catch (Exception e) {
            // Log error but don't throw - allow request to complete successfully
            log.error("Kafka send error - Topic: {}, Payload: {}, Error: {}", topic, payload, e.getMessage(), e);
            log.warn("Continuing request despite Kafka error");
        }
    }

    @Async
    public void sendLog(LogEvent log) {
        try {
            send("logging-service", log);
        } catch (Exception e) {
            // Silently ignore logging errors
            this.log.warn("Failed to send log to Kafka: {}", e.getMessage());
        }
    }

    public static <T> T objectMapping(String rawString, Class<T> clazz) {
        try {
            return objectMapper.readValue(rawString, clazz);
        } catch (JsonMappingException e) {
            throw new CustomExceptionHandler(ErrorCode.JSON_MAPPING_ERROR);
        } catch (JsonProcessingException e) {
            throw new CustomExceptionHandler(ErrorCode.JSON_PROCESSING_ERROR);
        }
    }
}
