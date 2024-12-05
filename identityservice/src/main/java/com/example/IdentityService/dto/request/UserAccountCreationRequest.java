package com.example.IdentityService.dto.request;

import lombok.*;
import lombok.experimental.FieldDefaults;

import java.time.LocalDate;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
@FieldDefaults(level = AccessLevel.PRIVATE)
public class UserAccountCreationRequest {
    String username;
    String password;
    String firstName;
    String lastName;
    LocalDate dob;
    String address;
    String gender;
    String email;
    String role;
}
