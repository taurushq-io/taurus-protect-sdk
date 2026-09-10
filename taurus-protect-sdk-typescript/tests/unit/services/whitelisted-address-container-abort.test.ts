/**
 * A container-level integrity failure must ABORT a whitelisted-address list call;
 * a row-level one must be EXCLUDED and reported.
 *
 * These are different failure classes and the distinction is load-bearing:
 *
 * - `ContainerIntegrityError` means this SDK cannot interpret the rules container the
 *   rows are being judged against. That invalidates EVERY row judged against it, not
 *   just the one being processed. Excluding them one by one would empty the whitelist
 *   and report success — a caller asking "is this destination approved?" would get a
 *   false negative on every address, with a 200-shaped result.
 * - A row-level `IntegrityError` / `WhitelistError` is a property of that one row, so
 *   the list survives and names what it withheld. Listing is how an operator finds the
 *   bad row, so failing the whole call would hide its own cause.
 * - Anything else is a defect in this SDK, and reporting a `TypeError` to the caller as
 *   "this address failed verification" hides it.
 *
 * Go re-throws via `errors.As`, Python via `except ContainerIntegrityError: raise`, and
 * Java by the exception being unchecked. This SDK caught everything with a blanket
 * `catch (error: unknown)` and turned all three classes into row exclusions — it was the
 * lone outlier, and nothing tested it.
 */

import { createHash } from "crypto";
import { WhitelistedAddressService } from "../../../src/services/whitelisted-address-service";
import {
  ContainerIntegrityError,
  IntegrityError,
  WhitelistError,
} from "../../../src/errors";
import type { AddressWhitelistingApi } from "../../../src/internal/openapi";
import type { WhitelistedAddressServiceConfig } from "../../../src/services/whitelisted-address-service";
import type {
  DecodedRulesContainer,
  RuleUserSignature,
} from "../../../src/models/governance-rules";
import { createEmptyRulesContainer } from "../../../src/models/governance-rules";
import type { WhitelistedAddress } from "../../../src/models/whitelisted-address";

// A valid P-256 public key (test key, no production value).
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

// Two rows carrying no rulesContainers array, so the per-page container cache stays
// empty and every row reaches the stubbed verifier.
const rows = [
  { id: "1", metadata: { hash: "h1", payloadAsString: '{"currency":"ETH"}' } },
  { id: "2", metadata: { hash: "h2", payloadAsString: '{"currency":"ETH"}' } },
];

const goodResult = {
  verifiedWhitelistedAddress: {
    id: 2,
    address: "0xabc",
  } as unknown as WhitelistedAddress,
  rulesContainer: createEmptyRulesContainer(),
  verifiedHash: "h2",
};

/**
 * Builds a service whose verifier fails on the FIRST row with `firstRowError` and
 * succeeds on every other row. Both list endpoints return the same two rows.
 */
function serviceThrowingOnFirstRow(firstRowError: Error): WhitelistedAddressService {
  const reply = { result: rows, totalItems: "2" };
  const api = {
    whitelistServiceGetWhitelistedAddresses: jest.fn().mockResolvedValue(reply),
    whitelistServiceGetWhitelistedAddressesForApproval: jest
      .fn()
      .mockResolvedValue(reply),
  } as unknown as AddressWhitelistingApi;

  const service = new WhitelistedAddressService(api, config);

  let call = 0;
  (
    service as unknown as { verifier: { verify: jest.Mock } }
  ).verifier.verify = jest.fn(() => {
    call += 1;
    if (call === 1) {
      throw firstRowError;
    }
    return goodResult;
  });

  return service;
}

describe("container-level integrity failures abort the list", () => {
  it("list() rejects with ContainerIntegrityError instead of excluding the row", async () => {
    const service = serviceThrowingOnFirstRow(
      new ContainerIntegrityError("line 0 carries a source cell this SDK cannot interpret")
    );

    await expect(service.list()).rejects.toBeInstanceOf(ContainerIntegrityError);
  });

  it("listForApproval() rejects with ContainerIntegrityError too", async () => {
    const service = serviceThrowingOnFirstRow(
      new ContainerIntegrityError("line 0 carries a source cell this SDK cannot interpret")
    );

    await expect(service.listForApproval()).rejects.toBeInstanceOf(
      ContainerIntegrityError
    );
  });

  it("does not let the surviving rows mask the failure", async () => {
    // The "rows returned but none survived" backstop only catches the homogeneous case.
    // Here row 2 verifies fine, so without the re-throw the call would return 1 item
    // plus 1 exclusion and look like a success.
    const service = serviceThrowingOnFirstRow(
      new ContainerIntegrityError("uninterpretable container")
    );

    await expect(service.list()).rejects.toThrow(/uninterpretable container/);
  });
});

