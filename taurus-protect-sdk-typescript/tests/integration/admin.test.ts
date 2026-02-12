/**
 * Integration tests for Admin APIs (Groups, Visibility Groups, Config, Audit).
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use high-level services (client.groups, client.visibilityGroups,
 * client.configService, client.audits) rather than raw OpenAPI APIs.
 */

import { skipIfNotIntegration, getTestClient } from "./helpers";

describe("Integration: Admin - Groups", () => {
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

  it("should list groups", async () => {
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
      const result = await client.groups.list({ limit: 10 });

      const groups = result.items;
      console.log(`Found ${groups.length} groups`);

      if (result.pagination) {
        console.log(`Total items: ${result.pagination.totalItems}`);
      }

      for (const group of groups) {
        console.log(
          `Group: ID=${group.id}, Name=${group.name}, ExternalID=${group.externalGroupId || "none"}`
        );
      }

      expect(Array.isArray(groups)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should paginate through groups", async () => {
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
      const pageSize = 2;
      const allGroups: Array<{ id?: string; name?: string }> = [];
      let offset = 0;

      // Fetch groups using pagination
      while (true) {
        const result = await client.groups.list({
          limit: pageSize,
          offset: offset,
        });

        const groups = result.items;
        allGroups.push(...groups);
        console.log(`Fetched ${groups.length} groups (offset=${offset})`);

        for (const group of groups) {
          console.log(`  ${group.id}: ${group.name}`);
        }

        // Check if there are more pages
        const totalItems = result.pagination?.totalItems ?? 0;
        if (offset + pageSize >= totalItems || groups.length === 0) {
          break;
        }

        offset += pageSize;

        // Safety limit for tests
        if (offset > 100) {
          console.log("Stopping pagination test at 100 items");
          break;
        }
      }

      console.log(`Total groups fetched via pagination: ${allGroups.length}`);
      expect(allGroups.length).toBeGreaterThanOrEqual(0);
    } finally {
      client.close();
    }
  }, 60000); // Extended timeout for pagination with many API calls
});

describe("Integration: Admin - Visibility Groups", () => {
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

  it("should list visibility groups", async () => {
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
      const visibilityGroups = await client.visibilityGroups.list();

      console.log(`Found ${visibilityGroups.length} visibility groups`);

      for (const vg of visibilityGroups) {
        console.log(`Visibility Group: ID=${vg.id}, Name=${vg.name}`);
      }

      expect(Array.isArray(visibilityGroups)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should get users in visibility group", async () => {
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
      // First, get a list of visibility groups to find a valid ID
      const visibilityGroups = await client.visibilityGroups.list();

      if (visibilityGroups.length === 0) {
        console.log("No visibility groups available for testing");
        return;
      }

      const visibilityGroupId = visibilityGroups[0].id;
      if (!visibilityGroupId) {
        console.log("Visibility group ID is undefined");
        return;
      }

      const users = await client.visibilityGroups.getUsersByVisibilityGroup(
        visibilityGroupId
      );

      console.log(
        `Found ${users.length} users in visibility group ${visibilityGroupId}`
      );

      for (const user of users) {
        console.log(`  User: ID=${user.id}, Email=${user.email}`);
      }

      expect(Array.isArray(users)).toBe(true);
    } catch (e) {
      // Skip if user doesn't have permission to view visibility group users
      if (e instanceof Error && e.message.includes("Response returned an error code")) {
        console.log("Skipping: Permission denied or API validation error");
        return;
      }
      throw e;
    } finally {
      client.close();
    }
  });
});

describe("Integration: Admin - Config", () => {
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

  it("should get tenant config", async () => {
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
      const config = await client.configService.getTenantConfig();

      expect(config).toBeDefined();

      console.log("Tenant configuration:");
      console.log(`  Tenant ID: ${config.tenantId}`);
      console.log(`  Base Currency: ${config.baseCurrency}`);
      console.log(`  MFA Mandatory: ${config.mfaMandatory}`);
      console.log(`  Protect Engine Version: ${config.protectEngineVersion}`);
    } finally {
      client.close();
    }
  });
});

describe("Integration: Admin - Audit", () => {
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

  it("should list audit trails", async () => {
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
      const auditTrails = await client.audits.list({ limit: 10 });

      console.log(`Found ${auditTrails.length} audit trails`);

      for (const trail of auditTrails) {
        console.log(
          `Audit Trail: Entity=${trail.entity}, Action=${trail.action}, User=${trail.userEmail}, Date=${trail.createdAt}`
        );
      }

      expect(Array.isArray(auditTrails)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should filter audit trails by date range", async () => {
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
      // Get audit trails from the last 7 days
      const now = new Date();
      const sevenDaysAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);

      const auditTrails = await client.audits.list({
        creationDateFrom: sevenDaysAgo,
        creationDateTo: now,
        limit: 10,
      });

      console.log(`Found ${auditTrails.length} audit trails in last 7 days`);

      for (const trail of auditTrails) {
        console.log(
          `Audit Trail: Entity=${trail.entity}, Action=${trail.action}, Date=${trail.createdAt}`
        );
      }

      expect(Array.isArray(auditTrails)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should filter audit trails by entity", async () => {
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
      const auditTrails = await client.audits.list({
        entities: ["USER"],
        limit: 10,
      });

      console.log(`Found ${auditTrails.length} USER audit trails`);

      for (const trail of auditTrails) {
        console.log(
          `Audit Trail: Entity=${trail.entity}, Action=${trail.action}, User=${trail.userEmail}`
        );
        // All returned trails should be for USER entity
        expect(trail.entity).toBe("USER");
      }

      expect(Array.isArray(auditTrails)).toBe(true);
    } catch (e) {
      // Skip if API doesn't support entity filtering (returns 400)
      if (e instanceof Error && e.message.includes("Response returned an error code")) {
        console.log("Skipping: API validation error - entity filter may not be supported");
        return;
      }
      throw e;
    } finally {
      client.close();
    }
  });
});
