package helper

import (
	"crypto/ecdsa"
	"fmt"
	"strings"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedAssetVerifier provides methods to verify whitelisted asset (contract) integrity.
type WhitelistedAssetVerifier struct {
	superAdminKeys     []*ecdsa.PublicKey
	minValidSignatures int
}

// NewWhitelistedAssetVerifier creates a new verifier with the given configuration.
func NewWhitelistedAssetVerifier(superAdminKeys []*ecdsa.PublicKey, minValidSignatures int) *WhitelistedAssetVerifier {
	return &WhitelistedAssetVerifier{
		superAdminKeys:     superAdminKeys,
		minValidSignatures: minValidSignatures,
	}
}

// AssetVerificationResult contains the result of verification and the decoded rules container.
type AssetVerificationResult struct {
	// RulesContainer is the decoded and verified rules container.
	RulesContainer *model.DecodedRulesContainer
	// VerifiedHash is the hash that was matched during verification.
	// This may differ from the input hash if a legacy hash format was used.
	VerifiedHash string
}

// VerifyWhitelistedAsset performs the complete 5-step verification of a whitelisted asset.
// This implements the same verification flow as the Java SDK.
//
// Steps:
// 1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
// 2. Verify rules container signatures (SuperAdmin signatures)
// 3. Decode rules container (base64 -> protobuf -> model)
// 4. Verify hash coverage (metadata.hash in at least one signature.hashes)
// 5. Verify whitelist signatures (user signatures meet governance thresholds)
//
// The function does not mutate the input asset. If a legacy hash was matched during
// verification, it is returned in the result's VerifiedHash field.
func (v *WhitelistedAssetVerifier) VerifyWhitelistedAsset(
	asset *model.WhitelistedAsset,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
) (*AssetVerificationResult, error) {
	if asset == nil {
		return nil, fmt.Errorf("whitelisted asset cannot be nil")
	}
	if asset.Metadata == nil {
		return nil, fmt.Errorf("metadata cannot be nil")
	}

	// Step 1: Verify metadata hash
	if err := v.verifyMetadataHash(asset); err != nil {
		return nil, err
	}

	// Step 2: Verify rules container signatures
	if err := v.verifyRulesContainerSignatures(asset, userSignaturesDecoder); err != nil {
		return nil, err
	}

	// Step 3: Decode rules container
	rulesContainer, err := v.decodeRulesContainer(asset, rulesContainerDecoder)
	if err != nil {
		return nil, err
	}

	// Step 4: Verify hash coverage
	// verifiedHash may differ from asset.Metadata.Hash if a legacy hash format was matched
	verifiedHash, err := v.verifyHashInSignedHashes(asset)
	if err != nil {
		return nil, err
	}

	// Step 5: Verify whitelist signatures using the verified hash
	if err := v.verifyWhitelistSignatures(asset, rulesContainer, verifiedHash); err != nil {
		return nil, err
	}

	return &AssetVerificationResult{
		RulesContainer: rulesContainer,
		VerifiedHash:   verifiedHash,
	}, nil
}

// verifyMetadataHash verifies that the computed hash matches the provided hash.
// Step 1 of the verification flow.
func (v *WhitelistedAssetVerifier) verifyMetadataHash(asset *model.WhitelistedAsset) error {
	if asset.Metadata.PayloadAsString == "" {
		return &model.IntegrityError{Message: "payloadAsString is empty"}
	}
	if asset.Metadata.Hash == "" {
		return &model.IntegrityError{Message: "metadata hash is empty"}
	}

	computedHash := crypto.CalculateHexHash(asset.Metadata.PayloadAsString)
	if !ConstantTimeCompare(computedHash, asset.Metadata.Hash) {
		return &model.IntegrityError{
			Message: "metadata hash verification failed",
		}
	}

	return nil
}

// verifyRulesContainerSignatures verifies SuperAdmin signatures on the rules container.
// Step 2 of the verification flow.
func (v *WhitelistedAssetVerifier) verifyRulesContainerSignatures(
	asset *model.WhitelistedAsset,
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
) error {
	if len(v.superAdminKeys) == 0 {
		return &model.IntegrityError{Message: "no SuperAdmin keys configured for verification"}
	}

	if asset.RulesContainer == "" {
		return &model.IntegrityError{Message: "rulesContainer is empty"}
	}
	if asset.RulesSignatures == "" {
		return &model.IntegrityError{Message: "rulesSignatures is empty"}
	}

	// Decode rules signatures (protobuf UserSignatures)
	signatures, err := userSignaturesDecoder(asset.RulesSignatures)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules signatures: %v", err),
		}
	}

	// Decode rules container data
	rulesData, err := DecodeBase64(asset.RulesContainer)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	// Verify signatures
	if err := VerifyGovernanceRulesSignatures(rulesData, signatures, v.superAdminKeys, v.minValidSignatures); err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("rules container signature verification failed: %v", err),
		}
	}

	return nil
}

