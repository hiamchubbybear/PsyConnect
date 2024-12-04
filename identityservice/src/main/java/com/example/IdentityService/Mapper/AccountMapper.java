package com.example.IdentityService.Mapper;

import com.example.IdentityService.dto.respone.AccountCreationRespone;
import com.example.IdentityService.model.UserAccount;
import org.mapstruct.Mapper;

@Mapper(componentModel = "spring")

public interface AccountMapper {
//    @Mapping(target = "role" , ignore = true)
    AccountCreationRespone toAccountRespone(UserAccount userAccount);

//    @Mapping(target = "password", ignore = true)
//    UserAccount toUserAccount(AccountCreationRequest request);
}