package dev.psyconnect.kafka.template;

import lombok.RequiredArgsConstructor;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Component;

@Component
public class KafkaProducerTemplate {
    private final KafkaTemplate<String, Object> kafkaTemplate;

    public KafkaProducerTemplate(KafkaTemplate<String, Object> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void sendMessage(String topic, Object message) {
        kafkaTemplate.send(topic, message);
    }
}
