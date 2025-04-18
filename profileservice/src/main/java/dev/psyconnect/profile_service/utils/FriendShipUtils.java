package dev.psyconnect.profile_service.utils;

import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.model.FriendRelationship;
import dev.psyconnect.profile_service.model.Profile;
import dev.psyconnect.profile_service.repository.ProfileRepository;

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
