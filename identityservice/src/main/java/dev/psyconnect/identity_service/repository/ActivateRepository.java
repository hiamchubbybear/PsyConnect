package dev.psyconnect.identity_service.repository;

import java.util.Optional;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import dev.psyconnect.identity_service.model.ActivationModel;

@Repository
public interface ActivateRepository extends JpaRepository<ActivationModel, String> {
    Optional<ActivationModel> findByToken(String token);
}
