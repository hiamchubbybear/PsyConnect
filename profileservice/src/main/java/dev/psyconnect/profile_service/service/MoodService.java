package dev.psyconnect.profile_service.service;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import dev.psyconnect.profile_service.model.Profile;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.dto.request.MoodCreateRequest;
import dev.psyconnect.profile_service.dto.request.MoodUpdateRequest;
import dev.psyconnect.profile_service.dto.response.*;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.model.Mood;
import dev.psyconnect.profile_service.repository.MoodRepository;
import dev.psyconnect.profile_service.repository.ProfileRepository;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@Service
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class MoodService {
    MoodRepository moodRepository;
    ProfileRepository profileRepository;

    @CacheEvict(key = "#profileId", value = "mood")
    public MoodCreateResponse createMoodByProfileId(String profileId, MoodCreateRequest request) {
        log.info("Create mood for user with profile id  {} ", profileId);
        log.info("Mood description {}", request.getMood());
        log.info("Mood description {}", request.getMoodDescription());
        log.info("Mood visibility {}", request.getVisibility());
        if (!profileRepository.existsById(profileId))
            throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        if (profileRepository.hasMood(profileId)) {
            throw new CustomExceptionHandler(ErrorCode.MOOD_ALREADY_EXISTS);
        }
        // Ho_Chi_Minh TimeZones
        final int TIME_ZONE = +7;
        Map<String, Long> dateTime = Time.MOOD_EXPIRES;
        long currentTime = dateTime.get("currentTimeMillis");
        long expiresTime = dateTime.get("expiresTimeMillis");
        LocalDateTime localDateTimeExpires = Time.fromTimeStampToLCD(TIME_ZONE, expiresTime);
        var response = moodRepository
                .createMood(
                        profileId,
                        request.getMood(),
                        request.getMoodDescription(),
                        request.getVisibility(),
                        currentTime,
                        expiresTime)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
        return MoodCreateResponse.builder()
                .profileId(profileId)
                .isSuccess(true)
                .moodDescription(request.getMoodDescription())
                .mood(request.getMood())
                .timeExpires(localDateTimeExpires)
                .build();
    }

    @CacheEvict(key = "#profileId", value = "mood")
    public MoodCreateResponse updateMoodByProfileId(String profileId, MoodUpdateRequest request) {
        if (!profileRepository.existsById(profileId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        var response = moodRepository
                .updateMood(profileId, request.getMood(), request.getMoodDescription(), request.getVisibility())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
        return MoodCreateResponse.builder()
                .isSuccess(true)
                .moodDescription(response.getDescription())
                .profileId(profileId)
                .build();
    }

    //    @Cacheable(key = "#profileId", value = "mood")
    public GetMoodResponse getMoodById(String profileId) {
        Profile profile = profileRepository.findById(profileId).orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        String fullName = profile.getFirstName() + profile.getLastName();
        var res = GetMoodResponse.builder()
                .profileId(profileId)
                .avatarUrl(profile.getAvatarUri())
                .fullName(fullName)
                .moodId(profile.getMoodList().getMoodId())
                .mood(profile.getMoodList().getMood())
                .description(profile.getMoodList().getDescription())
                .expiresAt(profile.getMoodList().getExpiresAt())
                .createdAt(profile.getMoodList().getCreatedAt())
                .visibility(profile.getMoodList().getVisibility())
                .build();
        log.info(res.toString());
        return res;
    }

    @CacheEvict(key = "#profileId", value = "mood")
    public DeleteMoodResponse deleteMood(String profileId) {
        if (!profileRepository.existsById(profileId)) throw new CustomExceptionHandler(ErrorCode.QUERY_FAILED);
        moodRepository.deleteById(profileId);
        return DeleteMoodResponse.builder()
                .isSuccess(true)
                .message("Delete success")
                .build();
    }

    public List<FriendMoodDTO> getFriendsMood(String profileId) {
        List<ProfileWithMood> foundObject = profileRepository.findFriendsWithMoodsByProfileId(profileId);
        List<FriendMoodDTO> response = new ArrayList<>();
        for (ProfileWithMood pwm : foundObject) {
            var profile = pwm.getProfile();
            var mood = pwm.getMood();
            response.add(FriendMoodDTO.builder()
                    .profileId(profile.getProfileId())
                    .fullName(profile.getFirstName() + profile.getLastName())
                    .avatarUrl(profile.getAvatarUri())
                    .moodId(mood.getMoodId())
                    .mood(mood.getMood())
                    .description(mood.getDescription())
                    .visibility(mood.getVisibility())
                    .createdAt(mood.getCreatedAt())
                    .expiresAt(mood.getExpiresAt())
                    .build());
        }
        return response;
    }
}
