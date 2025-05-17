package dev.psyconnect.identity_service.service;

import java.text.ParseException;
import java.util.Date;
import java.util.Map;
import java.util.UUID;
import java.util.function.Function;
import javax.naming.AuthenticationException;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.Authentication;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.stereotype.Service;

import com.nimbusds.jose.*;
import com.nimbusds.jose.crypto.MACSigner;
import com.nimbusds.jose.crypto.MACVerifier;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;

import dev.psyconnect.identity_service.dto.request.AuthenticationFilterRequest;
import dev.psyconnect.identity_service.dto.request.AuthenticationRequest;
import dev.psyconnect.identity_service.dto.request.GoogleAuthenticationRequest;
import dev.psyconnect.identity_service.dto.response.AuthenticationResponse;
import dev.psyconnect.identity_service.dto.response.LogoutRequest;
import dev.psyconnect.identity_service.dto.response.LogoutResponse;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.model.Account;
import dev.psyconnect.identity_service.model.BlackListToken;
import dev.psyconnect.identity_service.repository.BlackListTokenRepository;
import dev.psyconnect.identity_service.repository.RoleRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import lombok.extern.slf4j.Slf4j;

@Service
@Slf4j
public class AuthenticationService {

    private final BlackListTokenRepository blackListTokenRepository;

    @Value("${SIGNER_KEY}")
    private String SIGNER_KEY;
    // 30 Minutes
    static long TIME_EXPIRED = 30 * 60 * 60 * 100;

    final RoleRepository roleRepository;
    final UserAccountRepository userAccountRepository;

