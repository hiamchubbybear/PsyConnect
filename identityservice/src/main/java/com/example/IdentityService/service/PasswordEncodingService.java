package com.example.IdentityService.service;

import org.springframework.context.annotation.Bean;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

@Service
public class PasswordEncodingService {
    @Bean
    public PasswordEncoder encoder() {
        return new BCryptPasswordEncoder(10 );
    }
}
