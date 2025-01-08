package com.example.profileservice.repository;

import com.example.profileservice.dto.request.UserProfileUpdateRequest;
import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.stereotype.Repository;

import com.example.profileservice.model.UserProfile;

import java.util.Optional;

@Repository
public interface ProfileRepository extends Neo4jRepository<UserProfile, String> {
    Optional<UserProfile> findByEmail(String email);
}
