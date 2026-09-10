package protect

import (
	"context"
	"net/http"
	"testing"
)

func TestAPIKeyCredentialsApply(t *testing.T) {
	if _, _, err := APIKeyCredentials("", "deadbeef").apply("https://api.example.test", &http.Client{}); err == nil {
		t.Fatal("empty apiKey must fail")
	}
	if _, _, err := APIKeyCredentials("k", "").apply("https://api.example.test", &http.Client{}); err == nil {
		t.Fatal("empty apiSecret must fail")
	}
	if _, _, err := APIKeyCredentials("k", "not-hex").apply("https://api.example.test", &http.Client{}); err == nil {
		t.Fatal("non-hex apiSecret must fail")
	}

	hc, auth, err := APIKeyCredentials("k", "deadbeef0123456789abcdef01234567").apply("https://api.example.test", &http.Client{})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if auth == nil {
		t.Fatal("API-key apply must return a TPV1Auth to wipe on Close")
	}
	defer auth.Close()
	if _, ok := hc.Transport.(*TPV1Transport); !ok {
		t.Fatalf("transport = %T, want *TPV1Transport", hc.Transport)
	}
}

func TestBearerCredentialsApply(t *testing.T) {
	if _, _, err := BearerTokenProviderCredentials(nil).apply("https://api.example.test", &http.Client{}); err == nil {
		t.Fatal("nil provider must fail")
	}

	hc, auth, err := BearerTokenCredentials("tok").apply("https://api.example.test", &http.Client{})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}
	if auth != nil {
		t.Fatal("bearer apply must not return a TPV1Auth")
	}
	if _, ok := hc.Transport.(*bearerTransport); !ok {
		t.Fatalf("transport = %T, want *bearerTransport", hc.Transport)
	}
}

// An empty static bearer token must be rejected at construction, as in the other three
// SDKs: accepting it sends "Authorization: Bearer " on every request, which surfaces as
// an opaque 401 instead of naming the real cause.
func TestBearerTokenCredentialsRejectsEmptyToken(t *testing.T) {
	if _, _, err := BearerTokenCredentials("").apply("https://api.example.test", &http.Client{}); err == nil {
		t.Fatal("empty bearer token must fail at construction")
	}
	if _, _, err := BearerTokenCredentials("tok").apply("https://api.example.test", &http.Client{}); err != nil {
		t.Fatalf("non-empty bearer token must succeed: %v", err)
	}
}

// A provider is resolved per request, so an empty return can only be caught there. A
// failed refresh must not go out as an empty credential.
func TestBearerTransportRejectsEmptyProviderResult(t *testing.T) {
	hc, _, err := BearerTokenProviderCredentials(
		func(context.Context) (string, error) { return "", nil },
	).apply("https://api.example.test", &http.Client{})
	if err != nil {
		t.Fatalf("apply: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, "https://protect.example.com/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := hc.Transport.RoundTrip(req); err == nil {
		t.Fatal("an empty token from the provider must fail the request")
	}
}
