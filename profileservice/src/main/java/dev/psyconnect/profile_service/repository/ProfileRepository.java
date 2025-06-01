package dev.psyconnect.profile_service.repository;

import java.util.List;
import java.util.Optional;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.data.neo4j.repository.query.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import dev.psyconnect.profile_service.dto.response.ProfileWithMood;
import dev.psyconnect.profile_service.dto.response.ProfileWithRelationShipResponse;
import dev.psyconnect.profile_service.model.Profile;

@Repository
public interface ProfileRepository extends Neo4jRepository<Profile, String> {
    // Find all profile
    @Query(
            """
                    MATCH (userProfile:`user_profile`)
                    RETURN userProfile, elementId(userProfile)
                    AS __elementId__ SKIP $skip LIMIT $limit""")
    List<Profile> findAllProfilesPaged(@Param("skip") int skip, @Param("limit") int limit);

    @Query("""
                MATCH (u:user_profile {profileId: $profileId})-[:HAS_FRIEND]-(other:user_profile)
                OPTIONAL MATCH (other)-[:HAS_MOOD]->(m:mood)
                RETURN DISTINCT other AS profile, m AS mood
            """)
    List<ProfileWithRelationShipResponse> getProfileWithAllRelations(@Param("profileId") String profileId);


    @Query(
            """
                    	MATCH (u:user_profile {profileId: $profileId})
                    		-[:HAS_FRIEND]->(p:user_profile)
                    		-[:HAS_MOOD]-> (m:mood)
                    	RETURN p as profile, m as mood
                    """)
    List<ProfileWithMood> findFriendsWithMoodsByProfileId(String profileId);
}
