package dev.psyconnect.identity_service.enumeration;

import lombok.Getter;

@Getter
public enum PermissionEnum {
    // ADMIN Permissions
    ADMIN_CREATE_USER(1, "admin.user.create", "Admin creates a new user", "ADMIN"),
    ADMIN_EDIT_USER(2, "admin.user.edit", "Admin edits user information", "ADMIN"),
    ADMIN_DELETE_USER(3, "admin.user.delete", "Admin deletes a user", "ADMIN"),
    ADMIN_MANAGE_ROLE(4, "admin.role.manage", "Admin manages user roles", "ADMIN"),
    ADMIN_VIEW_USER_PROFILE(5, "admin.user.profile.view", "Admin views user profiles", "ADMIN"),
    ADMIN_MANAGE_CONTENT(6, "admin.content.manage", "Admin manages posts and content", "ADMIN"),
    ADMIN_VIEW_APPOINTMENTS(7, "admin.appointment.view", "Admin views all user appointments", "ADMIN"),
    ADMIN_VIEW_MESSAGES(8, "admin.message.view", "Admin views user messages", "ADMIN"),
    ADMIN_ACCESS_DASHBOARD(9, "admin.dashboard.access", "Admin accesses the dashboard", "ADMIN"),

    // THERAPIST Permissions
    THERAPIST_CREATE_POST(10, "therapist.post.create", "Therapist creates a counseling post", "THERAPIST"),
    THERAPIST_EDIT_POST(11, "therapist.post.edit", "Therapist edits a counseling post", "THERAPIST"),
    THERAPIST_DELETE_POST(12, "therapist.post.delete", "Therapist deletes a counseling post", "THERAPIST"),
    THERAPIST_RESPOND_TO_MESSAGE(13, "therapist.message.respond", "Therapist responds to messages", "THERAPIST"),
    THERAPIST_MANAGE_APPOINTMENTS(14, "therapist.appointment.manage", "Therapist manages appointments", "THERAPIST"),
    THERAPIST_ACCESS_DASHBOARD(15, "therapist.dashboard.access", "Therapist accesses the dashboard", "THERAPIST"),
    THERAPIST_VIEW_USER_PROFILE(16, "therapist.user.profile.view", "Therapist views user profiles", "THERAPIST"),
    THERAPIST_RATE_USER(17, "therapist.user.rate", "Therapist rates a user after a session", "THERAPIST"),

    // CLIENT/User Permissions
    USER_VIEW_POSTS(18, "user.post.view", "User views counseling posts", "CLIENT"),
    USER_SEND_MESSAGE(19, "user.message.send", "User sends counseling messages", "CLIENT"),
    USER_BOOK_APPOINTMENT(20, "user.appointment.book", "User books an appointment with a therapist", "CLIENT"),
    USER_COMMENT_ON_BLOG(21, "user.blog.comment", "User comments on blog posts", "CLIENT"),
    USER_READ_BLOG(22, "user.blog.read", "User reads counseling blog posts", "CLIENT"),
    USER_RATE_THERAPIST(23, "user.therapist.rate", "User rates a therapist after a session", "CLIENT"),
    USER_VIEW_PROFILE(24, "user.profile.view", "User views their profile", "CLIENT"),
    USER_UPDATE_PROFILE(25, "user.profile.update", "User updates their profile information", "CLIENT"),
    USER_ACCESS_DASHBOARD(26, "user.dashboard.access", "User accesses the dashboard", "CLIENT"),

    // Additional Admin Permissions
    MANAGE_NOTIFICATIONS(27, "admin.notification.manage", "Admin manages system notifications", "ADMIN"),
    VIEW_REPORTS(28, "admin.report.view", "Admin views activity reports", "ADMIN"),
    VIEW_ANALYTICS(29, "admin.analytics.view", "Admin views system analytics", "ADMIN");

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
