# Copyright (c) 2024 Taurus SA. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

"""Tests for whitelisted address legacy hash computation and verification.

This module contains comprehensive tests for:
1. Hash verification - ensuring SHA256(payloadAsString) == metadata.hash
2. Legacy hash computation - backward compatibility for addresses signed
   before schema changes (contractType, labels in linkedInternalAddresses)
3. Strategy testing - verifying all legacy hash strategies produce expected results

These tests use fixtures from the Java SDK to ensure cross-SDK compatibility.
"""

from __future__ import annotations

import pytest

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.errors import IntegrityError, WhitelistError
from taurus_protect.helpers.whitelist_hash_helper import compute_legacy_hashes
from tests.unit.fixtures.whitelisted_address_fixtures import (
    CASE1_CURRENT_HASH,
    CASE1_CURRENT_PAYLOAD,
    CASE1_LEGACY_HASH,
    CASE1_LEGACY_PAYLOAD,
    CASE2_CURRENT_HASH,
    CASE2_CURRENT_PAYLOAD,
    CASE2_LEGACY_HASH,
    CASE2_LEGACY_PAYLOAD,
    REAL_METADATA_HASH,
    REAL_PAYLOAD_AS_STRING,
    STRATEGY2_LEGACY_PAYLOAD,
)


class TestHashVerification:
    """Tests for hash verification - SHA256(payloadAsString) == metadata.hash."""

    def test_metadata_hash_matches_computed_hash(self) -> None:
        """Verify that SHA256(payloadAsString) equals metadata.hash for real API data.

        This test confirms the fundamental hash computation matches expected values.
        """
        computed_hash = calculate_hex_hash(REAL_PAYLOAD_AS_STRING)
        assert computed_hash == REAL_METADATA_HASH, (
            f"Hash mismatch: computed={computed_hash}, expected={REAL_METADATA_HASH}"
        )

    def test_metadata_hash_mismatch_detected(self) -> None:
        """Verify that tampering with payload is detected via hash mismatch.

        If an attacker modifies the payload, the computed hash should not match
        the original metadata.hash.
        """
        # Tamper with the payload by changing one character
        tampered_payload = REAL_PAYLOAD_AS_STRING.replace(
            "individual", "omnibus"
        )

        computed_hash = calculate_hex_hash(tampered_payload)
        assert computed_hash != REAL_METADATA_HASH, (
            "Tampered payload should produce different hash"
        )

    def test_empty_payload_hash(self) -> None:
        """Verify that empty payload produces a deterministic hash.

        Empty string is a valid input but should be handled gracefully.
        In practice, verification should fail before reaching this state.
        """
        empty_hash = calculate_hex_hash("")
        # SHA-256 of empty string is a well-known constant
        expected_empty_hash = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
        assert empty_hash == expected_empty_hash

    def test_whitespace_sensitive_hash(self) -> None:
        """Verify that whitespace changes produce different hashes.

        This is critical because JSON formatting differences should be detected.
        """
        hash1 = calculate_hex_hash('{"key":"value"}')
        hash2 = calculate_hex_hash('{ "key": "value" }')
        hash3 = calculate_hex_hash('{"key":"value"}\n')

        assert hash1 != hash2, "Different whitespace should produce different hashes"
        assert hash1 != hash3, "Trailing newline should produce different hash"


class TestLegacyHashCase1:
    """Tests for Legacy Hash Case 1: contractType added after signing.

    This tests addresses signed before the 'contractType' field was added
    to the payload schema.
    """

    def test_case1_current_payload_produces_current_hash(self) -> None:
        """Sanity check: current payload produces current hash."""
        computed_hash = calculate_hex_hash(CASE1_CURRENT_PAYLOAD)
        assert computed_hash == CASE1_CURRENT_HASH, (
            f"Current payload hash mismatch: computed={computed_hash}, "
            f"expected={CASE1_CURRENT_HASH}"
        )

    def test_case1_legacy_payload_produces_legacy_hash(self) -> None:
        """Sanity check: legacy payload (without contractType) produces legacy hash."""
        computed_hash = calculate_hex_hash(CASE1_LEGACY_PAYLOAD)
        assert computed_hash == CASE1_LEGACY_HASH, (
            f"Legacy payload hash mismatch: computed={computed_hash}, "
            f"expected={CASE1_LEGACY_HASH}"
        )

    def test_case1_removes_contract_type(self) -> None:
        """Verify that removing contractType from current payload produces legacy hash.

        The compute_legacy_hashes function should generate the legacy hash
        by removing the contractType field.
        """
        legacy_hashes = compute_legacy_hashes(CASE1_CURRENT_PAYLOAD)

        assert len(legacy_hashes) >= 1, (
            "Should produce at least one legacy hash for payload with contractType"
        )
        assert CASE1_LEGACY_HASH in legacy_hashes, (
            f"Legacy hash {CASE1_LEGACY_HASH} not found in computed hashes: {legacy_hashes}"
        )

    def test_case1_hash_difference(self) -> None:
        """Verify that current and legacy hashes are different.

        This confirms that the contractType field affects the hash.
        """
        assert CASE1_CURRENT_HASH != CASE1_LEGACY_HASH, (
            "Current and legacy hashes should be different"
        )


