package protect

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

// staticBearer is a constant BearerTokenProvider for tests.
func staticBearer(token string) BearerTokenProvider {
	return func(context.Context) (string, error) { return token, nil }
}

func testP256Key(t *testing.T) *ecdsa.PublicKey {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	return &priv.PublicKey
}

func TestValidateRequiresHost(t *testing.T) {
	c := &clientConfig{credentials: BearerTokenCredentials("tok")}
	if err := c.validate(); err == nil {
		t.Fatal("validate() with empty host must fail")
	}
}

func TestValidateRequiresCredentials(t *testing.T) {
	c := &clientConfig{host: "https://api.example.com"}
	if err := c.validate(); err == nil {
		t.Fatal("validate() with no credentials must fail")
	}
}

func TestValidateBearerRequiresSuperAdminKeys(t *testing.T) {
	noKeys := &clientConfig{host: "h", credentials: BearerTokenCredentials("tok")}
	if err := noKeys.validate(); err == nil {
		t.Fatal("validate() bearer mode without SuperAdmin keys must fail")
	}
	withKeys := &clientConfig{host: "h", credentials: BearerTokenCredentials("tok"), superAdminKeys: []*ecdsa.PublicKey{testP256Key(t)}, minValidSignatures: 1}
	if err := withKeys.validate(); err != nil {
		t.Fatalf("validate() bearer + keys = %v, want nil", err)
	}
}

func TestValidateAPIKeyRequiresSuperAdminKeys(t *testing.T) {
	c := &clientConfig{host: "h", credentials: APIKeyCredentials("k", "deadbeef")}
	if err := c.validate(); err == nil {
		t.Fatal("validate() API-key mode without SuperAdmin keys must fail")
	}
}

// TestValidateChecksMinSigsAgainstKeys verifies the SuperAdmin-key consistency
// checks run regardless of auth mode when keys ARE supplied.
func TestValidateChecksMinSigsAgainstKeys(t *testing.T) {
	key := testP256Key(t)

	tooMany := &clientConfig{host: "h", credentials: BearerTokenCredentials("tok"), superAdminKeys: []*ecdsa.PublicKey{key}, minValidSignatures: 2}
	if err := tooMany.validate(); err == nil {
		t.Fatal("minValidSignatures > len(superAdminKeys) must fail")
	}

	zeroSigs := &clientConfig{host: "h", credentials: BearerTokenCredentials("tok"), superAdminKeys: []*ecdsa.PublicKey{key}, minValidSignatures: 0}
	if err := zeroSigs.validate(); err == nil {
		t.Fatal("minValidSignatures == 0 with keys must fail")
	}

	ok := &clientConfig{host: "h", credentials: APIKeyCredentials("k", "deadbeef"), superAdminKeys: []*ecdsa.PublicKey{key}, minValidSignatures: 1}
	if err := ok.validate(); err != nil {
		t.Fatalf("validate() apikey+keys = %v, want nil", err)
	}
}
