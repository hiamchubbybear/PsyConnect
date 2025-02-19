package dev.psyconnect.profile_service.repository;

import java.util.List;
import java.util.Optional;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.data.neo4j.repository.query.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import dev.psyconnect.profile_service.model.UserProfile;

@Repository
public interface ProfileRepository extends Neo4jRepository<UserProfile, String> {
    // Find all profile
    @Query(
            "MATCH (userProfile:`user_profile`) RETURN userProfile, elementId(userProfile) AS __elementId__ SKIP $skip LIMIT $limit")
    List<UserProfile> findAllProfilesPaged(@Param("skip") int skip, @Param("limit") int limit);

    Optional<UserProfile> findByProfileId(String userId);
}
