/**
 * Whitelist hash computation utilities.
 *
 * This module provides functions for computing hashes of whitelisted addresses,
 * including legacy hash computation for backward compatibility with addresses
 * signed before schema changes.
 *
 * ## SECURITY DESIGN
 *
 * The API returns metadata with two representations of the same data:
 *   - `payload`: Raw JSON object (UNVERIFIED)
 *   - `payloadAsString`: JSON string that is cryptographically hashed (VERIFIED)
 *
 * The security model works as follows:
 *   1. Server computes: `metadata.hash = SHA256(payloadAsString)`
 *   2. The hash is signed by governance rules (SuperAdmin keys)
 *   3. Clients verify: `computed_hash(payloadAsString) === metadata.hash`
 *
 * ### ATTACK VECTOR (if using raw payload)
 *
 * An attacker intercepting API responses could:
 *   1. Modify the payload object (e.g., change destination address)
 *   2. Leave payloadAsString unchanged (hash still verifies)
 *   3. Client extracts data from modified payload → SECURITY BYPASS
 *
 * ### SOLUTION
 *
 * All parsing functions in this module (e.g., `parseWhitelistedAddressFromJson`)
 * expect a JSON string parameter, which should be the verified `payloadAsString`.
 * This ensures:
 *   - All extracted data comes from the cryptographically verified source
 *   - Any tampering with the raw payload object is ignored
 *   - The integrity chain: payloadAsString → hash → signature is preserved
 */

import { createHash } from "crypto";
import { calculateHexHash } from "../crypto";
import { constantTimeCompare } from "./constant-time";
import { IntegrityError } from "../errors";
import type {
  WhitelistedAddress,
  InternalAddress,
  InternalWallet,
} from "../models/whitelisted-address";

/**
 * Regular expression to match the contractType field in JSON.
 * Matches: ,"contractType":"..." (value in quotes)
 */
const CONTRACT_TYPE_PATTERN = /,"contractType":"[^"]*"/g;

/**
 * Regular expression to match label fields inside objects.
 * Matches: ,"label":"..."} (only labels followed by closing brace)
 * This only matches labels inside linkedInternalAddresses objects,
 * not the main address label which is followed by other fields.
 */
const LABEL_IN_OBJECT_PATTERN = /,"label":"[^"]*"}/g;

/**
 * Asset-specific legacy hash patterns.
 * Used for WhitelistedAsset backward compatibility.
 */

/**
 * Matches "isNFT":(true|false) with leading comma.
 */
const IS_NFT_PATTERN_LEADING_COMMA = /,"isNFT":(true|false)/g;

/**
 * Matches "isNFT":(true|false) with trailing comma.
 */
const IS_NFT_PATTERN_TRAILING_COMMA = /"isNFT":(true|false),/g;

/**
 * Matches "kindType":"..." with leading comma.
 */
const KIND_TYPE_PATTERN_LEADING_COMMA = /,"kindType":"[^"]*"/g;

/**
 * Matches "kindType":"..." with trailing comma.
 */
const KIND_TYPE_PATTERN_TRAILING_COMMA = /"kindType":"[^"]*",/g;

/**
 * Computes legacy hashes for backward compatibility.
 *
 * When a whitelisted address was signed before certain schema changes,
 * the hash may have been computed without those fields. This function
 * generates alternative hashes by removing fields that may not have
 * existed at signing time.
 *
 * Strategies:
 * 1. Remove contractType field (addresses signed before contractType was added)
 * 2. Remove labels from linkedInternalAddresses (after contractType but before labels were added)
 * 3. Remove both contractType and labels (before both fields were added)
 *
 * @param payloadAsString - The JSON payload string
 * @returns Array of possible hashes (may be empty if no legacy formats apply)
 */
export function computeLegacyHashes(payloadAsString: string): string[] {
  if (!payloadAsString) {
    return [];
  }

  const seen = new Set<string>();
  const hashes: string[] = [];

  const addHash = (payload: string): void => {
    const hash = calculateHexHash(payload);
    if (!seen.has(hash)) {
      seen.add(hash);
      hashes.push(hash);
    }
  };

  // Strategy 1: Remove contractType only
  // Handles addresses signed before contractType was added to schema
  const withoutContractType = payloadAsString.replace(
    CONTRACT_TYPE_PATTERN,
    ""
  );
  if (withoutContractType !== payloadAsString) {
    addHash(withoutContractType);
  }

  // Strategy 2: Remove labels from linkedInternalAddresses objects only (keep contractType)
  // Handles addresses signed after contractType was added but before labels were added
  const withoutLabels = payloadAsString.replace(LABEL_IN_OBJECT_PATTERN, "}");
  if (withoutLabels !== payloadAsString) {
    addHash(withoutLabels);
  }

  // Strategy 3: Remove BOTH contractType AND labels from linkedInternalAddresses
  // Handles addresses signed before both fields were added
  let withoutBoth = payloadAsString.replace(LABEL_IN_OBJECT_PATTERN, "}");
  withoutBoth = withoutBoth.replace(CONTRACT_TYPE_PATTERN, "");
  if (withoutBoth !== payloadAsString) {
    addHash(withoutBoth);
  }

  return hashes;
}

