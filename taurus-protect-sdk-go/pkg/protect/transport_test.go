package protect

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

func TestTPV1Transport_RoundTrip(t *testing.T) {
	// Create a test server that echoes request headers
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Authorization header not set")
		}
		if !strings.HasPrefix(auth, "TPV1-HMAC-SHA256") {
			t.Errorf("Authorization header should start with TPV1-HMAC-SHA256, got: %s", auth)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	auth, err := crypto.NewTPV1Auth("test-api-key", "deadbeef0123456789abcdef01234567")
	if err != nil {
		t.Fatalf("NewTPV1Auth() error = %v", err)
	}
	defer auth.Close()

	transport := &TPV1Transport{
		Base: http.DefaultTransport,
		Auth: auth,
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestTPV1Transport_RoundTrip_WithBody(t *testing.T) {
	var receivedBody []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedBody, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	auth, err := crypto.NewTPV1Auth("test-api-key", "deadbeef0123456789abcdef01234567")
	if err != nil {
		t.Fatalf("NewTPV1Auth() error = %v", err)
	}
	defer auth.Close()

	transport := &TPV1Transport{
		Base: http.DefaultTransport,
		Auth: auth,
	}

	client := &http.Client{Transport: transport}

	body := `{"test": "data"}`
	req, _ := http.NewRequest("POST", server.URL+"/test", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Do() error = %v", err)
	}
	defer resp.Body.Close()

	if string(receivedBody) != body {
		t.Errorf("Body = %v, want %v", string(receivedBody), body)
	}
}

func TestTPV1Transport_RoundTrip_DefaultTransport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	auth, err := crypto.NewTPV1Auth("test-api-key", "deadbeef0123456789abcdef01234567")
	if err != nil {
		t.Fatalf("NewTPV1Auth() error = %v", err)
	}
	defer auth.Close()

	// Base is nil, should use http.DefaultTransport
	transport := &TPV1Transport{
		Base: nil,
		Auth: auth,
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestTPV1Transport_RoundTrip_AuthHeaderFormat(t *testing.T) {
	var authHeader string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	auth, err := crypto.NewTPV1Auth("my-api-key", "0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatalf("NewTPV1Auth() error = %v", err)
	}
	defer auth.Close()

	transport := &TPV1Transport{
		Base: http.DefaultTransport,
		Auth: auth,
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Get(server.URL + "/test")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	defer resp.Body.Close()

	// Check header format
	if !strings.Contains(authHeader, "ApiKey=my-api-key") {
		t.Errorf("Authorization header should contain ApiKey=my-api-key, got: %s", authHeader)
	}
	if !strings.Contains(authHeader, "Nonce=") {
		t.Errorf("Authorization header should contain Nonce=, got: %s", authHeader)
	}
	if !strings.Contains(authHeader, "Timestamp=") {
		t.Errorf("Authorization header should contain Timestamp=, got: %s", authHeader)
	}
	if !strings.Contains(authHeader, "Signature=") {
		t.Errorf("Authorization header should contain Signature=, got: %s", authHeader)
	}
}

func TestNewHTTPClient(t *testing.T) {
	auth, err := crypto.NewTPV1Auth("test-api-key", "deadbeef0123456789abcdef01234567")
	if err != nil {
		t.Fatalf("NewTPV1Auth() error = %v", err)
	}
	defer auth.Close()

	baseClient := &http.Client{}
	client := newHTTPClient(auth, baseClient)

	if client == nil {
		t.Fatal("newHTTPClient() returned nil")
	}

	transport, ok := client.Transport.(*TPV1Transport)
	if !ok {
		t.Fatalf("Transport is %T, want *TPV1Transport", client.Transport)
	}

	if transport.Auth != auth {
		t.Error("Transport.Auth should be the provided auth")
	}
}

func TestNewHTTPClient_WithCustomTransport(t *testing.T) {
	auth, err := crypto.NewTPV1Auth("test-api-key", "deadbeef0123456789abcdef01234567")
	if err != nil {
		t.Fatalf("NewTPV1Auth() error = %v", err)
	}
	defer auth.Close()

	customTransport := &http.Transport{
		MaxIdleConns: 100,
	}
	baseClient := &http.Client{
		Transport: customTransport,
	}

	client := newHTTPClient(auth, baseClient)

	transport, ok := client.Transport.(*TPV1Transport)
	if !ok {
		t.Fatalf("Transport is %T, want *TPV1Transport", client.Transport)
	}

	if transport.Base != customTransport {
		t.Error("Transport.Base should be the custom transport")
	}
}
