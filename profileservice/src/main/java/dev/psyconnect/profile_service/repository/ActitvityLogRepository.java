package dev.psyconnect.profile_service.repository;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.stereotype.Repository;

import dev.psyconnect.profile_service.model.ActivityLog;

@Repository
public interface ActitvityLogRepository extends Neo4jRepository<ActivityLog, String> {}
