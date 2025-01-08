package com.example.profileservice.model;

import java.time.LocalDate;

import org.springframework.data.neo4j.core.schema.GeneratedValue;
import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.Node;
import org.springframework.data.neo4j.core.schema.Property;
import org.springframework.data.neo4j.core.support.UUIDStringGenerator;

import com.fasterxml.jackson.annotation.JsonFormat;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Node("user_profile")
public class UserProfile {
    @Id
    @GeneratedValue(UUIDStringGenerator.class)
    String profileId;

    @Property("userId")
    String userId;

    @Property("username")
    String username;

    @Property("firstName")
    String firstName;

    @Property("lastName")
    String lastName;

    @Property("dob")
    @JsonFormat(pattern = "yyyy-MM-dd")
    LocalDate dob;

    @Property("address")
    String address;

    @Property("gender")
    String gender;

    @Property("email")
    String email;

    @Property("avatarUri")
    String avatarUri;
}