// decodeRulesContainer decodes the base64 protobuf rules container.
// Step 3 of the verification flow.
func (v *WhitelistedAssetVerifier) decodeRulesContainer(
	asset *model.WhitelistedAsset,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
) (*model.DecodedRulesContainer, error) {
	if rulesContainerDecoder == nil {
		return nil, fmt.Errorf("rulesContainerDecoder is required")
	}

	container, err := rulesContainerDecoder(asset.RulesContainer)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	return container, nil
}

// verifyHashInSignedHashes verifies that the metadata hash is covered by at least one signature.
// Step 4 of the verification flow.
// Returns the hash that was found (may be a legacy hash for backward compatibility).
func (v *WhitelistedAssetVerifier) verifyHashInSignedHashes(asset *model.WhitelistedAsset) (string, error) {
	if asset.SignedContractAddress == nil {
		return "", &model.IntegrityError{Message: "signedContractAddress is nil"}
	}

	signatures := asset.SignedContractAddress.Signatures
	if len(signatures) == 0 {
		return "", &model.IntegrityError{Message: "no signatures in signedContractAddress"}
	}

	// Try the provided hash first
	providedHash := asset.Metadata.Hash
	if VerifyHashCoverage(providedHash, signatures) {
		return providedHash, nil
	}

	// Try legacy hashes for backward compatibility
	// This handles assets signed before schema changes (e.g., before isNFT or kindType was added)
	legacyHashes := ComputeAssetLegacyHashes(asset.Metadata.PayloadAsString)
	for _, legacyHash := range legacyHashes {
		if VerifyHashCoverage(legacyHash, signatures) {
			return legacyHash, nil
		}
	}

	return "", &model.IntegrityError{
		Message: "metadata hash is not covered by any signature",
	}
}

// verifyWhitelistSignatures verifies user signatures meet governance threshold requirements.
// Step 5 of the verification flow.
func (v *WhitelistedAssetVerifier) verifyWhitelistSignatures(
	asset *model.WhitelistedAsset,
	rulesContainer *model.DecodedRulesContainer,
	metadataHash string,
) error {

	// Find matching contract address whitelisting rules
	whitelistRules := rulesContainer.FindContractAddressWhitelistingRules(asset.Blockchain, asset.Network)
	if whitelistRules == nil {
		return &model.WhitelistError{
			Message: fmt.Sprintf("no contract address whitelisting rules found for blockchain=%s network=%s",
				asset.Blockchain, asset.Network),
		}
	}

	// Contract whitelisting uses parallelThresholds directly (no rule lines matching)
	parallelThresholds := whitelistRules.ParallelThresholds
	if len(parallelThresholds) == 0 {
		return &model.WhitelistError{Message: "no threshold rules defined"}
	}

	// Try to verify all paths (OR logic - only one needs to succeed)
	pathFailures := v.tryVerifyAllPaths(parallelThresholds, rulesContainer, asset.SignedContractAddress.Signatures, metadataHash)
	if len(pathFailures) > 0 {
		return &model.WhitelistError{
			Message: fmt.Sprintf("signature verification failed for whitelisted asset (ID: %s): "+
				"no approval path satisfied the threshold requirements. %s",
				asset.ID, strings.Join(pathFailures, "; ")),
		}
	}

	return nil
}

// tryVerifyAllPaths tries to verify all parallel threshold paths.
// Returns empty slice if verification passed, or list of failure messages if all paths failed.
func (v *WhitelistedAssetVerifier) tryVerifyAllPaths(
	parallelThresholds []*model.SequentialThresholds,
	rulesContainer *model.DecodedRulesContainer,
	signatures []model.WhitelistSignature,
	metadataHash string,
) []string {
	// Pre-compute JSON serialization of each signature's hashes array once,
	// so it can be reused across all group threshold checks.
	hashesJSONMap := precomputeHashesJSON(signatures)

	var pathFailures []string

	for i, seqThreshold := range parallelThresholds {
		err := v.verifySequentialThresholds(seqThreshold, rulesContainer, signatures, metadataHash, hashesJSONMap)
		if err == nil {
			return nil // Verification passed
		}
		pathFailures = append(pathFailures, fmt.Sprintf("Path %d: %s", i+1, err.Error()))
	}

	return pathFailures
}

