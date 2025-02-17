package dev.psyconnect.identity_service.interfaces;

import dev.psyconnect.identity_service.dto.request.*;
import dev.psyconnect.identity_service.dto.response.ActivateAccountResponse;
import dev.psyconnect.identity_service.dto.response.DeleteAccountResponse;
import dev.psyconnect.identity_service.dto.response.UpdateAccountResponse;
import dev.psyconnect.identity_service.dto.response.UserInfoResponse;
import dev.psyconnect.identity_service.globalexceptionhandle.CustomExceptionHandler;
import dev.psyconnect.identity_service.service.UserAccountService;
import jakarta.transaction.Transactional;

import java.util.UUID;
// Initialize method for circular bean  error
public interface IUserAccountService {
    public AuthenticationFilterRequest loadUserByUsername(String username) throws CustomExceptionHandler;
}