/**
 * The v2 asset-holders list returns an address only after the SDK verified it.
 *
 * The holder rows carry no signature, so `queryAssetAddresses` re-reads each internal row
 * through the managed-address list (HSM signature) and each whitelisted row through the
 * whitelisted-address list (6-step verification), keeping a row only when the verified
 * reader returns the same address under its id. Everything runs through the REAL
 * generated APIs over an injected `fetchApi`, so the deserializers and both verified
 * readers run exactly as in production.
 */

import * as crypto from "crypto";

import type { RulesContainerCache } from "../../../src/cache";
import { calculateHexHash, encodePublicKeyPem, signData } from "../../../src/crypto";
import {
  APIError,
  ContainerIntegrityError,
  IntegrityError,
} from "../../../src/errors";
import {
  AddressWhitelistingApi,
  AddressesApi,
  AssetV2Api,
  AssetsApi,
} from "../../../src/internal/openapi/apis";
import { Configuration } from "../../../src/internal/openapi/runtime";
import {
  createEmptyRulesContainer,
  type DecodedRulesContainer,
} from "../../../src/models/governance-rules";
import { AddressService } from "../../../src/services/address-service";
import { AssetService } from "../../../src/services/asset-service";
import { WhitelistedAddressService } from "../../../src/services/whitelisted-address-service";

const BASE_PATH = "https://protect.example.invalid";
const HOLDERS_PATH = "/api/rest/v2/assets/asset-1/addresses/query";
const ADDRESSES_PATH = "/api/rest/v1/addresses";
const WHITELIST_PATH = "/api/rest/v1/whitelists/addresses";

const INTERNAL = "ADDRESS_TYPE_V2_INTERNAL";
const WHITELISTED = "ADDRESS_TYPE_V2_WHITELISTED";
const EXTERNAL = "ADDRESS_TYPE_V2_EXTERNAL";

interface Recorded {
  readonly path: string;
  readonly query: Array<[string, string]>;
}

type Route = (request: Recorded) => { status?: number; body: unknown };

/** A fetch stand-in that answers by path and records every request. */
class RoutedTransport {
  readonly requests: Recorded[] = [];
  readonly configuration: Configuration;

  constructor(private readonly routes: Record<string, Route>) {
    this.configuration = new Configuration({
      basePath: BASE_PATH,
      fetchApi: async (input) => {
        const href =
          typeof input === "string" ? input : input instanceof URL ? input.href : input.url;
        const url = new URL(href);
        const query: Array<[string, string]> = [];
        url.searchParams.forEach((value, key) => query.push([key, value]));
        const request = { path: url.pathname, query };
        this.requests.push(request);
        const route = this.routes[url.pathname];
        if (route === undefined) {
          throw new Error(`unexpected request to ${url.pathname}`);
        }
        const { status = 200, body } = route(request);
        return new Response(JSON.stringify(body), {
          status,
          headers: { "Content-Type": "application/json" },
        });
      },
    });
  }

  to(path: string): Recorded[] {
    return this.requests.filter((r) => r.path === path);
  }
}

function values(request: Recorded, key: string): string[] {
  return request.query.filter(([k]) => k === key).map(([, v]) => v);
}

// --- managed addresses: signed by a throwaway HSM key -------------------------------

const hsm = crypto.generateKeyPairSync("ec", { namedCurve: "P-256" });
const otherKey = crypto.generateKeyPairSync("ec", { namedCurve: "P-256" });

function hsmContainer(): DecodedRulesContainer {
  return {
    ...createEmptyRulesContainer(),
    users: [
      {
        id: "hsm-slot",
        name: "HSM",
        publicKeyPem: encodePublicKeyPem(hsm.publicKey),
        roles: ["HSMSLOT"],
      },
    ],
  };
}

function rulesCacheOf(container: DecodedRulesContainer): RulesContainerCache {
  return {
    get: () => Promise.resolve(container),
    clear: () => undefined,
  } as unknown as RulesContainerCache;
}

/** A managed address row as GetAddresses returns it, HSM-signed unless told otherwise. */
function managedAddress(
  id: string,
  address: string,
  signer: crypto.KeyObject = hsm.privateKey
): Record<string, unknown> {
  return {
    id,
    address,
    walletId: "1",
    currency: "ETH",
    signature: signData(signer, Buffer.from(address, "utf-8")),
  };
}

/** Answers GetAddresses with the given rows, whatever ids were asked for. */
function addressesReply(...rows: Array<Record<string, unknown>>): Route {
  return () => ({ body: { result: rows, totalItems: String(rows.length) } });
}

// --- whitelisted addresses: a full 6-step chain on throwaway keys --------------------

