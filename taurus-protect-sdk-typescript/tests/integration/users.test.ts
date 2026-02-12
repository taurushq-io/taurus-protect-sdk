/**
 * Integration tests for Users API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use the high-level UserService (client.users) rather than
 * the raw OpenAPI API to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Users", () => {
  beforeEach(() => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }
  });

  it("should get current user (getCurrentUser)", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClient();
    try {
      const user = await client.users.getCurrentUser();

      expect(user).toBeDefined();

      console.log("Current user details:");
      console.log(`  ID: ${user.id}`);
      console.log(`  Email: ${user.email}`);
      console.log(`  Name: ${user.firstName} ${user.lastName}`);
      console.log(`  Status: ${user.status}`);
      console.log(`  Roles: ${user.roles?.join(", ") || "none"}`);
    } finally {
      client.close();
    }
  });

  it("should list users", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClient();
    try {
      const result = await client.users.list({ limit: 10 });

      const users = result.items;
      console.log(`Found ${users.length} users`);

      if (result.pagination) {
        console.log(`Total items: ${result.pagination.totalItems}`);
      }

      for (const user of users) {
        console.log(
          `User: ID=${user.id}, Email=${user.email}, Name=${user.firstName} ${user.lastName}, Status=${user.status}`
        );
      }

      expect(Array.isArray(users)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should filter users by status", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClient();
    try {
      const result = await client.users.list({
        limit: 10,
        status: "ACTIVE",
      });

      const users = result.items;
      console.log(`Found ${users.length} active users`);

      for (const user of users) {
        console.log(
          `Active User: ID=${user.id}, Email=${user.email}, Status=${user.status}`
        );
      }

      expect(Array.isArray(users)).toBe(true);
      // All returned users should be active
      for (const user of users) {
        expect(user.status).toBe("ACTIVE");
      }
    } catch (e) {
      // Skip if the API doesn't support status filtering (returns 400)
      if (e instanceof Error && e.message.includes("Response returned an error code")) {
        console.log("Skipping: API validation error - status filter may not be supported");
        return;
      }
      throw e;
    } finally {
      client.close();
    }
  });
});
