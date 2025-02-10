package dev.psyconnect.profile_service;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication(scanBasePackages = "dev.psyconnect.profile_service")
public class ProfileserviceApplication {

    public static void main(String[] args) {
        SpringApplication.run(ProfileserviceApplication.class, args);
    }
}
