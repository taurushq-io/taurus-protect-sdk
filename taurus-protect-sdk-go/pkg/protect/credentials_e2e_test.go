package protect

import (
	"context"
	"crypto/ecdsa"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// e2eAuthServer records the Authorization header of the last request and 200s.
func e2eAuthServer(gotAuth *string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
}

// TestE2E_APIKeyCredentialsSignRequests asserts a client built via
// WithCredentials(APIKeyCredentials(...)) signs a real round-trip with a
// TPV1-HMAC Authorization header.
func TestE2E_APIKeyCredentialsSignRequests(t *testing.T) {
	var gotAuth string
	srv := e2eAuthServer(&gotAuth)
	defer srv.Close()

	c, err := NewClient(srv.URL,
		WithCredentials(APIKeyCredentials("test-api-key", "deadbeef0123456789abcdef01234567")),
		WithSuperAdminKeys([]*ecdsa.PublicKey{testP256Key(t)}),
		WithMinValidSignatures(1))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer func() { _ = c.Close() }()

	resp, err := c.HTTPClient().Get(srv.URL + "/v1/health")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	_ = resp.Body.Close()

	if !strings.HasPrefix(gotAuth, "TPV1-HMAC-SHA256") {
		t.Fatalf("Authorization = %q, want a TPV1-HMAC-SHA256 signature", gotAuth)
	}
}

// TestE2E_BearerProviderCredentialsPerUser asserts one client built via
// WithCredentials(BearerTokenProviderCredentials(...)) serves distinct callers —
// each request carries that caller's token read from its context.
func TestE2E_BearerProviderCredentialsPerUser(t *testing.T) {
	var gotAuth string
	srv := e2eAuthServer(&gotAuth)
	defer srv.Close()

	c, err := NewClient(srv.URL,
		WithCredentials(BearerTokenProviderCredentials(func(ctx context.Context) (string, error) {
			token, _ := ctx.Value(ctxTokenKey{}).(string)
			return token, nil
		})),
		WithSuperAdminKeys([]*ecdsa.PublicKey{testP256Key(t)}),
		WithMinValidSignatures(1))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	defer func() { _ = c.Close() }()

	for _, token := range []string{"user-a-jwt", "user-b-jwt"} {
		req, err := http.NewRequestWithContext(
			context.WithValue(context.Background(), ctxTokenKey{}, token),
			http.MethodGet, srv.URL+"/v1/health", nil)
		if err != nil {
			t.Fatal(err)
		}
		resp, err := c.HTTPClient().Do(req)
		if err != nil {
			t.Fatalf("Do: %v", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()

		if want := "Bearer " + token; gotAuth != want {
			t.Fatalf("Authorization = %q, want %q", gotAuth, want)
		}
	}
}
