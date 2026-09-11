package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// recordingAssetServer answers every asset read with the same unverifiable row and
// remembers which paths it was asked for, so a test can assert the approve POST never
// happened rather than only that an error came back.
type recordingAssetServer struct {
	mu    sync.Mutex
	calls []string
}

func (r *recordingAssetServer) approvePosts() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	n := 0
	for _, c := range r.calls {
		if strings.HasPrefix(c, "POST") && strings.Contains(c, "approve") {
			n++
		}
	}
	return n
}

// The row's hash covers a different payload, so step 1 rejects it. That is enough for
// the batch semantics: what matters is that a row this SDK could not verify stops the
// call before anything is signed.
func approveAssetService(t *testing.T, rec *recordingAssetServer) *WhitelistedAssetService {
	t.Helper()

	good := `{"blockchain":"ETH","network":"mainnet","contractAddress":"0xA0b8","name":"USD Coin","symbol":"USDC","decimals":6}`
	row := map[string]any{
		"id": "1",
		"metadata": map[string]any{
			"hash":            crypto.CalculateHexHash(good),
			"payloadAsString": `{"contractAddress":"0xEVIL"}`,
		},
	}
	body, err := json.Marshal(map[string]any{"result": row})
	if err != nil {
		t.Fatal(err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rec.mu.Lock()
		rec.calls = append(rec.calls, req.Method+" "+req.URL.Path)
		rec.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, string(body))
	}))
	t.Cleanup(srv.Close)

	cfg := openapi.NewConfiguration()
	cfg.Servers = []openapi.ServerConfiguration{{URL: srv.URL}}
	cfg.HTTPClient = srv.Client()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return NewWhitelistedAssetServiceWithVerification(openapi.NewAPIClient(cfg), &WhitelistedAssetServiceConfig{
		SuperAdminKeys:     []*ecdsa.PublicKey{&key.PublicKey},
		MinValidSignatures: 1,
	})
}

func approveTestKey(t *testing.T) *ecdsa.PrivateKey {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	return key
}

// reviewedAssetSelection mints the content pin from a verified read, the only producer.
// A hand-built model.WhitelistedAssetApproval{} pins nothing, because its field is unexported.
func reviewedAssetSelection(t *testing.T, idToHash map[string]string) *model.WhitelistedAssetApproval {
	t.Helper()

	result := &model.WhitelistedAssetResult{}
	ids := make([]string, 0, len(idToHash))
	for id, hash := range idToHash {
		result.Assets = append(result.Assets, &model.WhitelistedAsset{
			ID:       id,
			Metadata: &model.WhitelistedAssetMetadata{Hash: hash},
		})
		ids = append(ids, id)
	}

	selection, err := result.Select(ids...)
	if err != nil {
		t.Fatalf("minting the reviewed selection: %v", err)
	}
	return selection
}

// 10A: one signature covers every hash in the batch, so a row that did not verify has to
// abort the whole call. Signing the survivors would tell the approver they approved less
// than they did — and the API takes one signature, so there is no partial submission.
func TestApproveWhitelistedAssetsSignsNothingWhenARowFailsVerification(t *testing.T) {
	rec := &recordingAssetServer{}
	svc := approveAssetService(t, rec)

	err := svc.ApproveWhitelistedAssets(context.Background(),
		reviewedAssetSelection(t, map[string]string{"1": "aaa", "2": "bbb"}),
		approveTestKey(t), "batch approval")
	if err == nil {
		t.Fatal("an unverifiable row must abort the approval")
	}
	if !strings.Contains(err.Error(), "refusing to sign") {
		t.Errorf("error should name the refusal, got: %v", err)
	}
	if n := rec.approvePosts(); n != 0 {
		t.Errorf("nothing may be signed or submitted, got %d approve request(s)", n)
	}
}

func TestApproveWhitelistedAssetsRejectsBadInput(t *testing.T) {
	rec := &recordingAssetServer{}
	svc := approveAssetService(t, rec)
	key := approveTestKey(t)

	one := reviewedAssetSelection(t, map[string]string{"1": "aaa"})

	cases := []struct {
		name      string
		selection *model.WhitelistedAssetApproval
		key       *ecdsa.PrivateKey
		comment   string
		want      string
	}{
		// An empty pin must be refused rather than treated as approve-nothing, or it
		// silently restores the unpinned behaviour the pin exists to remove.
		{"nil selection", nil, key, "c", "selection cannot be empty"},
		{"forged empty selection", &model.WhitelistedAssetApproval{}, key, "c", "selection cannot be empty"},
		{"nil key", one, nil, "c", "privateKey cannot be nil"},
		{"no comment", one, key, "", "comment is required"},
		// The IDs are sorted numerically before signing, so a non-numeric one has no
		// defined position in the signed array.
		{"non-numeric id",
			reviewedAssetSelection(t, map[string]string{"1": "aaa", "abc": "bbb"}),
			key, "c", "not a valid numeric ID"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := svc.ApproveWhitelistedAssets(context.Background(), tc.selection, tc.key, tc.comment)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("got %v, want an error containing %q", err, tc.want)
			}
			if n := rec.approvePosts(); n != 0 {
				t.Errorf("invalid input must not reach the API, got %d approve request(s)", n)
			}
		})
	}
}
