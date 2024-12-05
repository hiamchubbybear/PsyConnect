package com.example.IdentityService.dto.request;

import lombok.*;
import lombok.experimental.FieldDefaults;

import java.time.LocalDate;
import java.util.UUID;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE)
public class UserProfileCreationRequest {
    String username;
    String userId;
    String firstName;
    String lastName;
    LocalDate dob;
    String address;
    String gender;
    String email;
    String role;
}