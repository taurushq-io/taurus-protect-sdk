package helper

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Regular expressions for legacy hash computation.
var (
	// contractTypePattern matches ,\"contractType\":\"...\""
	contractTypePattern = regexp.MustCompile(`,\"contractType\":\"[^\"]*\"`)

	// labelInObjectPattern matches ,\"label\":\"...\"}
	// This only matches labels inside objects (followed by closing brace),
	// not the main address label which is followed by other fields.
	labelInObjectPattern = regexp.MustCompile(`,\"label\":\"[^\"]*\"}`)

	// isNFTPattern matches "isNFT":(true|false) with optional leading/trailing comma
	// Used for asset legacy hash computation.
	isNFTPatternWithLeadingComma  = regexp.MustCompile(`,\"isNFT\":(true|false)`)
	isNFTPatternWithTrailingComma = regexp.MustCompile(`\"isNFT\":(true|false),`)

	// kindTypePattern matches "kindType":"..." with optional leading/trailing comma
	// Used for asset legacy hash computation.
	kindTypePatternWithLeadingComma  = regexp.MustCompile(`,\"kindType\":\"[^\"]*\"`)
	kindTypePatternWithTrailingComma = regexp.MustCompile(`\"kindType\":\"[^\"]*\",`)
)

// ComputeLegacyHashes computes alternative hashes for backward compatibility.
// This handles addresses signed before schema changes by removing certain fields
// and recomputing the hash.
//
// Strategies:
// 1. Remove contractType field (addresses signed before contractType was added)
// 2. Remove labels from linkedInternalAddresses (after contractType but before labels were added)
// 3. Remove both contractType and labels (before both fields were added)
func ComputeLegacyHashes(payloadAsString string) []string {
	if payloadAsString == "" {
		return nil
	}

	// Use a map to track unique hashes (preserves insertion order with slice)
	seen := make(map[string]bool)
	var hashes []string

	addHash := func(payload string) {
		hash := crypto.CalculateHexHash(payload)
		if !seen[hash] {
			seen[hash] = true
			hashes = append(hashes, hash)
		}
	}

	// Strategy 1: Remove contractType only
	// Handles addresses signed before contractType was added to schema
	withoutContractType := contractTypePattern.ReplaceAllString(payloadAsString, "")
	if withoutContractType != payloadAsString {
		addHash(withoutContractType)
	}

	// Strategy 2: Remove labels from linkedInternalAddresses objects only (keep contractType)
	// Handles addresses signed after contractType was added but before labels were added
	withoutLabels := labelInObjectPattern.ReplaceAllString(payloadAsString, "}")
	if withoutLabels != payloadAsString {
		addHash(withoutLabels)
	}

	// Strategy 3: Remove BOTH contractType AND labels from linkedInternalAddresses
	// Handles addresses signed before both fields were added
	withoutBoth := labelInObjectPattern.ReplaceAllString(payloadAsString, "}")
	withoutBoth = contractTypePattern.ReplaceAllString(withoutBoth, "")
	if withoutBoth != payloadAsString {
		addHash(withoutBoth)
	}

	return hashes
}

// ComputeAssetLegacyHashes computes alternative hashes for backward compatibility with assets.
// This handles assets signed before schema changes by removing certain fields
// and recomputing the hash.
//
// Strategies (aligned with Java SDK WhitelistedAssetService.computeLegacyHashes):
// 1. Remove isNFT field (assets signed before isNFT was added)
// 2. Remove kindType field (assets signed before kindType was added)
// 3. Remove both isNFT and kindType (assets signed before both fields were added)
func ComputeAssetLegacyHashes(payloadAsString string) []string {
	if payloadAsString == "" {
		return nil
	}

	// Use a map to track unique hashes (preserves insertion order with slice)
	seen := make(map[string]bool)
	var hashes []string

	addHash := func(payload string) {
		hash := crypto.CalculateHexHash(payload)
		if !seen[hash] {
			seen[hash] = true
			hashes = append(hashes, hash)
		}
	}

	// Strategy 1: Remove isNFT only
	// Handles assets signed before isNFT was added to schema
	withoutIsNFT := isNFTPatternWithLeadingComma.ReplaceAllString(payloadAsString, "")
	withoutIsNFT = isNFTPatternWithTrailingComma.ReplaceAllString(withoutIsNFT, "")
	if withoutIsNFT != payloadAsString {
		addHash(withoutIsNFT)
	}

	// Strategy 2: Remove kindType only
	// Handles assets signed before kindType was added to schema
	withoutKindType := kindTypePatternWithLeadingComma.ReplaceAllString(payloadAsString, "")
	withoutKindType = kindTypePatternWithTrailingComma.ReplaceAllString(withoutKindType, "")
	if withoutKindType != payloadAsString {
		addHash(withoutKindType)
	}

	// Strategy 3: Remove BOTH isNFT AND kindType
	// Handles assets signed before both fields were added
	// Note: Order matches Java implementation - remove isNFT first, then kindType
	withoutBoth := isNFTPatternWithLeadingComma.ReplaceAllString(payloadAsString, "")
	withoutBoth = isNFTPatternWithTrailingComma.ReplaceAllString(withoutBoth, "")
	withoutBoth = kindTypePatternWithLeadingComma.ReplaceAllString(withoutBoth, "")
	withoutBoth = kindTypePatternWithTrailingComma.ReplaceAllString(withoutBoth, "")
	if withoutBoth != payloadAsString {
		addHash(withoutBoth)
	}

	return hashes
}

