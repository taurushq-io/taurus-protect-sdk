/**
 * Whitelisted addresses had no approval workflow at all: both the for-approval read and
 * the approve endpoint are generated in every SDK client and were wrapped by none, so an
 * approver could not act on a whitelisted destination through the SDK — verified or not.
 *
 * The batch re-read is the part that needs guarding: it goes through the verifying path
 * filtered by ids, so a page that omits a row must abort rather than become an approval
 * of fewer rows than the caller asked for.
 */

import * as crypto from "crypto";
import { WhitelistedAddressService } from "../../../src/services/whitelisted-address-service";
import { IntegrityError, ValidationError } from "../../../src/errors";
import type { AddressWhitelistingApi } from "../../../src/internal/openapi";
import type { WhitelistedAddressServiceConfig } from "../../../src/services/whitelisted-address-service";
import type {
  DecodedRulesContainer,
  RuleUserSignature,
} from "../../../src/models/governance-rules";
import { createEmptyRulesContainer } from "../../../src/models/governance-rules";

const TEST_SUPER_ADMIN_KEY_PEM = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvW
L+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==
-----END PUBLIC KEY-----`;

const config: WhitelistedAddressServiceConfig = {
  superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
  minValidSignatures: 1,
  rulesContainerDecoder: (_b64: string): DecodedRulesContainer =>
    createEmptyRulesContainer(),
  userSignaturesDecoder: (_b64: string): RuleUserSignature[] => [],
};

function setup(): {
  api: jest.Mocked<AddressWhitelistingApi>;
  svc: WhitelistedAddressService;
} {
  const api = {
    whitelistServiceGetWhitelistedAddresses: jest.fn().mockResolvedValue({
      result: [],
      totalItems: "0",
    }),
    whitelistServiceGetWhitelistedAddressesForApproval: jest.fn(),
    whitelistServiceApproveWhitelistedAddress: jest.fn().mockResolvedValue({}),
  } as unknown as jest.Mocked<AddressWhitelistingApi>;
  return { api, svc: new WhitelistedAddressService(api, config) };
}

const { privateKey } = crypto.generateKeyPairSync("ec", {
  namedCurve: "P-256",
});

describe("WhitelistedAddressService.approve", () => {
  it("aborts when the verified read omits a requested row", async () => {
    const { api, svc } = setup();

    // Empty page: no row verified, so every requested id is missing.
    await expect(svc.approve(["1", "2"], privateKey, "ok")).rejects.toThrow(
      /was not returned by the verified read/
    );
    expect(api.whitelistServiceApproveWhitelistedAddress).not.toHaveBeenCalled();
  });

  it("reads ONE id-filtered page rather than one GET per id", async () => {
    const { api, svc } = setup();

    await expect(svc.approve(["7", "3"], privateKey, "ok")).rejects.toThrow(
      IntegrityError
    );

    expect(api.whitelistServiceGetWhitelistedAddresses).toHaveBeenCalledTimes(1);
    const args = api.whitelistServiceGetWhitelistedAddresses.mock.calls[0]?.[0] as {
      ids?: string[];
      includeForApproval?: boolean;
    };
    // Ascending: the endpoint requires it, and it makes the signed order independent
    // of the order the caller passed.
    expect(args.ids).toEqual(["3", "7"]);
    expect(args.includeForApproval).toBe(true);
  });

  it.each([
    ["no ids", [] as string[], "ok"],
    ["no comment", ["1"], ""],
    ["non-numeric id", ["abc"], "ok"],
  ])("rejects bad input: %s", async (_name, ids, comment) => {
    const { api, svc } = setup();
    await expect(svc.approve(ids, privateKey, comment)).rejects.toThrow(
      ValidationError
    );
    expect(api.whitelistServiceApproveWhitelistedAddress).not.toHaveBeenCalled();
  });

  it("requires a private key", async () => {
    const { svc } = setup();
    await expect(
      svc.approve(["1"], undefined as unknown as crypto.KeyObject, "ok")
    ).rejects.toThrow(ValidationError);
  });
});

describe("WhitelistedAddressService.listForApproval", () => {
  it("verifies rather than returning rows as-is", async () => {
    const { api, svc } = setup();
    api.whitelistServiceGetWhitelistedAddressesForApproval.mockResolvedValue({
      result: [
        { id: "1", metadata: { hash: "deadbeef", payloadAsString: "{}" } },
      ],
      totalItems: "1",
    } as never);

    // Nothing in that row can verify, and a page where no row survives must not read
    // as a complete queue.
    await expect(svc.listForApproval({ limit: 10 })).rejects.toThrow(
      IntegrityError
    );
  });

  it("rejects a non-positive limit", async () => {
    const { svc } = setup();
    await expect(svc.listForApproval({ limit: 0 })).rejects.toThrow(
      ValidationError
    );
  });
});
