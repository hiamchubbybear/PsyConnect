package dev.psyconnect.identity_service.service;

import java.text.ParseException;
import java.util.Date;
import java.util.UUID;
import java.util.function.Function;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import com.nimbusds.jose.JOSEException;
import com.nimbusds.jose.JWSAlgorithm;
import com.nimbusds.jose.JWSHeader;
import com.nimbusds.jose.JWSObject;
import com.nimbusds.jose.JWSVerifier;
import com.nimbusds.jose.Payload;
import com.nimbusds.jose.crypto.MACSigner;
import com.nimbusds.jose.crypto.MACVerifier;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;

import dev.psyconnect.identity_service.dto.request.AuthenticationFilterRequest;
import dev.psyconnect.identity_service.dto.request.AuthenticationRequest;
import dev.psyconnect.identity_service.dto.request.Oauth2AuthenticationRequest;
import dev.psyconnect.identity_service.dto.response.AuthenticationResponse;
import dev.psyconnect.identity_service.dto.response.LogoutRequest;
import dev.psyconnect.identity_service.dto.response.LogoutResponse;
import dev.psyconnect.identity_service.dto.response.TokenClaimsResponse;
import dev.psyconnect.identity_service.dto.response.v1.AuthenticationResponseV1;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.globalexceptionhandle.ErrorCode;
import dev.psyconnect.identity_service.model.Account;
import dev.psyconnect.identity_service.model.BlackListToken;
import dev.psyconnect.identity_service.repository.BlackListTokenRepository;
import dev.psyconnect.identity_service.repository.RoleRepository;
import dev.psyconnect.identity_service.repository.TokenRepository;
import dev.psyconnect.identity_service.repository.UserAccountRepository;
import lombok.extern.slf4j.Slf4j;

@Service
@Slf4j
public class AuthenticationService {

    private final TokenRepository tokenRepository;

    private final BlackListTokenRepository blackListTokenRepository;
    private final TokenService tokenService;

    @Value("${SIGNER_KEY}")
    private String SIGNER_KEY;

    static long TIME_EXPIRED = 30 * 60 * 60 * 100;

    final RoleRepository roleRepository;
    final UserAccountRepository userAccountRepository;

