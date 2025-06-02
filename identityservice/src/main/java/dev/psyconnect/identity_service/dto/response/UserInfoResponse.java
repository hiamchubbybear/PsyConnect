package dev.psyconnect.identity_service.dto.response;

import java.time.LocalDate;
import java.util.Set;
import java.util.UUID;

import jakarta.persistence.*;

import com.fasterxml.jackson.annotation.JsonIgnore;

import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.model.RoleEntity;
import lombok.*;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class UserInfoResponse {
    private UUID accountId;
    private String username;
    private String email;

    @Enumerated(EnumType.STRING)
    private Provider provider;

    private @Builder.Default boolean isActivated = false;
    private LocalDate createdAt;

    @JsonIgnore
    @ManyToMany(fetch = FetchType.LAZY)
    Set<RoleEntity> role;

    @JsonIgnore
    UUID session;
}
