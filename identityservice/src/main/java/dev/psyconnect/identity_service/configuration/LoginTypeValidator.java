package dev.psyconnect.identity_service.configuration;

import java.util.List;

import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;

public class LoginTypeValidator implements ConstraintValidator<ValidateLoginType, String> {

    private final List<String> allowedLoginTypes = List.of("GOOGLE", "FACEBOOK", "NORMAL" ,"MOBILE");

    @Override
    public boolean isValid(String value, ConstraintValidatorContext context) {
        return value != null && allowedLoginTypes.contains(value.toUpperCase());
    }
}
