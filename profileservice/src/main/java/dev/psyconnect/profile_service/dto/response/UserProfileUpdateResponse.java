package dev.psyconnect.profile_service.dto.response;

import java.util.UUID;

import org.springframework.format.annotation.DateTimeFormat;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.google.type.DateTime;

import lombok.*;

@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class UserProfileUpdateResponse {
    UUID userId;
    String username;
    String firstName;
    String lastName;

    @JsonFormat(pattern = "yyyy-MM-dd")
    @DateTimeFormat(pattern = "yyyy-MM-dd")
    DateTime dob;

    String address;
    String gender;
    String avatarUri;
}
