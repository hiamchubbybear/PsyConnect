package dev.psyconnect.api_service.globalexceptionhandle;

import dev.psyconnect.api_service.apiresponse.ApiResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.HttpRequestMethodNotSupportedException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.MissingServletRequestParameterException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;

import java.nio.file.AccessDeniedException;

@ControllerAdvice
public class GlobalExceptionHandler {

    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ApiResponse<String>> handleGeneralException(Exception ex) {
        log.warn("General Exception: ", ex);
        ApiResponse<String> response =
                new ApiResponse<>(HttpStatus.INTERNAL_SERVER_ERROR.value(), "Server Error", null);
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }

    @ExceptionHandler(NullPointerException.class)
    public ResponseEntity<ApiResponse<String>> handleNullPointerException(NullPointerException ex) {
        log.warn("NullPointerException: ", ex);
        ApiResponse<String> response =
                new ApiResponse<>(HttpStatus.BAD_REQUEST.value(), "Null value encountered", null);
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    @ExceptionHandler(IllegalArgumentException.class)
    public ResponseEntity<ApiResponse<String>> handleIllegalArgumentException(IllegalArgumentException ex) {
        log.warn("IllegalArgumentException: ", ex);
        ApiResponse<String> response = new ApiResponse<>(HttpStatus.BAD_REQUEST.value(), "Illegal Argument", null);
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ApiResponse<String>> handleValidationException(MethodArgumentNotValidException ex) {
        log.warn("Validation Exception: ", ex);
        String errorMessage = ex.getBindingResult().getFieldErrors().stream()
                .map(fieldError -> fieldError.getField() + ": " + fieldError.getDefaultMessage())
                .findFirst()
                .orElse("Validation error");

        ApiResponse<String> response = new ApiResponse<>(HttpStatus.BAD_REQUEST.value(), "Validation Failed", null);
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    @ExceptionHandler(AccessDeniedException.class)
    public ResponseEntity<ApiResponse<String>> handleAccessDeniedException(AccessDeniedException ex) {
        log.warn("AccessDeniedException: ", ex);
        ApiResponse<String> response = new ApiResponse<>(HttpStatus.FORBIDDEN.value(), "Access denied", null);
        return ResponseEntity.status(HttpStatus.FORBIDDEN).body(response);
    }

    @ExceptionHandler(CustomExceptionHandler.class)
    public ResponseEntity<ApiResponse<?>> handleCustomException(CustomExceptionHandler customExceptionHandler) {
        log.warn("CustomExceptionHandler: ", customExceptionHandler);
        ErrorCode errorCode = customExceptionHandler.getErrorCode();
        return ResponseEntity.status(errorCode.getHttpStatus())
                .body(ApiResponse.builder()
                        .code(errorCode.getCode())
                        .message(errorCode.getMessage())
                        .build());
    }
}
