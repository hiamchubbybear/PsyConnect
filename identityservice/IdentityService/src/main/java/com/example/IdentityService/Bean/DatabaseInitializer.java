package com.example.IdentityService.Bean;

import jakarta.annotation.PostConstruct;
import org.springframework.data.neo4j.core.Neo4jClient;
import org.springframework.stereotype.Component;

@Component
public class DatabaseInitializer {

    private final Neo4jClient neo4jClient;

    public DatabaseInitializer(Neo4jClient neo4jClient) {
        this.neo4jClient = neo4jClient;
    }

    /**
 * This method is invoked after the bean initialization to reset the Neo4j database state.
 *
 * Specifically:
 * - Deletes all nodes and their relationships in the database (MATCH (n) DETACH DELETE n).
 * - Checks and removes the unique constraint on the `roleId` property if it exists
 *   (DROP CONSTRAINT unique_role_id IF EXISTS).
 * - Creates a new unique constraint to ensure the `roleId` property is unique for nodes
 *   with the `Role` label (CREATE CONSTRAINT unique_role_id FOR (r:Role) REQUIRE r.roleId IS UNIQUE).
 *
 * Note: This operation will erase all existing data in the database.
 * It should only be used in development environments or when data cleanup is required.
 */
    @PostConstruct
    public void initializeDatabase() {
        neo4jClient.query("MATCH (n) DETACH DELETE n").run();
        neo4jClient.query("DROP CONSTRAINT unique_role_id IF EXISTS").run();
        neo4jClient.query("CREATE CONSTRAINT unique_role_id FOR (r:Role) REQUIRE r.roleId IS UNIQUE").run();
    }

}
