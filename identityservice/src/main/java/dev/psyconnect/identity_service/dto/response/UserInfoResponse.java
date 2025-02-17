package dev.psyconnect.identity_service.dto.response;

import com.fasterxml.jackson.annotation.JsonIgnore;
import dev.psyconnect.identity_service.enumeration.Provider;
import dev.psyconnect.identity_service.model.RoleEntity;
import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.CreationTimestamp;

import java.sql.Timestamp;
import java.time.LocalDate;
import java.util.Locale;
import java.util.Set;
import java.util.UUID;

@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder

public class UserInfoResponse {
    private UUID userId;
    private String username;
    private String email;
    @Enumerated(EnumType.STRING)
    private Provider provider;
    private @Builder.Default boolean isActivated = false;
    private  LocalDate createdAt;
    @JsonIgnore
    @ManyToMany(fetch = FetchType.LAZY)
    Set<RoleEntity> role;
    UUID session;

}
