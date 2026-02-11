package helper

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// WhitelistedAddressVerifier provides methods to verify whitelisted address integrity.
type WhitelistedAddressVerifier struct {
	superAdminKeys     []*ecdsa.PublicKey
	minValidSignatures int
}

// NewWhitelistedAddressVerifier creates a new verifier with the given configuration.
func NewWhitelistedAddressVerifier(superAdminKeys []*ecdsa.PublicKey, minValidSignatures int) *WhitelistedAddressVerifier {
	return &WhitelistedAddressVerifier{
		superAdminKeys:     superAdminKeys,
		minValidSignatures: minValidSignatures,
	}
}

// VerificationResult contains the result of verification and the decoded rules container.
type VerificationResult struct {
	// RulesContainer is the decoded and verified rules container.
	RulesContainer *model.DecodedRulesContainer
	// VerifiedAddress is the address parsed from the verified payload.
	VerifiedAddress *model.WhitelistedAddress
	// VerifiedHash is the hash that was matched during verification.
	// This may differ from the input hash if a legacy hash format was used.
	VerifiedHash string
}

// VerifyWhitelistedAddress performs the complete 6-step verification of a whitelisted address.
// This implements the same verification flow as the Java SDK.
//
// Steps:
// 1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
// 2. Verify rules container signatures (SuperAdmin signatures)
// 3. Decode rules container (base64 -> protobuf -> model)
// 4. Verify hash coverage (metadata.hash in at least one signature.hashes)
// 5. Verify whitelist signatures (user signatures meet governance thresholds)
//
// The function does not mutate the input addr. If a legacy hash was matched during
// verification, it is returned in the result's VerifiedHash field.
// VerifyWhitelistedAddress performs the complete 6-step verification of a whitelisted address.
// If cachedRulesContainer is non-nil, steps 2-3 are skipped (already done during cache building).
func (v *WhitelistedAddressVerifier) VerifyWhitelistedAddress(
	addr *model.WhitelistedAddress,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
	cachedRulesContainer ...*model.DecodedRulesContainer,
) (*VerificationResult, error) {
	if addr == nil {
		return nil, fmt.Errorf("whitelisted address cannot be nil")
	}
	if addr.Metadata == nil {
		return nil, fmt.Errorf("metadata cannot be nil")
	}

	// Step 1: Verify metadata hash
	if err := v.verifyMetadataHash(addr); err != nil {
		return nil, err
	}

	var rulesContainer *model.DecodedRulesContainer
	if len(cachedRulesContainer) > 0 && cachedRulesContainer[0] != nil {
		// Steps 2-3 already done during cache building
		rulesContainer = cachedRulesContainer[0]
	} else {
		// Step 2: Verify rules container signatures
		if err := v.verifyRulesContainerSignatures(addr, userSignaturesDecoder); err != nil {
			return nil, err
		}

		// Step 3: Decode rules container
		var err error
		rulesContainer, err = v.decodeRulesContainer(addr, rulesContainerDecoder)
		if err != nil {
			return nil, err
		}
	}

	// Step 4: Verify hash coverage
	// verifiedHash may differ from addr.Metadata.Hash if a legacy hash format was matched
	verifiedHash, err := v.verifyHashInSignedHashes(addr)
	if err != nil {
		return nil, err
	}

	// Step 5: Verify whitelist signatures using the verified hash
	if err := v.verifyWhitelistSignatures(addr, rulesContainer, verifiedHash); err != nil {
		return nil, err
	}

	// Step 6: Parse WhitelistedAddress from verified payload
	verifiedAddr, err := ParseWhitelistedAddressFromJSON(addr.Metadata.PayloadAsString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse verified address: %w", err)
	}

	return &VerificationResult{
		RulesContainer:  rulesContainer,
		VerifiedAddress: verifiedAddr,
		VerifiedHash:    verifiedHash,
	}, nil
}

