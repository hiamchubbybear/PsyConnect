package dev.psyconnect.profile_service.service;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

import dev.psyconnect.profile_service.dto.response.*;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.dto.request.MoodCreateRequest;
import dev.psyconnect.profile_service.dto.request.MoodUpdateRequest;
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
        if (!profileRepository.existsById(profileId) || moodRepository.existsById(profileId))
            throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        log.info("Profile Id {}", profileId);
        final int TIME_ZONE = +7;
        Map<String, Long> dateTime = Time.MOOD_EXPIRES;
        long currentTime = dateTime.get("currentTimeMillis");
        long expiresTime = dateTime.get("expiresTimeMillis");
        LocalDateTime localDateTimeExpires = Time.fromTimeStampToLCD(TIME_ZONE, expiresTime);
        var response = moodRepository.createMood(profileId, request.getMood(), request.getMoodDescription(), request.getVisibility(), currentTime, expiresTime).orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
        log.info("Creating mood for profile {}", response.getMood());
        return MoodCreateResponse.builder().profileId(profileId).isSuccess(true).moodDescription(request.getMoodDescription()).mood(request.getMood()).timeExpires(localDateTimeExpires).build();
    }

    @CacheEvict(key = "#profileId", value = "mood")
    public MoodCreateResponse updateMoodByProfileId(String profileId, MoodUpdateRequest request) {
        if (!profileRepository.existsById(profileId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
        var response = moodRepository.updateMood(profileId, request.getMood(), request.getMoodDescription(), request.getVisibility()).orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
        log.info("Updating mood for profile {}", response.getMood());
        return MoodCreateResponse.builder().isSuccess(true).moodDescription(response.getDescription()).profileId(profileId).build();
    }

    @Cacheable(key = "#profileId", value = "mood")
    public GetMoodResponse getMoodById(String profileId) {
        Mood mood = moodRepository.getMood(profileId).orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
        return GetMoodResponse.builder().moodId(mood.getMoodId()).mood(mood.getMood()).description(mood.getDescription()).expiresAt(mood.getExpiresAt()).createdAt(mood.getCreatedAt()).visibility(mood.getVisibility()).build();
    }

    @CacheEvict(key = "#profileId", value = "mood")
    public DeleteMoodResponse deleteMood(String profileId) {
        if (!profileRepository.existsById(profileId)) throw new CustomExceptionHandler(ErrorCode.QUERY_FAILED);
        moodRepository.deleteById(profileId);
        return DeleteMoodResponse.builder().isSuccess(true).message("Delete success").build();
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
