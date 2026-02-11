/**
 * Air gap models for Taurus-PROTECT SDK.
 *
 * These models represent the data structures used for air-gapped signing
 * workflows with cold HSM (Hardware Security Module) devices.
 */

/**
 * Options for exporting requests for cold HSM signing.
 */
export interface GetOutgoingAirGapOptions {
  /**
   * List of request IDs to export for cold HSM signing.
   * These requests must be in HSM-ready state.
   */
  requestIds: string[];

  /**
   * Optional ECDSA signature of the request hashes.
   * Format: base64(ecdsa_sign(sha256([hex(sha256(req1)),hex(sha256(req2)),...]))
   */
  signature?: string;
}

/**
 * Options for exporting addresses for cold HSM signing.
 */
export interface GetOutgoingAirGapAddressOptions {
  /**
   * List of address IDs to export for cold HSM signing.
   */
  addressIds: string[];
}

/**
 * Options for submitting signed responses from the cold HSM.
 */
export interface SubmitIncomingAirGapOptions {
  /**
   * The signed payload from the cold HSM (base64-encoded).
   * This is the IncomingAirGapFile_SignedPayload from the HSM.
   */
  payload: string;

  /**
   * Optional signature of the air-gap importer (base64-encoded).
   */
  signature?: string;
}
