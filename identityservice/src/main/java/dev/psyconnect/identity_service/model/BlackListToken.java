package dev.psyconnect.identity_service.model;

import jakarta.persistence.*;
import lombok.*;

@Entity
@NoArgsConstructor
@AllArgsConstructor
@Builder
@Setter
@Getter
public class BlackListToken {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private Long id;

    @Lob
    @Column(columnDefinition = "TEXT", nullable = false)
    private String token;
}
