package com.example.IdentityService.model;

import java.util.Set;

import com.example.IdentityService.num.Provider;
import jakarta.persistence.*;

import com.fasterxml.jackson.annotation.JsonIgnore;

import lombok.*;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Entity
@Builder
@Table(name = "user_account", uniqueConstraints = @UniqueConstraint(columnNames = "username"))
public class UserAccount {
    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private String id;

    private String username;
    private String password;
    private String email;
    @Enumerated(EnumType.STRING)
    private Provider provider;
    @JsonIgnore
    @ManyToMany(fetch = FetchType.LAZY)
    Set<RoleEntity> role;
}
