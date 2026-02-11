package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// =============================================================================
// Helpers
// =============================================================================

func generateAddrTestKeyPair(t *testing.T) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key pair: %v", err)
	}
	return priv, &priv.PublicKey
}

func buildRulesContainerWithHSM(t *testing.T, hsmPub *ecdsa.PublicKey) *model.DecodedRulesContainer {
	t.Helper()
	return &model.DecodedRulesContainer{
		Users: []*model.RuleUser{
			{
				ID:        "hsm-slot-1",
				PublicKey: hsmPub,
				Roles:     []string{"HSMSLOT"},
			},
		},
	}
}

func signAddress(t *testing.T, privKey *ecdsa.PrivateKey, addressStr string) string {
	t.Helper()
	sig, err := crypto.SignData(privKey, []byte(addressStr))
	if err != nil {
		t.Fatalf("failed to sign address: %v", err)
	}
	return sig
}

// =============================================================================
// VerifyAddressSignature Tests
// =============================================================================

func TestVerifyAddressSignature_NilInput(t *testing.T) {
	t.Run("nil address", func(t *testing.T) {
		err := VerifyAddressSignature(nil, nil)
		if err == nil {
			t.Error("expected error for nil address")
		}
		assertIntegrityError(t, err, "address cannot be nil")
	})

	t.Run("nil rules container", func(t *testing.T) {
		addr := &model.Address{ID: "1", Address: "0xABC", Signature: "sig"}
		err := VerifyAddressSignature(addr, nil)
		if err == nil {
			t.Error("expected error for nil rules container")
		}
		assertIntegrityError(t, err, "rulesContainer cannot be nil")
	})
}

func TestVerifyAddressSignature_MissingFields(t *testing.T) {
	_, hsmPub := generateAddrTestKeyPair(t)
	container := buildRulesContainerWithHSM(t, hsmPub)

	t.Run("empty signature", func(t *testing.T) {
		addr := &model.Address{ID: "1", Address: "0xABC", Signature: ""}
		err := VerifyAddressSignature(addr, container)
		if err == nil {
			t.Error("expected error for empty signature")
		}
		assertIntegrityError(t, err, "has no signature")
	})

	t.Run("empty address string", func(t *testing.T) {
		addr := &model.Address{ID: "1", Address: "", Signature: "somesig"}
		err := VerifyAddressSignature(addr, container)
		if err == nil {
			t.Error("expected error for empty address string")
		}
		assertIntegrityError(t, err, "has no blockchain address")
	})
}

func TestVerifyAddressSignature_NoHSMKey(t *testing.T) {
	container := &model.DecodedRulesContainer{
		Users: []*model.RuleUser{
			{
				ID:    "regular-user",
				Roles: []string{"USER", "OPERATOR"},
			},
		},
	}

	addr := &model.Address{ID: "1", Address: "0xABC", Signature: "somesig"}
	err := VerifyAddressSignature(addr, container)
	if err == nil {
		t.Error("expected error when no HSMSLOT user found")
	}
	assertIntegrityError(t, err, "HSMSLOT")
}

