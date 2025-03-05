package dev.psyconnect.kafka.template;

import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

@Component
@Slf4j
public class KafkaConsumerTemplate {

    @KafkaListener(topics = "my-topic", groupId = "group-id")
    public void listen(String message) {
        log.info("Received Message: {}", message);
    }
}

