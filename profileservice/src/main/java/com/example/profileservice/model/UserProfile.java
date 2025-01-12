package com.example.profileservice.model;

import com.fasterxml.jackson.annotation.JsonFormat;
import jakarta.validation.constraints.NotNull;
import lombok.*;
import org.springframework.data.annotation.Id;
import org.springframework.data.annotation.Transient;
import org.springframework.data.neo4j.core.schema.GeneratedValue;
import org.springframework.data.neo4j.core.schema.Node;
import org.springframework.data.neo4j.core.schema.Property;
import org.springframework.data.neo4j.core.support.UUIDStringGenerator;

import java.time.LocalDate;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor()
@Node("user_profile")
public class UserProfile {

    @Id
    @GeneratedValue(UUIDStringGenerator.class)
    private String profileId;

    @Property("userId")
    @NotNull
    private String userId;

    @Property("username")
    private String username;

    @Property("firstName")
    private String firstName;

    @Property("lastName")
    private String lastName;

    @Property("dob")
    @JsonFormat(pattern = "yyyy-MM-dd")
    private LocalDate dob;

    @Property("address")
    private String address;

    @Property("gender")
    private String gender;

    @Property("avatarUri")
    private String avatarUri;

    @Transient
    private String elementId;
}