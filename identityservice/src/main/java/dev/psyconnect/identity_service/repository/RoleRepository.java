package dev.psyconnect.identity_service.repository;

import java.util.Optional;
import java.util.Set;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import dev.psyconnect.identity_service.model.Permission;
import dev.psyconnect.identity_service.model.RoleEntity;

@Repository
public interface RoleRepository extends JpaRepository<RoleEntity, String> {
    boolean existsByRoleId(String id);

    Optional<RoleEntity> findByName(String name);

    Set<RoleEntity> findAllByRoleId(String client);

    Set<RoleEntity> findAllByPermissionsContaining(Permission permission);
}
