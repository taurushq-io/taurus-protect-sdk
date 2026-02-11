/**
 * Test fixtures for whitelisted address unit tests.
 *
 * Contains real API response data and legacy hash test cases for backward
 * compatibility testing. These fixtures mirror the Java SDK test cases
 * to ensure cross-SDK alignment.
 *
 * @see Java SDK: WhitelistIntegrityHelperTest.java
 * @see Python SDK: test_whitelist_helper.py
 * @see Go SDK: whitelist_helper_test.go
 */

// ============================================================================
// REAL API RESPONSE DATA
// ============================================================================

/**
 * Real address captured from API integration test.
 *
 * ID: 36663
 * Blockchain: ALGO (Algorand)
 * Type: individual
 *
 * This payload represents an actual address returned from the Taurus-PROTECT API.
 * It includes all standard fields plus the tnParticipantID for Taurus Network
 * participant association.
 */
export const REAL_PAYLOAD_AS_STRING =
  '{"currency":"ALGO","addressType":"individual","address":"P4QCJV2YYLAEULGLJQAW4XTU3EBOHWL5C46I5SPLH2H7AJEE367ZDACV5A","memo":"","label":"TN_Bank ACC Cockroach_WTRTest","customerId":"","exchangeAccountId":"","linkedInternalAddresses":[],"contractType":"","tnParticipantID":"84dc35e3-0af8-4b6b-be75-785f4b149d16"}';

/**
 * SHA-256 hash of REAL_PAYLOAD_AS_STRING.
 *
 * This is the hash that would be signed by SuperAdmin keys and stored
 * in the envelope's signatures array.
 */
export const REAL_METADATA_HASH =
  "830063cfa8c1dbd696d670fc8360e85fbc57c3ffa66d22358b9a7d6befabb2f0";

// ============================================================================
// LEGACY HASH TEST CASES
// ============================================================================

/**
 * Case 1: contractType field added after original signing.
 *
 * This test case represents addresses that were signed before the contractType
 * field was added to the schema. When these addresses are retrieved, they now
 * include an empty contractType field (""), but the signature was computed
 * without this field.
 *
 * Verification must:
 * 1. First try matching the current payload hash
 * 2. If that fails, compute legacy hash (without contractType) and retry
 *
 * Schema evolution timeline:
 * - Original: {"currency", "addressType", "address", "memo", "label", ...}
 * - Current:  {"currency", "addressType", "address", "memo", "label", ..., "contractType"}
 */
export const CASE1_LEGACY_PAYLOAD =
  '{"currency":"ETH","addressType":"individual","address":"0x012566A179a935ACF1d81d4D237495DE933D12E6","memo":"","label":"CMTA20-KYC - 0x012566A179a935ACF1d81d4D237495DE933D12E6 (request 6826)","customerId":"","exchangeAccountId":"","linkedInternalAddresses":[]}';

export const CASE1_CURRENT_PAYLOAD =
  '{"currency":"ETH","addressType":"individual","address":"0x012566A179a935ACF1d81d4D237495DE933D12E6","memo":"","label":"CMTA20-KYC - 0x012566A179a935ACF1d81d4D237495DE933D12E6 (request 6826)","customerId":"","exchangeAccountId":"","linkedInternalAddresses":[],"contractType":""}';

/**
 * Hash of CASE1_LEGACY_PAYLOAD (without contractType).
 * This is what the original signature covers.
 */
export const CASE1_LEGACY_HASH =
  "cda66e821ec26f2432a717feaa1ef49be39a7ad9e93b6b8fcdce606659e964df";

/**
 * Hash of CASE1_CURRENT_PAYLOAD (with empty contractType).
 * This hash will NOT match the signature because contractType was added later.
 */
export const CASE1_CURRENT_HASH =
  "d95ae4359bea509c2542acf410649f1e361233da5e1ac7c7a198b6d6a2bbbe1f";

