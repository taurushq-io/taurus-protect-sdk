# Cryptographic Operations -- Go SDK

## Critical Files

- `tpv1.go` -- TPV1 auth, HMAC, ECDSA signing/verification, hash computation

## TPV1 Authentication

All API requests signed using TPV1-HMAC-SHA256, handled by HTTP transport middleware.
- `TPV1Transport` wraps `http.RoundTripper`
- Signs requests automatically before forwarding

## ECDSA Signing & Verification

- P-256 (secp256r1) curve -- validate before use
- Raw r||s signature format (64 bytes)
- `SignData(privateKey, data)` for request approval signing
- `CalculateHexHash(payload)` for SHA-256 hash computation
- `CalculateBase64Signature(privateKey, data)` for base64-encoded signatures

## Explicit Crypto Randomness

```go
// WRONG
ecdsa.Sign(nil, privateKey, hash[:])

// CORRECT
ecdsa.Sign(rand.Reader, privateKey, hash[:])
```
Always use explicit `crypto/rand.Reader`.

## Prevent Compiler Optimization of Secret Wiping

```go
func Wipe(data []byte) {
    for i := range data { data[i] = 0 }
    runtime.KeepAlive(data)
}
```

## Request Body Preservation

Read and restore body before cloning for signing -- prevents breaking retry middleware.
