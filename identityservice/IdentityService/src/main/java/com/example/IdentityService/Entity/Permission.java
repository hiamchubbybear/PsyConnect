package com.example.IdentityService.Entity;

import com.fasterxml.jackson.annotation.JsonIgnore;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.Node;
import org.springframework.data.neo4j.core.schema.Property;
import org.springframework.data.neo4j.core.schema.Relationship;

import javax.management.relation.Role;

@Node("permission")
@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
public class Permission {
    @Id
    String id;
    @Property("description")
    String description;
    @Property("name")
    String name;
    @Relationship(type = "BELONG_TO_ROLE" , direction = Relationship.Direction.OUTGOING)
    RoleEntity role;

}
