package dev.psyconnect.identity_service.model;

import java.sql.Timestamp;
import java.util.Set;
import java.util.UUID;

import jakarta.persistence.*;

import org.hibernate.annotations.CreationTimestamp;

import com.fasterxml.jackson.annotation.JsonIgnore;

import dev.psyconnect.identity_service.enumeration.Provider;
import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Builder
@Table(name = "user_account", uniqueConstraints = @UniqueConstraint(columnNames = "email"))
public class Account {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID accountId;
    private UUID profileId;
    private String username;
    private String password;
    private String email;

    @Enumerated(EnumType.STRING)
    private Provider provider;

    private boolean isActivated = false;

    private @CreationTimestamp Timestamp createdAt;

    @JsonIgnore
    @ManyToMany(fetch = FetchType.EAGER)
    Set<RoleEntity> role;

    UUID session;

    @OneToOne
    private Token token;
}
