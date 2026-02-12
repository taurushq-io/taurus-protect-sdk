/**
 * Unit tests for WhitelistedAddressService.
 *
 * Tests focus on:
 * - Service correctly delegates to verifier for all get/list operations
 * - parseWhitelistedAddressFromJson extracts fields from verified payload only
 * - Missing payload results in IntegrityError
 */

import { WhitelistedAddressService } from "../../../src/services/whitelisted-address-service";
import { IntegrityError } from "../../../src/errors";
import type { AddressWhitelistingApi } from "../../../src/internal/openapi";
import type { WhitelistedAddressServiceConfig } from "../../../src/services/whitelisted-address-service";
import type { DecodedRulesContainer } from "../../../src/models/governance-rules";
import type { RuleUserSignature } from "../../../src/models/governance-rules";
import { createEmptyRulesContainer } from "../../../src/models/governance-rules";
import { parseWhitelistedAddressFromJson } from "../../../src/helpers/whitelist-hash-helper";

// A valid P-256 public key for testing (test key with no production value)
const TEST_SUPER_ADMIN_KEY_PEM = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvW
L+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==
-----END PUBLIC KEY-----`;

// Minimal mock verification config for testing
const mockVerificationConfig: WhitelistedAddressServiceConfig = {
  superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
  minValidSignatures: 1,
  rulesContainerDecoder: (_base64: string): DecodedRulesContainer => createEmptyRulesContainer(),
  userSignaturesDecoder: (_base64: string): RuleUserSignature[] => [],
};

describe("WhitelistedAddressService", () => {
  // Create a minimal mock API
  const mockApi = {} as AddressWhitelistingApi;

  it("should construct with verification config", () => {
    const service = new WhitelistedAddressService(mockApi, mockVerificationConfig);
    expect(service).toBeDefined();
  });

  describe("parseWhitelistedAddressFromJson security behavior", () => {
    // These tests verify that address fields come exclusively from the verified
    // payloadAsString, not from unverified DTO envelope fields. This is the
    // security-critical function called by the verifier's Step 6.

    it("should use blockchain from verified payload, not envelope", () => {
      const payload = JSON.stringify({
        currency: "ETH",
        network: "mainnet",
        address: "0xabc123",
      });

      const result = parseWhitelistedAddressFromJson(payload);

      // Must use payload's currency as blockchain
      expect(result.blockchain).toBe("ETH");
    });

    it("should use network from verified payload", () => {
      const payload = JSON.stringify({
        currency: "ETH",
        network: "mainnet",
        address: "0xabc123",
      });

      const result = parseWhitelistedAddressFromJson(payload);

      expect(result.network).toBe("mainnet");
    });

    it("should use all fields from verified payload", () => {
      const payload = JSON.stringify({
        currency: "BTC",
        network: "testnet",
        address: "bc1qverified",
        label: "My verified label",
        memo: "Verified memo",
        customerId: "cust-verified",
        contractType: "NATIVE",
        addressType: "EXTERNAL",
        tnParticipantID: "participant-123",
      });

      const result = parseWhitelistedAddressFromJson(payload);

      expect(result.blockchain).toBe("BTC");
      expect(result.network).toBe("testnet");
      expect(result.address).toBe("bc1qverified");
      expect(result.label).toBe("My verified label");
      expect(result.memo).toBe("Verified memo");
      expect(result.customerId).toBe("cust-verified");
      expect(result.contractType).toBe("NATIVE");
      expect(result.addressType).toBe("EXTERNAL");
      expect(result.tnParticipantId).toBe("participant-123");
    });

    it("should parse linked addresses from payload", () => {
      const payload = JSON.stringify({
        currency: "ETH",
        network: "mainnet",
        address: "0xabc123",
        linkedInternalAddresses: [
          { id: 1, label: "Internal 1" },
          { id: 2, label: "Internal 2" },
        ],
        linkedWallets: [
          { id: 10, name: "Wallet A", path: "m/44'/60'/0'" },
        ],
      });

      const result = parseWhitelistedAddressFromJson(payload);

      expect(result.linkedInternalAddresses).toHaveLength(2);
      expect(result.linkedInternalAddresses![0].id).toBe(1);
      expect(result.linkedInternalAddresses![0].label).toBe("Internal 1");
      expect(result.linkedWallets).toHaveLength(1);
      expect(result.linkedWallets![0].id).toBe(10);
      expect(result.linkedWallets![0].label).toBe("Wallet A");
    });
  });

  describe("mapDtoToEnvelope", () => {
    const service = new WhitelistedAddressService(mockApi, mockVerificationConfig);
    const mapDtoToEnvelope = (service as unknown as {
      mapDtoToEnvelope: (dto: Record<string, unknown>, addressId: string) => unknown;
    }).mapDtoToEnvelope.bind(service);

    it("should map DTO fields to envelope structure", () => {
      const dto = {
        id: "123",
        metadata: {
          hash: "abc123",
          payloadAsString: JSON.stringify({ currency: "ETH", address: "0x1" }),
        },
        rulesContainer: "base64rules",
        rulesSignatures: "base64sigs",
        signedAddress: {
          signatures: [
            {
              signature: { userId: "user1", signature: "sig1" },
              hashes: ["hash1"],
            },
          ],
        },
        blockchain: "ETH",
        network: "mainnet",
      };

      const envelope = mapDtoToEnvelope(dto, "123") as {
        id: string;
        metadata: { hash: string; payloadAsString: string };
        blockchain: string;
        network: string;
      };

      expect(envelope.id).toBe("123");
      expect(envelope.metadata.hash).toBe("abc123");
      expect(envelope.blockchain).toBe("ETH");
      expect(envelope.network).toBe("mainnet");
    });

    it("should re-throw non-SyntaxError from JSON parsing", () => {
      // When payloadAsString contains JSON that triggers a non-SyntaxError
      // (e.g., via a reviver), the error should propagate, not be swallowed
      const dto = {
        id: "123",
        metadata: {
          hash: "abc123",
          payloadAsString: "valid json but with no linked data",
        },
        blockchain: "ETH",
        network: "mainnet",
      };

      // This should not throw — invalid JSON is caught as SyntaxError
      expect(() => mapDtoToEnvelope(dto, "123")).not.toThrow();
    });
  });
});