    @Autowired
    public AuthenticationService(
            RoleRepository roleRepository,
            UserAccountRepository userAccountRepository,
            BlackListTokenRepository blackListTokenRepository) {
        this.roleRepository = roleRepository;
        this.userAccountRepository = userAccountRepository;
        this.blackListTokenRepository = blackListTokenRepository;
    }

    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(10);
    }

    public String generateToken(AuthenticationRequest authenticationRequest, String loginType) {
        Account account = userAccountRepository
                .findByUsername(authenticationRequest.getUsername())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));

        String role = null;
        if (account.getRole() != null && !account.getRole().isEmpty()) {
            role = account.getRole().iterator().next().getName();
        } else {
            throw new CustomExceptionHandler(ErrorCode.ROLE_NOT_FOUND);
        }
        return createJwtToken(
                authenticationRequest.getUsername(),
                account.getAccountId().toString(),
                account.getProfileId().toString(),
                role,
                loginType);
    }

    public String generateGoogleAuthToken(GoogleAuthenticationRequest request, String loginType) {
        Account account = userAccountRepository
                .findByEmail(request.getEmail())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));

        String role = null;
        if (account.getRole() != null && !account.getRole().isEmpty()) {
            role = account.getRole().iterator().next().getName();
        } else {
            throw new CustomExceptionHandler(ErrorCode.ROLE_NOT_FOUND);
        }
        return createJwtToken(
                request.getEmail(),
                account.getAccountId().toString(),
                account.getProfileId().toString(),
                role,
                loginType);
    }

    private String createJwtToken(String subject, String accountId, String profileId, String role, String loginType) {
        long expiration = loginType.equalsIgnoreCase("mobile") ? TIME_EXPIRED * 2 * 24 * 30L : TIME_EXPIRED;

        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject(subject)
                .expirationTime(new Date(System.currentTimeMillis() + expiration))
                .issueTime(new Date())
                .jwtID(UUID.randomUUID().toString())
                .issuer("PsyConnect Authentication Service")
                .claim("scope", buildScope(role))
                .claim("type", loginType)
                .claim("accountId", accountId)
                .claim("profileId", profileId)
                .build();

        JWSObject jwsObject = new JWSObject(new JWSHeader(JWSAlgorithm.HS512), new Payload(claimsSet.toJSONObject()));

        try {
            jwsObject.sign(new MACSigner(SIGNER_KEY));
            return jwsObject.serialize();
        } catch (JOSEException e) {
            throw new RuntimeException("Sign Jwt failed", e);
        }
    }

    public SignedJWT verifyToken(String token) throws JOSEException, ParseException {
        JWSVerifier verifier = new MACVerifier(SIGNER_KEY);
        SignedJWT signedJWT = SignedJWT.parse(token);
        var verifyResult = signedJWT.verify(verifier);
        Date expirationDate = signedJWT.getJWTClaimsSet().getExpirationTime();
        if (expirationDate.before(new Date())) {
            throw new ParseException("Expired JWT", 0);
        }
        return signedJWT;
    }

    public AuthenticationResponse authenticate(AuthenticationRequest authenticationRequest, String loginType)
            throws AuthenticationException {
        var user = userAccountRepository
                .findByUsername(authenticationRequest.getUsername())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        log.debug("User request token is {}", authenticationRequest.getUsername());
        var password = authenticationRequest.getPassword();

        if (password == null) {
            throw new IllegalArgumentException("Username and password are required");
        } else if (!passwordEncoder().matches(password, user.getPassword())) {
            throw new IllegalArgumentException("Password does not match");
        } else {
            var response = AuthenticationResponse.builder()
                    .isSuccessful(true)
                    .token(generateToken(authenticationRequest, loginType))
                    .build();
            log.debug("Token is {}", response.getToken());
            return response;
        }
    }

    private String buildScope(String role) {
        StringBuilder builder = new StringBuilder();

        builder.append("role.").append(role).append(" ");

        roleRepository
                .findByName(role)
                .ifPresentOrElse(
                        roleEntity -> {
                            roleEntity.getPermissions().stream()
                                    .map(permission -> permission.getName())
                                    .forEach(permission ->
                                            builder.append(permission).append(" "));
                        },
                        () -> {
                            throw new IllegalArgumentException("Invalid role: " + role);
                        });

        log.debug("Built scope for role {}: {}", role, builder.toString().trim());

        return builder.toString().trim();
    }

    public LogoutResponse logout(LogoutRequest request) throws JOSEException, ParseException {
        if (request.getToken() == null || request.getToken().isEmpty()) {
            return new LogoutResponse(false);
        }
        try {
            SignedJWT signedJWT = verifyToken(request.getToken());
            Date expirationDate = signedJWT.getJWTClaimsSet().getExpirationTime();

            if (expirationDate.after(new Date())) {
                blackListTokenRepository.save(
                        BlackListToken.builder().token(request.getToken()).build());
                return new LogoutResponse(true);
            }
        } catch (ParseException | JOSEException e) {

            blackListTokenRepository.save(
                    BlackListToken.builder().token(request.getToken()).build());
            return new LogoutResponse(true);
        }
        return new LogoutResponse(false);
    }

    public String extractUsername(String token) throws ParseException, JOSEException {
        return extractClaim(token, JWTClaimsSet::getSubject);
    }

    public Date extractExpiration(String token) throws ParseException, JOSEException {
        return extractClaim(token, JWTClaimsSet::getExpirationTime);
    }

    // Extract a specific claim from the token
    public <T> T extractClaim(String token, Function<JWTClaimsSet, T> claimsResolver)
            throws ParseException, JOSEException {
        final JWTClaimsSet claims = extractAllClaims(token);
        return claimsResolver.apply(claims);
    }

    // Extract all claims from the token
    public JWTClaimsSet extractAllClaims(String token) throws ParseException, JOSEException {
        return verifyToken(token).getJWTClaimsSet();
    }

    // Check if the token is expired
    private Boolean isTokenExpired(String token) throws ParseException, JOSEException {
        return extractExpiration(token).before(new Date());
    }

    // Validate the token against user details and expiration
    public Boolean validateToken(String token, String username) throws ParseException, JOSEException {
        final String extractedUsername = extractUsername(token);
        return (extractedUsername.equals(username) && !isTokenExpired(token));
    }

    // Validate the token against user details and expiration
    public Boolean validateToken(String token, AuthenticationFilterRequest authenticationRequest)
            throws ParseException, JOSEException {
        final String username = extractUsername(token);
        return (username.equals(authenticationRequest.getUsername()) && !isTokenExpired(token));
    }

    public static String extractUsernameAuthenticationObject(Authentication authentication) {
        if (authentication == null) {
            throw new CustomExceptionHandler(ErrorCode.UNAUTHORIZED);
        }
        return authentication.getName();
    }

    public boolean isTokenInvalid(String token) {
        return blackListTokenRepository.existsByToken(token);
    }

    public static String extractUUIDClaim(Authentication authentication) {
        if (authentication == null) {
            throw new CustomExceptionHandler(ErrorCode.UNAUTHORIZED);
        }
        Object principal = authentication.getPrincipal();
        if (principal instanceof org.springframework.security.oauth2.jwt.Jwt) {
            org.springframework.security.oauth2.jwt.Jwt jwt = (Jwt) principal;
            Map<String, Object> returnValue = (jwt.getClaims());
            String uuid = returnValue.get("uuid").toString();
            log.debug("UUID is {}", uuid);
            return uuid;
        }
        throw new CustomExceptionHandler(ErrorCode.USER_UNAUTHENTICATED);
    }
}
