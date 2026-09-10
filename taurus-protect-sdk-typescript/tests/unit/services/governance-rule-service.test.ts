/**
 * Unit tests for GovernanceRuleService.
 *
 * These tests verify the critical security feature:
 * - ECDSA signature verification for governance rules
 */

import * as crypto from "crypto";
import { GovernanceRuleService } from "../../../src/services/governance-rule-service";
import { ConfigurationError, IntegrityError } from "../../../src/errors";
import type { GovernanceRules, RuleUserSignature } from "../../../src/models/governance-rules";
import type { GovernanceRulesApi } from "../../../src/internal/openapi/apis/GovernanceRulesApi";
import { signData } from "../../../src/crypto";

// Mock GovernanceRulesApi
function createMockGovernanceRulesApi(): jest.Mocked<GovernanceRulesApi> {
  return {
    ruleServiceGetRules: jest.fn(),
    ruleServiceGetRulesByID: jest.fn(),
    ruleServiceGetRulesProposal: jest.fn(),
    ruleServiceGetRulesHistory: jest.fn(),
    ruleServiceGetPublicKeys: jest.fn(),
  } as unknown as jest.Mocked<GovernanceRulesApi>;
}

// Generate a test ECDSA key pair
function generateTestKeyPair(): {
  privateKey: crypto.KeyObject;
  publicKey: crypto.KeyObject;
} {
  return crypto.generateKeyPairSync("ec", {
    namedCurve: "P-256",
  });
}