/**
 * Case 2: Both contractType AND labels in linkedInternalAddresses added after signing.
 *
 * This is the most complex legacy case. The original signing covered a payload that:
 * 1. Did not have the contractType field
 * 2. Did not have labels in linkedInternalAddresses objects
 *
 * The linkedInternalAddresses array originally contained only {id, address} objects.
 * Later, {label} was added to each linked address, and {contractType} was added
 * to the root payload.
 *
 * Legacy verification must try multiple strategies:
 * - Strategy 1: Remove only contractType
 * - Strategy 2: Remove only labels from linkedInternalAddresses
 * - Strategy 3: Remove both contractType AND labels
 *
 * Schema evolution:
 * - Original: {"...", "linkedInternalAddresses":[{"id":"10","address":"0x..."}]}
 * - Current:  {"...", "linkedInternalAddresses":[{"id":"10","address":"0x...","label":"..."}],"contractType":""}
 */
export const CASE2_CURRENT_PAYLOAD =
  '{"currency":"ETH","addressType":"individual","address":"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380","memo":"","label":"20200324 test address 2","customerId":"1","exchangeAccountId":"","linkedInternalAddresses":[{"id":"10","address":"0x589ef3d7585f54f0539e24253050887c691c9bd8","label":"client 0 ETH "},{"id":"13","address":"0x669805f31178faf0dca39c8a5c49ecc531b5156e","label":"ETH internal client 02.02"},{"id":"20","address":"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d","label":"LBR 07.02"},{"id":"98","address":"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e","label":"ETH LBR internal client 26.02"},{"id":"25","address":"0x9bc28e6710f5bb2511372987f613a436618e28ad","label":"LBR IC 13.02"}],"contractType":""}';

/**
 * Payload with BOTH contractType AND labels removed.
 * This is what the original signature covers (Strategy 3 - full legacy).
 */
export const CASE2_LEGACY_PAYLOAD =
  '{"currency":"ETH","addressType":"individual","address":"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380","memo":"","label":"20200324 test address 2","customerId":"1","exchangeAccountId":"","linkedInternalAddresses":[{"id":"10","address":"0x589ef3d7585f54f0539e24253050887c691c9bd8"},{"id":"13","address":"0x669805f31178faf0dca39c8a5c49ecc531b5156e"},{"id":"20","address":"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d"},{"id":"98","address":"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e"},{"id":"25","address":"0x9bc28e6710f5bb2511372987f613a436618e28ad"}]}';

/**
 * Hash of CASE2_LEGACY_PAYLOAD (without both contractType AND labels).
 * This is what the original signature covers.
 */
export const CASE2_LEGACY_HASH =
  "88e4e456f7ca1fc4ca415c6c571f828c0eb047e9f15f36d547c103b2ea0def9b";

/**
 * Hash of CASE2_CURRENT_PAYLOAD (with both contractType AND labels).
 * This hash will NOT match the signature.
 */
export const CASE2_CURRENT_HASH =
  "7d62d7f78ed55c716ea1278473d6cac5b60a31e1df941873118932822df32b03";

/**
 * Strategy 2 legacy payload: labels removed, contractType kept.
 *
 * This intermediate payload tests the case where:
 * - contractType is present (empty string)
 * - labels in linkedInternalAddresses are removed
 *
 * This is used to test that the legacy hash computation correctly
 * handles both fields independently.
 */
export const STRATEGY2_LEGACY_PAYLOAD =
  '{"currency":"ETH","addressType":"individual","address":"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380","memo":"","label":"20200324 test address 2","customerId":"1","exchangeAccountId":"","linkedInternalAddresses":[{"id":"10","address":"0x589ef3d7585f54f0539e24253050887c691c9bd8"},{"id":"13","address":"0x669805f31178faf0dca39c8a5c49ecc531b5156e"},{"id":"20","address":"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d"},{"id":"98","address":"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e"},{"id":"25","address":"0x9bc28e6710f5bb2511372987f613a436618e28ad"}],"contractType":""}';

// ============================================================================
// HELPER TYPES
// ============================================================================

/**
 * Represents a legacy hash test case for parameterized testing.
 */
