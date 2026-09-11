package helper

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// Step 5 picks WHICH governance rule judges a row from a (blockchain, network) key, and the
// network half is not always signed: governance carries a per-rule includeNetworkInPayload flag,
// and when it is off the signed payload has no `network` member at all — the common case in
// captured production data. The key then falls back to the network on the unsigned response DTO,
// which hands a response-controlling server the choice of which quorum the row must meet.
//
// The flag cannot be consulted: it has no proto backing in any SDK, so it is never part of the
// SuperAdmin-signed container. The answer is therefore to enforce every tier the unsigned value
// could have selected.
//
// Payload carries no `network`, which is what makes the DTO value load-bearing.
const tierPayloadWithoutNetwork = `{"currency":"ETH","address":"0xdead","memo":"",` +
	`"label":"Pending counterparty","customerId":"c-1","linkedInternalAddresses":[]}`

// buildTierFixture builds a container with a STRICT mainnet tier (2 distinct signers) and a WEAK
// goerli tier (1 signer), and a row carrying exactly one signature. Under the goerli tier alone
// the row passes; under mainnet it does not.
func buildTierFixture(t *testing.T, dtoNetwork string) (
	*WhitelistedAddressVerifier,
	*model.WhitelistedAddress,
	func(string) (*model.DecodedRulesContainer, error),
	func(string) ([]*model.RuleUserSignature, error),
) {
	t.Helper()

	saPriv, saPub := generateAssetTestKeyPair(t)
	opsPriv, opsPub := generateAssetTestKeyPair(t)
	_, compliancePub := generateAssetTestKeyPair(t)

	hash := crypto.CalculateHexHash(tierPayloadWithoutNetwork)
	hashes := []string{hash}
	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		t.Fatalf("marshal hashes: %v", err)
	}
	opsSig, err := crypto.SignData(opsPriv, hashesJSON)
	if err != nil {
		t.Fatalf("sign hashes: %v", err)
	}

	rulesB64 := encodeRulesContainerJSON(t, map[string]interface{}{"placeholder": true})
	rulesData, err := base64.StdEncoding.DecodeString(rulesB64)
	if err != nil {
		t.Fatalf("decode rules: %v", err)
	}
	saSig, err := crypto.SignData(saPriv, rulesData)
	if err != nil {
		t.Fatalf("sign rules: %v", err)
	}

	addr := &model.WhitelistedAddress{
		ID:         "addr-tier",
		Blockchain: "ETH",
		Network:    dtoNetwork, // unsigned, attacker-controlled
		Metadata: &model.WhitelistedAssetMetadata{
			Hash:            hash,
			PayloadAsString: tierPayloadWithoutNetwork,
		},
		RulesContainer:  rulesB64,
		RulesSignatures: base64.StdEncoding.EncodeToString([]byte("dummy")),
		SignedAddress: &model.SignedWhitelistedAddress{
			Signatures: []model.WhitelistSignature{
				{
					UserSignature: &model.WhitelistUserSignature{UserID: "ops@bank.com", Signature: opsSig},
					Hashes:        hashes,
				},
			},
		},
	}

	rcDecoder := func(string) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{ID: "ops@bank.com", PublicKey: opsPub, Roles: []string{"USER"}},
				{ID: "compliance@bank.com", PublicKey: compliancePub, Roles: []string{"USER"}},
			},
			Groups: []*model.RuleGroup{
				{ID: "ops", UserIDs: []string{"ops@bank.com"}},
				{ID: "compliance", UserIDs: []string{"compliance@bank.com"}},
			},
			AddressWhitelistingRules: []*model.AddressWhitelistingRules{
				{
					// STRICT: ops AND compliance.
					Currency: "ETH",
					Network:  "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "ops", MinimumSignatures: 1},
								{GroupID: "compliance", MinimumSignatures: 1},
							},
						},
					},
				},
				{
					// WEAK: ops alone.
					Currency: "ETH",
					Network:  "goerli",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "ops", MinimumSignatures: 1},
							},
						},
					},
				},
			},
		}, nil
	}
	usDecoder := func(string) ([]*model.RuleUserSignature, error) {
		return []*model.RuleUserSignature{{UserID: "sa@bank.com", Signature: saSig}}, nil
	}

	return NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{saPub}, 1), addr, rcDecoder, usDecoder
}

