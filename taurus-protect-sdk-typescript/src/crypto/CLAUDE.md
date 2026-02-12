# Cryptographic Operations — TypeScript SDK

## Critical Files

- `signing.ts` — ECDSA signing and verification (P-256, raw r||s format)
- `tpv1.ts` — TPV1 authentication middleware
- `hashing.ts` — SHA-256 and HMAC functions
- `keys.ts` — Key loading and validation

## TPV1 Authentication

- `createTPV1Middleware()` applied to all requests
- Handles HMAC-SHA256 request signing automatically

## ECDSA Signing & Verification

- P-256 (secp256r1) curve ONLY
- Raw r||s signature format (64 bytes), NOT DER
- Base64 encoding for output
- Must match Java's "SHA256withPLAIN-ECDSA" format
- `signData()` for request approval signing
- `verifySignature()` for signature verification

## MUST Throw for Invalid Curve

Both `signData()` and `verifySignature()` MUST throw `IntegrityError` for non-P-256 curves. Returning `false` for curve errors is a security issue — callers can't distinguish "bad signature" from "bad key."

## Client Error Types

Internal errors (like "No rules container returned from API") throw `ServerError` (not generic `Error`), so callers can catch with `instanceof APIError`.

## Request Hash Verification

Error messages MUST include both computed and provided hash values:
```typescript
throw new IntegrityError(`Request hash mismatch: computed=${computedHash}, provided=${providedHash}`);
```
When hash exists but payloadAsString is undefined, MUST throw IntegrityError.
