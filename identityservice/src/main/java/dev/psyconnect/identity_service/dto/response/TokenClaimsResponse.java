package dev.psyconnect.identity_service.dto.response;

import lombok.Builder;
import lombok.Data;

@Data
@Builder
public class TokenClaimsResponse {
    private String subject;
    private String accountId;
    private String profileId;
    private String scope;
    private String provider;
    private String platform;
    private Long issuedAt;
    private Long expiresAt;
}
