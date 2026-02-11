/**
 * User device pairing models for Taurus-PROTECT SDK.
 *
 * Device pairing is used to register and authenticate user devices
 * for secure access to the Taurus-PROTECT platform.
 */

/**
 * Status of a user device pairing request.
 */
export const UserDevicePairingStatus = {
  /** Waiting for pairing to start */
  Waiting: 'WAITING',
  /** Pairing process is in progress */
  Pairing: 'PAIRING',
  /** Pairing has been approved */
  Approved: 'APPROVED',
} as const;

export type UserDevicePairingStatus =
  (typeof UserDevicePairingStatus)[keyof typeof UserDevicePairingStatus];

/**
 * Represents a user device pairing request.
 *
 * A pairing request is created to initiate the device pairing process.
 * The pairing ID is used to track the pairing through its lifecycle.
 */
export interface UserDevicePairing {
  /** Unique identifier for the pairing request */
  readonly pairingId: string;
}

/**
 * Represents the status and information of a user device pairing.
 *
 * Contains the current status of the pairing process and, upon successful
 * completion, the API key for the paired device.
 */
export interface UserDevicePairingInfo {
  /** Unique identifier for the pairing request */
  readonly pairingId: string;
  /** Current status of the pairing process */
  readonly status: UserDevicePairingStatus;
  /** API key for the paired device (only available after approval) */
  readonly apiKey?: string;
}

/**
 * Options for starting a device pairing.
 */
export interface StartPairingOptions {
  /** The nonce for verification (6 digits) */
  nonce: string;
  /** The public key for the device */
  publicKey: string;
}

/**
 * Options for approving a device pairing.
 */
export interface ApprovePairingOptions {
  /** The nonce for verification */
  nonce: string;
}
