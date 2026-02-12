/**
 * Unit tests for GovernanceRuleService.
 *
 * These tests verify the critical security feature:
 * - ECDSA signature verification for governance rules
 */

import * as crypto from "crypto";
import { GovernanceRuleService } from "../../../src/services/governance-rule-service";
import { IntegrityError } from "../../../src/errors";
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
        "Insufficient valid signatures: found 0, required 1"
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
          { userId: "same-user", signature: signature2 }, // Duplicate user should be ignored
        ],
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      // Should fail because only 1 distinct user signature is counted
      expect(() => service.verifyGovernanceRules(rules)).toThrow(IntegrityError);
      expect(() => service.verifyGovernanceRules(rules)).toThrow(
        "Insufficient valid signatures: found 1, required 2"
      );
    });

    it("should skip verification when minValidSignatures is 0", () => {
      const mockApi = createMockGovernanceRulesApi();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [],
        minValidSignatures: 0, // Verification disabled
      });

      const rules: GovernanceRules = {
        rulesContainer: Buffer.from("{}").toString("base64"),
        rulesSignatures: [], // No signatures
        locked: true,
        creationDate: new Date(),
        updateDate: new Date(),
        trails: [],
      };

      // Should not throw when verification is disabled
      expect(() => service.verifyGovernanceRules(rules)).not.toThrow();
    });

    it("should throw when no SuperAdmin keys configured", () => {
      const mockApi = createMockGovernanceRulesApi();

      const service = new GovernanceRuleService(mockApi, {
        superAdminKeys: [],
        minValidSignatures: 1, // Verification enabled but no keys
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
        "No SuperAdmin keys configured for verification"
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
});
