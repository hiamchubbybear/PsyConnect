package com.example.IdentityService.Service;

import com.example.IdentityService.DTO.Request.AccountRequest;
import com.example.IdentityService.DTO.Respone.AccountRespone;
import com.example.IdentityService.Entity.RoleEntity;
import com.example.IdentityService.Entity.UserAccount;
import com.example.IdentityService.Mapper.AccountMapper;
import com.example.IdentityService.Repository.RoleRepository;
import com.example.IdentityService.Repository.UserAccountRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.security.crypto.bcrypt.BCrypt;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.stereotype.Service;

import java.util.Collection;
import java.util.HashSet;
import java.util.Set;
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

    public UserAccount createAccount(AccountRequest request) {
        UserAccount account = new UserAccount();
        account.setUsername(request.getUsername());
        account.setEmail(request.getEmail());
        account.setPassword(
                new BCryptPasswordEncoder().encode(request.getPassword()));
        account.setId(UUID.randomUUID());
        account.setRole(roleRepository.findById("002").get());
        accountRepository.save(account);
        return account;
    }

}
