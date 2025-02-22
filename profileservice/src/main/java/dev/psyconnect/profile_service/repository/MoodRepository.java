package dev.psyconnect.profile_service.repository;

import java.util.Optional;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.data.neo4j.repository.query.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import dev.psyconnect.profile_service.model.Mood;

@Repository
public interface MoodRepository extends Neo4jRepository<Mood, String> {
    @Query(
            """
                    MATCH (u:user_profile {profileId: $profileId})
                    CREATE (m:mood {
                    	moodId: $profileId,
                    	mood: $mood,
                    	description: $description,
                    	createdAt: $currentTimestamp,
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
            @Param("currenTimestamp") long currentTimestamp,
            @Param("expiresTimestamp") long expiresTimestamp);

    @Query(
            """
                    MATCH (u:user_profile {profileId: $profileId})-[:HAS_MOOD]->(m)
                    SET m.mood = $mood,
                    	m.description = $description,
                    	m.visibility = $visibility
                    RETURN m
                    """)
    Optional<Mood> updateMood(
            @Param("profileId") String profileId,
            @Param("mood") String mood,
            @Param("description") String description,
            @Param("visibility") String visibility
    );

    @Query(
            """
                    MATCH( u:user_profile {profileId: $profileId})-[:HAS_MOOD]->(m:mood)
                    RETURN m
                    """
    )
    Optional<Mood> getMood(@Param("profileId") String profileId);

}
