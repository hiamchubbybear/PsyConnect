package com.example.IdentityService.dto.response;
import java.time.LocalDate;
import java.util.UUID;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.format.annotation.DateTimeFormat;

@Setter
@Getter
@AllArgsConstructor
@NoArgsConstructor
public class UserCreationResponse {
    UUID userId;
    String firstName;
    String lastName;
    @JsonFormat(pattern = "yyyy-MM-dd")
    @DateTimeFormat(pattern = "yyyy-MM-dd")
    LocalDate dob;
    String address;
    String gender;
}
