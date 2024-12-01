package com.example.IdentityService.Repository;

import com.example.IdentityService.Entity.RoleEntity;
import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface RoleRepository extends Neo4jRepository<RoleEntity,String> {
    boolean existsByRoleId(String id);
    RoleEntity findByRoleId(String roleId);
}
