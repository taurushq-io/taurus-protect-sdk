package helper

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// The legacy-hash tolerance in step 4 accepts a hash computed over a REGEX-STRIPPED rewrite of
// the delivered payload. The strips are not injective, so a response-controlling server can
// append a member the strip removes to a genuinely signed payload: the residue is the signed
// bytes exactly, every signature check passes, and a step 6 that parsed the DELIVERED payload
// would hand back the appended value as verified.
//
// These are end-to-end tests through VerifyWhitelistedAddress with REAL signatures. Before this
// file, the full address flow was only ever exercised with nil inputs (see the two nil cases in
// whitelisted_address_verifier_test.go), so nothing in the Go suite could tell "closed" from
// "narrowed" — every legacy test asserted hashes, which the attack does not disturb.
//
// Two shapes, because the two defences cover different ones:
//
//	appended duplicate `label`      -> a duplicate key; caught by the variant parse AND, on the
//	                                   service path, by rejectDuplicateObjectKeys in the mapper
//	appended `contractType`         -> NOT a duplicate when the signed payload lacks the field,
//	                                   so ONLY the variant parse catches it

type addressInjectionFixture struct {
	verifier  *WhitelistedAddressVerifier
	addr      *model.WhitelistedAddress
	rcDecoder func(string) (*model.DecodedRulesContainer, error)
	usDecoder func(string) ([]*model.RuleUserSignature, error)
}

