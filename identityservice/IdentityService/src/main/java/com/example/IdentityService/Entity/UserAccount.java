package com.example.IdentityService.Entity;

import com.fasterxml.jackson.annotation.ObjectIdGenerators;
import jdk.jfr.Enabled;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.boot.autoconfigure.domain.EntityScan;
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

//    @Relationship(type = "ACTED_IN",direction = Relationship.Direction.INCOMING  )
//    Set<Role> role;
}
