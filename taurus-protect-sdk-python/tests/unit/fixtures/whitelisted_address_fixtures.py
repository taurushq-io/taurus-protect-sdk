# Copyright (c) 2024 Taurus SA. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

"""
Test fixtures for whitelisted address unit tests.

This module contains test data captured from real API responses and legacy hash
test cases from the Java SDK. These fixtures are used to verify:

1. Hash computation for current payload formats
2. Legacy hash computation for backward compatibility
3. Signature verification against known-good hashes

The legacy hash support is critical for verifying addresses that were signed
before certain payload fields were added (e.g., contractType, tnParticipantID,
labels in linkedInternalAddresses).
"""

# =============================================================================
# REAL API RESPONSE DATA
# =============================================================================
# Captured from integration test - represents actual API response format.
# Address ID: 36663, Blockchain: ALGO (Algorand)

REAL_PAYLOAD_AS_STRING = (
    '{"currency":"ALGO","addressType":"individual",'
    '"address":"P4QCJV2YYLAEULGLJQAW4XTU3EBOHWL5C46I5SPLH2H7AJEE367ZDACV5A",'
    '"memo":"","label":"TN_Bank ACC Cockroach_WTRTest","customerId":"",'
    '"exchangeAccountId":"","linkedInternalAddresses":[],"contractType":"",'
    '"tnParticipantID":"84dc35e3-0af8-4b6b-be75-785f4b149d16"}'
)
"""
Real payload from API response for an Algorand address.

Fields present in current payload format:
- currency: Blockchain identifier (ALGO = Algorand)
- addressType: "individual" or "omnibus"
- address: The actual blockchain address
- memo: Optional memo/tag for chains that support it
- label: Human-readable label
- customerId: Optional customer identifier
- exchangeAccountId: Optional exchange account reference
- linkedInternalAddresses: Array of linked internal addresses
- contractType: Type of contract (empty for non-contract addresses)
- tnParticipantID: Taurus Network participant ID (added later)
"""

REAL_METADATA_HASH = "830063cfa8c1dbd696d670fc8360e85fbc57c3ffa66d22358b9a7d6befabb2f0"
"""
SHA-256 hash of REAL_PAYLOAD_AS_STRING.

This hash is used for signature verification. The hash is computed as:
    sha256(payloadAsString.encode('utf-8')).hexdigest()
"""


# =============================================================================
# LEGACY HASH TEST CASE 1: contractType added after signing
# =============================================================================
# This case tests addresses that were signed before the 'contractType' field
# was added to the payload schema. The signature was computed against the
# legacy payload (without contractType), but the API now returns payloads
# with contractType included.

CASE1_LEGACY_PAYLOAD = (
    '{"currency":"ETH","addressType":"individual",'
    '"address":"0x012566A179a935ACF1d81d4D237495DE933D12E6",'
    '"memo":"","label":"CMTA20-KYC - 0x012566A179a935ACF1d81d4D237495DE933D12E6 (request 6826)",'
    '"customerId":"","exchangeAccountId":"","linkedInternalAddresses":[]}'
)
"""
Legacy payload format - BEFORE contractType field was added.

This represents the exact payload that was used when the address was
originally signed. The signature in the envelope was computed against
this payload (without contractType).
"""

CASE1_CURRENT_PAYLOAD = (
    '{"currency":"ETH","addressType":"individual",'
    '"address":"0x012566A179a935ACF1d81d4D237495DE933D12E6",'
    '"memo":"","label":"CMTA20-KYC - 0x012566A179a935ACF1d81d4D237495DE933D12E6 (request 6826)",'
    '"customerId":"","exchangeAccountId":"","linkedInternalAddresses":[],"contractType":""}'
)
"""
Current payload format - WITH contractType field added.

This is what the API now returns. The contractType field was added to the
schema after this address was signed, so the signature won't match this
payload's hash. The SDK must compute the legacy hash (without contractType)
to verify the signature.
"""

CASE1_LEGACY_HASH = "cda66e821ec26f2432a717feaa1ef49be39a7ad9e93b6b8fcdce606659e964df"
"""
Hash of CASE1_LEGACY_PAYLOAD - the hash that matches the signature.

When verification fails with the current hash, the SDK should try this
legacy hash to verify the signature.
"""

CASE1_CURRENT_HASH = "d95ae4359bea509c2542acf410649f1e361233da5e1ac7c7a198b6d6a2bbbe1f"
"""
Hash of CASE1_CURRENT_PAYLOAD - will NOT match the signature.

This hash is computed from the current payload but won't match because
the address was signed before contractType was added.
"""


# =============================================================================
# LEGACY HASH TEST CASE 2: contractType AND labels in linkedInternalAddresses
# =============================================================================
# This is a more complex case where two fields changed after signing:
# 1. contractType was added to the root payload
# 2. label field was added to each linkedInternalAddresses entry
#
# This requires trying multiple legacy hash strategies:
# - Strategy 1: Remove both contractType AND labels (original signing format)
# - Strategy 2: Remove only labels (intermediate format)

