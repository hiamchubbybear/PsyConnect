package dev.psyconnect.identity_service.dto.response;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
public class UserProfileCreationResponse {
    String profileId;
    String firstName;
    String lastName;
    String dob;
    String address;
    String gender;
    String avatarUri;
}
