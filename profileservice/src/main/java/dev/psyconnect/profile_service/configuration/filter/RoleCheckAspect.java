package dev.psyconnect.profile_service.configuration.filter;

import dev.psyconnect.profile_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.profile_service.globalexceptionhandle.ErrorCode;
import jakarta.servlet.http.HttpServletRequest;
import lombok.extern.slf4j.Slf4j;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.annotation.Before;
import org.springframework.stereotype.Component;
import org.springframework.web.context.request.RequestContextHolder;
import org.springframework.web.context.request.ServletRequestAttributes;

@Aspect
@Component
@Slf4j
public class RoleCheckAspect {
    @Before("@annotation(allowedRoles)")
    public void checkRoles(AllowedRoles allowedRoles) {
        HttpServletRequest request = ((ServletRequestAttributes) RequestContextHolder.currentRequestAttributes()).getRequest();
        String header = request.getHeader("X-Roles");
        String role = header.substring(0, header.indexOf(':'));
        boolean isAllowed = false;
        // If match one of all break the loop and return true
        for (String r : allowedRoles.value()) {
            if (role.equals("role." + r.toLowerCase().trim())) {
                log.info("User role {}", "role." + r.toLowerCase().trim());
                isAllowed = true;
                break;
            }
        }
        if (!isAllowed) throw new CustomExceptionHandler(ErrorCode.UNAUTHORIZED);
    }
}