const superAdmin = crypto.generateKeyPairSync("ec", { namedCurve: "P-256" });
const approver = crypto.generateKeyPairSync("ec", { namedCurve: "P-256" });
const APPROVER_ID = "approver@example.invalid";
const RULES_B64 = Buffer.from(JSON.stringify({ rules: "throwaway" })).toString("base64");
const RULES_SIGNATURE = signData(superAdmin.privateKey, Buffer.from(RULES_B64, "base64"));

function whitelistContainer(): DecodedRulesContainer {
  return {
    ...createEmptyRulesContainer(),
    users: [
      {
        id: APPROVER_ID,
        name: "Approver",
        publicKeyPem: encodePublicKeyPem(approver.publicKey),
        roles: ["USER"],
      },
    ],
    groups: [{ id: "approvers", name: "Approvers", userIds: [APPROVER_ID] }],
    addressWhitelistingRules: [
      {
        currency: "ETH",
        network: "mainnet",
        parallelThresholds: [
          { thresholds: [{ groupId: "approvers", minimumSignatures: 1, threshold: 0 }] },
        ],
        lines: [],
      },
    ],
  } as DecodedRulesContainer;
}

/** A whitelisted address envelope whose 6-step chain verifies, unless `tamper` is set. */
function whitelistedEnvelope(
  id: string,
  address: string,
  tamper = false
): Record<string, unknown> {
  const payloadAsString = JSON.stringify({
    currency: "ETH",
    addressType: "individual",
    address,
    memo: "",
    label: `holder ${id}`,
    customerId: "",
    exchangeAccountId: "",
    linkedInternalAddresses: [],
    contractType: "",
  });
  const hash = calculateHexHash(payloadAsString);
  const hashes = [hash];
  return {
    id,
    blockchain: "ETH",
    network: "mainnet",
    // A hash that no longer matches the payload fails step 1.
    metadata: { hash: tamper ? calculateHexHash("tampered") : hash, payloadAsString },
    rulesContainer: RULES_B64,
    rulesSignatures: Buffer.from("signatures").toString("base64"),
    signedAddress: {
      signatures: [
        {
          signature: {
            userId: APPROVER_ID,
            signature: signData(approver.privateKey, Buffer.from(JSON.stringify(hashes), "utf-8")),
          },
          hashes,
        },
      ],
    },
  };
}

function whitelistReply(...rows: Array<Record<string, unknown>>): Route {
  return () => ({ body: { result: rows, totalItems: String(rows.length) } });
}

// --- the service under test -----------------------------------------------------------

function holdersReply(
  rows: Array<Record<string, unknown>>,
  cursor?: Record<string, unknown>
): Route {
  return () => ({ body: { result: rows, ...(cursor ? { cursor } : {}) } });
}

function assetService(
  transport: RoutedTransport,
  container: DecodedRulesContainer = hsmContainer()
): AssetService {
  const c = transport.configuration;
  const rulesCache = rulesCacheOf(container);
  const addresses = new AddressService(new AddressesApi(c), rulesCache);
  const whitelisted = new WhitelistedAddressService(new AddressWhitelistingApi(c), {
    superAdminKeysPem: [encodePublicKeyPem(superAdmin.publicKey)],
    minValidSignatures: 1,
    rulesContainerDecoder: () => whitelistContainer(),
    userSignaturesDecoder: () => [{ userId: "superadmin", signature: RULES_SIGNATURE }],
  });
  return new AssetService(new AssetsApi(c), rulesCache, new AssetV2Api(c), {
    addresses: () => addresses,
    whitelistedAddresses: () => whitelisted,
  });
}

/** An external holder: on-chain data that keeps a page from being all-excluded. */
const EXTERNAL_ROW = { address: "0xEXT", addressType: EXTERNAL, balance: "1" };

