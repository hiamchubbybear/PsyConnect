package com.example.IdentityService.configuration;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.fasterxml.jackson.databind.ObjectMapper;

import feign.codec.Encoder;
import feign.jackson.JacksonEncoder;

@Configuration
public class FeignConfig {
    @Bean
    public Encoder feignEncoder() {
        return new JacksonEncoder(new ObjectMapper());
    }
}
