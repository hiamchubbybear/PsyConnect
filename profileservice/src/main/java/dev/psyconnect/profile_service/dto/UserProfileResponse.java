package dev.psyconnect.profile_service.dto;

import java.io.Serializable;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.google.type.DateTime;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class UserProfileResponse implements Serializable {
    String userId;
    String username;
    String firstName;
    String lastName;

    @JsonFormat(pattern = "yyyy-MM-dd")
    DateTime dob;

    String address;
    String gender;
    String email;
    String avatarUri;
}
