package dev.psyconnect.identity_service.model;

import jakarta.persistence.Column;
import jakarta.persistence.Entity;
import jakarta.persistence.Id;

import lombok.*;

@Entity
@NoArgsConstructor
@AllArgsConstructor
@Builder
@Setter
@Getter
public class BlackListToken {
    @Id
    @Column(length = 512)
    String token;
}
