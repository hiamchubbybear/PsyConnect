package com.example.IdentityService.configuration;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import feign.form.FormEncoder;

@Configuration
public class FeignConfiguration {

    @Bean
    public FormEncoder feignEncoder() {
        return new FormEncoder();
    }
}
