package dev.psyconnect.identity_service.interfaces;

import org.springframework.data.domain.Page;

import dev.psyconnect.identity_service.dto.request.*;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.model.Account;

// Initialize method for circular bean  error
public interface IUserAccountService {
    public AuthenticationFilterRequest loadUserByUsername(String username) throws CustomExceptionHandler;

    Page<Account> getAllAccount(int page);
}
