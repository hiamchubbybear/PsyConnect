package dev.psyconnect.identity_service.dto.request;

import java.util.UUID;

import lombok.*;

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
