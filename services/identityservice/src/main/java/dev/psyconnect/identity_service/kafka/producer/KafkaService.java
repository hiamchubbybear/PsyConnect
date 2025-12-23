package dev.psyconnect.identity_service.kafka.producer;

import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.scheduling.annotation.Async;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.databind.ObjectMapper;

import dev.psyconnect.identity_service.dto.LogEvent;
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
            log.info("[KAFKA] Attempting to send message to topic: {}, payload size: {} bytes", topic, json.length());

            kafkaTemplate.send(topic, json).whenComplete((result, ex) -> {
                if (ex != null) {
                    log.error(
                            "[KAFKA] ❌ Send FAILED - Topic: {}, Error: {}, Cause: {}",
                            topic,
                            ex.getMessage(),
                            ex.getCause() != null ? ex.getCause().getMessage() : "N/A");
                } else {
                    log.info(
                            "[KAFKA] ✅ Send SUCCESS - Topic: {}, Partition: {}, Offset: {}",
                            topic,
                            result.getRecordMetadata().partition(),
                            result.getRecordMetadata().offset());
                }
            });

        } catch (Exception e) {
            log.error(
                    "[KAFKA] ❌ Serialization FAILED - Payload: {}, Error: {}, Stack: {}",
                    payload.getClass().getSimpleName(),
                    e.getMessage(),
                    e.getStackTrace()[0]);
        }
    }

    @Async
    public void sendLog(LogEvent event) {
        send("logging-service", event);
    }
}
