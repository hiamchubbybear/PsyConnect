package com.example.IdentityService.configuration;

import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;

import java.util.List;

public class LoginTypeValidator implements ConstraintValidator<ValidateLoginType, String> {

    private final List<String> allowedLoginTypes = List.of("GOOGLE", "FACEBOOK", "NORMAL");

    @Override
    public boolean isValid(String value, ConstraintValidatorContext context) {
        return value != null && allowedLoginTypes.contains(value.toUpperCase());
    }
}
