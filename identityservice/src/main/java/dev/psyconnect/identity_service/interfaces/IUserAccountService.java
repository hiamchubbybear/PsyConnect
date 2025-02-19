package dev.psyconnect.identity_service.interfaces;

import dev.psyconnect.identity_service.dto.request.*;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;

// Initialize method for circular bean  error
public interface IUserAccountService {
    public AuthenticationFilterRequest loadUserByUsername(String username) throws CustomExceptionHandler;
}
