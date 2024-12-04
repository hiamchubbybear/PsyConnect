package com.example.IdentityService.num;

import com.fasterxml.jackson.databind.annotation.EnumNaming;
import jakarta.persistence.Entity;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;

import java.util.Arrays;
import java.util.List;
import java.util.Set;

public enum RoleEnum {
    ADMIN("Admin", "Admin default permission", Arrays.asList(PermissionEnum.values())),
    THERAPIST("Therapist", "Therapist default permission", Arrays.asList(
            PermissionEnum.THERAPIST_CREATE_POST,
            PermissionEnum.THERAPIST_EDIT_POST,
            PermissionEnum.THERAPIST_RESPOND_TO_MESSAGE,
            PermissionEnum.THERAPIST_MANAGE_APPOINTMENTS
    )),
    USER("Client", "Client default permission", Arrays.asList(
            PermissionEnum.USER_VIEW_POSTS,
            PermissionEnum.USER_SEND_MESSAGE,
            PermissionEnum.USER_BOOK_APPOINTMENT
    ));

    private String id;
    private  List<PermissionEnum> permissions;
    private  String name;


    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public List<PermissionEnum> getPermissions() {
        return permissions;
    }

    public String getName() {
        return name;
    }

    RoleEnum() {
    }

    RoleEnum(String id, String name , List<PermissionEnum> permissions) {
        this.id = id;
        this.permissions = permissions;
        this.name = name;
    }
}