describe("GovernanceRuleService", () => {
  describe("verifyGovernanceRules", () => {
    it("should pass verification with valid signatures", () => {
      const mockApi = createMockGovernanceRulesApi();
      const keyPair1 = generateTestKeyPair();
      const keyPair2 = generateTestKeyPair();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [keyPair1.publicKey, keyPair2.publicKey],
        minValidSignatures: 2,
      });

      // Create a rules container
      const rulesContainerData = JSON.stringify({
        users: [],
        groups: [],
        minimumDistinctUserSignatures: 2,
      });
      const rulesContainerBase64 = Buffer.from(rulesContainerData).toString("base64");
      const rulesContainerBuffer = Buffer.from(rulesContainerBase64, "base64");

      // Sign with both keys
      const signature1 = signData(keyPair1.privateKey, rulesContainerBuffer);
      const signature2 = signData(keyPair2.privateKey, rulesContainerBuffer);

      const rules: GovernanceRules = {
        rulesContainer: rulesContainerBase64,
        rulesSignatures: [
          { userId: "user1", signature: signature1 },
          { userId: "user2", signature: signature2 },
        ],
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      // Should not throw
      expect(() => service.verifyGovernanceRules(rules)).not.toThrow();
    });

    it("should fail verification with invalid signatures", () => {
      const mockApi = createMockGovernanceRulesApi();
      const keyPair = generateTestKeyPair();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [keyPair.publicKey],
        minValidSignatures: 1,
      });

      // Create a rules container
      const rulesContainerData = JSON.stringify({ users: [] });
      const rulesContainerBase64 = Buffer.from(rulesContainerData).toString("base64");

      // Use a completely fake signature
      const fakeSignature = Buffer.from("invalid-signature-data").toString("base64");

      const rules: GovernanceRules = {
        rulesContainer: rulesContainerBase64,
        rulesSignatures: [{ userId: "user1", signature: fakeSignature }],
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      expect(() => service.verifyGovernanceRules(rules)).toThrow(IntegrityError);
      expect(() => service.verifyGovernanceRules(rules)).toThrow(
        "insufficient distinct valid SuperAdmin signers: got 0, need 1"
      );
    });

    it("should fail verification with signature from wrong key", () => {
      const mockApi = createMockGovernanceRulesApi();
      const keyPairTrusted = generateTestKeyPair();
      const keyPairUntrusted = generateTestKeyPair();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [keyPairTrusted.publicKey],
        minValidSignatures: 1,
      });

      // Create a rules container
      const rulesContainerData = JSON.stringify({ users: [] });
      const rulesContainerBase64 = Buffer.from(rulesContainerData).toString("base64");
      const rulesContainerBuffer = Buffer.from(rulesContainerBase64, "base64");

      // Sign with untrusted key (not in superAdminKeys)
      const signatureFromUntrusted = signData(
        keyPairUntrusted.privateKey,
        rulesContainerBuffer
      );

      const rules: GovernanceRules = {
        rulesContainer: rulesContainerBase64,
        rulesSignatures: [{ userId: "user1", signature: signatureFromUntrusted }],
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      expect(() => service.verifyGovernanceRules(rules)).toThrow(IntegrityError);
    });

    it("should prevent duplicate signatures from same user", () => {
      const mockApi = createMockGovernanceRulesApi();
      const keyPair = generateTestKeyPair();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [keyPair.publicKey],
        minValidSignatures: 2, // Require 2 signatures
      });

      // Create a rules container
      const rulesContainerData = JSON.stringify({ users: [] });
      const rulesContainerBase64 = Buffer.from(rulesContainerData).toString("base64");
      const rulesContainerBuffer = Buffer.from(rulesContainerBase64, "base64");

      // Sign twice with the same key and same user ID
      const signature1 = signData(keyPair.privateKey, rulesContainerBuffer);
      const signature2 = signData(keyPair.privateKey, rulesContainerBuffer);

      const rules: GovernanceRules = {
        rulesContainer: rulesContainerBase64,
        rulesSignatures: [
          { userId: "same-user", signature: signature1 },
          { userId: "same-user", signature: signature2 }, // Same signing key counts once
        ],
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      // Should fail because only 1 distinct signing key is counted
      expect(() => service.verifyGovernanceRules(rules)).toThrow(IntegrityError);
      expect(() => service.verifyGovernanceRules(rules)).toThrow(
        "insufficient distinct valid SuperAdmin signers: got 1, need 2"
      );
    });

    it("should not let one key satisfy the threshold under different user IDs", () => {
      const mockApi = createMockGovernanceRulesApi();
      const keyPair = generateTestKeyPair();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [keyPair.publicKey],
        minValidSignatures: 2,
      });

      const rulesContainerBase64 = Buffer.from(
        JSON.stringify({ users: [] })
      ).toString("base64");
      const rulesContainerBuffer = Buffer.from(rulesContainerBase64, "base64");

      // ECDSA is randomized, so one key can emit as many distinct valid signatures as
      // the threshold demands. userId is server-supplied, so it cannot gate the count.
      const rules: GovernanceRules = {
        rulesContainer: rulesContainerBase64,
        rulesSignatures: [
          {
            userId: "attacker-label-1",
            signature: signData(keyPair.privateKey, rulesContainerBuffer),
          },
          {
            userId: "attacker-label-2",
            signature: signData(keyPair.privateKey, rulesContainerBuffer),
          },
        ],
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      expect(() => service.verifyGovernanceRules(rules)).toThrow(
        "insufficient distinct valid SuperAdmin signers: got 1, need 2"
      );
    });

    it("should count two distinct keys as two signers", () => {
      const mockApi = createMockGovernanceRulesApi();
      const keyPair1 = generateTestKeyPair();
      const keyPair2 = generateTestKeyPair();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [keyPair1.publicKey, keyPair2.publicKey],
        minValidSignatures: 2,
      });

      const rulesContainerBase64 = Buffer.from(
        JSON.stringify({ users: [] })
      ).toString("base64");
      const rulesContainerBuffer = Buffer.from(rulesContainerBase64, "base64");

      const rules: GovernanceRules = {
        rulesContainer: rulesContainerBase64,
        rulesSignatures: [
          {
            userId: "user1",
            signature: signData(keyPair1.privateKey, rulesContainerBuffer),
          },
          {
            userId: "user2",
            signature: signData(keyPair2.privateKey, rulesContainerBuffer),
          },
        ],
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      expect(() => service.verifyGovernanceRules(rules)).not.toThrow();
    });

    // This test used to assert the opposite — that a zero threshold made
    // verifyGovernanceRules return the rules untouched. That escape hatch is the
    // defect: a caller reaching this method is asking for verification, so answering
    // "fine" without checking anything is the one thing it must never do.
    it("refuses construction with a zero threshold and no keys", () => {
      const mockApi = createMockGovernanceRulesApi();

      // This was "the explicitly unverified service": constructible, and every read
      // then gated on `minValidSignatures > 0` and returned the UNVERIFIED trust root
      // with no error. The class is exported from the package barrel, so it was
      // reachable by any consumer.
      expect(
        () =>
          new GovernanceRuleService(mockApi, {
            superAdminKeys: [],
            minValidSignatures: 0,
          })
      ).toThrow(ConfigurationError);
    });

    it("rejects at construction when keys are supplied with a non-positive threshold", () => {
      const mockApi = createMockGovernanceRulesApi();

      // Keys present + threshold 0 reads as "verify with these", but every call site
      // gates on `minValidSignatures > 0`, so nothing would ever be verified.
      expect(
        () =>
          new GovernanceRuleService(mockApi, {
            superAdminKeys: [generateTestKeyPair().publicKey],
            minValidSignatures: 0,
          })
      ).toThrow(ConfigurationError);
    });

    it("should refuse construction when no SuperAdmin keys are configured", () => {
      const mockApi = createMockGovernanceRulesApi();

      // The throw moved to the CONSTRUCTOR: a keyless service used to be
      // constructible and then returned the unverified trust root on every read.
      expect(
        () =>
          new GovernanceRuleService(mockApi, {
            superAdminKeys: [],
            minValidSignatures: 1,
          })
      ).toThrow(ConfigurationError);
    });

    it("still rejects a ruleset carrying no signatures", () => {
      const mockApi = createMockGovernanceRulesApi();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [generateTestKeyPair().publicKey],
        minValidSignatures: 1,
      });

      const rules: GovernanceRules = {
        rulesContainer: Buffer.from("{}").toString("base64"),
        rulesSignatures: [{ userId: "user1", signature: "somesig" }],
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      expect(() => service.verifyGovernanceRules(rules)).toThrow(IntegrityError);
      expect(() => service.verifyGovernanceRules(rules)).toThrow(
        "insufficient distinct valid SuperAdmin signers"
      );
    });

    it("should throw when rules container is empty", () => {
      const mockApi = createMockGovernanceRulesApi();
      const keyPair = generateTestKeyPair();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [keyPair.publicKey],
        minValidSignatures: 1,
      });

      const rules: GovernanceRules = {
        rulesContainer: undefined, // Empty container
        rulesSignatures: [{ userId: "user1", signature: "somesig" }],
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      expect(() => service.verifyGovernanceRules(rules)).toThrow(IntegrityError);
      expect(() => service.verifyGovernanceRules(rules)).toThrow(
        "Rules container is empty, cannot verify"
      );
    });

    it("should throw when no signatures present", () => {
      const mockApi = createMockGovernanceRulesApi();
      const keyPair = generateTestKeyPair();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [keyPair.publicKey],
        minValidSignatures: 1,
      });

      const rules: GovernanceRules = {
        rulesContainer: Buffer.from("{}").toString("base64"),
        rulesSignatures: [], // No signatures
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      expect(() => service.verifyGovernanceRules(rules)).toThrow(IntegrityError);
      expect(() => service.verifyGovernanceRules(rules)).toThrow(
        "No signatures found on rules"
      );
    });
  });

  // Go, Java and Python all expose this; TypeScript had no getPublicKeys and no
  // SuperAdminPublicKey model at all, even though tg-protect-mcpd calls it. Without it
  // a caller cannot check that the keys it verifies against are the ones the server
  // enforces — the usual cause of a valid container failing verification.
  describe("getPublicKeys", () => {
    it("maps the configured SuperAdmin keys", async () => {
      const mockApi = createMockGovernanceRulesApi();
      mockApi.ruleServiceGetPublicKeys.mockResolvedValue({
        publicKeys: [
          { userID: "u1", publicKey: "-----BEGIN PUBLIC KEY-----k1" },
          { userID: "u2", publicKey: "-----BEGIN PUBLIC KEY-----k2" },
        ],
      });

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [generateTestKeyPair().publicKey],
        minValidSignatures: 1,
      });

      const keys = await service.getPublicKeys();

      expect(keys).toEqual([
        { userId: "u1", publicKey: "-----BEGIN PUBLIC KEY-----k1" },
        { userId: "u2", publicKey: "-----BEGIN PUBLIC KEY-----k2" },
      ]);
    });

    it("returns an empty list when the server reports no keys", async () => {
      const mockApi = createMockGovernanceRulesApi();
      mockApi.ruleServiceGetPublicKeys.mockResolvedValue({});

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [generateTestKeyPair().publicKey],
        minValidSignatures: 1,
      });

      await expect(service.getPublicKeys()).resolves.toEqual([]);
    });
  });
});
