package dev.psyconnect.profile_service.repository;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class DeleteResponse {
    private boolean isSuccess;
    private String message;
}
