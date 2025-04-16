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

    // PasswordEncoder is initialized with strength 10
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(10);
    }

    // This method generates a JWT token
    public String generateToken(AuthenticationRequest authenticationRequest, String loginType) {
        // If login is mobile set expire times about 1 month  else set it to 30 minutes
        TIME_EXPIRED = (loginType.equals("mobile")) ? TIME_EXPIRED * 2 * 24 * 30 : TIME_EXPIRED;
        // Get the username from the authentication request
        Account userAccountObject = userAccountRepository
                .findByUsername(authenticationRequest.getUsername())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        var username = authenticationRequest.getUsername();
        // Get the role of the user based on the username
        var roles = userAccountRepository
                .findByUsername(username)
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.ROLE_NOT_FOUND))
                .getRole()
                .stream()
                .findFirst()
                .get()
                .getName();

        // Create the JWT header using the RS256 algorithm
        JWSHeader header = new JWSHeader(JWSAlgorithm.HS512);

        /*
        Build the claims for the JWT (user-related information)
        -> Jwt Structure :
        -  Header
        -  Payload
        -  Signature
        -> Header :
        	Contain 2 parts : signing algorithm being used ( HMAC SHA256 or RSA )
        	and which type of token.
        Example :
        {
        	"alg": "HS256",
        	"type": "JWT"

        }
        Build the claims for the JWT (user-related information)
        -> Jwt Structure :
        -  Header
        -  Payload
        -  Signature
        -> Header :
        	Contain 2 parts : signing algorithm being used ( HMAC SHA256 or RSA )
        	and which type of token.
        Example :
        {
        	"alg": "HS256",
        	"typ": "JWT"
        }
        -> Payload :
        	The second part of the token is the payload, which contains the claims.
        	Claims are statements about an entity (typically, the user) and additional data.
        	There are three types of claims: registered, public, and private claims.
        Example :
        	{
        	"sub": "1234567890",
        	"name": "John Doe",
        	"admin": true
        	}
        -> Signature :
        	To create the signature part you have to take the encoded header,
        	the encoded payload, a secret, the algorithm specified in the header, and sign that.
        Example :
        	HMACSHA256(
        	base64UrlEncode(header) + "." +
        	base64UrlEncode(payload),
        	secret)
         */

        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject(authenticationRequest.getUsername())
                .expirationTime(new Date(new Date().getTime() + TIME_EXPIRED))
                .issueTime(new Date())
                .jwtID(UUID.randomUUID().toString())
                .issuer("PsyConnect Authentication Service")
                .claim("scope", buildScope(roles))
                .claim("type", loginType)
                .claim("accountId", userAccountObject.getAccountId())
                .claim("profileId", userAccountObject.getProfileId())
                .build();
        Payload payload = new Payload(claimsSet.toJSONObject());
        JWSObject jwsObject = new JWSObject(header, payload);

        try {
            jwsObject.sign(new MACSigner(SIGNER_KEY));
        } catch (JOSEException e) {
            throw new RuntimeException("Sign Jwt failed with exception", e);
        }

        // Return the serialized JWT
        return jwsObject.serialize();
    }

    public String generateGoogleAuthToken(GoogleAuthenticationRequest googleAuthenticationRequest, String loginType) {
        var user = userAccountRepository
                .findByEmail(googleAuthenticationRequest.getEmail())
                .orElse(null);
        assert user != null;
        var roles = user.getRole().stream().findFirst().get().getName();
        UUID uuid = user.getAccountId();

        // Create the JWT header using the RS256 algorithm
        JWSHeader header = new JWSHeader(JWSAlgorithm.HS512);

        /*
        Build the claims for the JWT (user-related information)
        -> Jwt Structure :
        -  Header
        -  Payload
        -  Signature
        -> Header :
        	Contain 2 parts : signing algorithm being used ( HMAC SHA256 or RSA )
        	and which type of token.
        Example :
        {
        	"alg": "HS256",
        	"type": "JWT"

        }
        Build the claims for the JWT (user-related information)
        -> Jwt Structure :
        -  Header
        -  Payload
        -  Signature
        -> Header :
        	Contain 2 parts : signing algorithm being used ( HMAC SHA256 or RSA )
        	and which type of token.
        Example :
        {
        	"alg": "HS256",
        	"typ": "JWT"
        }
        -> Payload :
        	The second part of the token is the payload, which contains the claims.
        	Claims are statements about an entity (typically, the user) and additional data.
        	There are three types of claims: registered, public, and private claims.
        Example :
        	{
        	"sub": "1234567890",
        	"name": "John Doe",
        	"admin": true
        	}
        -> Signature :
        	To create the signature part you have to take the encoded header,
        	the encoded payload, a secret, the algorithm specified in the header, and sign that.
        Example :
        	HMACSHA256(
        	base64UrlEncode(header) + "." +
        	base64UrlEncode(payload),
        	secret)
         */

        JWTClaimsSet claimsSet = new JWTClaimsSet.Builder()
                .subject(googleAuthenticationRequest.getEmail())
                .expirationTime(new Date(new Date().getTime() + TIME_EXPIRED))
                .issueTime(new Date())
                .jwtID(UUID.randomUUID().toString())
                .issuer("PsyConnect Authentication Service")
                .claim("scope", buildScope(roles))
                .claim("type", loginType)
                .claim("accountId", uuid.toString())
                .claim("profileId", user.getProfileId().toString())
                .build();
        Payload payload = new Payload(claimsSet.toJSONObject());
        JWSObject jwsObject = new JWSObject(header, payload);

        try {
            jwsObject.sign(new MACSigner(SIGNER_KEY));
        } catch (JOSEException e) {
            throw new RuntimeException("Sign Jwt failed with exception", e);
        }

        // Return the serialized JWT
        return jwsObject.serialize();
    }

    // Use to verify token
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

    // Use to authenticate account
    public AuthenticationResponse authenticate(AuthenticationRequest authenticationRequest, String loginType)
            throws AuthenticationException {
        var user = userAccountRepository
                .findByUsername(authenticationRequest.getUsername())
                .orElseThrow(() -> new CustomExceptionHandler(ErrorCode.USER_NOT_FOUND));
        log.debug("User request token is {}", authenticationRequest.getUsername());
        var password = authenticationRequest.getPassword();
        // Password not found throw
        if (password == null) {
            throw new IllegalArgumentException("Username and password are required");
        }
        // Password does not match with user account
        else if (!passwordEncoder().matches(password, user.getPassword())) {
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

        // Change prefix to match with OpenID Standard & Oauth2
        builder.append("role.").append(role).append(" ");

        // Find the role in the repository and retrieve its permissions
        roleRepository
                .findByName(role)
                .ifPresentOrElse(
                        roleEntity -> {
                            // Map each permission to the Spring Security format and append to the builder
                            roleEntity.getPermissions().stream()
                                    .map(permission ->
                                             permission.getName()) // Prefix permissions for clarity
                                    .forEach(permission ->
                                            builder.append(permission).append(" "));
                        },
                        () -> {
                            throw new IllegalArgumentException("Invalid role: " + role);
                        });

        log.debug("Built scope for role {}: {}", role, builder.toString().trim());

        // Return the concatenated authorities as a trimmed string
        return builder.toString().trim();
    }

    public LogoutResponse logout(LogoutRequest request) throws JOSEException, ParseException {
        if (request.getToken() == null || request.getToken().isEmpty()) {
            return new LogoutResponse(false); // Ignore empty or null tokens
        }
        try {
            SignedJWT signedJWT = verifyToken(request.getToken());
            Date expirationDate = signedJWT.getJWTClaimsSet().getExpirationTime();

            // If the token is expired, add it to the blacklist
            if (expirationDate.after(new Date())) {
                blackListTokenRepository.save(
                        BlackListToken.builder().token(request.getToken()).build());
                return new LogoutResponse(true);
            }
        } catch (ParseException | JOSEException e) {
            // If the token is invalid (tampered or incorrectly formatted), blacklist it
            blackListTokenRepository.save(
                    BlackListToken.builder().token(request.getToken()).build());
            return new LogoutResponse(true);
        }
        return new LogoutResponse(false);
    }

    public String extractUsername(String token) throws ParseException, JOSEException {
        return extractClaim(token, JWTClaimsSet::getSubject);
    }

    // Extract the expiration date from the token
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
    private JWTClaimsSet extractAllClaims(String token) throws ParseException, JOSEException {
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
