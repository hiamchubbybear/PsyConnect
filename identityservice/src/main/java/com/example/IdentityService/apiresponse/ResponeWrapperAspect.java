package com.example.IdentityService.apiresponse;

import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Component;

@Aspect
@Component
public class ResponeWrapperAspect {

    @Around("@annotation(com.example.IdentityService.apiresponse.CustomResponseWrapper)")
    public Object customResponseWrapper(ProceedingJoinPoint joinPoint) throws Throwable {
        Object result = joinPoint.proceed();
        if (result instanceof ResponseEntity) {
            ResponseEntity<?> response = (ResponseEntity<?>) result;
            return ResponseEntity.status(response.getStatusCode())
                    .body(new ApiResponse<>(response.getStatusCode().value(), "Success", response.getBody()));
        } else if (result instanceof ApiResponse) {
            return result;
        }
        return new ApiResponse<>(200, "Success", result);
    }
}
