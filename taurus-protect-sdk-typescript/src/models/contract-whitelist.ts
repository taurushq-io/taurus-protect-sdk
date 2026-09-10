/**
 * Contract whitelist models for Taurus-PROTECT SDK.
 *
 * WRITE-side models plus attributes. The read-side models that lived here
 * (WhitelistedContract, SignedWhitelistedContractEnvelope, WhitelistedContractResult and
 * their nested types) described the same server entity as WhitelistedAsset and were
 * populated with no verification. The verified models are in `whitelisted-asset.ts`.
 */

/**
 * Represents an attribute on a whitelisted contract.
 *
 * Attributes are key-value metadata that can be attached to whitelisted contracts
 * for additional categorization or information.
 */
export interface WhitelistedContractAttribute {
  /** Attribute ID. */
  readonly id: string | undefined;
  /** Attribute key. */
  readonly key: string | undefined;
  /** Attribute value. */
  readonly value: string | undefined;
  /** Content type. */
  readonly contentType: string | undefined;
  /** Attribute type. */
  readonly type: string | undefined;
  /** Attribute subtype. */
  readonly subType: string | undefined;
  /** Whether this attribute is a file. */
  readonly isFile: boolean | undefined;
}

/**
 * Request to create a whitelisted contract.
 */
export interface CreateWhitelistedContractRequest {
  /** Blockchain identifier (e.g., "ETH", "MATIC", "XTZ"). */
  blockchain: string;
  /** Network identifier (e.g., "mainnet", "goerli"). */
  network: string;
  /** The smart contract address. */
  contractAddress?: string;
  /** Token symbol (e.g., "USDC", "WETH"). */
  symbol: string;
  /** Human-readable name. */
  name: string;
  /** Number of decimals (0 for NFTs). */
  decimals: number;
  /** Contract kind (e.g., "erc20", "erc721", "fa2"). */
  kind: string;
  /** Token ID for NFTs (null for fungible tokens). */
  tokenId?: string;
}

/**
 * Request to update a whitelisted contract.
 */
export interface UpdateWhitelistedContractRequest {
  /** New symbol. */
  symbol: string;
  /** New name. */
  name: string;
  /** New decimals value. */
  decimals: number;
}
