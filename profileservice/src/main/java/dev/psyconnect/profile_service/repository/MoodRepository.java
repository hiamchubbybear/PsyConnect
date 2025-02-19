package dev.psyconnect.profile_service.repository;

import java.util.Optional;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.data.neo4j.repository.query.Query;
import org.springframework.stereotype.Repository;

import dev.psyconnect.profile_service.model.Mood;
import io.lettuce.core.dynamic.annotation.Param;

@Repository
public interface MoodRepository extends Neo4jRepository<Mood, String> {
    @Query(
            """
					MATCH (u:user_profile {profileId: $profileId})
					CREATE (m:mood {
						moodId: randomUUID(),
						mood: $mood,
						description: $description,
						createdAt: $currenTimestamp,
						expiresAt: $expiresTimestamp,
						visibility: $visibility

					})
					CREATE (u)-[:HAS_MOOD]->(m)
					RETURN m
					""")
    Optional<Mood> createMood(
            @Param("profileId") String profileId,
            @Param("mood") String mood,
            @Param("description") String description,
            @Param("visibility") String visibility,
            @Param("currenTimestamp") long currenTimestamp,
            @Param("expiresTimestamp") long expiresTimestamp);
}