// ParseWhitelistedAddressFromJSON parses a WhitelistedAddress from verified JSON payload.
// This is used to extract the signed fields from the cryptographically verified payload.
func ParseWhitelistedAddressFromJSON(jsonPayload string) (*model.WhitelistedAddress, error) {
	if jsonPayload == "" {
		return nil, fmt.Errorf("JSON payload cannot be empty")
	}

	var payload whitelistPayload
	if err := json.Unmarshal([]byte(jsonPayload), &payload); err != nil {
		return nil, fmt.Errorf("failed to parse whitelist payload: %w", err)
	}

	addr := &model.WhitelistedAddress{
		Blockchain:    payload.Currency,
		Network:       payload.Network,
		Address:       payload.Address,
		Memo:          payload.Memo,
		Label:         payload.Label,
		CustomerId:    payload.CustomerID,
		ContractType:  payload.ContractType,
		TnParticipantID: payload.TnParticipantID,
		AddressType:   payload.AddressType,
	}

	// Parse exchangeAccountId
	if payload.ExchangeAccountID != "" {
		var id int64
		if _, err := fmt.Sscanf(payload.ExchangeAccountID, "%d", &id); err == nil {
			addr.ExchangeAccountId = id
		}
	}

	// Parse linkedInternalAddresses
	for _, lia := range payload.LinkedInternalAddresses {
		addr.LinkedInternalAddresses = append(addr.LinkedInternalAddresses, model.InternalAddress{
			ID:    lia.ID,
			Label: lia.Label,
		})
	}

	// Parse linkedWallets
	for _, lw := range payload.LinkedWallets {
		addr.LinkedWallets = append(addr.LinkedWallets, model.InternalWallet{
			ID:    lw.ID,
			Path:  lw.Path,
			Label: lw.Name, // Note: JSON field is "name" but model field is "Label"
		})
	}

	return addr, nil
}

// whitelistPayload represents the JSON structure of the signed whitelist payload.
type whitelistPayload struct {
	Currency                string                   `json:"currency"`
	Network                 string                   `json:"network"`
	Address                 string                   `json:"address"`
	Memo                    string                   `json:"memo"`
	Label                   string                   `json:"label"`
	CustomerID              string                   `json:"customerId"`
	ContractType            string                   `json:"contractType"`
	TnParticipantID         string                   `json:"tnParticipantID"`
	AddressType             string                   `json:"addressType"`
	ExchangeAccountID       string                   `json:"exchangeAccountId"`
	LinkedInternalAddresses []linkedInternalAddress  `json:"linkedInternalAddresses"`
	LinkedWallets           []linkedWallet           `json:"linkedWallets"`
}

type linkedInternalAddress struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
	Label   string `json:"label"`
}

// UnmarshalJSON handles both numeric and string ID values from the API.
func (l *linkedInternalAddress) UnmarshalJSON(data []byte) error {
	type alias linkedInternalAddress
	aux := &struct {
		ID json.Number `json:"id"`
		*alias
	}{
		alias: (*alias)(l),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if aux.ID.String() != "" {
		id, err := aux.ID.Int64()
		if err != nil {
			return fmt.Errorf("failed to parse linkedInternalAddress id %q: %w", aux.ID.String(), err)
		}
		l.ID = id
	}
	return nil
}

type linkedWallet struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// UnmarshalJSON handles both numeric and string ID values from the API.
func (l *linkedWallet) UnmarshalJSON(data []byte) error {
	type alias linkedWallet
	aux := &struct {
		ID json.Number `json:"id"`
		*alias
	}{
		alias: (*alias)(l),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	if aux.ID.String() != "" {
		id, err := aux.ID.Int64()
		if err != nil {
			return fmt.Errorf("failed to parse linkedWallet id %q: %w", aux.ID.String(), err)
		}
		l.ID = id
	}
	return nil
}

// CheckHashesSignature verifies a signature of a list of hashes.
func CheckHashesSignature(hashes []string, signature string, publicKey interface{}) error {
	if len(hashes) == 0 {
		return fmt.Errorf("hashes cannot be empty")
	}
	if signature == "" {
		return fmt.Errorf("signature cannot be empty")
	}

	// JSON encode the hashes
	jsonBytes, err := json.Marshal(hashes)
	if err != nil {
		return fmt.Errorf("failed to JSON encode hashes: %w", err)
	}

	// Verify using ECDSA
	switch pk := publicKey.(type) {
	case *model.RuleUser:
		if pk.PublicKey == nil {
			return fmt.Errorf("user public key is nil")
		}
		valid, err := crypto.VerifySignature(pk.PublicKey, jsonBytes, signature)
		if err != nil {
			return fmt.Errorf("signature verification failed: %w", err)
		}
		if !valid {
			return fmt.Errorf("invalid signature for hashes %s", strings.Join(hashes, ", "))
		}
		return nil
	default:
		return fmt.Errorf("unsupported public key type")
	}
}