export interface LegacyHashTestCase {
  /** Human-readable description of the test case */
  description: string;
  /** The current payload as returned by the API */
  currentPayload: string;
  /** The legacy payload (what was originally signed) */
  legacyPayload: string;
  /** Hash of the current payload */
  currentHash: string;
  /** Hash of the legacy payload (should match signature) */
  legacyHash: string;
}

/**
 * All legacy hash test cases for parameterized testing.
 *
 * Usage:
 * ```typescript
 * import { LEGACY_HASH_TEST_CASES } from './fixtures/whitelisted-address-fixtures';
 *
 * describe('legacy hash computation', () => {
 *   it.each(LEGACY_HASH_TEST_CASES)('$description', ({ currentPayload, legacyHash }) => {
 *     const legacyHashes = computeLegacyHashes(currentPayload);
 *     expect(legacyHashes).toContain(legacyHash);
 *   });
 * });
 * ```
 */
export const LEGACY_HASH_TEST_CASES: LegacyHashTestCase[] = [
  {
    description: "Case 1: contractType added after signing",
    currentPayload: CASE1_CURRENT_PAYLOAD,
    legacyPayload: CASE1_LEGACY_PAYLOAD,
    currentHash: CASE1_CURRENT_HASH,
    legacyHash: CASE1_LEGACY_HASH,
  },
  {
    description: "Case 2: contractType AND labels added after signing",
    currentPayload: CASE2_CURRENT_PAYLOAD,
    legacyPayload: CASE2_LEGACY_PAYLOAD,
    currentHash: CASE2_CURRENT_HASH,
    legacyHash: CASE2_LEGACY_HASH,
  },
];

// ============================================================================
// MOCK ENVELOPE DATA
// ============================================================================

/**
 * Creates a mock whitelisted address envelope for testing.
 *
 * @param overrides - Partial envelope properties to override defaults
 * @returns A mock envelope object compatible with API response format
 */
export function createMockEnvelope(overrides?: {
  payloadAsString?: string;
  hash?: string;
  signatures?: Array<{ hashes: string[] }>;
  trails?: Array<{ action: string; date: Date }>;
  id?: string;
  status?: string;
}): {
  payloadAsString: string;
  hash: string;
  signatures: Array<{ hashes: string[] }>;
  trails: Array<{ action: string; date: Date }>;
  id: string;
  status: string;
} {
  return {
    payloadAsString: overrides?.payloadAsString ?? REAL_PAYLOAD_AS_STRING,
    hash: overrides?.hash ?? REAL_METADATA_HASH,
    signatures: overrides?.signatures ?? [{ hashes: [REAL_METADATA_HASH] }],
    trails: overrides?.trails ?? [
      { action: "created", date: new Date("2024-01-15T10:30:00Z") },
    ],
    id: overrides?.id ?? "36663",
    status: overrides?.status ?? "active",
  };
}

/**
 * Creates a mock envelope that requires legacy hash verification.
 *
 * The hash field contains the CURRENT hash, but the signatures only contain
 * the LEGACY hash. This simulates an address that was signed before schema changes.
 *
 * @param testCase - Which legacy hash test case to use
 * @returns A mock envelope that requires legacy hash fallback
 */
export function createLegacyMockEnvelope(
  testCase: "case1" | "case2" = "case1"
): {
  payloadAsString: string;
  hash: string;
  signatures: Array<{ hashes: string[] }>;
  trails: Array<{ action: string; date: Date }>;
  id: string;
  status: string;
} {
  const cases = {
    case1: {
      payload: CASE1_CURRENT_PAYLOAD,
      currentHash: CASE1_CURRENT_HASH,
      legacyHash: CASE1_LEGACY_HASH,
    },
    case2: {
      payload: CASE2_CURRENT_PAYLOAD,
      currentHash: CASE2_CURRENT_HASH,
      legacyHash: CASE2_LEGACY_HASH,
    },
  };

  const c = cases[testCase];
  return {
    payloadAsString: c.payload,
    hash: c.currentHash,
    signatures: [{ hashes: [c.legacyHash] }], // Only legacy hash in signatures
    trails: [{ action: "created", date: new Date("2020-03-24T00:00:00Z") }],
    id: testCase === "case1" ? "6826" : "12345",
    status: "active",
  };
}
