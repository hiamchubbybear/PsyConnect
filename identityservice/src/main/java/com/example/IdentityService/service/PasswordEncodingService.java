package com.example.IdentityService.service;

import org.springframework.context.annotation.Bean;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

@Service
public class PasswordEncodingService {
    String SIGNER_KEY = "ZW5jb2RlX3Bhc3NrZXlfYW5fbG9uX2E=";

    public static String encoder(String rawPassword) {
        return new BCryptPasswordEncoder(10).encode(rawPassword);
    }
}