CASE2_CURRENT_PAYLOAD = (
    '{"currency":"ETH","addressType":"individual",'
    '"address":"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380",'
    '"memo":"","label":"20200324 test address 2","customerId":"1",'
    '"exchangeAccountId":"","linkedInternalAddresses":['
    '{"id":"10","address":"0x589ef3d7585f54f0539e24253050887c691c9bd8","label":"client 0 ETH "},'
    '{"id":"13","address":"0x669805f31178faf0dca39c8a5c49ecc531b5156e","label":"ETH internal client 02.02"},'
    '{"id":"20","address":"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d","label":"LBR 07.02"},'
    '{"id":"98","address":"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e","label":"ETH LBR internal client 26.02"},'
    '{"id":"25","address":"0x9bc28e6710f5bb2511372987f613a436618e28ad","label":"LBR IC 13.02"}'
    '],"contractType":""}'
)
"""
Current payload format - WITH contractType AND labels in linkedInternalAddresses.

This is the current API response format. Both contractType (root level) and
label (in each linkedInternalAddresses entry) were added after the original
signing.
"""

CASE2_LEGACY_PAYLOAD = (
    '{"currency":"ETH","addressType":"individual",'
    '"address":"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380",'
    '"memo":"","label":"20200324 test address 2","customerId":"1",'
    '"exchangeAccountId":"","linkedInternalAddresses":['
    '{"id":"10","address":"0x589ef3d7585f54f0539e24253050887c691c9bd8"},'
    '{"id":"13","address":"0x669805f31178faf0dca39c8a5c49ecc531b5156e"},'
    '{"id":"20","address":"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d"},'
    '{"id":"98","address":"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e"},'
    '{"id":"25","address":"0x9bc28e6710f5bb2511372987f613a436618e28ad"}'
    ']}'
)
"""
Legacy payload format - Strategy 1: NO contractType AND NO labels.

This is the original signing format before both contractType and labels
were added. Use this to compute the legacy hash for very old addresses.
"""

CASE2_LEGACY_HASH = "88e4e456f7ca1fc4ca415c6c571f828c0eb047e9f15f36d547c103b2ea0def9b"
"""
Hash of CASE2_LEGACY_PAYLOAD - Strategy 1 (no contractType, no labels).

This is the hash that was used for the original signature when neither
contractType nor labels existed in the payload schema.
"""

CASE2_CURRENT_HASH = "7d62d7f78ed55c716ea1278473d6cac5b60a31e1df941873118932822df32b03"
"""
Hash of CASE2_CURRENT_PAYLOAD - will NOT match the signature.

This hash is computed from the current payload format with all fields present.
"""

# Strategy 2: Labels removed only (keep contractType)
STRATEGY2_LEGACY_PAYLOAD = (
    '{"currency":"ETH","addressType":"individual",'
    '"address":"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380",'
    '"memo":"","label":"20200324 test address 2","customerId":"1",'
    '"exchangeAccountId":"","linkedInternalAddresses":['
    '{"id":"10","address":"0x589ef3d7585f54f0539e24253050887c691c9bd8"},'
    '{"id":"13","address":"0x669805f31178faf0dca39c8a5c49ecc531b5156e"},'
    '{"id":"20","address":"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d"},'
    '{"id":"98","address":"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e"},'
    '{"id":"25","address":"0x9bc28e6710f5bb2511372987f613a436618e28ad"}'
    '],"contractType":""}'
)
"""
Intermediate legacy payload format - Strategy 2: WITH contractType, NO labels.

This format represents an intermediate schema version where contractType
was present but labels in linkedInternalAddresses were not yet added.
Some addresses may have been signed with this intermediate format.

Note: The hash for this payload is not provided in the Java SDK test cases,
but this payload is included for completeness and testing intermediate
legacy hash strategies.
"""


# =============================================================================
# HELPER FUNCTIONS FOR TESTS
# =============================================================================

def get_all_legacy_test_cases():
    """
    Return all legacy hash test cases as a list of tuples.

    Each tuple contains:
        (case_name, current_payload, legacy_payload, legacy_hash, current_hash)

    Returns:
        List of test case tuples for parameterized testing.
    """
    return [
        (
            "case1_contractType_added",
            CASE1_CURRENT_PAYLOAD,
            CASE1_LEGACY_PAYLOAD,
            CASE1_LEGACY_HASH,
            CASE1_CURRENT_HASH,
        ),
        (
            "case2_contractType_and_labels_added",
            CASE2_CURRENT_PAYLOAD,
            CASE2_LEGACY_PAYLOAD,
            CASE2_LEGACY_HASH,
            CASE2_CURRENT_HASH,
        ),
    ]