/**
 * Computes asset-specific legacy hashes for backward compatibility.
 *
 * When a whitelisted asset was signed before certain schema changes,
 * the hash may have been computed without those fields. This function
 * generates alternative hashes by removing fields that may not have
 * existed at signing time.
 *
 * Strategies (aligned with Java SDK WhitelistedAssetService.computeLegacyHashes):
 * 1. Remove isNFT field (assets signed before isNFT was added)
 * 2. Remove kindType field (assets signed before kindType was added)
 * 3. Remove both isNFT and kindType (assets signed before both fields were added)
 *
 * @param payloadAsString - The JSON payload string
 * @returns Array of possible hashes (may be empty if no legacy formats apply)
 */
export function computeAssetLegacyHashes(payloadAsString: string): string[] {
  if (!payloadAsString) {
    return [];
  }

  const seen = new Set<string>();
  const hashes: string[] = [];

  const addHash = (payload: string): void => {
    const hash = calculateHexHash(payload);
    if (!seen.has(hash)) {
      seen.add(hash);
      hashes.push(hash);
    }
  };

  // Strategy 1: Remove isNFT only
  // Handles assets signed before isNFT was added to schema
  let withoutIsNFT = payloadAsString.replace(IS_NFT_PATTERN_LEADING_COMMA, "");
  withoutIsNFT = withoutIsNFT.replace(IS_NFT_PATTERN_TRAILING_COMMA, "");
  if (withoutIsNFT !== payloadAsString) {
    addHash(withoutIsNFT);
  }

  // Strategy 2: Remove kindType only
  // Handles assets signed before kindType was added to schema
  let withoutKindType = payloadAsString.replace(
    KIND_TYPE_PATTERN_LEADING_COMMA,
    ""
  );
  withoutKindType = withoutKindType.replace(KIND_TYPE_PATTERN_TRAILING_COMMA, "");
  if (withoutKindType !== payloadAsString) {
    addHash(withoutKindType);
  }

  // Strategy 3: Remove BOTH isNFT AND kindType
  // Handles assets signed before both fields were added
  // Note: Order matches Java implementation - remove isNFT first, then kindType
  let withoutBoth = payloadAsString.replace(IS_NFT_PATTERN_LEADING_COMMA, "");
  withoutBoth = withoutBoth.replace(IS_NFT_PATTERN_TRAILING_COMMA, "");
  withoutBoth = withoutBoth.replace(KIND_TYPE_PATTERN_LEADING_COMMA, "");
  withoutBoth = withoutBoth.replace(KIND_TYPE_PATTERN_TRAILING_COMMA, "");
  if (withoutBoth !== payloadAsString) {
    addHash(withoutBoth);
  }

  return hashes;
}

/**
 * JSON payload structure for whitelisted address.
 */
interface WhitelistPayload {
  currency?: string;
  network?: string;
  address?: string;
  memo?: string;
  label?: string;
  customerId?: string;
  contractType?: string;
  tnParticipantID?: string;
  addressType?: string;
  exchangeAccountId?: string;
  linkedInternalAddresses?: Array<{
    id?: number;
    address?: string;
    label?: string;
  }>;
  linkedWallets?: Array<{
    id?: number;
    name?: string;
    path?: string;
  }>;
}

/**
 * Parses a WhitelistedAddress from a verified JSON payload.
 *
 * This extracts the signed fields from the cryptographically verified payload.
 * The payload structure differs from the API response - it uses the signed
 * field names (e.g., "currency" instead of "blockchain").
 *
 * **SECURITY NOTE:** This function expects the verified `payloadAsString` from
 * metadata, NOT the raw `payload` object. The `payloadAsString` is the
 * cryptographically verified source (its hash is signed by governance rules).
 * Using the raw payload object would bypass integrity verification and could
 * allow an attacker to inject tampered data.
 *
 * @param jsonPayload - JSON string from metadata.payloadAsString (verified source)
 * @returns Parsed WhitelistedAddress
 * @throws Error if the payload cannot be parsed
 *
 * @example
 * ```typescript
 * // CORRECT: Use payloadAsString
 * const address = parseWhitelistedAddressFromJson(metadata.payloadAsString);
 *
 * // WRONG: Never use raw payload directly
 * // const address = parseWhitelistedAddressFromJson(JSON.stringify(metadata.payload));
 * ```
 */
