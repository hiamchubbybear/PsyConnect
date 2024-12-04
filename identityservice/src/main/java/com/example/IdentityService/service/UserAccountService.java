package com.example.IdentityService.service;

import com.example.IdentityService.dto.request.AccountCreationRequest;
import com.example.IdentityService.model.UserAccount;
import com.example.IdentityService.Mapper.AccountMapper;
import com.example.IdentityService.repository.RoleRepository;
import com.example.IdentityService.repository.UserAccountRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.UUID;

@Service
@RequiredArgsConstructor
public class UserAccountService {
    private final RoleRepository roleRepository;
    public AccountMapper mapper;
    public UserAccountRepository accountRepository;
    public PasswordEncodingService passwordEncodingService;
    @Autowired
    public UserAccountService(AccountMapper mapper, UserAccountRepository accountRepository , PasswordEncodingService passwordEncodingService, RoleRepository roleRepository) {
        this.mapper = mapper;
        this.accountRepository = accountRepository;
        this.passwordEncodingService = passwordEncodingService;
        this.roleRepository = roleRepository;
    }

    public UserAccount createAccount(AccountCreationRequest request) {
        UserAccount account = new UserAccount();
        account.setUsername(request.getUsername());
        account.setEmail(request.getEmail());
        account.setPassword(
                new BCryptPasswordEncoder().encode(request.getPassword()));
        account.setId(UUID.randomUUID());
        account.setRole(roleRepository.findById(request.getRole()).get());
        accountRepository.save(account);
        return account;
    }

}
