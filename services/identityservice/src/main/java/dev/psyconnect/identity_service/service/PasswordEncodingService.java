package dev.psyconnect.identity_service.service;

import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

@Service
public class PasswordEncodingService {

    public static String encoder(String rawPassword) {
        return new BCryptPasswordEncoder(10).encode(rawPassword);
    }

    public static BCryptPasswordEncoder getBCryptPasswordEncoder() {
        return new BCryptPasswordEncoder();
    }
}