// TestUnsignedDtoNetworkCannotSelectTheWeakerTier is the finding: the row holds only an Ops
// signature, which satisfies the goerli tier but not mainnet. Because the signed payload carries
// no network, the server names "goerli" on the DTO to have the weaker quorum applied.
func TestUnsignedDtoNetworkCannotSelectTheWeakerTier(t *testing.T) {
	v, addr, rcDecoder, usDecoder := buildTierFixture(t, "goerli")

	_, err := v.VerifyWhitelistedAddress(addr, rcDecoder, usDecoder)
	if err == nil {
		t.Fatal("an unsigned DTO network selected the weaker tier: the row met goerli's 1-of-1 " +
			"but not mainnet's ops+compliance, and it verified anyway")
	}
	var whitelistErr *model.WhitelistError
	if !asWhitelistError(err, &whitelistErr) {
		t.Errorf("want a WhitelistError, got %T: %v", err, err)
	}
}

// TestUnsignedDtoNetworkIsStillRefusedWhenItNamesTheStrictTier shows the refusal does not
// depend on which network the server named — with the network unsigned, both tiers are enforced.
func TestUnsignedDtoNetworkIsStillRefusedWhenItNamesTheStrictTier(t *testing.T) {
	v, addr, rcDecoder, usDecoder := buildTierFixture(t, "mainnet")

	if _, err := v.VerifyWhitelistedAddress(addr, rcDecoder, usDecoder); err == nil {
		t.Fatal("expected refusal: the row does not satisfy the mainnet tier")
	}
}

// TestUnsignedDtoNetworkStillVerifiesWhenEveryTierIsSatisfied guards against over-strictness:
// conservative enforcement must not reject a row that meets every candidate tier. Without this,
// "reject when the network is unsigned" would pass the test above for the wrong reason.
func TestUnsignedDtoNetworkStillVerifiesWhenEveryTierIsSatisfied(t *testing.T) {
	v, addr, rcDecoder, usDecoder := buildTierFixture(t, "goerli")

	// Add the compliance signature, so ops+compliance (mainnet) and ops (goerli) both hold.
	compliancePriv, compliancePub := generateAssetTestKeyPair(t)
	hashesJSON, err := json.Marshal([]string{addr.Metadata.Hash})
	if err != nil {
		t.Fatalf("marshal hashes: %v", err)
	}
	complianceSig, err := crypto.SignData(compliancePriv, hashesJSON)
	if err != nil {
		t.Fatalf("sign hashes: %v", err)
	}
	addr.SignedAddress.Signatures = append(addr.SignedAddress.Signatures, model.WhitelistSignature{
		UserSignature: &model.WhitelistUserSignature{
			UserID:    "compliance@bank.com",
			Signature: complianceSig,
		},
		Hashes: []string{addr.Metadata.Hash},
	})

	// Point the container's compliance user at the key that actually signed.
	wrapped := func(b64 string) (*model.DecodedRulesContainer, error) {
		c, err := rcDecoder(b64)
		if err != nil {
			return nil, err
		}
		for _, u := range c.Users {
			if u.ID == "compliance@bank.com" {
				u.PublicKey = compliancePub
			}
		}
		return c, nil
	}

	if _, err := v.VerifyWhitelistedAddress(addr, wrapped, usDecoder); err != nil {
		t.Errorf("a row satisfying every candidate tier must verify: %v", err)
	}
}

