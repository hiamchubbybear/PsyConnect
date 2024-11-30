package com.example.IdentityService.Enum;

public enum RoleEnum {
    ROLE_ADMIN("001", "Permit All" , "Admin has all permission to do "),
    ROLE_USER("002", "User Permission" , "User has user permission to do "),
    ROLE_THERAPIST("003", "Therapist Permission" , "Therapist has therapist permission to do ");
    String id ;
    String permission;
    String description;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getPermission() {
        return permission;
    }

    public void setPermission(String permission) {
        this.permission = permission;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    RoleEnum(String id, String permission, String description) {
        this.id = id;
        this.permission = permission;
        this.description = description;
    }
}
