package dev.psyconnect.identity_service.globalexceptionhandle;

import org.springframework.http.HttpStatus;

import lombok.*;

@Getter
@NoArgsConstructor
@AllArgsConstructor
public enum ErrorCode {
    // NOT FOUND _ RESPONSE
    USER_NOTFOUND(001, "User not found", HttpStatus.NOT_FOUND), // 404
    IMAGES_NOTFOUND(001, "Images not found", HttpStatus.NOT_FOUND), // 404
    TOKEN_NOTFOUND(001, "Token not found or expired", HttpStatus.NOT_FOUND), // 404
    RESOURCE_NOT_FOUND(001, "Resource not found", HttpStatus.NOT_FOUND), // 404
    ROLE_NOT_FOUND(001, "Role not exists or can't be find ", HttpStatus.NOT_FOUND), // 404

    // EXISTED
    EMAIL_ALREADY_EXISTS(002, "Email existed", HttpStatus.CONFLICT), // 409
    USERNAME_ALREADY_EXISTS(002, "Username existed", HttpStatus.CONFLICT), // 409

    // INVALID RESPONSE
    USERNAME_INVALID(
            003, "Username must be at least 4 characters and less than 20 characters", HttpStatus.BAD_REQUEST), // 400
    PASSWORD_INVALID(003, "Password must be at least 8 characters", HttpStatus.BAD_REQUEST), // 400
    TOKEN_INVALID(003, "Token is used or has expired", HttpStatus.NOT_ACCEPTABLE), // 407

    // NOT MATCH
    UUID_NOT_MATCH(004, "UUID you provided do not match with any user", HttpStatus.NO_CONTENT), // 205

    // AUTHENTICATED FAILED
    USER_UNAUTHENTICATED(005, "User authentication failed", HttpStatus.UNAUTHORIZED), // 401
    UNAUTHORIZED(005, "You don't have permission to access", HttpStatus.UNAUTHORIZED), // 401
    ACTIVATED_FAILED(005, "Activated failed", HttpStatus.UNAUTHORIZED), // 401

    // SYSTEM ERROR
    RUNTIME_ERROR(006, "Run time error", HttpStatus.TOO_MANY_REQUESTS), // 429
    PRODUCT_EXISTED_CART(006, "Product added into cart", HttpStatus.NOT_ACCEPTABLE),
    NOT_MATCH(006, "Account doesn't match with system", HttpStatus.NOT_ACCEPTABLE), // 406

    // SERVER
    SERVER_INTERNAL_ERROR(007, "Server Internal Error", HttpStatus.INTERNAL_SERVER_ERROR),

    // UNCATEGORIZED_EXCEPTION
    UNCATEGORIZED_EXCEPTION(007, "Uncategorize Exception", HttpStatus.I_AM_A_TEAPOT), // 418

    // NULL
    NULL_EXCEPTION(010, "Your input data is null", HttpStatus.GONE); // 410

    private int code;
    private String status;
    private HttpStatus statusCode;
}
