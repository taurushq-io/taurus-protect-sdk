# Whitelisted Asset Integrity Verification — Python SDK

Whitelisted assets (contract addresses) are verified with the same chain as whitelisted
addresses, against `ContractAddressWhitelistingRules` instead of `AddressWhitelistingRules`.
See [`docs/INTEGRITY_VERIFICATION.md`](../../docs/INTEGRITY_VERIFICATION.md) for the
cross-SDK description; this file records how it is implemented here.

## Overview

Every read path verifies. There is no unverified accessor: `get`, `list` and the envelope
getter all run the full pipeline, so a caller cannot hold asset data that was never checked.

```
  ┌──────────────────────────────────────────────────────────────────────┐
  │ STEP 1  SHA-256(payloadAsString) == metadata.hash   (constant-time)  │
  ├──────────────────────────────────────────────────────────────────────┤
  │ STEP 2  SuperAdmin signatures on rulesContainer >= minValidSignatures│
  │         counted by DISTINCT signing key, never by signature entry    │
  ├──────────────────────────────────────────────────────────────────────┤
  │ STEP 3  Decode the rules container (base64 -> model)                 │
  ├──────────────────────────────────────────────────────────────────────┤
  │ STEP 4  metadata.hash — or a LEGACY hash — covered by a signature     │
  │         returns the hash it matched                                  │
  ├──────────────────────────────────────────────────────────────────────┤
  │ STEP 5  User signatures meet the governance thresholds               │
  │         checked against the hash STEP 4 returned                     │
  ├──────────────────────────────────────────────────────────────────────┤
  │ STEP 6  Parse the WhitelistedAsset from the VERIFIED payload         │
  └──────────────────────────────────────────────────────────────────────┘
```

## Where it lives

| Concern | Here |
|---|---|
| Verifier | `taurus_protect/helpers/whitelisted_asset_verifier.py` |
| Service | `taurus_protect/services/whitelisted_asset_service.py` |
| Asset legacy hashes | `compute_asset_legacy_hashes` |
| Parse from verified payload | `_map_asset_from_dto (payload-only)` |
| Hash coverage | `verify_hash_coverage` |
| Rule-key resolution | `resolve_rule_key` |
| Constant-time compare | `hmac.compare_digest` |
| Errors | `IntegrityError` / `WhitelistError` |

## Step 4: legacy hashes, and why the matched hash is returned

An asset signed before a schema change is covered by a LEGACY hash, not by SHA-256 of today's
payload. Three strategies are tried in order, de-duplicated:

1. remove `isNFT`
2. remove `kindType`
3. remove both

Step 4 returns the hash that was actually covered, and step 5 checks signatures against *that*
hash. A step 4 that returns nothing and lets step 5 re-read `metadata.hash` makes a
legacy-signed asset pass step 4 and then fail step 5.

**Assets and addresses use DIFFERENT legacy-hash functions.** `compute_asset_legacy_hashes` is not the address
one. Do not merge them.

## Step 5: which rules apply

`ContractAddressWhitelistingRules` uses `parallelThresholds` directly — a list of
`SequentialThresholds`, not `GroupThreshold`. There are no rule *lines* to narrow the
thresholds, unlike the address flow.

The `(blockchain, network)` pair that selects the rules comes from the **signed payload**, via
`resolve_rule_key`. A payload omitting the chain is an error, never a wildcard; a payload disagreeing
with the DTO is an error. The network falls back to the DTO only when the payload carries none
(governance rules carry a per-rule `includeNetworkInPayload` flag).

A populated group with `minimumSignatures = 0` is rejected as malformed — see the shared doc.

## Field sourcing

`name`, `symbol`, `contractAddress`, `blockchain`, `network` and `decimals` come **only** from
the verified payload. If the payload omits one, the result is empty for that field — never
back-filled from the DTO. Non-security fields (`id`, `status`) may come from the DTO.

## Usage

```
client.whitelisted_assets.get(asset_id)
client.whitelisted_assets.list(limit=50)
```

`list` verifies every row. A row that fails is **excluded and reported**, not fatal — one bad
row must not deny access to every good one, and listing is how an operator finds it. Excluding
stays fail-closed: an omitted asset cannot be selected. If rows came back but none survived,
that is an error, so a filtered page never reads as an empty whitelist.

## Related

- [Integrity Verification](../../docs/INTEGRITY_VERIFICATION.md) — cross-SDK flows
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) — the 6-step peer
- [Concepts](CONCEPTS.md) — model types