export function parseWhitelistedAddressFromJson(
  jsonPayload: string
): WhitelistedAddress {
  if (!jsonPayload) {
    throw new Error("JSON payload cannot be empty");
  }

  let payload: WhitelistPayload;
  try {
    payload = JSON.parse(jsonPayload) as WhitelistPayload;
  } catch (error) {
    throw new Error(
      `Failed to parse whitelist payload: ${error instanceof Error ? error.message : "unknown error"}`
    );
  }

  // Parse linkedInternalAddresses
  const linkedInternalAddresses: InternalAddress[] = [];
  if (Array.isArray(payload.linkedInternalAddresses)) {
    for (const lia of payload.linkedInternalAddresses) {
      linkedInternalAddresses.push({
        id: typeof lia.id === "number" ? lia.id : 0,
        label: typeof lia.label === "string" ? lia.label : undefined,
      });
    }
  }

  // Parse linkedWallets
  const linkedWallets: InternalWallet[] = [];
  if (Array.isArray(payload.linkedWallets)) {
    for (const lw of payload.linkedWallets) {
      linkedWallets.push({
        id: typeof lw.id === "number" ? lw.id : 0,
        // Note: JSON field is "name" but model field is "label"
        label: typeof lw.name === "string" ? lw.name : undefined,
        path: typeof lw.path === "string" ? lw.path : undefined,
      });
    }
  }

  // Parse exchangeAccountId (string in JSON, number in model)
  let exchangeAccountId: number | undefined;
  if (payload.exchangeAccountId) {
    const parsed = parseInt(payload.exchangeAccountId, 10);
    if (!isNaN(parsed)) {
      exchangeAccountId = parsed;
    }
  }

  return {
    id: "", // ID is not in the signed payload
    blockchain: payload.currency ?? "",
    network: payload.network ?? "",
    address: payload.address ?? "",
    memo: payload.memo,
    label: payload.label,
    customerId: payload.customerId,
    contractType: payload.contractType,
    addressType: payload.addressType,
    tnParticipantId: payload.tnParticipantID,
    exchangeAccountId,
    linkedInternalAddresses,
    linkedWallets,
    createdAt: undefined,
    attributes: {}, // Attributes come from envelope, not payload
  };
}

/**
 * Verifies hash coverage - checks if a hash is covered by at least one signature.
 *
 * SECURITY: Uses constant-time comparison to prevent timing attacks.
 * The loop does NOT early return to maintain constant execution time
 * regardless of where in the list the match is found.
 *
 * @param metadataHash - The hash to find
 * @param signatures - List of signature entries with their hashes
 * @returns true if the hash is found in any signature's hashes list
 */
export function verifyHashCoverage(
  metadataHash: string,
  signatures: Array<{ hashes: string[] }>
): boolean {
  let found = false;
  for (const sig of signatures) {
    if (containsHash(sig.hashes, metadataHash)) {
      found = true;
      // No break: returning on a match would leak its position through timing.
    }
  }
  return found;
}

/**
 * Reports whether one signature's hashes list covers `hash`.
 *
 * The per-signature half of the pair; {@link verifyHashCoverage} asks the same
 * question across every signature. These two are the ONLY places this SDK compares
 * hash strings — each verifier used to carry its own private copy.
 *
 * SECURITY: constant-time, and the loop does not break early.
 *
 * @param hashes - the hashes a signature covers
 * @param hash - the hash to look for
 * @returns true when the list covers the hash
 */
export function containsHash(
  hashes: string[] | undefined,
  hash: string
): boolean {
  if (!hashes || !hash) {
    return false;
  }

  let found = false;
  for (const candidate of hashes) {
    if (constantTimeCompare(hash, candidate)) {
      found = true;
      // No break, as above.
    }
  }
  return found;
}

/**
 * Maximum size of a signed payload before it is parsed.
 *
 * The payload is hash-checked in step 1 but not AUTHENTICATED until step 5's
 * signatures verify, so anything parsed in between is still attacker-influenced.
 * A generous ceiling that only stops a hostile response consuming memory.
 */
