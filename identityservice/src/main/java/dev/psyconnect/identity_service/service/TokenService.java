package dev.psyconnect.identity_service.service;

import java.sql.Timestamp;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.UUID;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import dev.psyconnect.identity_service.model.Token;
import dev.psyconnect.identity_service.repository.TokenRepository;
import lombok.RequiredArgsConstructor;

@RequiredArgsConstructor
@Service
public class TokenService {
    @Autowired
    private TokenRepository tokenRepository;

    public static int FRESH_TOKEN_TIME_EXPIRES = 7;

    public String generateRefreshToken(String username) {
        String generatedUUID = generateUUID();
        Instant currentTime = Instant.now().plus(FRESH_TOKEN_TIME_EXPIRES, ChronoUnit.DAYS);
        Token responseToken = Token.builder()
                .username(username)
                .token(generatedUUID)
                .expires(Timestamp.from(currentTime))
                .issuedAt(Timestamp.from(currentTime))
                .build();

        return tokenRepository.save(responseToken).getToken();
    }

    public String generateUUID() {
        return UUID.randomUUID().toString();
    }
}
