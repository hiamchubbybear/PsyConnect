package dev.psyconnect.identity_service.enumeration;

import lombok.Getter;

@Getter
public enum PermissionEnum {
    ADMIN_CREATE_USER(1, "ADMIN_CREATE_USER", "Admin creates a new user", "ADMIN"),
    ADMIN_EDIT_USER(2, "ADMIN_EDIT_USER", "Admin edits user information", "ADMIN"),
    ADMIN_DELETE_USER(3, "ADMIN_DELETE_USER", "Admin deletes a user", "ADMIN"),
    ADMIN_MANAGE_ROLE(4, "ADMIN_MANAGE_ROLE", "Admin manages user roles", "ADMIN"),
    ADMIN_VIEW_USER_PROFILE(5, "ADMIN_VIEW_USER_PROFILE", "Admin views user profiles", "ADMIN"),
    ADMIN_MANAGE_CONTENT(6, "ADMIN_MANAGE_CONTENT", "Admin manages posts and content", "ADMIN"),
    ADMIN_VIEW_APPOINTMENTS(7, "ADMIN_VIEW_APPOINTMENTS", "Admin views all user appointments", "ADMIN"),
    ADMIN_VIEW_MESSAGES(8, "ADMIN_VIEW_MESSAGES", "Admin views user messages", "ADMIN"),
    ADMIN_ACCESS_DASHBOARD(9, "ADMIN_ACCESS_DASHBOARD", "Admin accesses the dashboard", "ADMIN"),

    THERAPIST_CREATE_POST(10, "THERAPIST_CREATE_POST", "Therapist creates a counseling post", "THERAPIST"),
    THERAPIST_EDIT_POST(11, "THERAPIST_EDIT_POST", "Therapist edits a counseling post", "THERAPIST"),
    THERAPIST_DELETE_POST(12, "THERAPIST_DELETE_POST", "Therapist deletes a counseling post", "THERAPIST"),
    THERAPIST_RESPOND_TO_MESSAGE(13, "THERAPIST_RESPOND_TO_MESSAGE", "Therapist responds to messages", "THERAPIST"),
    THERAPIST_MANAGE_APPOINTMENTS(14, "THERAPIST_MANAGE_APPOINTMENTS", "Therapist manages appointments", "THERAPIST"),
    THERAPIST_ACCESS_DASHBOARD(15, "THERAPIST_ACCESS_DASHBOARD", "Therapist accesses the dashboard", "THERAPIST"),
    THERAPIST_VIEW_USER_PROFILE(16, "THERAPIST_VIEW_USER_PROFILE", "Therapist views user profiles", "THERAPIST"),
    THERAPIST_RATE_USER(17, "THERAPIST_RATE_USER", "Therapist rates a user after a session", "THERAPIST"),

    USER_VIEW_POSTS(18, "USER_VIEW_POSTS", "User views counseling posts", "CLIENT"),
    USER_SEND_MESSAGE(19, "USER_SEND_MESSAGE", "User sends counseling messages", "CLIENT"),
    USER_BOOK_APPOINTMENT(20, "USER_BOOK_APPOINTMENT", "User books an appointment with a therapist", "CLIENT"),
    USER_COMMENT_ON_BLOG(21, "USER_COMMENT_ON_BLOG", "User comments on blog posts", "CLIENT"),
    USER_READ_BLOG(22, "USER_READ_BLOG", "User reads counseling blog posts", "CLIENT"),
    USER_RATE_THERAPIST(23, "USER_RATE_THERAPIST", "User rates a therapist after a session", "CLIENT"),
    USER_VIEW_PROFILE(24, "USER_VIEW_PROFILE", "User views their profile", "CLIENT"),
    USER_UPDATE_PROFILE(25, "USER_UPDATE_PROFILE", "User updates their profile information", "CLIENT"),
    USER_ACCESS_DASHBOARD(26, "USER_ACCESS_DASHBOARD", "User accesses the dashboard", "CLIENT"),

    MANAGE_NOTIFICATIONS(27, "MANAGE_NOTIFICATIONS", "Admin manages system notifications", "ADMIN"),
    VIEW_REPORTS(28, "VIEW_REPORTS", "Admin views activity reports", "ADMIN"),
    VIEW_ANALYTICS(29, "VIEW_ANALYTICS", "Admin views system analytics", "ADMIN");


    private final int id;
    private final String name;
    private final String description;
    private final String role;

    PermissionEnum(int id, String name, String description, String role) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.role = role;
    }
}
