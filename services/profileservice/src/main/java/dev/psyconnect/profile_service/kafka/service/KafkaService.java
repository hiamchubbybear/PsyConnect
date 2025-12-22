package dev.psyconnect.profile_service.kafka.service;

import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;

import dev.psyconnect.profile_service.dto.request.LogEvent;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Service
public class KafkaService {

    private final KafkaTemplate<String, String> kafkaTemplate;
    private static final ObjectMapper objectMapper = new ObjectMapper();

    public KafkaService(KafkaTemplate<String, String> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    @Async
    public void send(String topic, Object payload) {
        try {
            String json = objectMapper.writeValueAsString(payload);

            kafkaTemplate.send(topic, json).whenComplete((result, ex) -> {
                if (ex != null) {
                    log.error("Kafka send error - Topic: {}, Error: {}", topic, ex.getMessage());
                }
            });

        } catch (Exception e) {
            log.error("Kafka serialization error - Payload: {}, Error: {}", payload, e.getMessage());
        }
    }

    @Async
    public void sendLog(LogEvent event) {
        send("logging-service", event);
    }
}
