package dev.psyconnect.profile_service.dto.response;

import java.io.Serializable;

import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class UserSettingResponse implements Serializable {
    private String profileId;
    private boolean isSuccess;
}
