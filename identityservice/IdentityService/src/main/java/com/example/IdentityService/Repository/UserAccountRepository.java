package com.example.IdentityService.Repository;

import com.example.IdentityService.Entity.UserAccount;
import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface UserAccountRepository extends Neo4jRepository  <UserAccount, UUID> {

}
