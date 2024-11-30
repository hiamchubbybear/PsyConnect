package com.example.IdentityService.Bean;

import org.springframework.boot.ApplicationRunner;
import org.springframework.context.annotation.Bean;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

public class BeanConfiguration {
    ApplicationRunner runner () {
        return args -> {

        };
    }

}
