package dev.psyconnect.kafka.service;

import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.support.KafkaHeaders;
import org.springframework.messaging.handler.annotation.Header;
import org.springframework.messaging.handler.annotation.Payload;

public class KafkaConsumerService {
    @KafkaListener(topics = "identity.create")
    public void listenWithHeaders(@Payload String message,

                                  @Header(KafkaHeaders.RECEIVED_PARTITION) int partition) {
        System.out.println("Received Message: " + message + "from partition: " + partition);
    }

    @KafkaListener(topics = "kafka_identity_service", groupId = "identity-kafka-server")
    public void listenGroupFoo(String message) {
        System.out.println("Received Message in group foo: " + message);
    }
}
