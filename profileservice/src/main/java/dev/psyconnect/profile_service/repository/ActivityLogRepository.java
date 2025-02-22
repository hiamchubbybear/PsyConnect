package dev.psyconnect.profile_service.repository;

import io.lettuce.core.dynamic.annotation.Param;
import org.springframework.data.neo4j.repository.Neo4jRepository;
import org.springframework.data.neo4j.repository.query.Query;
import org.springframework.stereotype.Repository;

import dev.psyconnect.profile_service.model.ActivityLog;

import java.util.List;

@Repository
public interface ActivityLogRepository extends Neo4jRepository<ActivityLog, String> {

    @Query("""
            MATCH (a:ActivityLog {profileId: $profileId})
            RETURN a
            ORDER BY a.timestamp DESC
            """)
    List<ActivityLog> findByProfileId(@Param("profileId") String profileId);

    @Query("""
            CREATE (a:ActivityLog {
                id: randomUUID(),
                profileId: $profileId,
                action: $action,
                timestamp: $timestamp,
                targetId: $targetId,
                targetType: $targetType,
                extraData: $extraData
            })
            RETURN a
            """)
    ActivityLog createLog(
            @Param("profileId") String profileId,
            @Param("action") String action,
            @Param("timestamp") long timestamp,
            @Param("targetId") String targetId,
            @Param("targetType") String targetType,
            @Param("extraData") String extraData
    );

    @Query("""
            MATCH (a:ActivityLog {id: $id})
            DELETE a
            """)
    void deleteById(@Param("id") String id);
}