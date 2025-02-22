package dev.psyconnect.profile_service.service;

import java.time.Instant;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.time.temporal.ChronoUnit;
import java.util.Map;

public class Time {
    public static int MOOD_TIME_EXPIRES = 24;

    public static final Map<String, Long> MOOD_EXPIRES = Map.of(
            "currentTimeMillis",
            Instant.now().toEpochMilli(),
            "expiresTimeMillis",
            Instant.now().plus(24, ChronoUnit.HOURS).toEpochMilli());

    public static LocalDate parseFromString(String stringDateTime) {
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd");
        return LocalDate.parse(stringDateTime, formatter);
    }

    public static LocalDateTime fromTimeStampToLCD(int timeZone, long timestamp) {
        return LocalDateTime.ofEpochSecond(timestamp / 1000, 0, ZoneOffset.ofHours(timeZone));
    }
}
