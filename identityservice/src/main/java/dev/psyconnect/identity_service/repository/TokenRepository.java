package dev.psyconnect.identity_service.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import dev.psyconnect.identity_service.model.Token;

@Repository
public interface TokenRepository extends JpaRepository<Token, String> {}
