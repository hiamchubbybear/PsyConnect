package dev.psyconnect.profile_service.kafka.config;

import java.util.HashMap;
import java.util.Map;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.apache.kafka.common.serialization.StringSerializer;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.ConcurrentKafkaListenerContainerFactory;
import org.springframework.kafka.core.*;
import org.springframework.kafka.listener.DefaultErrorHandler;
import org.springframework.util.backoff.FixedBackOff;

import lombok.extern.slf4j.Slf4j;

@Slf4j
@Configuration
public class KafkaProducerConfig {
    @Value("${BOOTSTRAP_ADDRESS}")
    String bootstrapAddress;

    @Value("${spring.kafka.consumer.group-id:profile-service}")
    String consumerGroupId;

    // Producer Configuration (for sending messages)
    @Bean
    public ProducerFactory<String, Object> producerFactory() {
        Map<String, Object> configProps = new HashMap<>();
        configProps.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapAddress);
        configProps.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, StringSerializer.class);
        configProps.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, StringSerializer.class);
        return new DefaultKafkaProducerFactory<>(configProps);
    }

    @Bean
    public KafkaTemplate<String, Object> kafkaTemplate() {
        return new KafkaTemplate<>(producerFactory());
    }

    // Consumer Configuration (for receiving messages)
    @Bean
    public ConsumerFactory<String, String> consumerFactory() {
        return new DefaultKafkaConsumerFactory<>(Map.of(
                ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG,
                bootstrapAddress,
                ConsumerConfig.GROUP_ID_CONFIG,
                consumerGroupId,
                ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG,
                StringDeserializer.class,
                ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG,
                StringDeserializer.class));
    }

    @Bean
    public DefaultErrorHandler errorHandler() {
        FixedBackOff backOff = new FixedBackOff(1000L, 3);
        DefaultErrorHandler errorHandler = new DefaultErrorHandler(backOff);
        errorHandler.setRetryListeners((record, ex, deliveryAttempt) -> {
            log.error("Failed to consume record: {}, attempt {}, error: {}", record, deliveryAttempt, ex.getMessage());
        });

        return errorHandler;
    }

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, String> kafkaListenerContainerFactory() {
        ConcurrentKafkaListenerContainerFactory<String, String> factory =
                new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(consumerFactory());
        factory.setCommonErrorHandler(errorHandler());
        return factory;
    }
}