// verifySequentialThresholds verifies all group thresholds in a sequential threshold path.
func (v *WhitelistedAssetVerifier) verifySequentialThresholds(
	seqThreshold *model.SequentialThresholds,
	rulesContainer *model.DecodedRulesContainer,
	signatures []model.WhitelistSignature,
	metadataHash string,
	hashesJSONMap map[int][]byte,
) error {
	if seqThreshold == nil || len(seqThreshold.Thresholds) == 0 {
		return &model.IntegrityError{Message: "no group thresholds defined"}
	}

	// ALL group thresholds must be satisfied (AND logic)
	for _, groupThreshold := range seqThreshold.Thresholds {
		if err := v.verifyGroupThreshold(groupThreshold, rulesContainer, signatures, metadataHash, hashesJSONMap); err != nil {
			return err
		}
	}

	return nil
}

// verifyGroupThreshold verifies that a group threshold is met.
func (v *WhitelistedAssetVerifier) verifyGroupThreshold(
	groupThreshold *model.GroupThreshold,
	rulesContainer *model.DecodedRulesContainer,
	signatures []model.WhitelistSignature,
	metadataHash string,
	hashesJSONMap map[int][]byte,
) error {
	groupID := groupThreshold.GroupID
	minSigs := groupThreshold.MinimumSignatures

	group := rulesContainer.FindGroupByID(groupID)
	if group == nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("group '%s' not found in rules container", groupID),
		}
	}

	if len(group.UserIDs) == 0 {
		if minSigs > 0 {
			return &model.IntegrityError{
				Message: fmt.Sprintf("group '%s' has no users but requires %d signature(s)", groupID, minSigs),
			}
		}
		return nil // minSignatures == 0, so empty group is OK
	}

	// Build set for faster lookup
	groupUserIDSet := make(map[string]bool)
	for _, uid := range group.UserIDs {
		groupUserIDSet[uid] = true
	}

	// Count valid signatures from users in this group
	validCount := 0
	var skippedReasons []string

	for i, sig := range signatures {
		if sig.UserSignature == nil {
			skippedReasons = append(skippedReasons, "signature has nil userSig")
			continue
		}

		sigUserID := sig.UserSignature.UserID
		if !groupUserIDSet[sigUserID] {
			continue // Signer not in this group - not an error, just not relevant
		}

		// Check that metadata hash is covered by this signature
		if !containsHash(sig.Hashes, metadataHash) {
			skippedReasons = append(skippedReasons, fmt.Sprintf(
				"user '%s' signature does not cover metadata hash '%s' (signed hashes=%v)",
				sigUserID, metadataHash, sig.Hashes))
			continue
		}

		user := rulesContainer.FindUserByID(sigUserID)
		if user == nil {
			skippedReasons = append(skippedReasons, fmt.Sprintf("user '%s' not found in rules container", sigUserID))
			continue
		}
		if user.PublicKey == nil {
			skippedReasons = append(skippedReasons, fmt.Sprintf("user '%s' has no public key", sigUserID))
			continue
		}

		// Use pre-computed JSON-encoded hashes array
		hashesJSON, ok := hashesJSONMap[i]
		if !ok {
			skippedReasons = append(skippedReasons, fmt.Sprintf("failed to marshal hashes for user '%s'", sigUserID))
			continue
		}

		valid, err := crypto.VerifySignature(user.PublicKey, hashesJSON, sig.UserSignature.Signature)
		if err != nil || !valid {
			skippedReasons = append(skippedReasons, fmt.Sprintf("user '%s' signature verification failed", sigUserID))
			continue
		}

		validCount++
		if validCount >= minSigs {
			return nil // Threshold met
		}
	}

	// Threshold not met
	message := fmt.Sprintf("group '%s' requires %d signature(s) but only %d valid", groupID, minSigs, validCount)
	if len(skippedReasons) > 0 {
		message += " [" + strings.Join(skippedReasons, "; ") + "]"
	}
	return &model.IntegrityError{Message: message}
}
