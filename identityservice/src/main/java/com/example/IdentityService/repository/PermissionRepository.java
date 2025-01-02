package com.example.IdentityService.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import com.example.IdentityService.model.Permission;

@Repository
public interface PermissionRepository extends JpaRepository<Permission, String> {
    Permission findByName(String name);
}
