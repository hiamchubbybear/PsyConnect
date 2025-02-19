package dev.psyconnect.profile_service.service;

import java.time.Instant;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.Map;

public class DateTimeTemplate {
    public static int MOOD_TIME_EXPIRES = 24;

    public static LocalDate parseFromString(String stringDateTime) {
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd");
        return LocalDate.parse(stringDateTime, formatter);
    }

    public static final Map<String, Long> MOOD_EXPIRES = Map.of(
            "currentTimeMillis",
            Instant.now().toEpochMilli(),
            "expiresTimeMillis",
            Instant.now().toEpochMilli());
}
