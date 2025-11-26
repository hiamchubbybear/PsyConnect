package dev.psyconnect.profile_service.configuration;

import java.util.HashMap;
import java.util.Map;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

import dev.psyconnect.profile_service.service.KafkaLoggerService;
import lombok.RequiredArgsConstructor;

@Component
@RequiredArgsConstructor
public class LoggingInterceptor implements HandlerInterceptor {

    private final KafkaLoggerService kafkaLogger;
    private static final String START_TIME_ATTRIBUTE = "startTime";

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) {
        request.setAttribute(START_TIME_ATTRIBUTE, System.currentTimeMillis());
        return true;
    }

    @Override
    public void afterCompletion(
            HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex) {
        Long startTime = (Long) request.getAttribute(START_TIME_ATTRIBUTE);
        if (startTime != null) {
            long duration = System.currentTimeMillis() - startTime;

            Map<String, Object> fields = new HashMap<>();
            fields.put("ip", request.getRemoteAddr());
            fields.put("userAgent", request.getHeader("User-Agent"));

            // Get userId from request attribute (set by auth filter)
            String userId = (String) request.getAttribute("userId");
            if (userId != null) {
                fields.put("userId", userId);
            }

            kafkaLogger.logHttpRequest(
                    request.getMethod(), request.getRequestURI(), response.getStatus(), duration, fields);
        }
    }
}
