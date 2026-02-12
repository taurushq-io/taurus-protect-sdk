# Cryptographic Operations -- Python SDK

## Critical Files

- `tpv1.py` -- TPV1Auth class with sign_request() and parse_url()
- `authenticated_rest.py` -- AuthenticatedRESTClient extends OpenAPI's RESTClientObject
- `signing.py` -- ECDSA signing and verification (P-256, raw r||s format)
- `hashing.py` -- SHA-256 hash computation
- `keys.py` -- Key loading and validation

## TPV1 Authentication

- `AuthenticatedRESTClient` intercepts all HTTP requests, adds Authorization header
- `ProtectClient._get_api_client()` replaces `api_client.rest_client` with `AuthenticatedRESTClient`
- Flow: create() -> TPV1Auth -> _get_api_client() -> replace rest_client -> all requests signed

## ECDSA Signing & Verification

- P-256 (secp256r1) curve -- MUST validate curve type before use
- Raw r||s signature format (64 bytes), NOT DER
- Request approval: `crypto.sign_data(private_key, hashes_json.encode())`
- `json.dumps(hashes, separators=(",", ":"))` for compact JSON matching Java GSON

## Key Lifecycle / Secret Cleanup

- `TPV1Auth.close()` wipes secret with `_wipe_bytes()`
- `__del__` safety: initialize `_closed` and `_secret` at TOP of `__init__` BEFORE any potentially-failing operations
- This prevents `AttributeError` in `__del__` when `__init__` fails partway through

## P-256 Curve Validation

```python
from cryptography.hazmat.primitives.asymmetric.ec import SECP256R1

if not isinstance(public_key.curve, SECP256R1):
    raise ValueError("Only P-256 (secp256r1) curve is supported")
```
