package dev.psyconnect.profile_service.model;

import org.springframework.data.neo4j.core.schema.Relationship;
import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface MoodRepository extends Neo4jRepository<Mood, String> {
}
