package dev.psyconnect.identity_service.controller;

import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.springframework.kafka.annotation.KafkaListener;

public class KafkaIdentityController {
    @KafkaListener(topics = "identity.user-created")
    public void listen(ConsumerRecord<String, String> record) {

    }
}
