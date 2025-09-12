package dev.psyconnect.identity_service.service;

import java.sql.Timestamp;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.UUID;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.model.Token;
import dev.psyconnect.identity_service.repository.TokenRepository;
import lombok.RequiredArgsConstructor;

@RequiredArgsConstructor
@Service
public class TokenService {
    @Autowired
    private TokenRepository tokenRepository;

    public static int FRESH_TOKEN_TIME_EXPIRES = 7;

    public String generateRefreshToken(String username, String oldRefreshToken) {
        String newUUID = generateUUID();
        Instant newExpiry = Instant.now().plus(FRESH_TOKEN_TIME_EXPIRES, ChronoUnit.DAYS);

        Token responseToken = tokenRepository
                .findById(username)
                .map(existingToken -> {
                    if (!existingToken.getToken().equals(oldRefreshToken)) {
                        throw new CustomExceptionHandler(ErrorCode.TOKEN_INVALID);
                    }
                    if (existingToken.isRevoked()) {
                        throw new CustomExceptionHandler(ErrorCode.TOKEN_INVALID);
                    }
                    if (existingToken.getExpires().toInstant().isBefore(Instant.now())) {
                        throw new CustomExceptionHandler(ErrorCode.TOKEN_INVALID);
                    }
                    existingToken.setToken(newUUID);
                    existingToken.setExpires(Timestamp.from(newExpiry));
                    existingToken.setIssuedAt(Timestamp.from(Instant.now()));
                    return existingToken;
                })
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.TOKEN_INVALID));
        return tokenRepository.save(responseToken).getToken();
    }

    public String generateUUID() {
        return UUID.randomUUID().toString();
    }
}