// buildAddressInjectionFixture signs signedPayload with a group member's key, then serves
// deliveredPayload instead, with metadata.hash recomputed over the delivered bytes so step 1
// passes. That is exactly the attacker's position: it controls the response, not the keys.
func buildAddressInjectionFixture(
	t *testing.T,
	signedPayload string,
	deliveredPayload string,
) *addressInjectionFixture {
	t.Helper()

	saPriv, saPub := generateAssetTestKeyPair(t)
	userPriv, userPub := generateAssetTestKeyPair(t)

	// The signature covers the hash of the payload that was genuinely signed.
	signedHash := crypto.CalculateHexHash(signedPayload)
	hashes := []string{signedHash}
	hashesJSON, err := json.Marshal(hashes)
	if err != nil {
		t.Fatalf("failed to marshal hashes: %v", err)
	}
	userSig, err := crypto.SignData(userPriv, hashesJSON)
	if err != nil {
		t.Fatalf("failed to sign hashes: %v", err)
	}

	rulesB64 := encodeRulesContainerJSON(t, map[string]interface{}{
		"users": []map[string]interface{}{
			{"id": "user1@bank.com", "publicKey": publicKeyToPEM(t, userPub), "roles": []string{"USER"}},
		},
	})
	rulesData, err := base64.StdEncoding.DecodeString(rulesB64)
	if err != nil {
		t.Fatalf("failed to decode rules container: %v", err)
	}
	saSig, err := crypto.SignData(saPriv, rulesData)
	if err != nil {
		t.Fatalf("failed to sign rules container: %v", err)
	}

	addr := &model.WhitelistedAddress{
		ID:         "addr-1",
		Blockchain: "ALGO",
		Network:    "mainnet",
		Metadata: &model.WhitelistedAssetMetadata{
			// Recomputed over the DELIVERED bytes, so step 1 cannot catch this.
			Hash:            crypto.CalculateHexHash(deliveredPayload),
			PayloadAsString: deliveredPayload,
		},
		RulesContainer:  rulesB64,
		RulesSignatures: base64.StdEncoding.EncodeToString([]byte("dummy")),
		SignedAddress: &model.SignedWhitelistedAddress{
			Signatures: []model.WhitelistSignature{
				{
					UserSignature: &model.WhitelistUserSignature{
						UserID:    "user1@bank.com",
						Signature: userSig,
					},
					Hashes: hashes,
				},
			},
		},
	}

	rcDecoder := func(string) (*model.DecodedRulesContainer, error) {
		return &model.DecodedRulesContainer{
			Users: []*model.RuleUser{
				{ID: "user1@bank.com", PublicKey: userPub, Roles: []string{"USER"}},
			},
			Groups: []*model.RuleGroup{
				{ID: "approvers", UserIDs: []string{"user1@bank.com"}},
			},
			AddressWhitelistingRules: []*model.AddressWhitelistingRules{
				{
					Currency: "ALGO",
					Network:  "mainnet",
					ParallelThresholds: []*model.SequentialThresholds{
						{
							Thresholds: []*model.GroupThreshold{
								{GroupID: "approvers", MinimumSignatures: 1},
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

	return &addressInjectionFixture{
		verifier:  NewWhitelistedAddressVerifier([]*ecdsa.PublicKey{saPub}, 1),
		addr:      addr,
		rcDecoder: rcDecoder,
		usDecoder: usDecoder,
	}
}

// genuine payload shape, matching the captured production fixture: the top-level label is
// followed by other members, and linkedInternalAddresses is empty.
const injectionSignedPayload = `{"currency":"ALGO","network":"mainnet","address":"ALGOADDR",` +
	`"memo":"","label":"TN_Bank treasury","customerId":"cust-1","linkedInternalAddresses":[],` +
	`"tnParticipantID":"participant-7"}`

// TestLegacyStripDoesNotLetAnAppendedLabelReachTheCaller pins the primary attack: a duplicate
// top-level `label` appended as the LAST member matches `,"label":"[^"]*"}` and is stripped back
// out, so strategy 2 recovers the signed bytes and steps 1-5 all pass on genuine signatures.
func TestLegacyStripDoesNotLetAnAppendedLabelReachTheCaller(t *testing.T) {
	const attackerLabel = "Coinbase Prime custody"

	// Insert the duplicate immediately before the closing brace.
	delivered := strings.TrimSuffix(injectionSignedPayload, "}") +
		`,"label":"` + attackerLabel + `"}`
	if delivered == injectionSignedPayload {
		t.Fatal("test setup: delivered payload is unchanged")
	}

	// Setup check: the strip really does recover the signed bytes, i.e. this test exercises the
	// path it claims to. Without this the test could pass because verification failed.
	var recovered bool
	for _, v := range ComputeLegacyPayloadVariants(delivered) {
		if v.Payload == injectionSignedPayload {
			recovered = true
		}
	}
	if !recovered {
		t.Fatal("test setup: no legacy variant recovers the signed payload, so step 4 would " +
			"reject for the wrong reason and this test would prove nothing")
	}

	f := buildAddressInjectionFixture(t, injectionSignedPayload, delivered)

	result, err := f.verifier.VerifyWhitelistedAddress(f.addr, f.rcDecoder, f.usDecoder)
	if err != nil {
		t.Fatalf("verification should succeed on genuine signatures: %v", err)
	}

	if result.VerifiedAddress.Label == attackerLabel {
		t.Errorf("attacker-appended label reached the caller as verified: %q",
			result.VerifiedAddress.Label)
	}
	if result.VerifiedAddress.Label != "TN_Bank treasury" {
		t.Errorf("Label = %q, want the signed value %q",
			result.VerifiedAddress.Label, "TN_Bank treasury")
	}
	if result.VerifiedPayload != injectionSignedPayload {
		t.Errorf("VerifiedPayload is not the payload the signature covered:\n got  %q\n want %q",
			result.VerifiedPayload, injectionSignedPayload)
	}
	// The delivered payload is left intact on the model, because a caller needs it to
	// reproduce metadata.hash. It still carries the attacker's member -- which is the point:
	// what changed is which string step 6 parses, not what the server sent.
	if !strings.Contains(f.addr.Metadata.PayloadAsString, attackerLabel) {
		t.Error("test setup: delivered payload no longer carries the injected label")
	}
}

// TestLegacyStripDoesNotLetAnAppendedContractTypeReachTheCaller is the shape duplicate-key
// rejection CANNOT catch: the signed payload has no `contractType`, so the appended member is
// not a duplicate of anything. Only parsing the matched variant closes it.
func TestLegacyStripDoesNotLetAnAppendedContractTypeReachTheCaller(t *testing.T) {
	const attackerContractType = "CMTA20"

	if strings.Contains(injectionSignedPayload, "contractType") {
		t.Fatal("test setup: the signed payload must NOT carry contractType for this shape")
	}
	delivered := strings.TrimSuffix(injectionSignedPayload, "}") +
		`,"contractType":"` + attackerContractType + `"}`

	// This shape must parse cleanly -- no duplicate key -- or it would be caught by the other
	// defence and this test would not be exercising the variant path.
	if err := rejectDuplicateObjectKeys(delivered); err != nil {
		t.Fatalf("test setup: delivered payload must have no duplicate key, got %v", err)
	}

	f := buildAddressInjectionFixture(t, injectionSignedPayload, delivered)

	result, err := f.verifier.VerifyWhitelistedAddress(f.addr, f.rcDecoder, f.usDecoder)
	if err != nil {
		t.Fatalf("verification should succeed on genuine signatures: %v", err)
	}

	if result.VerifiedAddress.ContractType != "" {
		t.Errorf("attacker-appended contractType reached the caller as verified: %q",
			result.VerifiedAddress.ContractType)
	}
	if result.VerifiedPayload != injectionSignedPayload {
		t.Errorf("VerifiedPayload is not the payload the signature covered:\n got  %q\n want %q",
			result.VerifiedPayload, injectionSignedPayload)
	}
}

// TestLegacyStripStillAcceptsAGenuineLegacyRow guards the other direction: the fix must not
// narrow the tolerance the feature exists for. A row whose signature covers only the
// contractType-free hash must still verify, and its contractType must come back empty because
// no signature ever covered a value for it.
func TestLegacyStripStillAcceptsAGenuineLegacyRow(t *testing.T) {
	// The server adds contractType:"" to a row signed before the field existed. This is the
	// legitimate migration, byte-identical in shape to the attack above.
	delivered := strings.TrimSuffix(injectionSignedPayload, "}") + `,"contractType":""}`

	f := buildAddressInjectionFixture(t, injectionSignedPayload, delivered)

	result, err := f.verifier.VerifyWhitelistedAddress(f.addr, f.rcDecoder, f.usDecoder)
	if err != nil {
		t.Fatalf("a genuine legacy row must still verify: %v", err)
	}
	if result.VerifiedAddress.Address != "ALGOADDR" {
		t.Errorf("Address = %q, want %q", result.VerifiedAddress.Address, "ALGOADDR")
	}
	if result.VerifiedAddress.ContractType != "" {
		t.Errorf("ContractType = %q, want empty: no signature covered a value for it",
			result.VerifiedAddress.ContractType)
	}
}

// TestParseWhitelistedAddressRejectsDuplicateKeys pins the second defence directly, for the
// shapes the strip does not reach.
func TestParseWhitelistedAddressRejectsDuplicateKeys(t *testing.T) {
	dup := `{"currency":"ALGO","address":"A","label":"one","label":"two"}`

	if _, err := ParseWhitelistedAddressFromJSON(dup); err == nil {
		t.Fatal("expected a duplicate-key rejection; encoding/json would silently keep \"two\"")
	}

	// And the honest payload still parses.
	if _, err := ParseWhitelistedAddressFromJSON(injectionSignedPayload); err != nil {
		t.Errorf("a well-formed payload must still parse: %v", err)
	}
}

// TestParseWhitelistedAddressRejectsDuplicateKeysInNestedObjects covers the nested case: the
// per-object key set has to be per object, not global, or `id` appearing once in each linked
// address would read as a duplicate.
func TestParseWhitelistedAddressRejectsDuplicateKeysInNestedObjects(t *testing.T) {
	nestedDup := `{"currency":"ALGO","linkedInternalAddresses":[{"id":"1","label":"a","label":"b"}]}`
	if _, err := ParseWhitelistedAddressFromJSON(nestedDup); err == nil {
		t.Error("expected a duplicate-key rejection inside a nested object")
	}

	// Same key in two SIBLING objects is not a duplicate.
	siblings := `{"currency":"ALGO","linkedInternalAddresses":[{"id":"1","label":"a"},{"id":"2","label":"b"}]}`
	if _, err := ParseWhitelistedAddressFromJSON(siblings); err != nil {
		t.Errorf("repeated keys in sibling objects are not duplicates: %v", err)
	}
}
