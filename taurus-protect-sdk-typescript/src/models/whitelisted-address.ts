/**
 * Whitelisted address models for Taurus-PROTECT SDK.
 *
 * This module provides domain models for whitelisted addresses and
 * their cryptographic verification envelopes.
 */

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
