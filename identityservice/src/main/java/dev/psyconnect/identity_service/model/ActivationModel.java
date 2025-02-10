package dev.psyconnect.identity_service.model;

import java.sql.Timestamp;
import java.util.UUID;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;

import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Entity
public class ActivationModel {
    @Id
    private UUID activateId;

    private String username;
    private String token;
    private Timestamp issuedAt;
    private Timestamp expires;
    private boolean revoked;
}
