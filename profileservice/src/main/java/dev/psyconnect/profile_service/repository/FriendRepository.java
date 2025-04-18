package dev.psyconnect.profile_service.repository;

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
}
