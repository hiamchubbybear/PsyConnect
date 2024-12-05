package com.example.profileservice.dto.respone;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.time.LocalDate;

@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
public class UserProfileCreationRespone {
    String userId;
    String firstName;
    String lastName;
    LocalDate dob;
    String address;
    String gender;
}
