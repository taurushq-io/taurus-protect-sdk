package helper

import (
	"crypto/ecdsa"
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
//
// Accepts an empty key set on purpose: only step 2 uses the keys, so steps 1, 5 and 6
// are verifiable without them. A keyless verifier fails closed in
// VerifyGovernanceRulesSignatures; clientConfig.validate rejects one earlier still.
// Making this a panic would assert an invariant the type does not have.
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

// VerifyAndDecodeRulesContainer performs steps 2-3 for one container, so the service can
// build its per-page cache.
func (v *WhitelistedAddressVerifier) VerifyAndDecodeRulesContainer(
	rulesContainerBase64 string,
	rulesSignaturesBase64 string,
	rulesContainerDecoder func(base64Data string) (*model.DecodedRulesContainer, error),
	userSignaturesDecoder func(base64Data string) ([]*model.RuleUserSignature, error),
) (*model.DecodedRulesContainer, error) {
	return VerifyAndDecodeRulesContainer(
		rulesContainerBase64,
		rulesSignaturesBase64,
		v.superAdminKeys,
		v.minValidSignatures,
		rulesContainerDecoder,
		userSignaturesDecoder,
	)
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
	// Which rules judge this address is decided by the SIGNED payload, not by the
	// surrounding response. A DTO that set blockchain="" would select the
	// global-default tier — broader than the rule the address belongs to.
	blockchain, network, err := resolveRuleKeyFor(addr.Metadata, addr.Blockchain, addr.Network)
	if err != nil {
		return err
	}

	// Find matching address whitelisting rules
	whitelistRules := rulesContainer.FindAddressWhitelistingRules(blockchain, network)
	if whitelistRules == nil {
		return &model.WhitelistError{
			Message: fmt.Sprintf("no address whitelisting rules found for blockchain=%s network=%s",
				blockchain, network),
		}
	}

	// Determine which thresholds to use based on rule lines
	parallelThresholds, err := v.getApplicableThresholds(whitelistRules, addr)
	if err != nil {
		return err
	}
	if len(parallelThresholds) == 0 {
		return &model.WhitelistError{Message: "no threshold rules defined"}
	}

	// Try to verify all paths (OR logic - only one needs to succeed)
	pathFailures := tryVerifyAllPaths(parallelThresholds, rulesContainer, addr.SignedAddress.Signatures, metadataHash)
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
//
// Returns a ContainerIntegrityError when a line carries a source cell this SDK could
// not type. Falling through to the container defaults in that case is the dangerous
// path: the line might be the one that matches, and its thresholds might be stricter
// than the default, so the address would be verified against a quorum governance
// never granted it — silently, with no error.
//
//	line source typed ──┬─ matches wallet path ─▶ line thresholds
//	                    └─ no match ────────────▶ container defaults
//	line source RAW ─────────────────────────────▶ ContainerIntegrityError
func (v *WhitelistedAddressVerifier) getApplicableThresholds(
	rules *model.AddressWhitelistingRules,
	addr *model.WhitelistedAddress,
) ([]*model.SequentialThresholds, error) {
	hasLinkedAddresses := len(addr.LinkedInternalAddresses) > 0
	walletCount := len(addr.LinkedWallets)

	// Check rule lines only if: no linked addresses AND exactly 1 linked wallet
	shouldCheckRuleLines := !hasLinkedAddresses && walletCount == 1

	if shouldCheckRuleLines && len(rules.Lines) > 0 {
		walletPath := addr.LinkedWallets[0].Path

		// Find matching line by wallet path
		for i, line := range rules.Lines {
			if lineHasUntypedSource(line) {
				return nil, &model.ContainerIntegrityError{
					Message: fmt.Sprintf(
						"address whitelisting rules for blockchain=%s network=%s line %d carry a "+
							"source cell this SDK version cannot interpret; refusing to fall back to "+
							"the container default thresholds, which may be weaker than the line's. "+
							"Upgrade the SDK to match the validatord that signed this container",
						rules.Currency, rules.Network, i),
				}
			}
			if v.matchesWalletPath(line, walletPath) {
				return line.ParallelThresholds, nil
			}
		}
	}

	// Fallback to default thresholds
	return rules.ParallelThresholds, nil
}

// lineHasUntypedSource reports whether the line's source cell was preserved verbatim
// because this SDK could not decode it into the typed model. matchesWalletPath reads
// Cells[0], so that is the cell whose meaning must be known.
func lineHasUntypedSource(line *model.AddressWhitelistingLine) bool {
	if line == nil || len(line.Cells) == 0 {
		return false
	}
	source := line.Cells[0]
	return source != nil && len(source.Raw) > 0
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

// containsHash lives in signature_verifier.go, beside VerifyHashCoverage.
// The step-5 threshold walk lives in group_threshold.go, shared with the asset verifier.
