package dev.psyconnect.identity_service.repository;

import java.util.Optional;
import java.util.UUID;

import jakarta.transaction.Transactional;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import dev.psyconnect.identity_service.model.Account;

@Repository
public interface UserAccountRepository extends JpaRepository<Account, UUID> {
    @Query("update Account m set m.isActivated = true where m.email = :email")
    @Modifying
    @Transactional
    void activateUser(@Param("email") String email);

    boolean existsByUsername(String name);

    boolean existsByEmail(String email);

    Optional<Account> findByUsername(String username);

    Optional<Account> findByEmail(String email);

    @Query("SELECT u.isActivated AS boolean_value\n" + "FROM Account  u where u.username=?1")
    Boolean isActive(String username);

    @Query("SELECT u.isActivated AS boolean_value\n" + "FROM Account  u where u.email=?1")
    Boolean isActiveByEmail(String email);
}
