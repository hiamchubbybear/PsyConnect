package com.example.IdentityService.Bean;

import com.example.IdentityService.Entity.RoleEntity;
import com.example.IdentityService.Enum.RoleEnum;
import com.example.IdentityService.Repository.RoleRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.ApplicationRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.http.HttpStatus;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.web.client.HttpClientErrorException;

@Configuration
public class BeanConfiguration {
    private final RoleRepository roleRepository ;
    @Autowired
    public BeanConfiguration(RoleRepository roleRepository) {
        this.roleRepository = roleRepository;
    }

    @Bean
    ApplicationRunner runner () {
        return args -> {
                for( RoleEnum roleEnum : RoleEnum.values() ) {
                    if(roleRepository.existsByRoleId(roleEnum.getId())) {
                        roleRepository.save(new RoleEntity(roleEnum.getId(), roleEnum.getPermission(), roleEnum.getDescription()));
                    } else throw new IllegalArgumentException("NOT ACCEPT" );
                }
        };
    }

}
