package dev.psyconnect.profile_service;

import java.time.LocalDateTime;
import java.time.ZoneOffset;
import java.util.Map;

import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.dto.response.MoodCreateResponse;
import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.profile_service.repository.MoodRepository;
import dev.psyconnect.profile_service.repository.ProfileRepository;
import dev.psyconnect.profile_service.service.DateTimeTemplate;
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

    public MoodCreateResponse createMoodByProfileId(String profileId, MoodCreateRequest request) {
        log.info("Profile Id {}", profileId);

        Map<String, Long> dateTime = DateTimeTemplate.MOOD_EXPIRES;
        long currentTime = dateTime.get("currentTimeMillis");
        long expiresTime = dateTime.get("expiresTimeMillis");

        if (!profileRepository.existsById(profileId)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);

        var response = moodRepository
                .createMood(
                        profileId,
                        request.mood,
                        request.moodDescription,
                        request.getVisibility(),
                        currentTime,
                        expiresTime)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));

        log.info("Creating mood for profile {}", response.getMood());

        return MoodCreateResponse.builder()
                .profileId(profileId)
                .isSuccess(true)
                .moodDescription(request.moodDescription)
                .mood(request.mood)
                .timeExpires(LocalDateTime.ofEpochSecond(expiresTime / 1000, 0, ZoneOffset.UTC))
                .build();
    }

    //    public MoodUpdateResponse updateMood(String profileId, MoodUpdateRequest request) {
    //        log.info("Profile Id {}", profileId);
    //        if (!moodRepository.existsById(request.)) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
    //        var response = moodRepository
    //                .createMood(profileId, request.mood, request.moodDescription)
    //                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.QUERY_FAILED));
    //        log.info("Creating mood for profile {}", response.getMood());
    //        if (response == null) throw new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND);
    //        return MoodCreateResponse.builder()
    //                .profileId(profileId)
    //                .isSuccess(true)
    //                .moodDescription(request.moodDescription)
    //                .mood(request.mood)
    //                .build();
    //    }
}