// VerifyAndDecodeRulesContainer verifies SuperAdmin signatures on a rules container
// and decodes it. This performs steps 2-3 of the verification flow for a single
// rules container, used by the service to build the normalized cache.
func (v *WhitelistedAddressVerifier) VerifyAndDecodeRulesContainer(
	rulesContainerBase64 string,
	rulesSignaturesBase64 string,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
) (*model.DecodedRulesContainer, error) {
	if rulesContainerBase64 == "" {
		return nil, &model.IntegrityError{Message: "rulesContainer is empty"}
	}
	if rulesSignaturesBase64 == "" {
		return nil, &model.IntegrityError{Message: "rulesSignatures is empty"}
	}

	// Decode and verify signatures
	signatures, err := userSignaturesDecoder(rulesSignaturesBase64)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules signatures: %v", err),
		}
	}

	rulesData, err := DecodeBase64(rulesContainerBase64)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	if err := VerifyGovernanceRulesSignatures(rulesData, signatures, v.superAdminKeys, v.minValidSignatures); err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("rules container signature verification failed: %v", err),
		}
	}

	// Decode rules container
	container, err := rulesContainerDecoder(rulesContainerBase64)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	return container, nil
}

// verifyMetadataHash verifies that the computed hash matches the provided hash.
// Step 1 of the verification flow.
func (v *WhitelistedAddressVerifier) verifyMetadataHash(addr *model.WhitelistedAddress) error {
	if addr.Metadata.PayloadAsString == "" {
		return &model.IntegrityError{Message: "payloadAsString is empty"}
	}
	if addr.Metadata.Hash == "" {
		return &model.IntegrityError{Message: "metadata hash is empty"}
	}

	computedHash := crypto.CalculateHexHash(addr.Metadata.PayloadAsString)
	if !ConstantTimeCompare(computedHash, addr.Metadata.Hash) {
		return &model.IntegrityError{
			Message: "metadata hash verification failed",
		}
	}

	return nil
}

// verifyRulesContainerSignatures verifies SuperAdmin signatures on the rules container.
// Step 2 of the verification flow.
func (v *WhitelistedAddressVerifier) verifyRulesContainerSignatures(
	addr *model.WhitelistedAddress,
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
) error {
	if len(v.superAdminKeys) == 0 {
		return &model.IntegrityError{Message: "no SuperAdmin keys configured for verification"}
	}

	if addr.RulesContainer == "" {
		return &model.IntegrityError{Message: "rulesContainer is empty"}
	}
	if addr.RulesSignatures == "" {
		return &model.IntegrityError{Message: "rulesSignatures is empty"}
	}

	// Decode rules signatures (protobuf UserSignatures)
	signatures, err := userSignaturesDecoder(addr.RulesSignatures)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules signatures: %v", err),
		}
	}

	// Decode rules container data
	rulesData, err := DecodeBase64(addr.RulesContainer)
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
func (v *WhitelistedAddressVerifier) decodeRulesContainer(
	addr *model.WhitelistedAddress,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
) (*model.DecodedRulesContainer, error) {
	if rulesContainerDecoder == nil {
		return nil, fmt.Errorf("rulesContainerDecoder is required")
	}

	container, err := rulesContainerDecoder(addr.RulesContainer)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	return container, nil
}

