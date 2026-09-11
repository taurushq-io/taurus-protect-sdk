/**
 * Whitelisted address models for Taurus-PROTECT SDK.
 *
 * This module provides domain models for whitelisted addresses and
 * their cryptographic verification envelopes.
 */

import type { Verified } from "../helpers/verified";
import type { DecodedRulesContainer } from "./governance-rules";

/**
 * An internal address linked to a whitelisted address.
 */
export interface InternalAddress {
  /** Unique identifier for the internal address. */
  readonly id: number;
  /** Human-readable label for the address. */
  readonly label: string | undefined;
}

/**
 * An internal wallet linked to a whitelisted address.
 */
export interface InternalWallet {
  /** Unique identifier for the internal wallet. */
  readonly id: number;
  /** Wallet path (e.g., "BTC/wallet1"). */
  readonly path: string | undefined;
  /** Human-readable label for the wallet. */
  readonly label: string | undefined;
}

/**
 * A whitelisted external address.
 *
 * Whitelisted addresses are pre-approved destinations for withdrawals.
 * They must be verified with cryptographic signatures before use.
 */
export interface WhitelistedAddress {
  /** Unique identifier for the whitelisted address. */
  readonly id: string;
  /** Blockchain address string. */
  readonly address: string;
  /** Blockchain/currency identifier (e.g., "ETH", "BTC"). */
  readonly blockchain: string;
  /** Network type (e.g., "mainnet", "testnet"). */
  readonly network: string;
  /** Human-readable label for the address. */
  readonly label: string | undefined;
  /** Optional memo field for Stellar or destination tag for Ripple. */
  readonly memo: string | undefined;
  /** Customer ID for external reconciliation. */
  readonly customerId: string | undefined;
  /** Smart contract type (e.g., "CMTA20", "ERC20") for contract addresses. */
  readonly contractType: string | undefined;
  /** Address type (individual, exchange, contract, etc.). */
  readonly addressType: string | undefined;
  /** Taurus Network participant ID. */
  readonly tnParticipantId: string | undefined;
  /** Exchange account ID when the address belongs to an exchange. */
  readonly exchangeAccountId: number | undefined;
  /** Linked internal addresses. */
  readonly linkedInternalAddresses: InternalAddress[];
  /** Linked internal wallets that can send to this whitelisted address. */
  readonly linkedWallets: InternalWallet[];
  /** Creation timestamp. */
  readonly createdAt: Date | undefined;
  /** Custom attributes. */
  readonly attributes: Record<string, unknown>;
}

/**
 * Metadata containing the hash and signed payload for verification.
 */
export interface WhitelistMetadata {
  /** SHA-256 hash of the payloadAsString. */
  readonly hash: string;
  /** JSON string of the signed payload. */
  readonly payloadAsString: string;
}

/**
 * User signature details for a whitelist entry.
 */
export interface WhitelistUserSignature {
  /** ID of the signing user. */
  readonly userId: string | undefined;
  /** Base64-encoded cryptographic signature. */
  readonly signature: string | undefined;
  /** Optional comment from the signer. */
  readonly comment: string | undefined;
}

/**
 * A signature entry with the hashes it covers.
 */
export interface WhitelistSignatureEntry {
  /** User signature details. */
  readonly userSignature: WhitelistUserSignature | undefined;
  /** Hashes covered by this signature. */
  readonly hashes: string[];
}

/**
 * Signed whitelisted address data containing signatures.
 */
export interface SignedWhitelistedAddress {
  /** Base64-encoded signed payload. */
  readonly payload: string | undefined;
  /** List of signatures on this address. */
  readonly signatures: WhitelistSignatureEntry[];
}

/**
 * Envelope containing a whitelisted address with all data needed for verification.
 *
 * This envelope contains:
 * - Metadata with hash and payload for step 1 (hash verification)
 * - Rules container and signatures for steps 2-3 (SuperAdmin verification)
 * - Signed address with user signatures for steps 4-6 (hash coverage and threshold verification)
 */
export interface SignedWhitelistedAddressEnvelope {
  /** Unique identifier of the whitelisted address. */
  readonly id: string;
  /** Metadata with hash and payload. */
  readonly metadata: WhitelistMetadata;
  /** Base64-encoded rules container. */
  readonly rulesContainerBase64: string;
  /** Base64-encoded rules signatures. */
  readonly rulesSignaturesBase64: string;
  /** Signed address data with user signatures. */
  readonly signedAddress: SignedWhitelistedAddress;
  /** Blockchain identifier. */
  readonly blockchain: string;
  /** Network identifier. */
  readonly network: string;
  /** Linked internal addresses for rule line matching. */
  readonly linkedInternalAddresses: InternalAddress[];
  /** Linked wallets for rule line matching. */
  readonly linkedWallets: InternalWallet[];
  /** Hash of the rules container (for normalized caching). */
  readonly rulesContainerHash?: string;
}

