/**
 * Integration tests for Governance Rules API.
 *
 * These tests require a live API connection. Configure via environment variables
 * or hard-coded defaults in helpers.ts.
 *
 * Tests use the high-level GovernanceRuleService (client.governanceRules) rather than
 * the raw OpenAPI API to demonstrate proper SDK usage patterns.
 */

import { ProtectClient } from "../../src/client";
import {
  skipIfNotIntegration,
  getTestClient,
  getTestClientWithVerification,
  getConfig,
  DEFAULT_MIN_VALID_SIGNATURES,
} from "./helpers";

describe("Integration: Governance Rules", () => {
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

  it("should get governance rules", async () => {
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
      const rules = await client.governanceRules.get();

      console.log("Governance rules response received");

      // The response contains the rules
      if (rules) {
        if (rules.rulesContainer) {
          console.log(
            `  Rules container present: ${rules.rulesContainer.length > 0}`
          );
        }

        if (rules.rulesSignatures) {
          console.log(`  Signatures count: ${rules.rulesSignatures.length}`);
        }

        if (rules.locked !== undefined) {
          console.log(`  Locked: ${rules.locked}`);
        }

        if (rules.creationDate) {
          console.log(`  Creation date: ${rules.creationDate}`);
        }
      }

      // Basic validation - rules should be defined
      expect(rules).toBeDefined();
    } finally {
      client.close();
    }
  });

  it("should get governance rules history", async () => {
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
      const result = await client.governanceRules.getHistory({ limit: 10 });

      console.log("Governance rules history response received");

      const history = result.items;
      console.log(`  History entries: ${history.length}`);

      if (result.nextCursor) {
        console.log(`  Has more pages: yes (cursor: ${result.nextCursor})`);
      }

      for (const entry of history) {
        console.log(`  - Created: ${entry.creationDate}, Locked: ${entry.locked}`);
      }

      expect(Array.isArray(history)).toBe(true);
    } finally {
      client.close();
    }
  });

  it("should get decoded rules container", async () => {
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
      const decoded = await client.governanceRules.getDecodedRulesContainer();

      console.log("Decoded rules container response received");

      expect(decoded).toBeDefined();

      console.log(`  Users count: ${decoded.users?.length ?? 0}`);
      console.log(`  Groups count: ${decoded.groups?.length ?? 0}`);
    } finally {
      client.close();
    }
  });

  it("should verify governance rules with SuperAdmin signatures", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClientWithVerification();
    try {
      const rules = await client.governanceRules.get();

      if (!rules) {
        console.log("No governance rules available - skipping test");
        return;
      }

      console.log("Governance rules retrieved:");
      console.log(`  Rules container length: ${rules.rulesContainer?.length ?? 0} bytes`);
      console.log(`  Signatures count: ${rules.rulesSignatures?.length ?? 0}`);

      // Log signature user IDs
      if (rules.rulesSignatures) {
        rules.rulesSignatures.forEach((sig, i) => {
          console.log(`  Signature[${i}] userID: ${sig.userId}`);
        });
      }

      // getDecodedRulesContainer performs signature verification when SuperAdmin keys are configured
      const decoded = await client.governanceRules.getDecodedRulesContainer();

      console.log("Signature verification PASSED");
      console.log("Decoded rules container:");
      console.log(`  Users count: ${decoded.users?.length ?? 0}`);
      console.log(`  Groups count: ${decoded.groups?.length ?? 0}`);
      console.log(`  Address whitelisting rules: ${decoded.addressWhitelistingRules?.length ?? 0}`);
      console.log(`  Contract address whitelisting rules: ${decoded.contractAddressWhitelistingRules?.length ?? 0}`);
      // Note: requestRules not in current DecodedRulesContainer model

      expect(decoded).toBeDefined();
    } finally {
      client.close();
    }
  });

  it("should fail verification with invalid SuperAdmin keys", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const config = getConfig();

    // Create client with invalid SuperAdmin keys (should fail when keys are parsed)
    const invalidKeys = [
      `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==
-----END PUBLIC KEY-----`,
    ];

    const client = ProtectClient.create({
      host: config.host,
      apiKey: config.apiKey,
      apiSecret: config.apiSecret,
      superAdminKeysPem: invalidKeys,
      minValidSignatures: 1,
    });

    try {
      // Accessing governanceRules service should fail because the invalid
      // public key cannot be parsed (not a valid EC point)
      expect(() => client.governanceRules).toThrow("Failed to decode public key");

      console.log("Correctly rejected invalid SuperAdmin keys at parse time");
    } finally {
      client.close();
    }
  });

  it("should have SuperAdmin keys configured on verification client", async () => {
    try {
      skipIfNotIntegration();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_INTEGRATION") {
        return;
      }
      throw e;
    }

    const client = getTestClientWithVerification();
    try {
      // Verify client has SuperAdmin keys configured
      const keyCount = client.superAdminKeysPem?.length ?? 0;
      const minSigs = client.minValidSignatures ?? 0;

      console.log(`Client configured with ${keyCount} SuperAdmin keys`);
      console.log(`MinValidSignatures: ${minSigs}`);

      expect(keyCount).toBe(3);
      expect(minSigs).toBe(DEFAULT_MIN_VALID_SIGNATURES);
    } finally {
      client.close();
    }
  });
});
