package dev.psyconnect.profile_service.dto.request;

import org.springframework.format.annotation.DateTimeFormat;

import com.fasterxml.jackson.annotation.JsonFormat;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class UserProfileCreationRequest {
    String accountId;
    String profileId;
    String firstName;
    String lastName;
    String dob;
    String address;
    String gender;
    String role;
    String avatarUri;
}
