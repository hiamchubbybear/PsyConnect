package dev.psyconnect.identity_service.controller;

import com.nimbusds.jose.JOSEException;
import com.nimbusds.jwt.JWTClaimsSet;
import dev.psyconnect.identity_service.apiresponse.ApiResponse;
import dev.psyconnect.identity_service.dto.response.TokenClaimsResponse;
import dev.psyconnect.identity_service.service.AuthenticationService;
import dev.psyconnect.identity_service.service.ExtractClaimService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.*;

import java.text.ParseException;
import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/auth/token")
@RequiredArgsConstructor
public class ExtractClaimController {
    private final AuthenticationService authenticationService;
    private final ExtractClaimService extractClaimService;

    @GetMapping("/claims")
    public Map<String, Object> getAllClaims(@RequestParam("token") String token) throws ParseException, JOSEException {
        JWTClaimsSet claims = extractClaimService.extractAllClaims(token);
        Map<String, Object> result = new HashMap<>();
        result.put("subject", claims.getSubject());
        result.put("accountId", claims.getClaim("accountId"));
        result.put("profileId", claims.getClaim("profileId"));
        result.put("role", extractClaimService.extractRole(token));
        result.put("scope", claims.getClaim("scope"));
        result.put("type", claims.getClaim("type"));
        result.put("platform", extractClaimService.extractPlatform(token));
        result.put("issuer", claims.getIssuer());
        result.put("expiration", claims.getExpirationTime());
        result.put("issuedAt", claims.getIssueTime());
        result.put("jwtId", claims.getJWTID());
        return result;
    }

    @GetMapping("/username")
    public Map<String, String> getUsername(@RequestParam("token") String token) throws ParseException, JOSEException {
        String username = extractClaimService.extractUsername(token);
        return Map.of("username", username);
    }

    @GetMapping("/platform")
    public Map<String, String> getPlatform(@RequestParam("token") String token) throws ParseException, JOSEException {
        String platform = extractClaimService.extractPlatform(token);
        return Map.of("platform", platform);
    }

    @GetMapping("/role")
    public Map<String, String> getRole(@RequestParam("token") String token) throws ParseException, JOSEException {
        String role = extractClaimService.extractRole(token);
        return Map.of("role", role);
    }

    @PostMapping("/token/claims")
    public ApiResponse<TokenClaimsResponse> extractClaims(@RequestBody String token) throws ParseException, JOSEException {
        return new ApiResponse<>(authenticationService.extractClaims(token));
    }

}
