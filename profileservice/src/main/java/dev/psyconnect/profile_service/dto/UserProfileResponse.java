package dev.psyconnect.profile_service.dto;

import java.io.Serializable;
import java.util.UUID;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.google.type.DateTime;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.format.annotation.DateTimeFormat;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class UserProfileResponse implements Serializable {
    UUID userId;
    UUID profileId;
    String firstName;
    String lastName;
    @JsonFormat(pattern = "yyyy-MM-dd", shape = JsonFormat.Shape.STRING)
    @DateTimeFormat(pattern = "yyyy-MM-dd")
    String dob;
    String address;
    String gender;
    String role;
    String avatarUri;
}
