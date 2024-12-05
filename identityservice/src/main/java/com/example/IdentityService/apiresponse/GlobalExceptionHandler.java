package com.example.IdentityService.apiresponse;

import org.springframework.dao.DataAccessException;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.HttpRequestMethodNotSupportedException;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.MissingServletRequestParameterException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.servlet.NoHandlerFoundException;

import java.nio.file.AccessDeniedException;

@ControllerAdvice
public class GlobalExceptionHandler {

    // 1. Xử lý ngoại lệ chung (Exception)
    @ExceptionHandler(Exception.class)
    public ResponseEntity<ApiResponse<String>> handleGeneralException(Exception ex) {
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.INTERNAL_SERVER_ERROR.value(),
                ex.getMessage(),
                "Internal Server Error"
        );
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }

    // 2. Xử lý ngoại lệ NullPointerException
    @ExceptionHandler(NullPointerException.class)
    public ResponseEntity<ApiResponse<String>> handleNullPointerException(NullPointerException ex) {
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.BAD_REQUEST.value(),
                "Null value encountered",
                "Bad Request"
        );
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    // 3. Xử lý ngoại lệ IllegalArgumentException
    @ExceptionHandler(IllegalArgumentException.class)
    public ResponseEntity<ApiResponse<String>> handleIllegalArgumentException(IllegalArgumentException ex) {
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.BAD_REQUEST.value(),
                ex.getMessage(),
                "Illegal Argument"
        );
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    // 4. Xử lý ngoại lệ MethodArgumentNotValidException (Validation lỗi)
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ApiResponse<String>> handleValidationException(MethodArgumentNotValidException ex) {
        String errorMessage = ex.getBindingResult().getFieldErrors().stream()
                .map(fieldError -> fieldError.getField() + ": " + fieldError.getDefaultMessage())
                .findFirst()
                .orElse("Validation error");

        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.BAD_REQUEST.value(),
                errorMessage,
                "Validation Failed"
        );
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    // 5. Xử lý ngoại lệ HttpRequestMethodNotSupportedException
    @ExceptionHandler(HttpRequestMethodNotSupportedException.class)
    public ResponseEntity<ApiResponse<String>> handleMethodNotSupportedException(HttpRequestMethodNotSupportedException ex) {
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.METHOD_NOT_ALLOWED.value(),
                "HTTP method not supported: " + ex.getMethod(),
                "Method Not Allowed"
        );
        return ResponseEntity.status(HttpStatus.METHOD_NOT_ALLOWED).body(response);
    }

    // 6. Xử lý ngoại lệ MissingServletRequestParameterException
    @ExceptionHandler(MissingServletRequestParameterException.class)
    public ResponseEntity<ApiResponse<String>> handleMissingParamsException(MissingServletRequestParameterException ex) {
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.BAD_REQUEST.value(),
                "Missing parameter: " + ex.getParameterName(),
                "Bad Request"
        );
        return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(response);
    }

    // 7. Xử lý ngoại lệ DataAccessException (liên quan đến database)
    @ExceptionHandler(DataAccessException.class)
    public ResponseEntity<ApiResponse<String>> handleDataAccessException(DataAccessException ex) {
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.INTERNAL_SERVER_ERROR.value(),
                "Database error: " + ex.getMessage(),
                "Internal Server Error"
        );
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(response);
    }

    // 8. Xử lý ngoại lệ NoHandlerFoundException (404 Not Found)
    @ExceptionHandler(NoHandlerFoundException.class)
    public ResponseEntity<ApiResponse<String>> handleNoHandlerFoundException(NoHandlerFoundException ex) {
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.NOT_FOUND.value(),
                "Endpoint not found: " + ex.getRequestURL(),
                "Not Found"
        );
        return ResponseEntity.status(HttpStatus.NOT_FOUND).body(response);
    }

    // 9. Xử lý ngoại lệ AccessDeniedException (403 Forbidden)
    @ExceptionHandler(AccessDeniedException.class)
    public ResponseEntity<ApiResponse<String>> handleAccessDeniedException(AccessDeniedException ex) {
        ApiResponse<String> response = new ApiResponse<>(
                HttpStatus.FORBIDDEN.value(),
                "Access denied: " + ex.getMessage(),
                "Forbidden"
        );
        return ResponseEntity.status(HttpStatus.FORBIDDEN).body(response);
    }

//    // 10. Xử lý ngoại lệ CustomException (nếu có ngoại lệ tùy chỉnh)
//    @ExceptionHandler(CustomException.class)
//    public ResponseEntity<ApiResponse<String>> handleCustomException(CustomException ex) {
//        ApiResponse<String> response = new ApiResponse<>(
//                ex.getStatusCode(),
//                ex.getMessage(),
//                "Custom Error"
//        );
//        return ResponseEntity.status(ex.getStatusCode()).body(response);
//    }
}