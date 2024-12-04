package com.example.IdentityService.repository;

import com.example.IdentityService.model.RoleEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface RoleRepository extends JpaRepository<RoleEntity,String> {
    boolean existsByRoleId(String id);
    RoleEntity findByRoleId(String roleId);
    List<RoleEntity> findAllByName(String roleName);
    RoleEntity findByName(String role);

}