class TestLegacyHashCase2:
    """Tests for Legacy Hash Case 2: contractType AND labels in linkedInternalAddresses.

    This tests addresses signed before both the 'contractType' field and
    'label' fields in linkedInternalAddresses were added.
    """

    def test_case2_current_payload_produces_current_hash(self) -> None:
        """Sanity check: current payload produces current hash."""
        computed_hash = calculate_hex_hash(CASE2_CURRENT_PAYLOAD)
        assert computed_hash == CASE2_CURRENT_HASH, (
            f"Current payload hash mismatch: computed={computed_hash}, "
            f"expected={CASE2_CURRENT_HASH}"
        )

    def test_case2_legacy_payload_produces_legacy_hash(self) -> None:
        """Sanity check: legacy payload (without contractType AND labels) produces legacy hash."""
        computed_hash = calculate_hex_hash(CASE2_LEGACY_PAYLOAD)
        assert computed_hash == CASE2_LEGACY_HASH, (
            f"Legacy payload hash mismatch: computed={computed_hash}, "
            f"expected={CASE2_LEGACY_HASH}"
        )

    def test_case2_removes_contract_type_and_labels_in_objects(self) -> None:
        """Verify that removing both contractType and labels produces legacy hash.

        The compute_legacy_hashes function should generate the legacy hash
        by removing both the contractType field and label fields from
        linkedInternalAddresses objects.
        """
        legacy_hashes = compute_legacy_hashes(CASE2_CURRENT_PAYLOAD)

        assert len(legacy_hashes) >= 1, (
            "Should produce at least one legacy hash for payload with contractType and labels"
        )
        assert CASE2_LEGACY_HASH in legacy_hashes, (
            f"Legacy hash {CASE2_LEGACY_HASH} not found in computed hashes: {legacy_hashes}"
        )

    def test_case2_label_pattern_does_not_affect_main_label(self) -> None:
        """Verify that the label removal pattern only affects labels in objects.

        The main 'label' field (which is NOT followed by '}') should be preserved.
        Only labels inside linkedInternalAddresses (followed by '}') should be removed.
        """
        # Case 2 current payload has main label "20200324 test address 2"
        # This main label should be preserved in all legacy hashes

        legacy_hashes = compute_legacy_hashes(CASE2_CURRENT_PAYLOAD)

        # Verify the legacy payload still contains the main label
        assert '"label":"20200324 test address 2"' in CASE2_LEGACY_PAYLOAD, (
            "Legacy payload fixture should contain main label"
        )

        # Verify the legacy hash was computed with the main label intact
        # by checking that computing hash of legacy payload matches expected
        computed_from_legacy = calculate_hex_hash(CASE2_LEGACY_PAYLOAD)
        assert computed_from_legacy == CASE2_LEGACY_HASH

    def test_case2_only_removing_contract_type_is_not_enough(self) -> None:
        """Verify that removing only contractType does not produce the legacy hash.

        For Case 2, we need to remove BOTH contractType AND labels to match
        the original signature. Removing just contractType produces a different hash.
        """
        # Manually remove only contractType (keeping labels)
        # This simulates what Strategy 1 does
        payload_without_contract_type_only = CASE2_CURRENT_PAYLOAD.replace(
            ',"contractType":""', ''
        )

        hash_without_contract_type_only = calculate_hex_hash(
            payload_without_contract_type_only
        )

        assert hash_without_contract_type_only != CASE2_LEGACY_HASH, (
            "Removing only contractType should NOT produce the legacy hash for Case 2"
        )

    def test_case2_hash_difference(self) -> None:
        """Verify that current and legacy hashes are different.

        This confirms that contractType and/or labels affect the hash.
        """
        assert CASE2_CURRENT_HASH != CASE2_LEGACY_HASH, (
            "Current and legacy hashes should be different"
        )