describe("row-level integrity failures are excluded, not fatal", () => {
  it("list() excludes an IntegrityError row and returns the rest", async () => {
    const service = serviceThrowingOnFirstRow(new IntegrityError("bad row"));

    const result = await service.list();

    expect(result.items).toHaveLength(1);
    expect(result.excludedUnverified).toEqual([{ id: "1", reason: "bad row" }]);
    // totalItems reduced by the exclusion: the server counted 2, the caller gets 1.
    expect(result.pagination?.totalItems).toBe(1);
  });

  it("list() excludes a WhitelistError row and returns the rest", async () => {
    const service = serviceThrowingOnFirstRow(
      new WhitelistError("threshold not met")
    );

    const result = await service.list();

    expect(result.items).toHaveLength(1);
    expect(result.excludedUnverified).toEqual([
      { id: "1", reason: "threshold not met" },
    ]);
  });
});

describe("the row-to-container binding is authenticated, not taken on trust", () => {
  // `rulesContainerHash` is the LABEL a row uses to pick its container out of the
  // normalized `rulesContainers` array — and both the label and the container come from
  // the same response. validatord derives the label as
  // base64(SHA256(<the base64 container text>)), so the SDK can recompute it. Without
  // that check a server can file container A under container B's label and steer any row
  // to any OTHER validly-signed container — an older ruleset with a weaker group
  // threshold, say. Both containers pass the SuperAdmin check, so step 2 does not catch
  // it.
  it("rejects a container filed under a label that is not its own", async () => {
    const reply = {
      result: rows,
      totalItems: "2",
      rulesContainers: [
        {
          // A label that does not match these bytes.
          hash: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
          rulesContainer: "c3RhbGUtY29udGFpbmVy",
          rulesSignatures: "sigs",
        },
      ],
    };
    const api = {
      whitelistServiceGetWhitelistedAddresses: jest.fn().mockResolvedValue(reply),
    } as unknown as AddressWhitelistingApi;

    const service = new WhitelistedAddressService(api, config);

    await expect(service.list()).rejects.toThrow(/hash mismatch/);
  });

  it("accepts a container whose label matches its bytes", async () => {
    // Same shape, correct label: the mismatch check must not reject honest responses.
    // It still fails later (the stub container does not verify), so assert only that the
    // failure is no longer the label check.
    const containerBase64 = "c3RhbGUtY29udGFpbmVy";
    const reply = {
      result: rows,
      totalItems: "2",
      rulesContainers: [
        {
          hash: createHash("sha256").update(containerBase64, "utf-8").digest("base64"),
          rulesContainer: containerBase64,
          rulesSignatures: "sigs",
        },
      ],
    };
    const api = {
      whitelistServiceGetWhitelistedAddresses: jest.fn().mockResolvedValue(reply),
    } as unknown as AddressWhitelistingApi;

    const service = new WhitelistedAddressService(api, config);

    await expect(service.list()).rejects.not.toThrow(/hash mismatch/);
  });
});

describe("a defect is not a verification failure", () => {
  // The blanket catch reported a TypeError from a genuine bug to the caller as
  // "this address failed verification", where it was indistinguishable from a
  // tampered row. It must escape the row loop instead.
  //
  // It surfaces as ServerError, not TypeError: BaseService.execute() passes the known
  // integrity types through untouched (which is why the ContainerIntegrityError cases
  // above assert the exact class) and wraps anything unrecognised. The point of this
  // test is that the call REJECTS and the original message survives — not which
  // wrapper the base class chose.
  it("list() surfaces a defect instead of reporting it as an excluded row", async () => {
    const service = serviceThrowingOnFirstRow(
      new TypeError("cannot read property 'x' of undefined")
    );

    await expect(service.list()).rejects.toThrow(/cannot read property 'x'/);
  });
});
