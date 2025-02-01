package com.example.IdentityService.configuration;

import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;

import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.util.ContentCachingRequestWrapper;

public class LoggingInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler)
            throws Exception {
        ContentCachingRequestWrapper requestWrapper = new ContentCachingRequestWrapper(request);
        String uri = requestWrapper.getRequestURI();
        String method = requestWrapper.getMethod();
        String headers = requestWrapper.getHeader("User-Agent");

        // In ra thông tin về request
        System.out.println("Request Method: " + method);
        System.out.println("Request URI: " + uri);
        System.out.println("Request Headers: " + headers);

        // Nếu có body trong request, in luôn
        if (requestWrapper.getContentLength() > 0) {
            String body = new String(requestWrapper.getContentAsByteArray());
            System.out.println("Request Body: " + body);
        }

        return true;
    }
}
