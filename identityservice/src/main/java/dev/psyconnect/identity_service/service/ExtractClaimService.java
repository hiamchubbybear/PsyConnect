package dev.psyconnect.identity_service.service;

import com.nimbusds.jose.JOSEException;
import com.nimbusds.jwt.JWTClaimsSet;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;

import java.text.ParseException;

@Service
@RequiredArgsConstructor
public class ExtractClaimService {

    private final AuthenticationService authenticationService;

    public JWTClaimsSet extractAllClaims(String token) throws ParseException, JOSEException {
        return authenticationService.extractAllClaims(token);
    }

    public String extractUsername(String token) throws ParseException, JOSEException {
        return authenticationService.extractUsername(token);
    }

    public String extractScope(String token) throws ParseException, JOSEException {
        JWTClaimsSet claims = authenticationService.extractAllClaims(token);
        return (String) claims.getClaim("scope");
    }

    public String extractRole(String token) throws ParseException, JOSEException {
        String scope = extractScope(token);
        if (scope != null && scope.startsWith("role.")) {
            return scope.split(" ")[0].substring(5);
        }
        return null;
    }

    public String extractPlatform(String token) throws ParseException, JOSEException {
        JWTClaimsSet claims = authenticationService.extractAllClaims(token);
        return (String) claims.getClaim("platform");
    }
}
