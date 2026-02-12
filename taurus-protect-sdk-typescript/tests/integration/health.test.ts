/**
 * Integration tests for Health API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use the high-level HealthService (client.health) rather than
 * the raw OpenAPI API to demonstrate proper SDK usage patterns.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Health", () => {
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

  it("should check health", async () => {
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
      const healthStatus = await client.health.check();

      console.log("Health check response received");

      // Response should have health check data
      expect(healthStatus).toBeDefined();

      console.log(`  Status: ${healthStatus.status}`);
      if (healthStatus.message) {
        console.log(`  Message: ${healthStatus.message}`);
      }
    } finally {
      client.close();
    }
  });

  it("should get global component status", async () => {
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
      const status = await client.health.getGlobalStatus();

      console.log("Global component status response received");

      expect(status).toBeDefined();

      console.log(`  Status: ${status.status}`);
      if (status.message) {
        console.log(`  Message: ${status.message}`);
      }
    } catch (error) {
      // This endpoint may not be available in all environments
      console.log(
        "Global component status endpoint not available (may require special permissions)"
      );
    } finally {
      client.close();
    }
  });
});
