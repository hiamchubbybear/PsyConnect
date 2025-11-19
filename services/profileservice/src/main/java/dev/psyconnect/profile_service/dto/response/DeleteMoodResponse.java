package dev.psyconnect.profile_service.dto.response;

import lombok.*;

@AllArgsConstructor
@NoArgsConstructor
@Getter
@Setter
@Builder
public class DeleteMoodResponse {
    boolean isSuccess;
    String message;
}
