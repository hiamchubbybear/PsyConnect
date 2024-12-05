package com.example.profileservice.model;

import lombok.*;
import org.springframework.data.neo4j.core.schema.GeneratedValue;
import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.Node;
import org.springframework.data.neo4j.core.schema.Property;
import org.springframework.data.neo4j.core.support.UUIDStringGenerator;

import java.time.LocalDate;

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
    @Property("firstName")
    String firstName;
    @Property("lastName")
    String lastName;
    @Property("dob")
    LocalDate dob;
    @Property("address")
    String address;
    @Property("gender")
    String gender;
}
