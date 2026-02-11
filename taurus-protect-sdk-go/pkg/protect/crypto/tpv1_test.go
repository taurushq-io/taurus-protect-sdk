package crypto

import (
	"net/http"
	"strings"
	"testing"
)

func TestNewTPV1Auth(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		apiSecret string
		wantErr   bool
	}{
		{
			name:      "valid credentials",
			apiKey:    "test-api-key",
			apiSecret: "deadbeef",
			wantErr:   false,
		},
		{
			name:      "empty api key",
			apiKey:    "",
			apiSecret: "deadbeef",
			wantErr:   true,
		},
		{
			name:      "empty api secret",
			apiKey:    "test-api-key",
			apiSecret: "",
			wantErr:   true,
		},
		{
			name:      "invalid hex secret",
			apiKey:    "test-api-key",
			apiSecret: "not-hex",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth, err := NewTPV1Auth(tt.apiKey, tt.apiSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTPV1Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && auth == nil {
				t.Error("NewTPV1Auth() returned nil auth without error")
			}
			if auth != nil {
				auth.Close()
			}
		})
	}
}

func TestTPV1Auth_SignRequest(t *testing.T) {
	auth, err := NewTPV1Auth("test-api-key", "deadbeefcafe1234")
	if err != nil {
		t.Fatalf("Failed to create TPV1Auth: %v", err)
	}
	defer auth.Close()

	req, err := http.NewRequest("GET", "https://api.example.com/api/rest/v1/wallets", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	err = auth.SignRequest(req, nil)
	if err != nil {
		t.Errorf("SignRequest() error = %v", err)
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		t.Error("Authorization header not set")
	}

	if !strings.HasPrefix(authHeader, "TPV1-HMAC-SHA256 ") {
		t.Errorf("Authorization header has wrong prefix: %s", authHeader)
	}

	if !strings.Contains(authHeader, "ApiKey=test-api-key") {
		t.Errorf("Authorization header missing ApiKey: %s", authHeader)
	}

	if !strings.Contains(authHeader, "Nonce=") {
		t.Errorf("Authorization header missing Nonce: %s", authHeader)
	}

	if !strings.Contains(authHeader, "Timestamp=") {
		t.Errorf("Authorization header missing Timestamp: %s", authHeader)
	}

	if !strings.Contains(authHeader, "Signature=") {
		t.Errorf("Authorization header missing Signature: %s", authHeader)
	}
}

func TestTPV1Auth_SignRequestWithBody(t *testing.T) {
	auth, err := NewTPV1Auth("test-api-key", "deadbeefcafe1234")
	if err != nil {
		t.Fatalf("Failed to create TPV1Auth: %v", err)
	}
	defer auth.Close()

	body := []byte(`{"name":"test-wallet"}`)
	req, err := http.NewRequest("POST", "https://api.example.com/api/rest/v1/wallets", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	err = auth.SignRequest(req, body)
	if err != nil {
		t.Errorf("SignRequest() error = %v", err)
	}

	// Content-Type should be normalized
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("Content-Type not normalized: %s", contentType)
	}

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		t.Error("Authorization header not set")
	}
}

func TestCalculateSignedHeader(t *testing.T) {
	apiSecret := []byte{0xde, 0xad, 0xbe, 0xef}
	header := CalculateSignedHeader(
		"test-key",
		apiSecret,
		"test-nonce",
		1234567890,
		"GET",
		"api.example.com",
		"/api/v1/test",
		"",
		"",
		"",
	)

	if !strings.HasPrefix(header, "TPV1-HMAC-SHA256 ") {
		t.Errorf("Header has wrong prefix: %s", header)
	}

	if !strings.Contains(header, "ApiKey=test-key") {
		t.Errorf("Header missing ApiKey: %s", header)
	}

	if !strings.Contains(header, "Nonce=test-nonce") {
		t.Errorf("Header missing Nonce: %s", header)
	}

	if !strings.Contains(header, "Timestamp=1234567890") {
		t.Errorf("Header missing Timestamp: %s", header)
	}
}

func TestCalculateBase64HMAC(t *testing.T) {
	secret := []byte("test-secret")
	data := "test-data"

	hmac1 := CalculateBase64HMAC(secret, data)
	hmac2 := CalculateBase64HMAC(secret, data)

	if hmac1 != hmac2 {
		t.Error("HMAC should be deterministic")
	}

	if !CheckBase64HMAC(secret, data, hmac1) {
		t.Error("CheckBase64HMAC should return true for valid HMAC")
	}

	if CheckBase64HMAC(secret, data, "invalid") {
		t.Error("CheckBase64HMAC should return false for invalid HMAC")
	}
}

func TestCalculateHexHash(t *testing.T) {
	hash := CalculateHexHash("test")
	// SHA256 of "test"
	expected := "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	if hash != expected {
		t.Errorf("CalculateHexHash() = %s, want %s", hash, expected)
	}
}

func TestWipe(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	Wipe(data)

	for i, b := range data {
		if b != 0 {
			t.Errorf("Wipe() did not zero byte at index %d: got %d", i, b)
		}
	}
}

func TestWipe_EmptySlice(t *testing.T) {
	var data []byte
	// Should not panic
	Wipe(data)
}

func TestTPV1Auth_Close_WipesSecret(t *testing.T) {
	auth, err := NewTPV1Auth("test-api-key", "deadbeef0123456789abcdef01234567")
	if err != nil {
		t.Fatalf("NewTPV1Auth() error = %v", err)
	}

	// Capture secret length before close
	secretLen := len(auth.apiSecret)
	if secretLen == 0 {
		t.Fatal("apiSecret should not be empty")
	}

	auth.Close()

	// Verify secret is wiped (all zeros)
	for i, b := range auth.apiSecret {
		if b != 0 {
			t.Errorf("apiSecret[%d] = %v after Close(), want 0", i, b)
		}
	}
}

func TestDecodePublicKeyPEM_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		pem     string
		wantErr bool
	}{
		{"empty PEM", "", true},
		{"invalid PEM", "not a pem", true},
		{"invalid PEM header", "-----BEGIN SOMETHING-----\naW52YWxpZA==\n-----END SOMETHING-----", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodePublicKeyPEM(tt.pem)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodePublicKeyPEM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDecodePrivateKeyPEM_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		pem     string
		wantErr bool
	}{
		{"empty PEM", "", true},
		{"invalid PEM", "not a pem", true},
		{"invalid PEM content", "-----BEGIN EC PRIVATE KEY-----\naW52YWxpZA==\n-----END EC PRIVATE KEY-----", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecodePrivateKeyPEM(tt.pem)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodePrivateKeyPEM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewTPV1Auth_ErrorMessages(t *testing.T) {
	tests := []struct {
		name         string
		apiKey       string
		apiSecretHex string
		errContains  string
	}{
		{
			name:         "empty api key",
			apiKey:       "",
			apiSecretHex: "deadbeef",
			errContains:  "apiKey cannot be empty",
		},
		{
			name:         "empty api secret",
			apiKey:       "test-key",
			apiSecretHex: "",
			errContains:  "apiSecret cannot be empty",
		},
		{
			name:         "invalid hex",
			apiKey:       "test-key",
			apiSecretHex: "not-valid-hex!",
			errContains:  "invalid apiSecret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTPV1Auth(tt.apiKey, tt.apiSecretHex)
			if err == nil {
				t.Error("Expected error, got nil")
				return
			}
			if !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("Error = %v, want containing %q", err, tt.errContains)
			}
		})
	}
}
