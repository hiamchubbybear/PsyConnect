package dev.psyconnect.kafka.controller;

import dev.psyconnect.kafka.practice.PracticalAdvice;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.apache.kafka.common.header.Headers;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.messaging.handler.annotation.Payload;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.stream.IntStream;
import java.util.stream.StreamSupport;
@Slf4j
public class KafkaController {

    private static final Logger logger =
            LoggerFactory.getLogger(KafkaController.class);

    private final KafkaTemplate<String, Object> template;
    private final String topicName;
    private final int messagesPerRequest;
    private CountDownLatch latch;

    public KafkaController(
            final KafkaTemplate<String, Object> template,
            @Value("${tpd.topic-name}") final String topicName,
            @Value("${tpd.messages-per-request}") final int messagesPerRequest) {
        this.template = template;
        this.topicName = topicName;
        this.messagesPerRequest = messagesPerRequest;
    }


    @KafkaListener(topics = "advice-topic", clientIdPrefix = "json",
            containerFactory = "kafkaListenerContainerFactory")
    public void listenAsObject(ConsumerRecord<String, PracticalAdvice> cr,
                               @Payload PracticalAdvice payload) {
        logger.info("Logger 1 [JSON] received key {}: Type [{}] | Payload: {} | Record: {}", cr.key(),
                typeIdHeader(cr.headers()), payload, cr.toString());
        latch.countDown();
    }

    @KafkaListener(topics = "advice-topic", clientIdPrefix = "string",
            containerFactory = "kafkaListenerStringContainerFactory")
    public void listenasString(ConsumerRecord<String, String> cr,
                               @Payload String payload) {
        logger.info("Logger 2 [String] received key {}: Type [{}] | Payload: {} | Record: {}", cr.key(),
                typeIdHeader(cr.headers()), payload, cr.toString());
        latch.countDown();
    }

    @KafkaListener(topics = "advice-topic", clientIdPrefix = "bytearray",
            containerFactory = "kafkaListenerByteArrayContainerFactory")
    public void listenAsByteArray(ConsumerRecord<String, byte[]> cr,
                                  @Payload byte[] payload) {
        logger.info("Logger 3 [ByteArray] received key {}: Type [{}] | Payload: {} | Record: {}", cr.key(),
                typeIdHeader(cr.headers()), payload, cr.toString());
        latch.countDown();
    }

    private static String typeIdHeader(Headers headers) {
        return StreamSupport.stream(headers.spliterator(), false)
                .filter(header -> header.key().equals("__TypeId__"))
                .findFirst().map(header -> new String(header.value())).orElse("N/A");
    }
}
