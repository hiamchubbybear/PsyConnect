package com.example.IdentityService.Service;

import org.springframework.context.annotation.Bean;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import java.security.SecureRandom;
@Service
public class PasswordEncodingService {
    String SIGNER_KEY = "Jh6OEn0FD9bsClIQBqVU9889ElikEJnC";
    @Bean
    public PasswordEncoder encoder() {
        return new BCryptPasswordEncoder(10 );
    }
}
