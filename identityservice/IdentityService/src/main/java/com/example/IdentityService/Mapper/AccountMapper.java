package com.example.IdentityService.Mapper;

import com.example.IdentityService.DTO.Request.AccountRequest;
import com.example.IdentityService.DTO.Respone.AccountRespone;
import com.example.IdentityService.Entity.UserAccount;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;

@Mapper(componentModel = "spring")

public interface AccountMapper {
    AccountRespone toAccountRespone(UserAccount userAccount);

    @Mapping(target = "password", ignore = true)
    UserAccount toUserAccount(AccountRequest request);
}