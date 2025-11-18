package dev.psyconnect.profile_service.repository;

import java.util.List;
import java.util.Map;

import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.data.neo4j.repository.query.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import dev.psyconnect.profile_service.enums.FriendShipStatus;
import dev.psyconnect.profile_service.model.FriendRelationship;

@Repository
public interface FriendRepository extends Neo4jRepository<FriendRelationship, String> {

    @Query(
            """
			MATCH (a:user_profile {profileId : $id})
			-[f:HAS_FRIEND]->(b:user_profile {profileId: $toId})
			SET f.status= $status , f.createdAt = datetime()
			""")
    public void accept(
            @Param("id") String id, @Param("toId") String targetID, @Param("status") FriendShipStatus status);

    @Query(
            """
			MATCH  (a:user_profile {profileId: $id}),
			(b:user_profile {profileId: $toId})
			CREATE (a)-[r:HAS_FRIEND {
			status: $status,
			createdAt: datetime()
			}]->(b)
			""")
    public void create(
            @Param("id") String id, @Param("toId") String targetID, @Param("status") FriendShipStatus status);

    @Query(
            """
			MATCH(a:user_profile {profileId : $id})
			-[f:HAS_FRIEND {status : $status}]->
			(b:user_profile{profileId: $toId})
			RETURN COUNT(f) > 0 AS existsRelationShip
			""")
    public boolean exists(
            @Param("id") String id, @Param("toId") String targetID, @Param("status") FriendShipStatus status);

    @Query(
            """
			MATCH(a:user_profile {profileId : $id})
			-[r:HAS_FRIEND]-(b:user_profile {profileId : $toId})
			WHERE r.status = $status DELETE r
			""")
    public void delete(
            @Param("id") String id, @Param("toId") String targetID, @Param("status") FriendShipStatus status);

    @Query(
            """
			MATCH (a:user_profile {profileId: $profileId})-[r:HAS_FRIEND {status:'ACCEPTED'}]->(friend:user_profile)
			RETURN DISTINCT {
			profileId: friend.profileId,
			firstName: friend.firstName,
			lastName: friend.lastName,
			avatarUri: friend.avatarUri
			} AS friend
			""")
    List<Map<String, Object>> findAcceptedFriendsLightweight(String profileId);

    @Query(
            """
			MATCH (a:user_profile {profileId: $id})
			<-[:HAS_FRIEND {status:'PENDING'}]-
			(requester:user_profile)
			RETURN DISTINCT {
				profileId: requester.profileId,
				firstName: requester.firstName,
				lastName: requester.lastName,
				avatarUri: requester.avatarUri
			} AS profile
		""")
    List<Map<String, Object>> findReceivedFriendRequests(@Param("id") String profileId);

    @Query(
            """
			MATCH (a:user_profile {profileId: $id})
			-[:HAS_FRIEND {status:'PENDING'}]->
			(target:user_profile)
			RETURN DISTINCT {
				profileId: target.profileId,
				firstName: target.firstName,
				lastName: target.lastName,
				avatarUri: target.avatarUri
			} AS profile
		""")
    List<Map<String, Object>> findSentFriendRequests(@Param("id") String profileId);

    @Query(
            """
			MATCH (a:user_profile {profileId: $id})-[:HAS_FRIEND {status:'ACCEPTED'}]->(f:user_profile)
			MATCH (b:user_profile {profileId: $toId})-[:HAS_FRIEND {status:'ACCEPTED'}]->(f)
			RETURN DISTINCT {
				profileId: f.profileId,
				firstName: f.firstName,
				lastName: f.lastName,
				avatarUri: f.avatarUri
			} AS profile
		""")
    List<Map<String, Object>> findMutualFriends(@Param("id") String id, @Param("toId") String targetId);

    @Query(
            """
			MATCH (me:user_profile {profileId: $id})
			OPTIONAL MATCH (me)-[:HAS_FRIEND {status:'ACCEPTED'}]->(f1:user_profile)
			OPTIONAL MATCH (f1)-[:HAS_FRIEND {status:'ACCEPTED'}]->(suggested:user_profile)
			WHERE NOT (me)-[:HAS_FRIEND]->(suggested)
			AND me <> suggested
			WITH me, collect(DISTINCT suggested) AS suggestedFriends

			OPTIONAL MATCH (other:user_profile)
			WHERE NOT (me)-[:HAS_FRIEND]->(other)
			AND me <> other
			WITH me, suggestedFriends, collect(DISTINCT other)[..10] AS fallback

			WITH suggestedFriends + fallback AS allSuggestions
			UNWIND allSuggestions AS s
			WITH DISTINCT s
			RETURN {
				profileId: s.profileId,
				firstName: s.firstName,
				lastName: s.lastName,
				avatarUri: s.avatarUri
			} AS profile
		""")
    List<Map<String, Object>> suggestFriends(@Param("id") String profileId);
}
