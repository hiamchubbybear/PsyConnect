package com.example.IdentityService.model;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import lombok.*;

import java.util.List;
import java.util.Set;

@Entity
@Table(name = "permission")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
public class Permission {
    @Id
    private String id;
    @Column(nullable = false, unique = true)
    private String name;
    private String description;
    @ManyToMany(mappedBy = "permissions", fetch = FetchType.EAGER)
    @JsonIgnore
    private List<RoleEntity> roles;
}