/**
 * Blockchain models for Taurus-PROTECT SDK.
 */

import type { Currency } from './currency';

/**
 * Blockchain-specific information for Polkadot (DOT) networks.
 */
export interface DOTBlockchainInfo {
  /** SS58 address format prefix */
  readonly ss58Format?: number;
}

/**
 * Blockchain-specific information for EVM-compatible networks.
 */
export interface EVMBlockchainInfo {
  /** EIP-155 chain ID */
  readonly chainId?: string;
}

/**
 * Blockchain-specific information for Tezos (XTZ) networks.
 */
export interface XTZBlockchainInfo {
  /** Protocol hash */
  readonly protocolHash?: string;
}

/**
 * Represents blockchain network information.
 *
 * Contains details about a supported blockchain including its symbol, name,
 * network type, chain ID, and various blockchain-specific configuration.
 */
export interface Blockchain {
  /** Blockchain symbol (e.g., "BTC", "ETH", "SOL") */
  readonly symbol?: string;
  /** Blockchain name (e.g., "Bitcoin", "Ethereum") */
  readonly name?: string;
  /** Network type (e.g., "mainnet", "testnet") */
  readonly network?: string;
  /** Chain ID (for EVM-compatible chains) */
  readonly chainId?: string;
  /** Number of confirmations required for transactions */
  readonly confirmations?: string;
  /** Current block height of the blockchain */
  readonly blockHeight?: string;
  /** Blackhole/burn address for this blockchain (used for cancel requests) */
  readonly blackholeAddress?: string;
  /** Whether this is a Layer 2 chain */
  readonly isLayer2Chain?: boolean;
  /** Underlying Layer 1 blockchain network (only relevant when isLayer2Chain is true) */
  readonly layer1Network?: string;
  /** Base currency information for this blockchain */
  readonly baseCurrency?: Currency;
  /** Polkadot-specific blockchain information */
  readonly dotInfo?: DOTBlockchainInfo;
  /** EVM-specific blockchain information */
  readonly ethInfo?: EVMBlockchainInfo;
  /** Tezos-specific blockchain information */
  readonly xtzInfo?: XTZBlockchainInfo;
}

/**
 * Options for listing blockchains.
 */
export interface ListBlockchainsOptions {
  /** Filter by blockchain symbol (e.g., "ETH", "BTC") */
  blockchain?: string;
  /** Filter by network type (e.g., "mainnet", "testnet") */
  network?: string;
  /** Include current block height in response (requires blockchain/network to be specified) */
  includeBlockHeight?: boolean;
}
