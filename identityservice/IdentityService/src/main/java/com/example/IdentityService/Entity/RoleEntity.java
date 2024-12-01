package com.example.IdentityService.Entity;

import lombok.*;
import org.springframework.data.neo4j.core.schema.*;

import java.util.List;
import java.util.Set;

@Node("role")
@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class RoleEntity {
    @Id @Property("roleId")
    String roleId;
    @Relationship(type = "HAS_PERMISSION" , direction = Relationship.Direction.OUTGOING)
    List<Permission> permission;
    @Property("description")
    String description;
}
