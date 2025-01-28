package com.example.IdentityService.configuration;

import java.util.*;

import jakarta.transaction.Transactional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.ApplicationRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import com.example.IdentityService.model.Permission;
import com.example.IdentityService.model.RoleEntity;
import com.example.IdentityService.enumeration.PermissionEnum;
import com.example.IdentityService.enumeration.RoleEnum;
import com.example.IdentityService.repository.PermissionRepository;
import com.example.IdentityService.repository.RoleRepository;

import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@Configuration
@RequiredArgsConstructor
@FieldDefaults(makeFinal = true, level = AccessLevel.PRIVATE)
public class BeanConfiguration {
    private static final Logger log = LoggerFactory.getLogger(BeanConfiguration.class);

    RoleRepository roleRepository;
    PermissionRepository permissionRepository;

    /**
     * @return This bean run when initializing the application.
     * This checks roles, creates permissions from PermissionEnum, and associates roles with permissions.
     */
    @Bean
    @Transactional
    ApplicationRunner runner() {
        return args -> {
            for (PermissionEnum permissionEnum : PermissionEnum.values()) {
                Permission permission = permissionRepository.findByName(permissionEnum.getName());
                if (permission == null) {
                    permission = new Permission(
                            String.valueOf(permissionEnum.getId()),
                            permissionEnum.getName(),
                            permissionEnum.getDescription(),
                            new ArrayList<>());
                    permissionRepository.save(permission);
                    log.info("Created permission: {}", permissionEnum.getName());
                }
            }

            for (RoleEnum roleEnum : RoleEnum.values()) {
                if (!roleRepository.existsByRoleId(roleEnum.getId())) {
                    List<Permission> permissions = new ArrayList<>();
                    for (PermissionEnum permissionEnum : roleEnum.getPermissions()) {
                        Permission permission = permissionRepository.findByName(permissionEnum.getName());
                        if (permission != null) {
                            permissions.add(permission);
                        } else {
                            log.warn(
                                    "Permission {} not found for role {}",
                                    permissionEnum.getName(),
                                    roleEnum.getName());
                        }
                    }

                    RoleEntity roleEntity = new RoleEntity(roleEnum.getId(), permissions, null, roleEnum.getName());
                    roleRepository.save(roleEntity);
                    log.info("Created role: {}", roleEntity.getRoleId());
                } else {
                    log.info("Role {} already exists.", roleEnum.getName());
                }
            }

            for (PermissionEnum permissionEnum : PermissionEnum.values()) {
                Permission permission = permissionRepository.findByName(permissionEnum.getName());
                if (permission != null
                        && (permission.getRoles() == null
                                || permission.getRoles().isEmpty())) {
                    List<RoleEntity> roleEntities =
                            roleRepository.findAllById(Collections.singleton(permissionEnum.getRole()));
                    if (!roleEntities.isEmpty()) {
                        permission.setRoles(roleEntities);
                        permissionRepository.save(permission);
                        log.info("Updated permission {} with roles.", permissionEnum.getName());
                    } else {
                        log.warn("No role found for permission: {}", permissionEnum.getName());
                    }
                }
            }
        };
    }
}
