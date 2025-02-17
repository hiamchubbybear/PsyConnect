package dev.psyconnect.identity_service.enumeration;

import java.util.Arrays;
import java.util.List;

public enum RoleEnum {
    ADMIN("ADMIN", "ADMIN_PERMISSION", Arrays.asList(PermissionEnum.values())),
    THERAPIST(
            "THERAPIST",
            "THERAPIST_PERMISSION",
            Arrays.asList(
                    PermissionEnum.THERAPIST_CREATE_POST,
                    PermissionEnum.THERAPIST_EDIT_POST,
                    PermissionEnum.THERAPIST_RESPOND_TO_MESSAGE,
                    PermissionEnum.THERAPIST_MANAGE_APPOINTMENTS)),
    USER(
            "CLIENT",
            "CLIENT_PERMISSION",
            Arrays.asList(
                    PermissionEnum.USER_VIEW_POSTS,
                    PermissionEnum.USER_SEND_MESSAGE,
                    PermissionEnum.USER_BOOK_APPOINTMENT));

    private String id;
    private List<PermissionEnum> permissions;
    private String name;

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

    RoleEnum() {}

    RoleEnum(String id, String name, List<PermissionEnum> permissions) {
        this.id = id;
        this.permissions = permissions;
        this.name = name;
    }
}
