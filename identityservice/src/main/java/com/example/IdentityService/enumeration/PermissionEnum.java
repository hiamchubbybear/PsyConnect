package com.example.IdentityService.enumeration;

import lombok.Getter;

@Getter
public enum PermissionEnum {
    ADMIN_CREATE_USER(1, "admin_create_user", "Quản lý tạo người dùng", "Admin"),
    ADMIN_EDIT_USER(2, "admin_edit_user", "Quản lý chỉnh sửa thông tin người dùng", "Admin"),
    ADMIN_DELETE_USER(3, "admin_delete_user", "Quản lý xóa người dùng", "Admin"),
    ADMIN_MANAGE_ROLE(4, "admin_manage_role", "Quản lý vai trò người dùng", "Admin"),
    ADMIN_VIEW_USER_PROFILE(5, "admin_view_user_profile", "Xem hồ sơ người dùng", "Admin"),
    ADMIN_MANAGE_CONTENT(6, "admin_manage_content", "Quản lý bài đăng và nội dung", "Admin"),
    ADMIN_VIEW_APPOINTMENTS(7, "admin_view_appointments", "Xem lịch hẹn của tất cả người dùng", "Admin"),
    ADMIN_VIEW_MESSAGES(8, "admin_view_messages", "Xem tin nhắn của người dùng", "Admin"),
    ADMIN_ACCESS_DASHBOARD(9, "admin_access_dashboard", "Truy cập bảng điều khiển Admin", "Admin"),

    THERAPIST_CREATE_POST(10, "therapist_create_post", "Tạo bài đăng tư vấn tâm lý", "Therapist"),
    THERAPIST_EDIT_POST(11, "therapist_edit_post", "Chỉnh sửa bài đăng tư vấn tâm lý", "Therapist"),
    THERAPIST_DELETE_POST(12, "therapist_delete_post", "Xóa bài đăng tư vấn tâm lý", "Therapist"),
    THERAPIST_RESPOND_TO_MESSAGE(13, "therapist_respond_to_message", "Trả lời tin nhắn tư vấn", "Therapist"),
    THERAPIST_MANAGE_APPOINTMENTS(14, "therapist_manage_appointments", "Quản lý lịch hẹn trực tiếp", "Therapist"),
    THERAPIST_ACCESS_DASHBOARD(15, "therapist_access_dashboard", "Truy cập bảng điều khiển Therapist", "Therapist"),
    THERAPIST_VIEW_USER_PROFILE(16, "therapist_view_user_profile", "Xem hồ sơ người dùng", "Therapist"),
    THERAPIST_RATE_USER(17, "therapist_rate_user", "Đánh giá người dùng sau khi tư vấn", "Therapist"),

    USER_VIEW_POSTS(18, "user_view_posts", "Xem bài đăng tư vấn tâm lý", "Client"),
    USER_SEND_MESSAGE(19, "user_send_message", "Gửi tin nhắn tư vấn tâm lý", "Client"),
    USER_BOOK_APPOINTMENT(20, "user_book_appointment", "Đặt lịch hẹn trực tiếp với therapist", "Client"),
    USER_COMMENT_ON_BLOG(21, "user_comment_on_blog", "Bình luận trên blog", "Client"),
    USER_READ_BLOG(22, "user_read_blog", "Đọc bài blog tư vấn tâm lý", "Client"),
    USER_RATE_THERAPIST(23, "user_rate_therapist", "Đánh giá chuyên gia tâm lý sau khi tư vấn", "Client"),
    USER_VIEW_PROFILE(24, "user_view_profile", "Xem hồ sơ cá nhân", "Client"),
    USER_UPDATE_PROFILE(25, "user_update_profile", "Cập nhật thông tin cá nhân", "Client"),
    USER_ACCESS_DASHBOARD(26, "user_access_dashboard", "Truy cập bảng điều khiển Client", "Client"),

    MANAGE_NOTIFICATIONS(27, "manage_notifications", "Quản lý thông báo hệ thống", "Admin"),
    VIEW_REPORTS(28, "view_reports", "Xem báo cáo hoạt động", "Admin"),
    VIEW_ANALYTICS(29, "view_analytics", "Xem phân tích dữ liệu hệ thống", "Admin");

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
