/**
 * `WhitelistedAddressService.getEnvelope` must return what verification PRODUCED, not
 * the envelope it was handed.
 *
 * It used to run `verify()` and then return the unverified INPUT envelope, discarding
 * the result — so a caller reading the returned envelope's raw fields
 * (`signedAddress.payloadAsString`, the signature list) read attacker-controllable data
 * that had merely been *passed through* a verifier. The asset side has returned
 * `verifyEnvelope(...).verifiedEnvelope` since the 2026-09-04 pass; the address side was
 * the lone outlier and no test covered it, which is why it survived.
 *
 * Two independent gates here:
 *
 * 1. **Compile time** — the declared `Verified<SignedWhitelistedAddressEnvelope>`
 *    assignment below only type-checks if `getEnvelope` returns the branded type. An
 *    unbranded return is a `tsc` error under ts-jest, which shows up as
 *    "Test suite failed to run".
 * 2. **Runtime** — the verifier is stubbed to return an envelope that is a DIFFERENT
 *    object from its input, so returning the input fails by identity.
 */

import { WhitelistedAddressService } from "../../../src/services/whitelisted-address-service";
import type { WhitelistedAddressServiceConfig } from "../../../src/services/whitelisted-address-service";
import { WhitelistedAddressVerifier } from "../../../src/helpers/whitelisted-address-verifier";
import { NotFoundError, ValidationError } from "../../../src/errors";
import { createEmptyRulesContainer } from "../../../src/models/governance-rules";
import type {
  DecodedRulesContainer,
  RuleUserSignature,
} from "../../../src/models/governance-rules";
import type {
  SignedWhitelistedAddressEnvelope,
  WhitelistedAddressVerificationResult,
} from "../../../src/models/whitelisted-address";
import type { Verified } from "../../../src/helpers/verified";
import type { AddressWhitelistingApi } from "../../../src/internal/openapi";

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
    whitelistServiceGetWhitelistedAddress: jest.fn(),
  } as unknown as jest.Mocked<AddressWhitelistingApi>;
  return { api, svc: new WhitelistedAddressService(api, config) };
}

/** A DTO that maps to an envelope carrying the SERVER's values. */
function serverDto(): Record<string, unknown> {
  return {
    id: "1",
    metadata: {
      hash: "server-hash",
      payloadAsString: JSON.stringify({
        blockchain: "ETH",
        address: "0xSERVER",
      }),
    },
    signedAddress: {
      payloadAsString: JSON.stringify({
        blockchain: "ETH",
        address: "0xSERVER",
      }),
      signatures: [],
    },
    rulesContainer: "",
    rulesSignatures: [],
  };
}

describe("WhitelistedAddressService.getEnvelope", () => {
  afterEach(() => {
    jest.restoreAllMocks();
  });

  it("returns the envelope verification produced, not the one it was given", async () => {
    const { api, svc } = setup();
    api.whitelistServiceGetWhitelistedAddress.mockResolvedValue(
      { result: serverDto() } as never
    );

    // A sentinel that is deliberately NOT the mapped input envelope. If getEnvelope
    // returns its own `mapDtoToEnvelope(dto)` result, this identity check fails —
    // that is the whole gate.
    const verifiedSentinel = {
      id: 1,
      metadata: { hash: "verified-hash", payloadAsString: "{}" },
      signedAddress: { payloadAsString: "{}", signatures: [] },
      rulesContainer: "",
      rulesSignatures: [],
    } as unknown as Verified<SignedWhitelistedAddressEnvelope>;

    jest
      .spyOn(WhitelistedAddressVerifier.prototype, "verify")
      .mockReturnValue({
        verifiedEnvelope: verifiedSentinel,
      } as unknown as WhitelistedAddressVerificationResult);

    // Compile-time gate: this annotation does not type-check unless getEnvelope's
    // return type carries the `Verified<>` brand.
    const envelope: Verified<SignedWhitelistedAddressEnvelope> =
      await svc.getEnvelope("1");

    expect(envelope).toBe(verifiedSentinel);
    expect(envelope.metadata.hash).toBe("verified-hash");
    // Belt and braces: the server's value must not be what came back.
    expect(envelope.metadata.hash).not.toBe("server-hash");
  });

  it("propagates a verification failure instead of returning the raw envelope", async () => {
    const { api, svc } = setup();
    api.whitelistServiceGetWhitelistedAddress.mockResolvedValue(
      { result: serverDto() } as never
    );

    // Unstubbed verifier: an unsigned envelope cannot pass the 6 steps.
    await expect(svc.getEnvelope("1")).rejects.toThrow();
  });

  it("throws ValidationError when addressId is empty", async () => {
    const { svc } = setup();
    await expect(svc.getEnvelope("")).rejects.toThrow(ValidationError);
  });

  it("throws NotFoundError when the address is not found", async () => {
    const { api, svc } = setup();
    api.whitelistServiceGetWhitelistedAddress.mockResolvedValue(
      { result: undefined } as never
    );

    await expect(svc.getEnvelope("999")).rejects.toThrow(NotFoundError);
  });
});
