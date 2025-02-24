package dev.psyconnect.identity_service.configuration;

import java.util.*;

import jakarta.transaction.Transactional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.ApplicationRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import dev.psyconnect.identity_service.enumeration.PermissionEnum;
import dev.psyconnect.identity_service.enumeration.RoleEnum;
import dev.psyconnect.identity_service.model.Permission;
import dev.psyconnect.identity_service.model.RoleEntity;
import dev.psyconnect.identity_service.repository.PermissionRepository;
import dev.psyconnect.identity_service.repository.RoleRepository;
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

    @Bean
    @Transactional
    ApplicationRunner runner() {
        return args -> {
            // Tạo các Permission từ PermissionEnum
            List<Permission> permissionsToSave = new ArrayList<>();
            for (PermissionEnum permissionEnum : PermissionEnum.values()) {
                if (!permissionRepository.existsByName(permissionEnum.getName())) {
                    Permission permission = new Permission(
                            String.valueOf(permissionEnum.getId()),
                            permissionEnum.getName(),
                            permissionEnum.getDescription(),
                            new HashSet<>());
                    permissionsToSave.add(permission);
                    log.info("Created permission: {}", permissionEnum.getName());
                }
            }
            permissionRepository.saveAll(permissionsToSave);

            // Tạo các RoleEntity từ RoleEnum và liên kết với Permission
            Set<RoleEntity> rolesToSave = new HashSet<>();
            for (RoleEnum roleEnum : RoleEnum.values()) {
                if (!roleRepository.existsByRoleId(roleEnum.getId())) {
                    Set<Permission> permissions = new HashSet<>();
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
                    rolesToSave.add(roleEntity);
                    log.info("Created role: {}", roleEntity.getRoleId());
                } else {
                    log.info("Role {} already exists.", roleEnum.getName());
                }
            }
            roleRepository.saveAll(rolesToSave);

            // Cập nhật quan hệ giữa Permission và RoleEntity
            for (PermissionEnum permissionEnum : PermissionEnum.values()) {
                Permission permission = permissionRepository.findByName(permissionEnum.getName());
                if (permission != null
                        && (permission.getRoles() == null
                                || permission.getRoles().isEmpty())) {
                    Set<RoleEntity> roleEntities = roleRepository.findAllByPermissionsContaining(permission);
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
