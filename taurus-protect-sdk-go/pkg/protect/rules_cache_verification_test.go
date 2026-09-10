package protect

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// The rules cache supplies the HSM public key that address signature verification
// trusts, so the container it holds must have its SuperAdmin signatures verified
// before anything reads it. An unverified container makes that check pass against
// whatever key the container carries — attacker-chosen addresses verify clean.
//
// This is a cross-SDK invariant and TypeScript had drifted off it: its cache fetched
// through the generated API and the raw mapper, skipping verification entirely. Nothing
// in any SDK tested the path, which is why it survived. Guard for this SDK.
func TestRulesCacheRefusesAnUnverifiedContainer(t *testing.T) {
	unsigned, err := mapper.RulesContainerToBase64(&model.DecodedRulesContainer{
		MinimumDistinctUserSignatures: 2,
	})
	if err != nil {
		t.Fatalf("encode container: %v", err)
	}

	for _, tc := range []struct {
		name string
		body string
	}{
		{
			// Wire-valid container with NO signatures over it.
			name: "no signatures",
			body: `{"result":{"rulesContainer":"` + unsigned + `","rulesSignatures":[]}}`,
		},
		{
			// Signatures present, produced by a key this client does not trust, so the
			// distinct-signer threshold cannot be met.
			name: "signature from an unconfigured key",
			body: `{"result":{"rulesContainer":"` + unsigned +
				`","rulesSignatures":[{"userId":"attacker","signature":"YmFkLXNpZy1ub3QtZWNkc2E="}]}}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			client, err := NewClient(srv.URL,
				WithCredentials(APIKeyCredentials("test-key", "abcdef0123456789")),
				WithSuperAdminKeysPEM([]string{rulesCacheTestKeyPEM(t)}),
				WithMinValidSignatures(1),
			)
			if err != nil {
				t.Fatalf("NewClient: %v", err)
			}
			defer func() { _ = client.Close() }()

			container, err := client.RulesCache().Get(context.Background())
			if err == nil {
				t.Fatalf("rules cache accepted an unverified container: %+v", container)
			}
			if container != nil {
				t.Errorf("a container was returned alongside the error: %+v", container)
			}
		})
	}
}

func rulesCacheTestKeyPEM(t *testing.T) string {
	t.Helper()
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	der, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("marshal key: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
}
