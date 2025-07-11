package dev.psyconnect.identity_service.configuration;

import java.util.List;

import jakarta.validation.ConstraintValidator;
import jakarta.validation.ConstraintValidatorContext;

public class ProviderTypeValidator implements ConstraintValidator<ValidateProviderType, String> {

    private final List<String> allowedLoginTypes = List.of("GOOGLE", "FACEBOOK", "NORMAL", "APPLEID");

    @Override
    public boolean isValid(String value, ConstraintValidatorContext context) {
        return value != null && allowedLoginTypes.contains(value.toUpperCase());
    }
}
