/**
 * Fee payer models for Taurus-PROTECT SDK.
 *
 * Fee payers are accounts used to pay transaction fees on behalf of other
 * addresses, commonly used for sponsored transactions on EVM-compatible
 * blockchains like Ethereum.
 */

/**
 * Represents a fee payer configuration in the Taurus-PROTECT system.
 *
 * A fee payer is an account that pays transaction fees on behalf of other
 * addresses, enabling sponsored transactions on blockchains like Ethereum.
 */
export interface FeePayer {
  /** Unique fee payer identifier */
  readonly id?: string;
  /** Tenant identifier */
  readonly tenantId?: string;
  /** Blockchain type (e.g., "ETH") */
  readonly blockchain?: string;
  /** Network (e.g., "mainnet", "goerli") */
  readonly network?: string;
  /** Fee payer name */
  readonly name?: string;
  /** Creation date */
  readonly creationDate?: Date;
  /** Blockchain-specific fee payer information */
  readonly feePayerInfo?: FeePayerInfo;
}

/**
 * Blockchain-specific information for a fee payer.
 */
export interface FeePayerInfo {
  /** Blockchain type */
  readonly blockchain?: string;
  /** Ethereum-specific fee payer configuration */
  readonly eth?: FeePayerEth;
}

/**
 * Ethereum-specific fee payer configuration.
 */
export interface FeePayerEth {
  /** Kind of fee payer (e.g., "local", "remote") */
  readonly kind?: string;
  /** Local fee payer configuration */
  readonly local?: FeePayerEthLocal;
  /** Remote fee payer configuration */
  readonly remote?: FeePayerEthRemote;
  /** Encrypted remote configuration */
  readonly remoteEncrypted?: string;
}

/**
 * Local (internally managed) Ethereum fee payer configuration.
 */
export interface FeePayerEthLocal {
  /** Address ID of the fee payer account */
  readonly addressId?: string;
  /** Address ID of the forwarder contract */
  readonly forwarderAddressId?: string;
  /** Whether requests are auto-approved */
  readonly autoApprove?: boolean;
  /** Address ID of the creator account */
  readonly creatorAddressId?: string;
  /** Kind of forwarder (e.g., "OpenZeppelinForwarder") */
  readonly forwarderKind?: string;
  /** Domain separator for EIP-712 signatures */
  readonly domainSeparator?: string;
}

/**
 * Remote (externally managed) Ethereum fee payer configuration.
 */
export interface FeePayerEthRemote {
  /** URL of the remote fee payer service */
  readonly url?: string;
  /** Username for authentication */
  readonly username?: string;
  /** Address ID of the sender account */
  readonly fromAddressId?: string;
  /** Address of the forwarder contract */
  readonly forwarderAddress?: string;
  /** Address ID of the forwarder contract */
  readonly forwarderAddressId?: string;
  /** Address of the creator account */
  readonly creatorAddress?: string;
  /** Address ID of the creator account */
  readonly creatorAddressId?: string;
  /** Kind of forwarder (e.g., "OpenZeppelinForwarder") */
  readonly forwarderKind?: string;
  /** Domain separator for EIP-712 signatures */
  readonly domainSeparator?: string;
}

/**
 * Options for listing fee payers.
 */
export interface ListFeePayersOptions {
  /** Maximum number of results to return */
  limit?: number;
  /** Number of results to skip for pagination */
  offset?: number;
  /** List of specific IDs to filter by */
  ids?: string[];
  /** Blockchain to filter by */
  blockchain?: string;
  /** Network to filter by */
  network?: string;
}
