package dev.psyconnect.identity_service.globalexceptionhandle;

import org.springframework.http.HttpStatus;

import lombok.Getter;
import lombok.RequiredArgsConstructor;

@Getter
@RequiredArgsConstructor
public enum ErrorCode {
    // NOT FOUND - 404
    USER_NOT_FOUND(101, "User not found", HttpStatus.NOT_FOUND),
    IMAGES_NOT_FOUND(102, "Images not found", HttpStatus.NOT_FOUND),
    TOKEN_NOT_FOUND(103, "Token not found or expired", HttpStatus.NOT_FOUND),
    RESOURCE_NOT_FOUND(104, "Resource not found", HttpStatus.NOT_FOUND),
    ROLE_NOT_FOUND(105, "Role does not exist or cannot be found", HttpStatus.NOT_FOUND),

    // ALREADY EXISTS - 409
    EMAIL_ALREADY_EXISTS(201, "Email already exists", HttpStatus.CONFLICT),
    USERNAME_ALREADY_EXISTS(202, "Username already exists", HttpStatus.CONFLICT),

    // INVALID REQUEST - 401
    USERNAME_INVALID(301, "Username must be between 4 and 20 characters", HttpStatus.BAD_REQUEST),
    PASSWORD_INVALID(302, "Password must be at least 8 characters", HttpStatus.BAD_REQUEST),
    TOKEN_INVALID(303, "Token is used or has expired", HttpStatus.BAD_REQUEST),

    // NOT MATCH - 404
    UUID_NOT_MATCH(401, "UUID provided does not match any user", HttpStatus.NOT_FOUND),

    // AUTHENTICATION FAILURE - 401
    USER_UNAUTHENTICATED(501, "User authentication failed", HttpStatus.UNAUTHORIZED),
    UNAUTHORIZED(502, "You don't have permission to access this resource", HttpStatus.UNAUTHORIZED),
    ACTIVATION_FAILED(503, "Activation failed", HttpStatus.UNAUTHORIZED),

    // SYSTEM ERROR - 429 / 406
    RUNTIME_ERROR(601, "Runtime error", HttpStatus.INTERNAL_SERVER_ERROR),
    PRODUCT_ALREADY_IN_CART(602, "Product already added to cart", HttpStatus.NOT_ACCEPTABLE),
    ACCOUNT_NOT_MATCH(603, "Account does not match the system", HttpStatus.NOT_ACCEPTABLE),

    // SERVER ERROR - 500
    SERVER_INTERNAL_ERROR(701, "Internal server error", HttpStatus.INTERNAL_SERVER_ERROR),

    // GENERAL EXCEPTIONS - 500
    UNCATEGORIZED_EXCEPTION(801, "Uncategorized exception", HttpStatus.INTERNAL_SERVER_ERROR),
    DELETE_ACCOUNT_FAILED(802, "An error occurred while deleting the account", HttpStatus.INTERNAL_SERVER_ERROR),

    // NULL INPUT - 410
    NULL_EXCEPTION(901, "Your input data is null", HttpStatus.GONE);

    private final int code;
    private final String message;
    private final HttpStatus httpStatus;
}
