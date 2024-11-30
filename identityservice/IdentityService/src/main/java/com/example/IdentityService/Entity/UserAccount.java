package com.example.IdentityService.Entity;

import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.neo4j.core.schema.*;

import java.util.Set;
import java.util.UUID;

@Node("useraccount")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class UserAccount {
    @Id
    private UUID id;
    @Property("username")
    private String username;
    @Property("password")
    private String password;
    @Property("email")
    private String email;
    @Relationship(type = "HAS_ROLE",direction = Relationship.Direction.OUTGOING  )
    RoleEntity role;
}
