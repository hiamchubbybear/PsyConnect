package com.example.IdentityService.Bean;

import com.example.IdentityService.Entity.Permission;
import com.example.IdentityService.Entity.RoleEntity;
import com.example.IdentityService.Enum.PermissionEnum;
import com.example.IdentityService.Enum.RoleEnum;
import com.example.IdentityService.Repository.PermissionRepository;
import com.example.IdentityService.Repository.RoleRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.ApplicationRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.List;
import java.util.stream.Collectors;

@Configuration
public class BeanConfiguration {
    private final RoleRepository roleRepository;

    @Autowired
    public BeanConfiguration(RoleRepository roleRepository) {
        this.roleRepository = roleRepository;
    }

    /**
     *
     * @param permissionRepository
     * @return This bean run when initialize the applications
     * Check role if it doesn't exist create a role get all permission from PermissionEnum
     * Then loop the permission enum if it doesn't have permission role . Role will be find
     * and add  into Permission
     *
     */
    @Bean
    ApplicationRunner runner(PermissionRepository permissionRepository) {
        return args -> {
            for (RoleEnum roleEnum : RoleEnum.values()) {
                if (!roleRepository.existsByRoleId(roleEnum.getId())) {
                    List<Permission> permissions = roleEnum.getPermissions().stream()
                            .map(permissionEnum -> new Permission(
                                    String.valueOf(permissionEnum.getId()),
                                    permissionEnum.getDescription(),
                                    permissionEnum.getName(),
                                    null
                            ))
                            .collect(Collectors.toList());
                    RoleEntity roleEntity = new RoleEntity(roleEnum.getId(), permissions, roleEnum.getDescription());
                    roleRepository.save(roleEntity);
                } else {
                    System.out.println("Role " + roleEnum.getPermissions() + " already exists.");
                }
            }
            for (PermissionEnum permissionEnum : PermissionEnum.values()) {
                Permission permission = permissionRepository.findByName(permissionEnum.getName());
                if (permission == null) {
                    permission = new Permission(
                            String.valueOf(permissionEnum.getId()),
                            permissionEnum.getDescription(),
                            permissionEnum.getName(),
                            roleRepository.findByRoleId(permissionEnum.getRole())
                    );
                    var rolee = permissionEnum.getRole();
                    permissionRepository.save(permission);
                    System.out.println("Created permission: " + permissionEnum.getName());
                }
                if (permission.getRole() == null) {
                    RoleEntity roleEntity = roleRepository.findById(permissionEnum.getRole()).orElse(null);
                    if (roleEntity != null) {
                        permission.setRole(roleEntity);
                        permissionRepository.save(permission);
                    } else {
                        System.out.println("Role not found for permission " + permissionEnum.getName());
                    }
                }
            }
        };
    }
}