/**
 * Result of whitelisted address verification.
 */
export interface WhitelistedAddressVerificationResult {
  /** The verified whitelisted address parsed from the payload. */
  readonly verifiedWhitelistedAddress: WhitelistedAddress;
  /** The decoded rules container used during verification. */
  readonly verifiedRulesContainer?: DecodedRulesContainer;
  /** The hash that was verified (may be a legacy hash). */
  readonly verifiedHash: string;
  /**
   * The payload {@link verifiedHash} covers, and the bytes
   * {@link verifiedWhitelistedAddress} was parsed from.
   *
   * When a legacy variant matched, this is that variant rather than
   * `metadata.payloadAsString` — the delivered payload carries members no signature
   * covered. `metadata.payloadAsString` is deliberately left untouched, because a caller
   * needs it to reproduce `metadata.hash`; read this field instead when the question is
   * "what was actually signed".
   */
  readonly verifiedPayload: string;
  /**
   * The envelope, marked as having passed verification.
   *
   * Anything that reads raw envelope fields takes this type rather than the bare
   * envelope, so a read path that skipped `verify()` will not compile. The asset side
   * has carried this since the 2026-09-04 pass; the address side returned the
   * unverified INPUT envelope from `getEnvelope` and discarded the verification result
   * entirely.
   */
  readonly verifiedEnvelope: Verified<SignedWhitelistedAddressEnvelope>;
}

/**
 * The set of rows an approver reviewed, carrying the metadata hash each one had AT
 * REVIEW TIME. It is the content pin the approval path signs against.
 *
 * Why this type exists rather than a plain `string[]` of ids. The approval API accepts
 * only ids: the SDK re-reads them and signs whatever the server returns under those ids.
 * Nothing bound the approver's intent to the bytes signed, so a response-controlling
 * server could answer the id-filtered re-read with a different row — one whose existing
 * signatures already satisfy the container it presents — and harvest a genuine approver
 * signature over content the approver never saw. This is the same shape
 * `GovernanceRuleService.approveRulesProposal` was hardened against with its mandatory
 * `expectedContainerHash`.
 *
 * The pin lives on the RESULT rather than on the row, because this SDK's
 * `WhitelistedAddress` deliberately carries no `metadata` — the hash it would pin is on
 * the envelope, and the verifier's result, not the model. `list()` and
 * `listForApproval()` therefore capture it while they still hold the DTO.
 *
 *	verified read ──▶ result.select(ids) ──▶ WhitelistedAddressApproval
 *	                                                 │
 *	                       approve(selection, key, comment) ◀┘
 *	                                                 │
 *	                          re-read ──▶ hash == pinned ? sign : abort
 *
 * `#pinned` is a true private field, which is the closest TypeScript comes to Go's
 * unexported-field idiom: a hand-built object literal is not an instance, and an
 * instance built with an empty map pins nothing and is refused by `approve`. What it
 * proves is that the approver went through a verified read — not that the read was
 * against the right keys. It is a structural guard against forgetting to pin.
 */
export class WhitelistedAddressApproval {
  /** Row id -> the metadata hash that row carried when it was reviewed. */
  readonly #pinned: ReadonlyMap<string, string>;

  /**
   * Minted by a verified read. Prefer `result.select(...)` / `result.selectAll()` on a
   * {@link WhitelistedAddress} listing over calling this directly — those are the
   * paths that guarantee the hashes came from rows this SDK verified.
   *
   * @param pinned - row id -> reviewed metadata hash
   */
  constructor(pinned: ReadonlyMap<string, string>) {
    this.#pinned = new Map(pinned);
  }

  /**
   * The pinned row ids, in no particular order.
   *
   * The approval path sorts them itself, because the endpoint requires ascending order
   * and the signed array must not depend on the order the caller happened to select in.
   */
  ids(): string[] {
    return [...this.#pinned.keys()];
  }

  /**
   * The reviewed metadata hash for `id`, or `undefined` when that id was not pinned.
   *
   * @param id - the row id
   */
  pinnedHash(id: string): string | undefined {
    return this.#pinned.get(id);
  }

  /** True when this selection pins nothing. */
  isEmpty(): boolean {
    return this.#pinned.size === 0;
  }
}

/**
 * Creates an empty WhitelistedAddress with required fields.
 *
 * @returns An empty WhitelistedAddress object
 */
export function createEmptyWhitelistedAddress(): WhitelistedAddress {
  return {
    id: "",
    address: "",
    blockchain: "",
    network: "",
    label: undefined,
    memo: undefined,
    customerId: undefined,
    contractType: undefined,
    addressType: undefined,
    tnParticipantId: undefined,
    exchangeAccountId: undefined,
    linkedInternalAddresses: [],
    linkedWallets: [],
    createdAt: undefined,
    attributes: {},
  };
}