// TestSignedNetworkStillSelectsExactlyOneTier: when the payload DOES carry the network, the key
// is fully signed and only that tier applies — conservative enforcement must not leak into the
// normal case and start demanding quorums governance did not assign.
func TestSignedNetworkStillSelectsExactlyOneTier(t *testing.T) {
	v, addr, rcDecoder, usDecoder := buildTierFixture(t, "goerli")

	// Re-sign with a payload that names goerli explicitly.
	signed := `{"currency":"ETH","network":"goerli","address":"0xdead","memo":"",` +
		`"label":"Pending counterparty","customerId":"c-1","linkedInternalAddresses":[]}`
	hash := crypto.CalculateHexHash(signed)
	addr.Metadata.Hash = hash
	addr.Metadata.PayloadAsString = signed

	opsPriv, opsPub := generateAssetTestKeyPair(t)
	hashesJSON, err := json.Marshal([]string{hash})
	if err != nil {
		t.Fatalf("marshal hashes: %v", err)
	}
	opsSig, err := crypto.SignData(opsPriv, hashesJSON)
	if err != nil {
		t.Fatalf("sign hashes: %v", err)
	}
	addr.SignedAddress.Signatures = []model.WhitelistSignature{
		{
			UserSignature: &model.WhitelistUserSignature{UserID: "ops@bank.com", Signature: opsSig},
			Hashes:        []string{hash},
		},
	}
	wrapped := func(b64 string) (*model.DecodedRulesContainer, error) {
		c, err := rcDecoder(b64)
		if err != nil {
			return nil, err
		}
		for _, u := range c.Users {
			if u.ID == "ops@bank.com" {
				u.PublicKey = opsPub
			}
		}
		return c, nil
	}

	if _, err := v.VerifyWhitelistedAddress(addr, wrapped, usDecoder); err != nil {
		t.Errorf("a signed network must select exactly its own tier: %v", err)
	}
}

// TestFindAddressWhitelistingRuleCandidates_Reachability pins the candidate set against the tier
// walk it has to mirror. Getting this wrong in either direction is a real bug: too few candidates
// leaves the steering open, too many rejects rows governance would accept.
func TestFindAddressWhitelistingRuleCandidates_Reachability(t *testing.T) {
	rule := func(currency, network string) *model.AddressWhitelistingRules {
		return &model.AddressWhitelistingRules{Currency: currency, Network: network}
	}

	t.Run("exact tiers plus the global default when the chain has no wildcard network", func(t *testing.T) {
		c := &model.DecodedRulesContainer{AddressWhitelistingRules: []*model.AddressWhitelistingRules{
			rule("ETH", "mainnet"), rule("ETH", "goerli"), rule("BTC", "mainnet"), rule("", ""),
		}}
		got := c.FindAddressWhitelistingRuleCandidates("ETH")
		if len(got) != 3 {
			t.Fatalf("want ETH/mainnet, ETH/goerli and the global default, got %d", len(got))
		}
	})

	t.Run("global default is unreachable when the chain has a wildcard-network rule", func(t *testing.T) {
		c := &model.DecodedRulesContainer{AddressWhitelistingRules: []*model.AddressWhitelistingRules{
			rule("ETH", "mainnet"), rule("ETH", ""), rule("", ""),
		}}
		got := c.FindAddressWhitelistingRuleCandidates("ETH")
		if len(got) != 2 {
			t.Fatalf("the tier walk stops at priority 2, so the global default is not a candidate; got %d", len(got))
		}
		for _, r := range got {
			if r.Currency == "" {
				t.Error("global default must not be a candidate when ETH/* exists")
			}
		}
	})

	t.Run("a single tier means unchanged behaviour", func(t *testing.T) {
		c := &model.DecodedRulesContainer{AddressWhitelistingRules: []*model.AddressWhitelistingRules{
			rule("ETH", "mainnet"),
		}}
		if got := c.FindAddressWhitelistingRuleCandidates("ETH"); len(got) != 1 {
			t.Fatalf("want exactly the one reachable tier, got %d", len(got))
		}
	})
}

// asWhitelistError is errors.As without importing errors into every assertion.
func asWhitelistError(err error, target **model.WhitelistError) bool {
	if e, ok := err.(*model.WhitelistError); ok {
		*target = e
		return true
	}
	return false
}
