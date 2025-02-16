package dev.psyconnect.profile_service.service;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;

public class DateTimeFormatting {
    public static LocalDate parseFromString(String stringDateTime) {
        DateTimeFormatter formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd");
        return LocalDate.parse(stringDateTime, formatter);
    }
}
