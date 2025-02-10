package dev.psyconnect.identity_service.globalexceptionhandle;

import java.nio.file.AccessDeniedException;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.dao.DataAccessException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.HttpRequestMethodNotSupportedException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.MissingServletRequestParameterException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.servlet.NoHandlerFoundException;

import dev.psyconnect.identity_service.apiresponse.ApiResponse;

@ControllerAdvice
public class GlobalExceptionHandler {

    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    // 1. Xử lý ngoại lệ chung (Exception)
    @ExceptionHandler(Exception.class)
    public ResponseEntity<ApiResponse<String>> handleGeneralException(Exception ex) {
        log.warn("General Exception: ", ex);
        ApiResponse<String> response =
                new ApiResponse<>(HttpStatus.INTERNAL_SERVER_ERROR.value(), ex.getMessage(), "Internal Server Error");
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }

    // 2. Xử lý ngoại lệ NullPointerException
    @ExceptionHandler(NullPointerException.class)
    public ResponseEntity<ApiResponse<String>> handleNullPointerException(NullPointerException ex) {
        log.warn("NullPointerException: ", ex);
        ApiResponse<String> response =
                new ApiResponse<>(HttpStatus.BAD_REQUEST.value(), "Null value encountered", "Bad Request");
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    // 3. Xử lý ngoại lệ IllegalArgumentException
    @ExceptionHandler(IllegalArgumentException.class)
    public ResponseEntity<ApiResponse<String>> handleIllegalArgumentException(IllegalArgumentException ex) {
        log.warn("IllegalArgumentException: ", ex);
        ApiResponse<String> response =
                new ApiResponse<>(HttpStatus.BAD_REQUEST.value(), ex.getMessage(), "Illegal Argument");
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    // 4. Xử lý ngoại lệ MethodArgumentNotValidException (Validation lỗi)
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ApiResponse<String>> handleValidationException(MethodArgumentNotValidException ex) {
        log.warn("Validation Exception: ", ex);
        String errorMessage = ex.getBindingResult().getFieldErrors().stream()
                .map(fieldError -> fieldError.getField() + ": " + fieldError.getDefaultMessage())
                .findFirst()
                .orElse("Validation error");

        ApiResponse<String> response =
                new ApiResponse<>(HttpStatus.BAD_REQUEST.value(), errorMessage, "Validation Failed");
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    // 5. Xử lý ngoại lệ HttpRequestMethodNotSupportedException
    @ExceptionHandler(HttpRequestMethodNotSupportedException.class)
    public ResponseEntity<ApiResponse<String>> handleMethodNotSupportedException(
            HttpRequestMethodNotSupportedException ex) {
        log.warn("HttpRequestMethodNotSupportedException: ", ex);
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.METHOD_NOT_ALLOWED.value(),
                "HTTP method not supported: " + ex.getMethod(),
                "Method Not Allowed");
        return ResponseEntity.status(HttpStatus.METHOD_NOT_ALLOWED).body(response);
    }

    // 6. Xử lý ngoại lệ MissingServletRequestParameterException
    @ExceptionHandler(MissingServletRequestParameterException.class)
    public ResponseEntity<ApiResponse<String>> handleMissingParamsException(
            MissingServletRequestParameterException ex) {
        log.warn("MissingServletRequestParameterException: ", ex);
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.BAD_REQUEST.value(), "Missing parameter: " + ex.getParameterName(), "Bad Request");
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    // 7. Xử lý ngoại lệ DataAccessException (liên quan đến database)
    @ExceptionHandler(DataAccessException.class)
    public ResponseEntity<ApiResponse<String>> handleDataAccessException(DataAccessException ex) {
        log.warn("DataAccessException: ", ex);
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.INTERNAL_SERVER_ERROR.value(),
                "Database error: " + ex.getMessage(),
                "Internal Server Error");
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }

    // 8. Xử lý ngoại lệ NoHandlerFoundException (404 Not Found)
    @ExceptionHandler(NoHandlerFoundException.class)
    public ResponseEntity<ApiResponse<String>> handleNoHandlerFoundException(NoHandlerFoundException ex) {
        log.warn("NoHandlerFoundException: ", ex);
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.NOT_FOUND.value(), "Endpoint not found: " + ex.getRequestURL(), "Not Found");
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(response);
    }

    // 9. Xử lý ngoại lệ AccessDeniedException (403 Forbidden)
    @ExceptionHandler(AccessDeniedException.class)
    public ResponseEntity<ApiResponse<String>> handleAccessDeniedException(AccessDeniedException ex) {
        log.warn("AccessDeniedException: ", ex);
        ApiResponse<String> response =
                new ApiResponse<>(HttpStatus.FORBIDDEN.value(), "Access denied: " + ex.getMessage(), "Forbidden");
        return ResponseEntity.status(HttpStatus.FORBIDDEN).body(response);
    }

    // 10. Xử lý ngoại lệ CustomExceptionHandler (ngoại lệ tuỳ chỉnh)
    @ExceptionHandler(CustomExceptionHandler.class)
    public ResponseEntity<ApiResponse<?>> handleCustomException(CustomExceptionHandler customExceptionHandler) {
        log.warn("CustomExceptionHandler: ", customExceptionHandler);
        ErrorCode errorCode = customExceptionHandler.getErrorCode();
        return ResponseEntity.status(errorCode.getStatusCode())
                .body(new ApiResponse<>(
                        errorCode.getCode(), "An error occurred while processing the request", errorCode.getStatus()));
    }
}
