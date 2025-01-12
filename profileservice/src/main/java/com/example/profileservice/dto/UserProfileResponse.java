package com.example.profileservice.dto;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.neo4j.core.schema.Property;

import java.time.LocalDate;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class UserProfileResponse {
    String userId;
    String username;
    String firstName;
    String lastName;
    @JsonFormat(pattern = "yyyy-MM-dd")
    LocalDate dob;
    String address;
    String gender;
    String email;
    String avatarUri;
}
