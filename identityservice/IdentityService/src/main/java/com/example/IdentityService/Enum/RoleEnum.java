package com.example.IdentityService.Enum;

import com.example.IdentityService.Entity.Permission;

import java.util.Arrays;
import java.util.List;

public enum RoleEnum {
        ADMIN("Admin","Admin default permission", Arrays.asList(PermissionEnum.values())),
        THERAPIST("Therapist","Therapist default permission", Arrays.asList(
                PermissionEnum.THERAPIST_CREATE_POST,
                PermissionEnum.THERAPIST_EDIT_POST,
                PermissionEnum.THERAPIST_RESPOND_TO_MESSAGE,
                PermissionEnum.THERAPIST_MANAGE_APPOINTMENTS
        )),
        USER("Client","Client default permission", Arrays.asList(
                PermissionEnum.USER_VIEW_POSTS,
                PermissionEnum.USER_SEND_MESSAGE,
                PermissionEnum.USER_BOOK_APPOINTMENT
        ));
        private String id;
        private final List<PermissionEnum> permissions;
        private final String description;
    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getDescription() {
        return description;
    }

    RoleEnum(String id, String description , List<PermissionEnum> permissions) {
        this.id = id;
        this.permissions = permissions;
        this.description = description;
    }

    public List<PermissionEnum> getPermissions() {
            return permissions;
        }
    }