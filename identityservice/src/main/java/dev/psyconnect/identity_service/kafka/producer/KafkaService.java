package dev.psyconnect.identity_service.kafka.producer;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.core.KafkaTemplate;
import org.springframework.stereotype.Service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;

import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;

@Service
public class KafkaService {
    private final KafkaTemplate<String, Object> kafkaTemplate;
    private static final ObjectMapper objectMapper = new ObjectMapper();

    @Autowired
    public KafkaService(KafkaTemplate<String, Object> kafkaTemplate) {
        this.kafkaTemplate = kafkaTemplate;
    }

    public void send(String topic, Object payload ) {
        try {
            String objectMapper = new ObjectMapper().writeValueAsString(payload);
            kafkaTemplate.send(topic, objectMapper);
        } catch (Exception e) {
            throw new CustomExceptionHandler(ErrorCode.KAFKA_SERVER_ERROR);
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
