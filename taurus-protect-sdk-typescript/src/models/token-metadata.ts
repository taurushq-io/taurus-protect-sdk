/**
 * Token metadata models for Taurus-PROTECT SDK.
 *
 * These models represent token metadata for various token standards
 * including ERC-20, ERC-721, ERC-1155 (Ethereum), FA tokens (Tezos), and CryptoPunks.
 */

/**
 * Token metadata for ERC-20/ERC-721/ERC-1155 or FA tokens.
 *
 * Contains details about a token including name, description, decimals,
 * and optional media data for NFTs.
 */
export interface TokenMetadata {
  /** The token name */
  name?: string;
  /** The token symbol (primarily for FA tokens) */
  symbol?: string;
  /** The token description */
  description?: string;
  /** The number of decimals for the token */
  decimals?: string;
  /** The data type for NFT media (e.g., "image/png") */
  dataType?: string;
  /** Base64 encoded data for NFT media */
  base64Data?: string;
  /** The token metadata URI */
  uri?: string;
}

/**
 * CryptoPunk metadata.
 *
 * Contains details specific to CryptoPunk NFTs including
 * punk ID, attributes, and image data.
 */
export interface CryptoPunkMetadata {
  /** The punk ID (0-9999) */
  punkId?: string;
  /** The punk attributes (e.g., "Male 2, Earring, Bandana, ...") */
  punkAttributes?: string;
  /** Base64 encoded image data */
  image?: string;
}

/**
 * Options for getting ERC token metadata.
 */
export interface GetERCTokenMetadataOptions {
  /** The network (e.g., "mainnet", "goerli") */
  network: string;
  /** The contract address */
  contract: string;
  /** The token ID (required for ERC-721/1155, optional for ERC-20) */
  tokenId?: string;
  /** Whether to include base64 data (for NFTs) */
  withData?: boolean;
  /** The blockchain symbol (e.g., "ETH") */
  blockchain?: string;
}

/**
 * Options for getting EVM ERC token metadata.
 */
export interface GetEVMERCTokenMetadataOptions {
  /** The network (e.g., "mainnet") */
  network: string;
  /** The contract address */
  contract: string;
  /** The token ID (required for ERC-721/1155, optional for ERC-20) */
  tokenId?: string;
  /** Whether to include base64 data (for NFTs) */
  withData?: boolean;
  /** The blockchain symbol (e.g., "MATIC", "AVAX") - required */
  blockchain: string;
}

/**
 * Options for getting FA token metadata.
 */
export interface GetFATokenMetadataOptions {
  /** The network (e.g., "mainnet") */
  network: string;
  /** The contract address */
  contract: string;
  /** The token ID (must be "0" for FA1.2, any existing token for FA2) */
  tokenId?: string;
  /** Whether to include base64 data */
  withData?: boolean;
}

/**
 * Options for getting CryptoPunk metadata.
 */
export interface GetCryptoPunkMetadataOptions {
  /** The network (e.g., "mainnet") */
  network: string;
  /** The contract address (should be 0xb47e3cd837dDF8e4c57F05d70Ab865de6e193BBB for ETH/mainnet) */
  contract: string;
  /** The punk ID (0-9999) */
  punkId: string;
  /** The blockchain symbol */
  blockchain?: string;
}