// verifyHashInSignedHashes verifies that the metadata hash is covered by at least one signature.
// Step 4 of the verification flow.
// Returns the hash that was found (may be a legacy hash).
func (v *WhitelistedAddressVerifier) verifyHashInSignedHashes(addr *model.WhitelistedAddress) (string, error) {
	if addr.SignedAddress == nil {
		return "", &model.IntegrityError{Message: "signedAddress is nil"}
	}

	signatures := addr.SignedAddress.Signatures
	if len(signatures) == 0 {
		return "", &model.IntegrityError{Message: "no signatures in signedAddress"}
	}

	// Try the provided hash first
	providedHash := addr.Metadata.Hash
	if VerifyHashCoverage(providedHash, signatures) {
		return providedHash, nil
	}

	// Try legacy hashes for backward compatibility
	legacyHashes := ComputeLegacyHashes(addr.Metadata.PayloadAsString)
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
func (v *WhitelistedAddressVerifier) verifyWhitelistSignatures(
	addr *model.WhitelistedAddress,
	rulesContainer *model.DecodedRulesContainer,
	metadataHash string,
) error {
	// Find matching address whitelisting rules
	whitelistRules := rulesContainer.FindAddressWhitelistingRules(addr.Blockchain, addr.Network)
	if whitelistRules == nil {
		return &model.WhitelistError{
			Message: fmt.Sprintf("no address whitelisting rules found for blockchain=%s network=%s",
				addr.Blockchain, addr.Network),
		}
	}

	// Determine which thresholds to use based on rule lines
	parallelThresholds := v.getApplicableThresholds(whitelistRules, addr)
	if len(parallelThresholds) == 0 {
		return &model.WhitelistError{Message: "no threshold rules defined"}
	}

	// Try to verify all paths (OR logic - only one needs to succeed)
	pathFailures := v.tryVerifyAllPaths(parallelThresholds, rulesContainer, addr.SignedAddress.Signatures, metadataHash)
	if len(pathFailures) > 0 {
		return &model.WhitelistError{
			Message: fmt.Sprintf("signature verification failed for whitelisted address (ID: %s): "+
				"no approval path satisfied the threshold requirements. %s",
				addr.ID, strings.Join(pathFailures, "; ")),
		}
	}

	return nil
}

// getApplicableThresholds determines which thresholds to use based on rule lines.
// Checks rule lines only when: NO linked addresses AND exactly 1 linked wallet.
func (v *WhitelistedAddressVerifier) getApplicableThresholds(
	rules *model.AddressWhitelistingRules,
	addr *model.WhitelistedAddress,
) []*model.SequentialThresholds {
	hasLinkedAddresses := len(addr.LinkedInternalAddresses) > 0
	walletCount := len(addr.LinkedWallets)

	// Check rule lines only if: no linked addresses AND exactly 1 linked wallet
	shouldCheckRuleLines := !hasLinkedAddresses && walletCount == 1

	if shouldCheckRuleLines && len(rules.Lines) > 0 {
		walletPath := addr.LinkedWallets[0].Path

		// Find matching line by wallet path
		for _, line := range rules.Lines {
			if v.matchesWalletPath(line, walletPath) {
				return line.ParallelThresholds
			}
		}
	}

	// Fallback to default thresholds
	return rules.ParallelThresholds
}

// matchesWalletPath checks if a rule line matches the given wallet path.
func (v *WhitelistedAddressVerifier) matchesWalletPath(line *model.AddressWhitelistingLine, walletPath string) bool {
	if len(line.Cells) == 0 {
		return false
	}

	source := line.Cells[0]
	if source.Type != model.RuleSourceTypeInternalWallet {
		return false
	}

	if source.InternalWallet == nil {
		return false
	}

	return walletPath != "" && walletPath == source.InternalWallet.Path
}

// precomputeHashesJSON pre-computes JSON serialization of each signature's hashes array.
// This avoids redundant json.Marshal calls when the same signature is checked
// across multiple group thresholds in the verification loops.
// Returns a map from signature index to marshaled JSON bytes.
// If marshaling fails for a signature, that index is absent from the map.
func precomputeHashesJSON(signatures []model.WhitelistSignature) map[int][]byte {
	result := make(map[int][]byte, len(signatures))
	for i, sig := range signatures {
		hashesJSON, err := json.Marshal(sig.Hashes)
		if err == nil {
			result[i] = hashesJSON
		}
	}
	return result
}

// tryVerifyAllPaths tries to verify all parallel threshold paths.
// Returns empty slice if verification passed, or list of failure messages if all paths failed.
func (v *WhitelistedAddressVerifier) tryVerifyAllPaths(
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
func (v *WhitelistedAddressVerifier) verifySequentialThresholds(
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
func (v *WhitelistedAddressVerifier) verifyGroupThreshold(
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

// containsHash checks if a hash is in the list (using constant-time comparison).
func containsHash(hashes []string, hash string) bool {
	for _, h := range hashes {
		if ConstantTimeCompare(h, hash) {
			return true
		}
	}
	return false
}
