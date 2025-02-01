package com.example.profileservice.model;

import java.time.LocalDate;

import jakarta.persistence.Lob;
import jakarta.validation.constraints.NotNull;

import org.springframework.data.annotation.Id;
import org.springframework.data.annotation.Transient;
import org.springframework.data.neo4j.core.schema.GeneratedValue;
import org.springframework.data.neo4j.core.schema.Node;
import org.springframework.data.neo4j.core.schema.Property;
import org.springframework.data.neo4j.core.support.UUIDStringGenerator;

import com.fasterxml.jackson.annotation.JsonFormat;

import lombok.*;

@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor()
@Node("user_profile")
public class UserProfile {
    private @GeneratedValue(UUIDStringGenerator.class) @Id String profileId;
    private @Property("userId") @NotNull String userId;
    private @Property("username") String username;
    private @Property("firstName") String firstName;
    private @Property("lastName") String lastName;
    private @Property("dob") @JsonFormat(pattern = "yyyy-MM-dd") LocalDate dob;
    private @Property("address") String address;
    private @Property("gender") String gender;
    private @Lob @Property("avatarUri") byte avatarUri;
    private @Transient String elementId;
}
