package dev.psyconnect.profile_service.dto.request;

import java.util.List;
import java.util.stream.Collectors;

import org.springframework.stereotype.Service;

import dev.psyconnect.profile_service.model.ActivityLog;
import dev.psyconnect.profile_service.repository.ActivityLogRepository;
import lombok.AccessLevel;
import lombok.RequiredArgsConstructor;
import lombok.experimental.FieldDefaults;

@Service
@RequiredArgsConstructor
@FieldDefaults(level = AccessLevel.PRIVATE, makeFinal = true)
public class ActivityLogService {
    ActivityLogRepository activityLogRepository;

    public ActivityLogResponse createLog(String profileId, ActivityLogRequest request) {
        long timestamp = System.currentTimeMillis();
        ActivityLog log = activityLogRepository.createLog(
                profileId,
                request.getAction(),
                timestamp,
                request.getTargetId(),
                request.getTargetType(),
                request.getExtraData());

        return ActivityLogResponse.builder()
                .id(log.getId())
                .profileId(log.getProfileId())
                .action(log.getAction())
                .timestamp(log.getTimestamp())
                .targetId(log.getTargetId())
                .targetType(log.getTargetType())
                .extraData(log.getExtraData())
                .build();
    }

    public List<ActivityLogResponse> getLogsByProfileId(String profileId) {
        return activityLogRepository.findByProfileId(profileId).stream()
                .map(log -> ActivityLogResponse.builder()
                        .id(log.getId())
                        .profileId(log.getProfileId())
                        .action(log.getAction())
                        .timestamp(log.getTimestamp())
                        .targetId(log.getTargetId())
                        .targetType(log.getTargetType())
                        .extraData(log.getExtraData())
                        .build())
                .collect(Collectors.toList());
    }

    public void deleteLog(String logId) {
        activityLogRepository.deleteById(logId);
    }
}
