package dev.psyconnect.profile_service.repository;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.stereotype.Repository;

import dev.psyconnect.profile_service.model.UserSetting;

@Repository
public interface UserSettingRepository extends Neo4jRepository<UserSetting, String> {}
