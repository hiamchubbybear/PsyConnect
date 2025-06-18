package dev.psyconnect.profile_service.model;

import java.io.Serializable;

import org.springframework.data.neo4j.core.schema.Id;
import org.springframework.data.neo4j.core.schema.Node;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Node("user_setting")
@Data
@Builder
@AllArgsConstructor
@NoArgsConstructor
public class Setting implements Serializable {
    @Id
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
