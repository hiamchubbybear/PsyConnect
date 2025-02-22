package dev.psyconnect.profile_service.repository;

import java.util.List;
import java.util.Optional;

import dev.psyconnect.profile_service.dto.response.ProfileWithRelationShipResponse;
import dev.psyconnect.profile_service.model.Profile;
import org.neo4j.graphdb.Path;
import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.data.neo4j.repository.query.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

@Repository
public interface ProfileRepository extends Neo4jRepository<Profile, String> {
    // Find all profile
    @Query(
            "MATCH (userProfile:`user_profile`) RETURN userProfile, elementId(userProfile) AS __elementId__ SKIP $skip LIMIT $limit")
    List<Profile> findAllProfilesPaged(@Param("skip") int skip, @Param("limit") int limit);

    @Query("""
    MATCH (u:user_profile {profileId: $profileId})
    OPTIONAL MATCH (u)-[:HAS_MOOD]->(m:mood)
    OPTIONAL MATCH (u)-[:HAS_SETTING]->(s:user_setting)
    RETURN u as profile, m AS mood, s AS setting
""")
//    @Query("""
//    MATCH (u:user_profile {profileId: $profileId})
//    OPTIONAL MATCH (u)-[:HAS_MOOD]->(m:mood)
//    OPTIONAL MATCH (u)-[:HAS_SETTING]->(s:user_setting)
//    OPTIONAL MATCH (u)-[:FRIEND]->(f:user_profile)
//    RETURN u, collect(m) AS moods, collect(s) AS settings, collect(f) AS friends
//""")
    Optional<ProfileWithRelationShipResponse> getProfileWithAllRelations(@Param("profileId") String profileId);


}
