package com.example.IdentityService.configuration;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

import jakarta.validation.Constraint;

import com.nimbusds.jose.Payload;

@Constraint(validatedBy = LoginTypeValidator.class)
@Target({ElementType.PARAMETER, ElementType.FIELD})
@Retention(RetentionPolicy.RUNTIME)
public @interface ValidateLoginType {
    String message() default "Invalid Login Type : Only Contain  { GOOGLE , FACEBOOK , NORMAL }";

    Class<?>[] groups() default {};

    Class<? extends Payload>[] payload() default {};
}
