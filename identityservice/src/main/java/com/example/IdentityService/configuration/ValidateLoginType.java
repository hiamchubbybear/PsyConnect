package com.example.IdentityService.configuration;

import com.nimbusds.jose.Payload;
import jakarta.validation.Constraint;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

@Constraint(validatedBy = LoginTypeValidator.class)
@Target({ElementType.PARAMETER, ElementType.FIELD})
@Retention(RetentionPolicy.RUNTIME)
public @interface ValidateLoginType {
    String message() default "Invalid Login Type : Only Contain  { GOOGLE , FACEBOOK , NORMAL }";
    Class<?>[] groups() default {};
    Class<? extends Payload>[] payload() default {};
}
