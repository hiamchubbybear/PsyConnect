package dev.psyconnect.profile_service.model;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.stereotype.Repository;

import java.util.UUID;

@Repository
public interface ActitvityLogRepository extends Neo4jRepository<ActivityLog, UUID> {
}
