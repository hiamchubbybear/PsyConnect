package com.example.IdentityService.Mapper;

import com.example.IdentityService.DTO.Request.AccountRequest;
import com.example.IdentityService.DTO.Respone.AccountRespone;
import com.example.IdentityService.Entity.UserAccount;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.springframework.data.neo4j.core.schema.Id;

@Mapper(componentModel = "spring")

public interface AccountMapper {
//    @Mapping(target = "role" , ignore = true)
    AccountRespone toAccountRespone(UserAccount userAccount);

//    @Mapping(target = "password", ignore = true)
//    UserAccount toUserAccount(AccountRequest request);
}