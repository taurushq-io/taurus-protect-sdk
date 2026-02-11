/**
 * Helper utilities for the Taurus-PROTECT SDK.
 *
 * This module exports utility functions for security-sensitive operations.
 */

export {
  constantTimeCompare,
  constantTimeCompareBytes,
} from "./constant-time";

export { isValidSignature, verifyGovernanceRules } from "./signature-verifier";

export {
  verifyAddressSignature,
  verifyAddressSignatures,
} from "./address-signature-verifier";

export {
  computeLegacyHashes,
  computeAssetLegacyHashes,
  parseWhitelistedAddressFromJson,
  verifyHashCoverage,
} from "./whitelist-hash-helper";

export {
  WhitelistedAddressVerifier,
  type WhitelistedAddressVerifierConfig,
  type RulesContainerDecoder,
  type UserSignaturesDecoder,
} from "./whitelisted-address-verifier";

export {
  WhitelistedAssetVerifier,
  type WhitelistedAssetVerifierConfig,
  type RulesContainerDecoder as AssetRulesContainerDecoder,
  type UserSignaturesDecoder as AssetUserSignaturesDecoder,
} from "./whitelisted-asset-verifier";

export { getSourceAddress, getDestinationAddress, getAmount } from "./metadata-utils";
