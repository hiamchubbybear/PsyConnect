package dev.psyconnect.identity_service.model;

import java.sql.Time;
import java.sql.Timestamp;
import java.time.Instant;
import java.util.Date;
import java.util.Set;
import java.util.UUID;

import com.google.type.DateTime;
import jakarta.persistence.*;

import org.hibernate.annotations.CreationTimestamp;

import com.fasterxml.jackson.annotation.JsonIgnore;

import dev.psyconnect.identity_service.enumeration.Provider;
import lombok.*;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

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


    public Account createAdminAccount(RoleEntity roleEntity) {
        Account admin = new Account();
        admin.accountId = UUID.randomUUID();
        admin.profileId = UUID.randomUUID();
        admin.username = "chessy1603";
        admin.password = new BCryptPasswordEncoder().encode("16032004");
        admin.email = "tranvanhuy160304@gmail.com";
        admin.provider = Provider.GOOGLE;
        admin.isActivated = true;
        admin.createdAt = Timestamp.from(Instant.now());
        admin.role = Set.of(roleEntity);
        admin.session = UUID.randomUUID();
        admin.token = null;
        return admin;
    }
}