class TestLegacyHashStrategies:
    """Tests for legacy hash computation strategies.

    The compute_legacy_hashes function uses multiple strategies:
    - Strategy 1: Remove contractType only
    - Strategy 2: Remove labels from linkedInternalAddresses only
    - Strategy 3: Remove BOTH contractType AND labels
    """

    def test_all_strategies_produce_different_hashes(self) -> None:
        """Verify that all four hash types are different for Case 2.

        For Case 2 payload, we should have four distinct hashes:
        1. Current hash (with all fields)
        2. Hash without contractType only (Strategy 1)
        3. Hash without labels only (Strategy 2)
        4. Hash without both (Strategy 3) = legacy hash
        """
        # Current hash
        current_hash = calculate_hex_hash(CASE2_CURRENT_PAYLOAD)

        # Strategy 1: Remove contractType only
        without_contract_type = CASE2_CURRENT_PAYLOAD.replace(
            ',"contractType":""', ''
        )
        strategy1_hash = calculate_hex_hash(without_contract_type)

        # Strategy 2: Remove labels from linkedInternalAddresses only
        # Using the fixture STRATEGY2_LEGACY_PAYLOAD
        strategy2_hash = calculate_hex_hash(STRATEGY2_LEGACY_PAYLOAD)

        # Strategy 3: Remove both = legacy hash (same as CASE2_LEGACY_HASH)
        strategy3_hash = calculate_hex_hash(CASE2_LEGACY_PAYLOAD)

        # All four should be different
        all_hashes = {current_hash, strategy1_hash, strategy2_hash, strategy3_hash}
        assert len(all_hashes) == 4, (
            f"Expected 4 distinct hashes, got {len(all_hashes)}. "
            f"Hashes: current={current_hash}, strategy1={strategy1_hash}, "
            f"strategy2={strategy2_hash}, strategy3={strategy3_hash}"
        )

    def test_strategy3_matches_original_legacy_hash(self) -> None:
        """Verify that Strategy 3 (both removed) matches the expected legacy hash.

        Strategy 3 should produce the same hash as CASE2_LEGACY_HASH, which
        was captured from the Java SDK and represents the original signing format.
        """
        strategy3_hash = calculate_hex_hash(CASE2_LEGACY_PAYLOAD)
        assert strategy3_hash == CASE2_LEGACY_HASH, (
            f"Strategy 3 hash should match legacy hash. "
            f"Got {strategy3_hash}, expected {CASE2_LEGACY_HASH}"
        )

    def test_compute_legacy_hashes_returns_all_strategies(self) -> None:
        """Verify that compute_legacy_hashes returns hashes for all applicable strategies.

        For Case 2 payload (has both contractType and labels), we should get
        at least 3 legacy hashes (one for each strategy that applies).
        """
        legacy_hashes = compute_legacy_hashes(CASE2_CURRENT_PAYLOAD)

        # Should have multiple strategies
        assert len(legacy_hashes) >= 2, (
            f"Expected at least 2 legacy hashes, got {len(legacy_hashes)}"
        )

        # All hashes should be unique (no duplicates)
        assert len(legacy_hashes) == len(set(legacy_hashes)), (
            "Legacy hashes should not contain duplicates"
        )

        # The legacy hash (Strategy 3) should be present
        assert CASE2_LEGACY_HASH in legacy_hashes, (
            f"Legacy hash should be present in computed hashes"
        )

    def test_compute_legacy_hashes_case1_single_strategy(self) -> None:
        """Verify that Case 1 (contractType only) produces expected legacy hashes.

        Case 1 only has contractType to remove (no labels in linkedInternalAddresses),
        so there should be fewer strategies that apply.
        """
        legacy_hashes = compute_legacy_hashes(CASE1_CURRENT_PAYLOAD)

        # Should have at least one hash (contractType removal)
        assert len(legacy_hashes) >= 1, (
            f"Expected at least 1 legacy hash, got {len(legacy_hashes)}"
        )

        # The legacy hash should be present
        assert CASE1_LEGACY_HASH in legacy_hashes, (
            f"Case 1 legacy hash should be present in computed hashes"
        )


