package com.example.IdentityService.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.example.IdentityService.model.RoleEntity;

@Repository
public interface RoleRepository extends JpaRepository<RoleEntity, String> {
    boolean existsByRoleId(String id);
}
