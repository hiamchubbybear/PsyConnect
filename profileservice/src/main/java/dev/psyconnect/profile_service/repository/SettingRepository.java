package dev.psyconnect.profile_service.repository;

import java.util.Optional;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.data.neo4j.repository.query.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import dev.psyconnect.profile_service.model.Setting;

@Repository
public interface SettingRepository extends Neo4jRepository<Setting, String> {

    @Query("""
            MATCH (u:user_profile {profileId: $profileId})-[:HAS_SETTING]->(s:user_setting)
            RETURN s
            """)
    Optional<Setting> getUserSetting(@Param("profileId") String profileId);

    @Query(
            """
                    MATCH (u:user_profile {profileId: $profileId})
                    CREATE (s:user_setting {
                    	profileId: $profileId,
                    	privacyLevel: $privacyLevel,
                    	showLastSeen: $showLastSeen,
                    	showProfilePicture: $showProfilePicture,
                    	showMood: $showMood,
                    	notificationsEnabled: $notificationsEnabled,
                    	emailNotifications: $emailNotifications,
                    	pushNotifications: $pushNotifications,
                    	smsNotifications: $smsNotifications,
                    	twoFactorAuth: $twoFactorAuth,
                    	allowLoginAlerts: $allowLoginAlerts,
                    	trustedDevices: $trustedDevices,
                    	language: $language,
                    	theme: $theme,
                    	autoDeleteOldMoods: $autoDeleteOldMoods
                    })
                    CREATE (u)-[:HAS_SETTING]->(s)
                    RETURN s
                    """)
    Optional<Setting> createUserSetting(
            @Param("profileId") String profileId,
            @Param("privacyLevel") String privacyLevel,
            @Param("showLastSeen") boolean showLastSeen,
            @Param("showProfilePicture") boolean showProfilePicture,
            @Param("showMood") boolean showMood,
            @Param("notificationsEnabled") boolean notificationsEnabled,
            @Param("emailNotifications") boolean emailNotifications,
            @Param("pushNotifications") boolean pushNotifications,
            @Param("smsNotifications") boolean smsNotifications,
            @Param("twoFactorAuth") boolean twoFactorAuth,
            @Param("allowLoginAlerts") boolean allowLoginAlerts,
            @Param("trustedDevices") String trustedDevices,
            @Param("language") String language,
            @Param("theme") String theme,
            @Param("autoDeleteOldMoods") boolean autoDeleteOldMoods);

    @Query(
            """
                    MATCH (u:user_profile {profileId: $profileId})-[:HAS_SETTING]->(s:user_setting)
                    SET s.privacyLevel = $privacyLevel,
                    	s.showLastSeen = $showLastSeen,
                    	s.showProfilePicture = $showProfilePicture,
                    	s.showMood = $showMood,
                    	s.notificationsEnabled = $notificationsEnabled,
                    	s.emailNotifications = $emailNotifications,
                    	s.pushNotifications = $pushNotifications,
                    	s.smsNotifications = $smsNotifications,
                    	s.twoFactorAuth = $twoFactorAuth,
                    	s.allowLoginAlerts = $allowLoginAlerts,
                    	s.trustedDevices = $trustedDevices,
                    	s.language = $language,
                    	s.theme = $theme,
                    	s.autoDeleteOldMoods = $autoDeleteOldMoods
                    RETURN s
                    """)
    Optional<Setting> updateUserSetting(
            @Param("profileId") String profileId,
            @Param("privacyLevel") String privacyLevel,
            @Param("showLastSeen") boolean showLastSeen,
            @Param("showProfilePicture") boolean showProfilePicture,
            @Param("showMood") boolean showMood,
            @Param("notificationsEnabled") boolean notificationsEnabled,
            @Param("emailNotifications") boolean emailNotifications,
            @Param("pushNotifications") boolean pushNotifications,
            @Param("smsNotifications") boolean smsNotifications,
            @Param("twoFactorAuth") boolean twoFactorAuth,
            @Param("allowLoginAlerts") boolean allowLoginAlerts,
            @Param("trustedDevices") String trustedDevices,
            @Param("language") String language,
            @Param("theme") String theme,
            @Param("autoDeleteOldMoods") boolean autoDeleteOldMoods);

    @Query("""
                MATCH (u:user_profile {profileId: $profileId})-[:HAS_SETTING]->(s:user_setting)
                DETACH DELETE s
            """)
    void deleteUserSetting(@Param("profileId") String profileId);

}