describe("queryAssetAddresses: internal holders", () => {
  it("keeps a holder whose HSM-verified address matches, with the verified address", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xINT", addressType: INTERNAL, addressID: "11", balance: "5" },
      ]),
      [ADDRESSES_PATH]: addressesReply(managedAddress("11", "0xINT")),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.items).toEqual([
      {
        address: "0xINT",
        addressType: INTERNAL,
        addressId: "11",
        balance: "5",
        kycStatus: undefined,
        whitelistedAddressId: undefined,
        verified: true,
      },
    ]);
    expect(page.excludedUnverified).toEqual([]);
    const [read] = transport.to(ADDRESSES_PATH);
    expect(values(read, "addressIds")).toEqual(["11"]);
    expect(values(read, "limit")).toEqual(["1"]);
  });

  it("excludes a holder whose address differs from the verified one", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xATTACKER", addressType: INTERNAL, addressID: "12" },
        EXTERNAL_ROW,
      ]),
      [ADDRESSES_PATH]: addressesReply(managedAddress("12", "0xGENUINE")),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.items.map((h) => h.address)).toEqual(["0xEXT"]);
    expect(page.excludedUnverified).toHaveLength(1);
    expect(page.excludedUnverified[0].id).toBe("12");
    expect(page.excludedUnverified[0].reason).toMatch(/differs from the verified/);
  });

  it("excludes a holder whose id the verified reader did not return", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xINT", addressType: INTERNAL, addressID: "13" },
        EXTERNAL_ROW,
      ]),
      [ADDRESSES_PATH]: addressesReply(),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.excludedUnverified).toEqual([
      { id: "13", reason: expect.stringMatching(/was not returned by the verified/) },
    ]);
  });

  it("excludes an internal holder without an addressID, by its address, with no read", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([{ address: "0xNOID", addressType: INTERNAL }, EXTERNAL_ROW]),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.excludedUnverified).toEqual([
      { id: "0xNOID", reason: expect.stringMatching(/carries no addressID/) },
    ]);
    expect(transport.requests).toHaveLength(1);
  });

  it("excludes a holder whose HSM signature fails in the verified reader", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xBAD", addressType: INTERNAL, addressID: "14" },
        { address: "0xGOOD", addressType: INTERNAL, addressID: "15" },
      ]),
      [ADDRESSES_PATH]: addressesReply(
        managedAddress("14", "0xBAD", otherKey.privateKey),
        managedAddress("15", "0xGOOD")
      ),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.items.map((h) => [h.address, h.verified])).toEqual([["0xGOOD", true]]);
    expect(page.excludedUnverified).toEqual([
      { id: "14", reason: expect.stringMatching(/signature verification failed/) },
    ]);
  });

  it("excludes an id the verified reader returned twice", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xDUP", addressType: INTERNAL, addressID: "16" },
        EXTERNAL_ROW,
      ]),
      [ADDRESSES_PATH]: addressesReply(
        managedAddress("16", "0xDUP"),
        managedAddress("16", "0xOTHER")
      ),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.excludedUnverified).toEqual([
      { id: "16", reason: expect.stringMatching(/returned more than once/) },
    ]);
  });

  it("re-reads 51 internal ids as exactly two requests, 50 + 1", async () => {
    const ids = Array.from({ length: 51 }, (_v, i) => String(i + 1));
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply(
        ids.map((id) => ({ address: `0xA${id}`, addressType: INTERNAL, addressID: id }))
      ),
      // Echo a signed row for every id the request names.
      [ADDRESSES_PATH]: (request) => ({
        body: {
          result: values(request, "addressIds").map((id) => managedAddress(id, `0xA${id}`)),
        },
      }),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    const reads = transport.to(ADDRESSES_PATH);
    expect(reads.map((r) => values(r, "addressIds").length)).toEqual([50, 1]);
    expect(reads.map((r) => values(r, "limit"))).toEqual([["50"], ["1"]]);
    expect(values(reads[0], "addressIds")).toEqual(ids.slice(0, 50));
    expect(values(reads[1], "addressIds")).toEqual(["51"]);
    expect(page.items).toHaveLength(51);
    expect(page.items.every((h) => h.verified)).toBe(true);
  });

  it("aborts when the rules container cannot verify anything (no HSM key)", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xINT", addressType: INTERNAL, addressID: "17" },
        EXTERNAL_ROW,
      ]),
      [ADDRESSES_PATH]: addressesReply(managedAddress("17", "0xINT")),
    });

    await expect(
      assetService(transport, createEmptyRulesContainer()).queryAssetAddresses("asset-1")
    ).rejects.toBeInstanceOf(ContainerIntegrityError);
  });

  it("propagates an error from the verified reader instead of excluding", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xINT", addressType: INTERNAL, addressID: "18" },
        EXTERNAL_ROW,
      ]),
      [ADDRESSES_PATH]: () => ({ status: 500, body: { message: "down" } }),
    });

    await expect(
      assetService(transport).queryAssetAddresses("asset-1")
    ).rejects.toBeInstanceOf(APIError);
  });
});