func TestVerifyAddressSignature_ValidSignature(t *testing.T) {
	hsmPriv, hsmPub := generateAddrTestKeyPair(t)
	container := buildRulesContainerWithHSM(t, hsmPub)

	addressStr := "0x1234567890abcdef"
	sig := signAddress(t, hsmPriv, addressStr)

	addr := &model.Address{
		ID:        "addr-1",
		Address:   addressStr,
		Signature: sig,
	}

	err := VerifyAddressSignature(addr, container)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerifyAddressSignature_InvalidSignature(t *testing.T) {
	hsmPriv, hsmPub := generateAddrTestKeyPair(t)
	container := buildRulesContainerWithHSM(t, hsmPub)

	// Sign a different address
	sig := signAddress(t, hsmPriv, "0xDIFFERENT_ADDRESS")

	addr := &model.Address{
		ID:        "addr-1",
		Address:   "0x1234567890abcdef",
		Signature: sig,
	}

	err := VerifyAddressSignature(addr, container)
	if err == nil {
		t.Error("expected error for invalid signature (signed different data)")
	}
}

func TestVerifyAddressSignature_WrongKey(t *testing.T) {
	// Sign with one key, verify with another
	hsmPriv, _ := generateAddrTestKeyPair(t)
	_, otherPub := generateAddrTestKeyPair(t)
	container := buildRulesContainerWithHSM(t, otherPub)

	addressStr := "0x1234567890abcdef"
	sig := signAddress(t, hsmPriv, addressStr)

	addr := &model.Address{
		ID:        "addr-1",
		Address:   addressStr,
		Signature: sig,
	}

	err := VerifyAddressSignature(addr, container)
	if err == nil {
		t.Error("expected error when verifying with wrong key")
	}
}

func TestVerifyAddressSignature_MalformedSignature(t *testing.T) {
	_, hsmPub := generateAddrTestKeyPair(t)
	container := buildRulesContainerWithHSM(t, hsmPub)

	addr := &model.Address{
		ID:        "addr-1",
		Address:   "0x1234567890abcdef",
		Signature: "not-a-valid-base64-signature!!!",
	}

	err := VerifyAddressSignature(addr, container)
	if err == nil {
		t.Error("expected error for malformed signature")
	}
}

func TestVerifyAddressSignature_EmptyRulesContainerUsers(t *testing.T) {
	container := &model.DecodedRulesContainer{
		Users: []*model.RuleUser{},
	}

	addr := &model.Address{ID: "1", Address: "0xABC", Signature: "somesig"}
	err := VerifyAddressSignature(addr, container)
	if err == nil {
		t.Error("expected error when no HSM user in empty users list")
	}
}

// =============================================================================
// VerifyAddressSignatures (batch) Tests
// =============================================================================

func TestVerifyAddressSignatures_EmptyList(t *testing.T) {
	container := &model.DecodedRulesContainer{}
	err := VerifyAddressSignatures(nil, container)
	if err != nil {
		t.Errorf("unexpected error for nil address list: %v", err)
	}

	err = VerifyAddressSignatures([]*model.Address{}, container)
	if err != nil {
		t.Errorf("unexpected error for empty address list: %v", err)
	}
}

func TestVerifyAddressSignatures_MultipleValid(t *testing.T) {
	hsmPriv, hsmPub := generateAddrTestKeyPair(t)
	container := buildRulesContainerWithHSM(t, hsmPub)

	addr1Str := "0xABC"
	addr2Str := "0xDEF"
	sig1 := signAddress(t, hsmPriv, addr1Str)
	sig2 := signAddress(t, hsmPriv, addr2Str)

	addresses := []*model.Address{
		{ID: "1", Address: addr1Str, Signature: sig1},
		{ID: "2", Address: addr2Str, Signature: sig2},
	}

	err := VerifyAddressSignatures(addresses, container)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerifyAddressSignatures_FailsOnFirstInvalid(t *testing.T) {
	hsmPriv, hsmPub := generateAddrTestKeyPair(t)
	container := buildRulesContainerWithHSM(t, hsmPub)

	addr1Str := "0xABC"
	sig1 := signAddress(t, hsmPriv, addr1Str)

	addresses := []*model.Address{
		{ID: "1", Address: addr1Str, Signature: sig1},
		{ID: "2", Address: "0xDEF", Signature: "invalid-sig"},
	}

	err := VerifyAddressSignatures(addresses, container)
	if err == nil {
		t.Error("expected error when one address has invalid signature")
	}
}

// =============================================================================
// Helper
// =============================================================================

func assertIntegrityError(t *testing.T, err error, contains string) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected IntegrityError containing %q, got nil", contains)
	}
	intErr, ok := err.(*model.IntegrityError)
	if !ok {
		t.Fatalf("expected IntegrityError, got %T: %v", err, err)
	}
	if contains != "" && !strings.Contains(intErr.Message, contains) {
		t.Errorf("expected error containing %q, got %q", contains, intErr.Message)
	}
}
