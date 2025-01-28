package com.example.IdentityService.model;

import java.util.Set;
import java.util.UUID;

import jakarta.persistence.*;

import com.example.IdentityService.enumeration.Provider;
import com.fasterxml.jackson.annotation.JsonIgnore;

import lombok.*;

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

    @JsonIgnore
    @ManyToMany(fetch = FetchType.LAZY)
    Set<RoleEntity> role;
}
