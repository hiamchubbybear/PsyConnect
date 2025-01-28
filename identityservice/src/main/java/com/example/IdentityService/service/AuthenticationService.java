package com.example.IdentityService.service;

import java.text.ParseException;
import java.util.Date;
import java.util.UUID;
import javax.naming.AuthenticationException;

import org.springframework.beans.InvalidPropertyException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;

import com.example.IdentityService.dto.request.AuthenticationRequest;
import com.example.IdentityService.dto.request.GoogleAuthenticationRequest;
import com.example.IdentityService.dto.response.AuthenticationResponse;
import com.example.IdentityService.globalexceptionhandle.CustomExceptionHandler;
import com.example.IdentityService.globalexceptionhandle.ErrorCode;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.repository.RoleRepository;
import com.example.IdentityService.repository.UserAccountRepository;
import com.nimbusds.jose.*;
import com.nimbusds.jose.crypto.MACSigner;
import com.nimbusds.jose.crypto.MACVerifier;
import com.nimbusds.jwt.JWTClaimsSet;
import com.nimbusds.jwt.SignedJWT;

import lombok.extern.slf4j.Slf4j;

@Service
@Slf4j
public class AuthenticationService {

    @Value("${SIGNER_KEY}")
    private String SIGNER_KEY;

    static long TIME_EXPIRED = 30 * 60 * 60 * 100;

    final RoleRepository roleRepository;
    final UserAccountRepository userAccountRepository;

    @Autowired
    public AuthenticationService(RoleRepository roleRepository, UserAccountRepository userAccountRepository) {
        this.roleRepository = roleRepository;
        this.userAccountRepository = userAccountRepository;
    }

    // PasswordEncoder is initialized with strength 10
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(10);
    }

    // This method generates a JWT token
    public String generateToken(AuthenticationRequest authenticationRequest, String loginType) {
        // Get the username from the authentication request
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
        UUID uuid = user.getUserId();

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
                .claim("uuid", uuid.toString())
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
        if (expirationDate.after(new Date())) {
            throw new ParseException("Expired JWT", 0);
        }
        return signedJWT;
    }

    // Use to authenticate account
    public AuthenticationResponse authenticate(AuthenticationRequest authenticationRequest, String loginType)
            throws AuthenticationException {
        var user = userAccountRepository
                .findByUsername(authenticationRequest.getUsername())
                .orElseThrow(() -> new InvalidPropertyException(
                        UserAccount.class, "User not found", authenticationRequest.getUsername()));

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

        // Prefix the role with "ROLE_" to match Spring Security conventions
        builder.append("ROLE_").append(role).append(" ");

        // Find the role in the repository and retrieve its permissions
        roleRepository
                .findByName(role)
                .ifPresentOrElse(
                        roleEntity -> {
                            // Map each permission to the Spring Security format and append to the builder
                            roleEntity.getPermissions().stream()
                                    .map(permission ->
                                            "PERMISSION_" + permission.getName()) // Prefix permissions for clarity
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
}
