package dev.psyconnect.identity_service.repository;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import dev.psyconnect.identity_service.model.Token;

@Repository
public interface ActivateRepository extends JpaRepository<Token, String> {
    Optional<Token> findByToken(String token);
}
