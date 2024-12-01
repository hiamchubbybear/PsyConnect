package com.example.IdentityService.Repository;

import com.example.IdentityService.Entity.Permission;
import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface PermissionRepository extends Neo4jRepository<Permission, String> {
    Permission findByName(String name);
}