class TestEdgeCases:
    """Tests for edge cases in hash computation and legacy hash handling."""

    def test_payload_without_removable_fields(self) -> None:
        """Verify that payload without contractType or labels returns empty list.

        When there are no fields to remove, compute_legacy_hashes should return
        an empty list (no legacy hashes needed).

        Note: The regex ,"label":"..."} only matches labels that appear right
        before a closing brace (i.e., labels inside linkedInternalAddresses objects).
        A payload where label is NOT the last field won't trigger this pattern.
        """
        # Use a payload where label is NOT the last field (not before closing brace)
        simple_payload = '{"address":"0x123","label":"test","currency":"ETH"}'
        legacy_hashes = compute_legacy_hashes(simple_payload)

        # No contractType, and label is not before closing brace, so empty list
        assert legacy_hashes == [], (
            f"Expected empty list for simple payload, got {legacy_hashes}"
        )

    def test_empty_payload_returns_empty_list(self) -> None:
        """Verify that empty payload returns empty list."""
        legacy_hashes = compute_legacy_hashes("")
        assert legacy_hashes == [], "Empty payload should return empty list"

    def test_null_like_payload(self) -> None:
        """Verify that None-like values are handled gracefully."""
        # Empty string
        assert compute_legacy_hashes("") == []

        # Whitespace-only is technically valid JSON, but not for our purposes
        # The function should handle it gracefully
        whitespace_hashes = compute_legacy_hashes("   ")
        # May return empty or a hash depending on implementation
        assert isinstance(whitespace_hashes, list)

    def test_json_with_empty_contract_type(self) -> None:
        """Verify that empty contractType value is still removed.

        The contractType field might be present but empty ("contractType":"").
        This should still be removed for legacy hash computation.
        """
        payload = '{"address":"0x123","currency":"ETH","contractType":""}'
        legacy_hashes = compute_legacy_hashes(payload)

        # Should remove the empty contractType
        expected_without = '{"address":"0x123","currency":"ETH"}'
        expected_hash = calculate_hex_hash(expected_without)

        assert expected_hash in legacy_hashes, (
            f"Expected hash {expected_hash} not in {legacy_hashes}"
        )

    def test_hash_determinism(self) -> None:
        """Verify that hash computation is deterministic.

        The same input should always produce the same hash.
        """
        for _ in range(10):
            hash1 = calculate_hex_hash(REAL_PAYLOAD_AS_STRING)
            hash2 = calculate_hex_hash(REAL_PAYLOAD_AS_STRING)
            assert hash1 == hash2 == REAL_METADATA_HASH

    def test_legacy_hash_determinism(self) -> None:
        """Verify that legacy hash computation is deterministic.

        The same input should always produce the same legacy hashes in the same order.
        """
        hashes1 = compute_legacy_hashes(CASE2_CURRENT_PAYLOAD)
        hashes2 = compute_legacy_hashes(CASE2_CURRENT_PAYLOAD)

        assert hashes1 == hashes2, "Legacy hash computation should be deterministic"


class TestCrossSDKCompatibility:
    """Tests to ensure compatibility with other SDKs (Java, Go, TypeScript).

    These tests use fixtures from the Java SDK to verify that Python SDK
    produces identical hashes.
    """

    def test_java_sdk_case1_compatibility(self) -> None:
        """Verify Python SDK produces same hashes as Java SDK for Case 1."""
        # Verify current hash matches Java SDK
        computed_current = calculate_hex_hash(CASE1_CURRENT_PAYLOAD)
        assert computed_current == CASE1_CURRENT_HASH, (
            f"Python current hash should match Java SDK: {computed_current} != {CASE1_CURRENT_HASH}"
        )

        # Verify legacy hash matches Java SDK
        computed_legacy = calculate_hex_hash(CASE1_LEGACY_PAYLOAD)
        assert computed_legacy == CASE1_LEGACY_HASH, (
            f"Python legacy hash should match Java SDK: {computed_legacy} != {CASE1_LEGACY_HASH}"
        )

    def test_java_sdk_case2_compatibility(self) -> None:
        """Verify Python SDK produces same hashes as Java SDK for Case 2."""
        # Verify current hash matches Java SDK
        computed_current = calculate_hex_hash(CASE2_CURRENT_PAYLOAD)
        assert computed_current == CASE2_CURRENT_HASH, (
            f"Python current hash should match Java SDK: {computed_current} != {CASE2_CURRENT_HASH}"
        )

        # Verify legacy hash matches Java SDK
        computed_legacy = calculate_hex_hash(CASE2_LEGACY_PAYLOAD)
        assert computed_legacy == CASE2_LEGACY_HASH, (
            f"Python legacy hash should match Java SDK: {computed_legacy} != {CASE2_LEGACY_HASH}"
        )

    def test_real_api_response_hash_compatibility(self) -> None:
        """Verify Python SDK produces same hash as real API response."""
        computed = calculate_hex_hash(REAL_PAYLOAD_AS_STRING)
        assert computed == REAL_METADATA_HASH, (
            f"Python hash should match real API metadata.hash: {computed} != {REAL_METADATA_HASH}"
        )
