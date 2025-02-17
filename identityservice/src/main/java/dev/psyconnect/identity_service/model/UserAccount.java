package dev.psyconnect.identity_service.model;

import java.sql.Timestamp;
import java.util.Collection;
import java.util.List;
import java.util.Set;
import java.util.UUID;

import jakarta.persistence.*;

import org.hibernate.annotations.CreationTimestamp;

import com.fasterxml.jackson.annotation.JsonIgnore;

import dev.psyconnect.identity_service.enumeration.Provider;
import lombok.*;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Builder
@Table(name = "user_account", uniqueConstraints = @UniqueConstraint(columnNames = "email"))
public class UserAccount {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID userId;

    private String username;
    private String password;
    private String email;

    @Enumerated(EnumType.STRING)
    private Provider provider;

    private @Builder.Default boolean isActivated = false;

    private @CreationTimestamp Timestamp createdAt;

    @JsonIgnore
    @ManyToMany(fetch = FetchType.EAGER)
    Set<RoleEntity> role;

    UUID session;

    @OneToOne
    private Token token;

}