    public AuthenticationService(
            RoleRepository roleRepository,
            UserAccountRepository userAccountRepository,
            BlackListTokenRepository blackListTokenRepository,
            TokenService tokenService,
            TokenRepository tokenRepository) {
        this.roleRepository = roleRepository;
        this.tokenService = tokenService;
        this.userAccountRepository = userAccountRepository;
        this.blackListTokenRepository = blackListTokenRepository;
        this.tokenRepository = tokenRepository;
    }

    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(10);
    }

    public String generateToken(String username, String provider, String platForm) {
        Account account = userAccountRepository
                .findByUsername(username)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));

        String role = null;
        if (account.getRole() != null && !account.getRole().isEmpty()) {
            role = account.getRole().iterator().next().getName();
        } else {
            throw new CustomExceptionHandler(ErrorCode.ROLE_NOT_FOUND);
        }
        return createJwtToken(
                username,
                account.getAccountId().toString(),
                account.getProfileId().toString(),
                role,
                provider,
                platForm);
    }

    public AuthenticationResponse generateOAuth2LoginToken(
            Oauth2AuthenticationRequest request, String provider, String platForm) {
        Account account = userAccountRepository
                .findByEmail(request.getEmail())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        log.info("Found session token {}, ", account.getSession().toString());
        if (!request.getSessionToken().trim().equals(account.getSession().trim()))
            throw new CustomExceptionHandler(ErrorCode.USER_UNAUTHENTICATED);
        if (!request.getProvider().toUpperCase().equals(account.getProvider().toString()))
            throw new CustomExceptionHandler(ErrorCode.AUTHENTICATE_REQUIRED_DENY);
        String role = null;
        if (account.getRole() != null && !account.getRole().isEmpty()) {
            role = account.getRole().iterator().next().getName();
        } else {
            throw new CustomExceptionHandler(ErrorCode.ROLE_NOT_FOUND);
        }
        String jwtToken = createJwtToken(
                request.getEmail(),
                account.getAccountId().toString(),
                account.getProfileId().toString(),
                role,
                provider,
                platForm);
        int effectRow =
                userAccountRepository.updateSessionIdPostLogin(account.getEmail(), request.getSessionToken(), "");
        if (effectRow <= 0) throw new CustomExceptionHandler(ErrorCode.DELETE_SESSION_FAILED);
        return AuthenticationResponse.builder()
                .token(jwtToken)
                .isSuccessful(true)
                .build();
    }

    private String createJwtToken(
            String subject, String accountId, String profileId, String role, String provider, String platForm) {
        long expiration =
                platForm.toLowerCase().equalsIgnoreCase("mobile") ? TIME_EXPIRED * 2 * 24 * 30L : TIME_EXPIRED;

        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject(subject)
                .expirationTime(new Date(System.currentTimeMillis() + expiration))
                .issueTime(new Date())
                .jwtID(UUID.randomUUID().toString())
                .issuer("PsyConnect Authentication Service")
                .claim("scope", buildScope(role))
                .claim("type", provider)
                .claim("accountId", accountId)
                .claim("profileId", profileId)
                .claim("platform", platForm)
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

    public AuthenticationResponse authenticate(
            AuthenticationRequest authenticationRequest, String provider, String clientPlatform) {

        var user = userAccountRepository
                .findByUsername(authenticationRequest.getUsername())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        log.debug("User request token is {}", authenticationRequest.getUsername());
        var password = authenticationRequest.getPassword();

        if (password == null) throw new CustomExceptionHandler(ErrorCode.PASSWORD_INVALID);
        else if (!passwordEncoder().matches(password, user.getPassword()))
            throw new CustomExceptionHandler(ErrorCode.PASSWORD_INVALID);
        else {
            var response = AuthenticationResponse.builder()
                    .isSuccessful(true)
                    .token(generateToken(authenticationRequest.getUsername(), provider, clientPlatform))
                    .build();
            return response;
        }
    }

    public AuthenticationResponseV1 authenticateV1(
            AuthenticationRequest authenticationRequest, String provider, String clientPlatform) {
        var password = authenticationRequest.getPassword();
        var user = userAccountRepository
                .findByUsername(authenticationRequest.getUsername())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        if (password == null) throw new CustomExceptionHandler(ErrorCode.PASSWORD_INVALID);
        else if (!passwordEncoder().matches(password, user.getPassword()))
            throw new CustomExceptionHandler(ErrorCode.PASSWORD_WRONG);
        String refreshToken = tokenService.generateRefreshToken(authenticationRequest.getUsername());
        var response = AuthenticationResponseV1.builder()
                .isSuccessful(true)
                .token(generateToken(authenticationRequest.getUsername(), provider, clientPlatform))
                .refreshToken(refreshToken)
                .build();
        return response;
    }

    public AuthenticationResponseV1 refreshTokenV1(
            String username, String refreshToken, String provider, String clientPlatform) {
        var token = tokenRepository
                .findById(username)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.TOKEN_INVALID));
        String newRefreshToken = tokenService.checkAndReGenerateRefreshToken(username, refreshToken);
        var response = AuthenticationResponseV1.builder()
                .isSuccessful(true)
                .token(generateToken(username, provider, clientPlatform))
                .refreshToken(newRefreshToken)
                .build();
        return response;
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

    public <T> T extractClaim(String token, Function<JWTClaimsSet, T> claimsResolver)
            throws ParseException, JOSEException {
        final JWTClaimsSet claims = extractAllClaims(token);
        return claimsResolver.apply(claims);
    }

    public JWTClaimsSet extractAllClaims(String token) throws ParseException, JOSEException {
        return verifyToken(token).getJWTClaimsSet();
    }

    private Boolean isTokenExpired(String token) throws ParseException, JOSEException {
        return extractExpiration(token).before(new Date());
    }

    public Boolean validateToken(String token, AuthenticationFilterRequest authenticationRequest)
            throws ParseException, JOSEException {
        final String username = extractUsername(token);
        return (username.equals(authenticationRequest.getUsername()) && !isTokenExpired(token));
    }

    public TokenClaimsResponse extractClaims(String token) throws ParseException, JOSEException {
        var claimsSet = verifyToken(token).getJWTClaimsSet();

        return TokenClaimsResponse.builder()
                .subject(claimsSet.getSubject())
                .accountId((String) claimsSet.getClaim("accountId"))
                .profileId((String) claimsSet.getClaim("profileId"))
                .scope((String) claimsSet.getClaim("scope"))
                .provider((String) claimsSet.getClaim("type"))
                .platform((String) claimsSet.getClaim("platform"))
                .issuedAt(claimsSet.getIssueTime().getTime())
                .expiresAt(claimsSet.getExpirationTime().getTime())
                .build();
    }

    public boolean isTokenValid(String token) {
        return !blackListTokenRepository.existsByToken(token);
    }

    public String generateTokenV1(AuthenticationRequest authenticationRequest, String provider, String platForm) {
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
                provider,
                platForm);
    }

    public AuthenticationResponse generateOAuth2LoginTokenV1(
            Oauth2AuthenticationRequest request, String provider, String platForm) {
        Account account = userAccountRepository
                .findByEmail(request.getEmail())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        log.info("Found session token {}, ", account.getSession().toString());
        if (!request.getSessionToken().trim().equals(account.getSession().trim()))
            throw new CustomExceptionHandler(ErrorCode.USER_UNAUTHENTICATED);
        if (!request.getProvider().toUpperCase().equals(account.getProvider().toString()))
            throw new CustomExceptionHandler(ErrorCode.AUTHENTICATE_REQUIRED_DENY);
        String role = null;
        if (account.getRole() != null && !account.getRole().isEmpty()) {
            role = account.getRole().iterator().next().getName();
        } else {
            throw new CustomExceptionHandler(ErrorCode.ROLE_NOT_FOUND);
        }
        String jwtToken = createJwtToken(
                request.getEmail(),
                account.getAccountId().toString(),
                account.getProfileId().toString(),
                role,
                provider,
                platForm);
        int effectRow =
                userAccountRepository.updateSessionIdPostLogin(account.getEmail(), request.getSessionToken(), "");
        if (effectRow <= 0) throw new CustomExceptionHandler(ErrorCode.DELETE_SESSION_FAILED);
        return AuthenticationResponse.builder()
                .token(jwtToken)
                .isSuccessful(true)
                .build();
    }
}