export const MAX_PAYLOAD_BYTES = 1 << 20;

/**
 * Returns the (blockchain, network) pair that selects the governance rules, taken
 * from the SIGNED payload rather than the surrounding DTO.
 *
 * The DTO is free-floating: nothing binds it to the signatures, so a response that
 * set blockchain to empty would steer verification to the global-default rule tier,
 * which is broader than the rule the entity belongs to. The payload pair is at
 * least hash-bound to the material step 5 checks.
 *
 * Addresses name the chain `currency`; assets name it `blockchain`. Both are
 * accepted so one helper serves both flows.
 *
 * The chain is always in the payload. The NETWORK is not: governance rules carry a
 * per-rule `includeNetworkInPayload` flag, and real signed payloads omit `network`
 * when it is off. Requiring it would reject correctly-signed addresses, so the DTO
 * network is used only when the payload has none.
 *
 * @param payloadAsString - the signed payload
 * @param dtoBlockchain - the blockchain the response claims
 * @param dtoNetwork - the network the response claims
 * @returns the pair to look rules up with
 * @throws {@link IntegrityError} if the payload omits the chain — never treated as
 *   a wildcard, because an empty value matches every tier — or if the payload and
 *   the DTO disagree on a field the payload does carry
 */
export function resolveRuleKey(
  payloadAsString: string | undefined,
  dtoBlockchain: string | undefined,
  dtoNetwork: string | undefined
): { blockchain: string; network: string } {
  if (!payloadAsString) {
    throw new IntegrityError("cannot resolve governance rule key: payload is empty");
  }
  if (payloadAsString.length > MAX_PAYLOAD_BYTES) {
    throw new IntegrityError(
      `cannot resolve governance rule key: payload exceeds ${MAX_PAYLOAD_BYTES} bytes`
    );
  }

  let payload: Record<string, unknown>;
  try {
    const parsed: unknown = JSON.parse(payloadAsString);
    if (typeof parsed !== "object" || parsed === null || Array.isArray(parsed)) {
      throw new IntegrityError(
        "cannot resolve governance rule key: payload is not an object"
      );
    }
    payload = parsed as Record<string, unknown>;
  } catch (error) {
    if (error instanceof IntegrityError) {
      throw error;
    }
    throw new IntegrityError(
      `cannot resolve governance rule key: ${error instanceof Error ? error.message : "unknown error"}`
    );
  }

  const blockchain = String(payload.blockchain ?? payload.currency ?? "");
  let network = String(payload.network ?? "");

  if (!blockchain) {
    throw new IntegrityError(
      "signed payload does not carry a blockchain; " +
        "refusing to fall back to a wildcard rule"
    );
  }

  // Only reachable when includeNetworkInPayload is off for this rule.
  const networkFromPayload = network !== "";
  if (!networkFromPayload) {
    network = dtoNetwork ?? "";
  }

  if (dtoBlockchain && dtoBlockchain.toLowerCase() !== blockchain.toLowerCase()) {
    throw new IntegrityError(
      `blockchain disagrees between the signed payload (${blockchain}) and the response (${dtoBlockchain})`
    );
  }
  if (networkFromPayload && dtoNetwork && dtoNetwork.toLowerCase() !== network.toLowerCase()) {
    throw new IntegrityError(
      `network disagrees between the signed payload (${network}) and the response (${dtoNetwork})`
    );
  }

  return { blockchain, network };
}

/**
 * Recomputes the label validatord files a normalized rules container under, so a row's
 * `rulesContainerHash` resolves only to bytes that really hash to it.
 *
 * The convention is validatord's, at `internal/api/v1/whitelist-controller.go`:
 *
 * ```go
 * h := base64.StdEncoding.EncodeToString(crypto.Sha256([]byte(e.GetRulesContainer())))
 * ```
 *
 * Note what is hashed. `GetRulesContainer()` is ALREADY a base64 string there, so the
 * digest is over the base64 TEXT, not over the decoded protobuf, and the output is
 * base64 rather than hex. Decoding the container first and hashing the protobuf yields a
 * different value and would reject every container, breaking all list calls.
 *
 * Do not confuse it with `enforcedRulesHash`, which is `base64(SHA256(raw protobuf))` and
 * is a backlink to a ruleset's predecessor rather than its own identity. Both are 44-char
 * base64 SHA-256 digests, so mixing them up is silent.
 */
export function containerHashLabel(containerBase64: string): string {
  return createHash("sha256").update(containerBase64, "utf-8").digest("base64");
}
