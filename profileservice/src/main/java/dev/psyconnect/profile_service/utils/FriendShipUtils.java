package dev.psyconnect.profile_service.utils;

import dev.psyconnect.profile_service.enums.FriendShipStatus;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.model.FriendRelationship;
import dev.psyconnect.profile_service.model.Profile;
import dev.psyconnect.profile_service.repository.ProfileRepository;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class FriendShipUtils {
    private final ProfileRepository profileRepository;
    ProfileRepository repository;

    public FriendShipUtils(ProfileRepository profileRepository) {
        this.profileRepository = profileRepository;
    }

    public boolean hasAnyRelationship(Profile profile, String targetId) {
        for (FriendRelationship fr : profile.getFriends()) {
            if (fr.getTarget().getProfileId().equals(targetId)) {
                return true;
            }
        }
        return false;
    }

    public boolean existByUser(String u1, String u2) {
        return profileRepository.existsById(u1) || profileRepository.existsById(u2);
    }
}
