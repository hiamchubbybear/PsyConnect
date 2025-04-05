package dev.psyconnect.profile_service.dto.request;

import lombok.*;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
@Setter
public class UserSettingRequest {
    private String profileId;
    private String privacyLevel;
    private boolean showLastSeen;
    private boolean showProfilePicture;
    private boolean showMood;
    private boolean notificationsEnabled;
    private boolean emailNotifications;
    private boolean pushNotifications;
    private boolean smsNotifications;
    private boolean twoFactorAuth;
    private boolean allowLoginAlerts;
    private String trustedDevices;
    private String language;
    private String theme;
    private boolean autoDeleteOldMoods;
}
