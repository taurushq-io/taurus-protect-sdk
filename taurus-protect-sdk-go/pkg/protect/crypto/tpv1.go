// Package crypto provides cryptographic utilities for the Taurus-PROTECT SDK.
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

// TPV1Auth holds credentials for TPV1-HMAC-SHA256 authentication.
type TPV1Auth struct {
	APIKey    string
	apiSecret []byte
}

// NewTPV1Auth creates a new TPV1Auth with the given API key and hex-encoded secret.
// The secret is wiped from memory on any initialization failure to prevent leakage.
func NewTPV1Auth(apiKey, apiSecretHex string) (*TPV1Auth, error) {
	if apiKey == "" {
		return nil, errors.New("apiKey cannot be empty")
	}
	if apiSecretHex == "" {
		return nil, errors.New("apiSecret cannot be empty")
	}

	secret, err := hex.DecodeString(apiSecretHex)
	if err != nil {
		// Wipe any partially decoded bytes (defense in depth)
		Wipe(secret)
		return nil, fmt.Errorf("invalid apiSecret: must be valid hex encoding: %w", err)
	}

	return &TPV1Auth{
		APIKey:    apiKey,
		apiSecret: secret,
	}, nil
}

// Close securely wipes the API secret from memory.
func (a *TPV1Auth) Close() {
	Wipe(a.apiSecret)
}

// SignRequest signs an HTTP request with TPV1-HMAC-SHA256 authentication.
// It sets the Authorization header on the request.
func (a *TPV1Auth) SignRequest(req *http.Request, body []byte) error {
	nonce := uuid.New().String()
	timestamp := time.Now().UnixMilli()

	host := getHost(req.URL)
	path := req.URL.Path
	query := req.URL.RawQuery
	contentType := req.Header.Get("Content-Type")

	// Normalize content-type for JSON
	if contentType == "application/json" {
		contentType = "application/json; charset=utf-8"
		req.Header.Set("Content-Type", contentType)
	}

	header := CalculateSignedHeader(
		a.APIKey,
		a.apiSecret,
		nonce,
		timestamp,
		req.Method,
		host,
		path,
		query,
		contentType,
		string(body),
	)

	req.Header.Set("Authorization", header)
	return nil
}

// CalculateSignedHeader computes the TPV1-HMAC-SHA256 Authorization header value.
func CalculateSignedHeader(apiKey string, apiSecret []byte, nonce string, timestamp int64, method, host, path, query, contentType, body string) string {
	// Build message by joining non-empty parts with spaces
	parts := []string{"TPV1", apiKey, nonce, fmt.Sprintf("%d", timestamp), method, host, path, query, contentType, body}
	var nonEmpty []string
	for _, p := range parts {
		if p != "" {
			nonEmpty = append(nonEmpty, p)
		}
	}
	msg := strings.Join(nonEmpty, " ")

	signature := CalculateBase64HMAC(apiSecret, msg)

	return fmt.Sprintf("TPV1-HMAC-SHA256 ApiKey=%s Nonce=%s Timestamp=%d Signature=%s",
		apiKey, nonce, timestamp, signature)
}

// CalculateBase64HMAC computes HMAC-SHA256 and returns base64-encoded result.
func CalculateBase64HMAC(secret []byte, data string) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// CheckBase64HMAC verifies an HMAC-SHA256 signature in constant time.
func CheckBase64HMAC(secret []byte, data, base64HMAC string) bool {
	expected := CalculateBase64HMAC(secret, data)
	return hmac.Equal([]byte(expected), []byte(base64HMAC))
}

// CalculateHexHash computes SHA256 and returns hex-encoded result.
func CalculateHexHash(data string) string {
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

// getHost returns the host (with port if non-standard) for signing.
func getHost(u *url.URL) string {
	port := u.Port()
	if port == "" {
		return u.Hostname()
	}
	// Omit standard ports
	if (u.Scheme == "http" && port == "80") || (u.Scheme == "https" && port == "443") {
		return u.Hostname()
	}
	return u.Host
}

// Wipe securely zeros a byte slice to remove sensitive data from memory.
// Uses runtime.KeepAlive to prevent the compiler from optimizing away the wipe operation.
func Wipe(data []byte) {
	for i := range data {
		data[i] = 0
	}
	// Ensure the wipe operation is not optimized away by the compiler
	runtime.KeepAlive(data)
}

// DecodePublicKeyPEM decodes a PEM-encoded ECDSA public key.
// Only P-256 (secp256r1) curve is supported for security reasons.
func DecodePublicKeyPEM(pemData string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not an ECDSA public key")
	}

	// Validate curve is P-256 (secp256r1) - reject weaker curves
	if ecdsaPub.Curve != elliptic.P256() {
		curveName := "unknown"
		if ecdsaPub.Curve != nil && ecdsaPub.Curve.Params() != nil {
			curveName = ecdsaPub.Curve.Params().Name
		}
		return nil, fmt.Errorf("only P-256 curve is supported, got %s", curveName)
	}

	return ecdsaPub, nil
}

// DecodePrivateKeyPEM decodes a PEM-encoded ECDSA private key.
func DecodePrivateKeyPEM(pemData string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		// Try PKCS8 format
		pkcs8Key, err2 := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err2 != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		ecdsaKey, ok := pkcs8Key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.New("not an ECDSA private key")
		}
		return ecdsaKey, nil
	}

	return key, nil
}

// VerifySignature verifies an ECDSA signature (SHA256withPLAIN-ECDSA format).
// The signature is expected to be base64-encoded raw r||s format.
func VerifySignature(publicKey *ecdsa.PublicKey, data []byte, base64Signature string) (bool, error) {
	// Validate P-256 curve for defense-in-depth
	if publicKey.Curve != elliptic.P256() {
		return false, fmt.Errorf("only P-256 curve is supported for signature verification")
	}

	sig, err := base64.StdEncoding.DecodeString(base64Signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature: %w", err)
	}

	// Plain ECDSA signature is r||s concatenated
	keySize := (publicKey.Curve.Params().BitSize + 7) / 8
	if len(sig) != 2*keySize {
		return false, fmt.Errorf("invalid signature length: expected %d, got %d", 2*keySize, len(sig))
	}

	r := new(big.Int).SetBytes(sig[:keySize])
	s := new(big.Int).SetBytes(sig[keySize:])

	hash := sha256.Sum256(data)
	return ecdsa.Verify(publicKey, hash[:], r, s), nil
}

// SignData signs data using ECDSA with SHA256 (plain format).
// Returns base64-encoded raw r||s signature.
func SignData(privateKey *ecdsa.PrivateKey, data []byte) (string, error) {
	// Validate P-256 curve for defense-in-depth
	if privateKey.Curve != elliptic.P256() {
		return "", fmt.Errorf("only P-256 curve is supported for signing")
	}

	hash := sha256.Sum256(data)
	// Use explicit crypto/rand.Reader for randomness source
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign: %w", err)
	}

	keySize := (privateKey.Curve.Params().BitSize + 7) / 8
	sig := make([]byte, 2*keySize)

	rBytes := r.Bytes()
	sBytes := s.Bytes()

	// Pad r and s to keySize bytes
	copy(sig[keySize-len(rBytes):keySize], rBytes)
	copy(sig[2*keySize-len(sBytes):], sBytes)

	return base64.StdEncoding.EncodeToString(sig), nil
}
