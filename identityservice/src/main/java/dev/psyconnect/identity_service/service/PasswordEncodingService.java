package dev.psyconnect.identity_service.service;

import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

@Service
public class PasswordEncodingService {
    String SIGNER_KEY = "ZW5jb2RlX3Bhc3NrZXlfYW5fbG9uX2E=";

    public static String encoder(String rawPassword) {
        return new BCryptPasswordEncoder(10).encode(rawPassword);
    }

    public static BCryptPasswordEncoder getBCryptPasswordEncoder() {
        return new BCryptPasswordEncoder();
    }
}