describe("queryAssetAddresses: whitelisted holders", () => {
  it("keeps a holder whose 6-step-verified address matches, with the verified address", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xWL", addressType: WHITELISTED, whitelistedAddressID: "501" },
      ]),
      [WHITELIST_PATH]: whitelistReply(whitelistedEnvelope("501", "0xWL")),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.items.map((h) => [h.address, h.whitelistedAddressId, h.verified])).toEqual([
      ["0xWL", "501", true],
    ]);
    const [read] = transport.to(WHITELIST_PATH);
    expect(values(read, "ids")).toEqual(["501"]);
    expect(values(read, "limit")).toEqual(["1"]);
    expect(values(read, "rulesContainerNormalized")).toEqual(["true"]);
    expect(transport.to(ADDRESSES_PATH)).toHaveLength(0);
  });

  it("excludes a holder whose address differs from the verified entry", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xATTACKER", addressType: WHITELISTED, whitelistedAddressID: "502" },
        EXTERNAL_ROW,
      ]),
      [WHITELIST_PATH]: whitelistReply(whitelistedEnvelope("502", "0xGENUINE")),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.excludedUnverified).toEqual([
      { id: "502", reason: expect.stringMatching(/differs from the verified/) },
    ]);
  });

  it("excludes a holder whose id the verified reader did not return", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xWL", addressType: WHITELISTED, whitelistedAddressID: "503" },
        EXTERNAL_ROW,
      ]),
      [WHITELIST_PATH]: whitelistReply(),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.excludedUnverified).toEqual([
      { id: "503", reason: expect.stringMatching(/was not returned by the verified/) },
    ]);
  });

  it("excludes a whitelisted holder without a whitelistedAddressID, with no read", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([{ address: "0xWLNOID", addressType: WHITELISTED }, EXTERNAL_ROW]),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.excludedUnverified).toEqual([
      { id: "0xWLNOID", reason: expect.stringMatching(/carries no whitelistedAddressID/) },
    ]);
    expect(transport.requests).toHaveLength(1);
  });

  it("excludes a holder whose entry fails the 6-step verification", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xTAMPERED", addressType: WHITELISTED, whitelistedAddressID: "504" },
        { address: "0xWL", addressType: WHITELISTED, whitelistedAddressID: "505" },
      ]),
      [WHITELIST_PATH]: whitelistReply(
        whitelistedEnvelope("504", "0xTAMPERED", true),
        whitelistedEnvelope("505", "0xWL")
      ),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.items.map((h) => h.whitelistedAddressId)).toEqual(["505"]);
    expect(page.excludedUnverified.map((e) => e.id)).toEqual(["504"]);
  });
});

describe("queryAssetAddresses: the page", () => {
  it("returns external and untyped holders unverified, with no extra request", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        // An id on a non-internal row is not a reason to read it.
        { address: "0xEXT", addressType: EXTERNAL, addressID: "99" },
        { address: "0xUNTYPED" },
      ]),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page.items.map((h) => [h.address, h.verified])).toEqual([
      ["0xEXT", false],
      ["0xUNTYPED", false],
    ]);
    expect(page.excludedUnverified).toEqual([]);
    expect(transport.requests.map((r) => r.path)).toEqual([HOLDERS_PATH]);
  });

  it("throws when rows came back and none survived", async () => {
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply([
        { address: "0xA", addressType: INTERNAL, addressID: "21" },
        { address: "0xB", addressType: WHITELISTED, whitelistedAddressID: "22" },
      ]),
      [ADDRESSES_PATH]: addressesReply(),
      [WHITELIST_PATH]: whitelistReply(),
    });

    const call = assetService(transport).queryAssetAddresses("asset-1");
    await expect(call).rejects.toBeInstanceOf(IntegrityError);
    await expect(call).rejects.toThrow(/all 2 asset holder\(s\) failed verification/);
  });

  it("returns an empty page for an empty reply, with no extra request", async () => {
    const transport = new RoutedTransport({ [HOLDERS_PATH]: holdersReply([]) });

    const page = await assetService(transport).queryAssetAddresses("asset-1");

    expect(page).toEqual({
      items: [],
      excludedUnverified: [],
      pagination: { pageSize: 20, nextCursor: "", hasMore: false },
    });
    expect(transport.requests).toHaveLength(1);
  });

  it("never moves the cursor for an exclusion", async () => {
    const token = Buffer.from("holders page 2").toString("base64");
    const transport = new RoutedTransport({
      [HOLDERS_PATH]: holdersReply(
        [{ address: "0xINT", addressType: INTERNAL, addressID: "31" }, EXTERNAL_ROW],
        { currentPage: token, hasNext: true }
      ),
      [ADDRESSES_PATH]: addressesReply(),
    });

    const page = await assetService(transport).queryAssetAddresses("asset-1", { pageSize: 2 });

    expect(page.excludedUnverified).toHaveLength(1);
    expect(page.pagination).toEqual({ pageSize: 2, nextCursor: token, hasMore: true });
  });
});
