package com.example.IdentityService.repository;

import com.example.IdentityService.model.RoleEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import javax.management.relation.Role;
import java.util.List;
import java.util.Optional;

@Repository
public interface RoleRepository extends JpaRepository<RoleEntity,String> {
    boolean existsByRoleId(String id);

    Optional<RoleEntity> findByName(String name);
}
