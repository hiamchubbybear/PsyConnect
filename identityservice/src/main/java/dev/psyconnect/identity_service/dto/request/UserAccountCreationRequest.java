package dev.psyconnect.identity_service.dto.request;

import org.springframework.format.annotation.DateTimeFormat;

import com.fasterxml.jackson.annotation.JsonFormat;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class UserAccountCreationRequest {
    String accountId;
    String profileId;
    String username;
    String password;
    String firstName;
    String lastName;
    String oauth2Session;

    @JsonFormat(pattern = "yyyy-MM-dd")
    @DateTimeFormat(pattern = "yyyy-MM-dd")
    String dob;

    String address;
    String gender;
    String email;
    String avatarUri;
    String role;
}
