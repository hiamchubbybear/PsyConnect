package com.example.profileservice.repository;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.stereotype.Repository;

import com.example.profileservice.model.UserProfile;

@Repository
public interface ProfileRepository extends Neo4jRepository<UserProfile, String> {}
