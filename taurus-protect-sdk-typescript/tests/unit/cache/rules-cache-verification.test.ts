/**
 * Gate: the rules container feeding address signature verification must have its
 * SuperAdmin signatures verified before anything reads the HSM public key out of it.
 *
 * This SDK fetched that container through the generated API and the raw mapper, so it
 * was decoded unverified and `rulesSignatures` was discarded — AddressService then
 * verified addresses against an HSM key nothing had authenticated. Go, Java and Python
 * all fetch through their governance service; only TypeScript did not, and nothing
 * tested the path, which is why it survived.
 *
 * Drives the real client wiring rather than an injected mock, because the defect WAS
 * the wiring: every service-level test passed throughout.
 */

import { generateKeyPairSync } from "crypto";

import { ProtectClient } from "../../../src/client";
import { IntegrityError } from "../../../src/errors";
import { AddressesApi } from "../../../src/internal/openapi/apis/AddressesApi";
import { GovernanceRulesApi } from "../../../src/internal/openapi/apis/GovernanceRulesApi";
import { rulesContainerToBase64 } from "../../../src/mappers/protobuf-rules-container-encode";
import { createEmptyRulesContainer } from "../../../src/models/governance-rules";

function superAdminKeyPem(): string {
  const { publicKey } = generateKeyPairSync("ec", { namedCurve: "prime256v1" });
  return publicKey.export({ type: "spki", format: "pem" }).toString();
}

function createClient(): ProtectClient {
  return ProtectClient.create({
    host: "https://protect.example.com",
    apiKey: "test-api-key-12345",
    apiSecret: "aabbccdd11223344556677889900aabbccddeeff",
    superAdminKeysPem: [superAdminKeyPem()],
    minValidSignatures: 1,
  });
}

function stubAddress(): void {
  jest
    .spyOn(AddressesApi.prototype, "walletServiceGetAddress")
    .mockResolvedValue({
      result: {
        id: "1",
        walletId: "100",
        address: "0xaaa",
        currency: "ETH",
        signature: "c2ln",
      },
    } as never);
}

describe("rules container verification on the address path", () => {
  let client: ProtectClient | undefined;

  afterEach(() => {
    if (client && !client.isClosed) {
      client.close();
    }
    client = undefined;
    jest.restoreAllMocks();
  });

  it("refuses an unsigned rules container instead of trusting its HSM key", async () => {
    // Wire-valid container carrying NO SuperAdmin signatures — the shape an on-path
    // attacker or a compromised endpoint supplies to swap in their own HSM key.
    jest
      .spyOn(GovernanceRulesApi.prototype, "ruleServiceGetRules")
      .mockResolvedValue({
        result: {
          rulesContainer: rulesContainerToBase64(createEmptyRulesContainer()),
          rulesSignatures: [],
        },
      } as never);
    stubAddress();

    client = createClient();

    await expect(client.addresses.get(1)).rejects.toThrow(IntegrityError);
  });

  it("refuses a container whose signatures verify against none of the configured keys", async () => {
    // Signatures present but produced by a key the client does not trust, so the
    // distinct-signer threshold cannot be met.
    jest
      .spyOn(GovernanceRulesApi.prototype, "ruleServiceGetRules")
      .mockResolvedValue({
        result: {
          rulesContainer: rulesContainerToBase64(createEmptyRulesContainer()),
          rulesSignatures: [
            { userId: "attacker", signature: "YmFkLXNpZ25hdHVyZS1ub3QtZWNkc2E=" },
          ],
        },
      } as never);
    stubAddress();

    client = createClient();

    await expect(client.addresses.get(1)).rejects.toThrow(IntegrityError);
  });

  it("reaches the governance service rather than the generated rules API", async () => {
    // The cache must not call ruleServiceGetRules on its own: that path skips
    // verification entirely. It reaches it only THROUGH GovernanceRuleService, which
    // verifies first — so a rejection here proves the verifying path ran.
    const rulesSpy = jest
      .spyOn(GovernanceRulesApi.prototype, "ruleServiceGetRules")
      .mockResolvedValue({
        result: {
          rulesContainer: rulesContainerToBase64(createEmptyRulesContainer()),
          rulesSignatures: [],
        },
      } as never);
    stubAddress();

    client = createClient();

    await expect(client.addresses.get(1)).rejects.toThrow(IntegrityError);
    expect(rulesSpy).toHaveBeenCalled();
  });
});
