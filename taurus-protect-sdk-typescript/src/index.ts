/**
 * Taurus-PROTECT TypeScript SDK
 *
 * This SDK provides a TypeScript client for interacting with the Taurus-PROTECT API.
 *
 * @example
 * ```typescript
 * import { ProtectClient } from '@taurushq/protect-sdk';
 *
 * const client = ProtectClient.create({
 *   host: 'https://your-protect-instance.example.com',
 *   apiKey: 'your-api-key',
 *   apiSecret: 'your-hex-encoded-secret',
 * });
 *
 * // Use the client
 * const wallets = await client.walletsApi.walletServiceGetWalletsV2();
 *
 * // Access TaurusNetwork APIs
 * const participants = await client.taurusNetwork.participantApi.getAllParticipants();
 *
 * // Clean up
 * client.close();
 * ```
 *
 * @packageDocumentation
 */

// Main client
export {
  ProtectClient,
  TaurusNetworkNamespace,
  type ProtectClientConfig,
} from "./client";

// Error types
export {
  APIError,
  ValidationError,
  AuthenticationError,
  AuthorizationError,
  NotFoundError,
  RateLimitError,
  ServerError,
  IntegrityError,
  WhitelistError,
  ConfigurationError,
  RequestMetadataError,
  mapHttpError,
} from "./errors";

// Crypto utilities (for advanced use cases)
export {
  calculateHexHash,
  calculateSha256Bytes,
  calculateBase64Hmac,
  verifyBase64Hmac,
  constantTimeCompare,
  constantTimeCompareBytes,
  signData,
  verifySignature,
  decodePrivateKeyPem,
  decodePublicKeyPem,
  decodePublicKeysPem,
  encodePublicKeyPem,
  getPublicKeyFromPrivate,
  TPV1Auth,
  calculateSignedHeader,
} from "./crypto";

// Models
export * from "./models";

// Mappers (for advanced use cases)
export * from "./mappers";

// Services
export * from "./services";

// Transport (for custom middleware)
export { createTPV1Middleware } from "./transport";
