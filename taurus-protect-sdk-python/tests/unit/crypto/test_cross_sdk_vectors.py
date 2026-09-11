"""Cross-SDK cryptographic test vectors.

These tests verify that Python SDK cryptographic functions produce
identical results to the Java, Go, and TypeScript SDKs. All SDKs
read the same test vectors from docs/test-vectors/crypto-test-vectors.json.
"""

from __future__ import annotations

import json
import os
from typing import Any, Dict, List

import pytest

from taurus_protect.crypto.hashing import calculate_hex_hash, constant_time_compare
from taurus_protect.helpers.whitelist_hash_helper import (
    compute_asset_legacy_hashes,
    compute_legacy_hashes,
)

VECTORS_PATH = os.path.join(
    os.path.dirname(__file__),
    "..", "..", "..", "..", "docs", "test-vectors", "crypto-test-vectors.json",
)


@pytest.fixture(scope="module")
def vectors() -> Dict[str, Any]:
    with open(VECTORS_PATH) as f:
        return json.load(f)["vectors"]


class TestCrossSdkHexHash:
    """Verify SHA-256 hex hash matches all SDKs."""

    def test_hex_hash_vectors(self, vectors: Dict[str, Any]) -> None:
        for vec in vectors["hex_hash"]:
            result = calculate_hex_hash(vec["input"])
            assert result == vec["expected"], (
                f"SHA-256 mismatch for {vec['description']}: "
                f"got {result}, expected {vec['expected']}"
            )


class TestCrossSdkConstantTimeCompare:
    """Verify constant-time comparison matches all SDKs."""

    def test_constant_time_compare_vectors(self, vectors: Dict[str, Any]) -> None:
        for vec in vectors["constant_time_compare"]:
            result = constant_time_compare(vec["a"], vec["b"])
            assert result == vec["expected"], (
                f"Constant-time compare mismatch for {vec['description']}: "
                f"got {result}, expected {vec['expected']}"
            )


class TestCrossSdkLegacyAddressHash:
    """Verify legacy address hash computation matches all SDKs."""

    def test_legacy_address_hashes(self, vectors: Dict[str, Any]) -> None:
        for vec in vectors["legacy_hash_address"]:
            legacy_hashes = compute_legacy_hashes(vec["payload"])
            expected_count = vec["expected_legacy_count"]

            assert len(legacy_hashes) == expected_count, (
                f"Legacy hash count mismatch for {vec['description']}: "
                f"got {len(legacy_hashes)}, expected {expected_count}"
            )

            if expected_count > 0:
                assert vec["expected_without_contract_type"] in legacy_hashes, (
                    f"Missing without_contract_type hash for {vec['description']}"
                )
                assert vec["expected_without_labels"] in legacy_hashes, (
                    f"Missing without_labels hash for {vec['description']}"
                )
                assert vec["expected_without_both"] in legacy_hashes, (
                    f"Missing without_both hash for {vec['description']}"
                )

    def test_legacy_address_original_hash(self, vectors: Dict[str, Any]) -> None:
        """Verify original hashes match across SDKs."""
        for vec in vectors["legacy_hash_address"]:
            result = calculate_hex_hash(vec["payload"])
            assert result == vec["original_hash"], (
                f"Original hash mismatch for {vec['description']}: "
                f"got {result}, expected {vec['original_hash']}"
            )


class TestCrossSdkLegacyAssetHash:
    """Verify legacy asset hash computation matches all SDKs."""

    def test_legacy_asset_hashes(self, vectors: Dict[str, Any]) -> None:
        for vec in vectors["legacy_hash_asset"]:
            legacy_hashes = compute_asset_legacy_hashes(vec["payload"])
            expected_count = vec["expected_legacy_count"]

            assert len(legacy_hashes) == expected_count, (
                f"Legacy hash count mismatch for {vec['description']}: "
                f"got {len(legacy_hashes)}, expected {expected_count}"
            )

            if expected_count > 0:
                assert vec["expected_without_is_nft"] in legacy_hashes, (
                    f"Missing without_is_nft hash for {vec['description']}"
                )
                assert vec["expected_without_kind_type"] in legacy_hashes, (
                    f"Missing without_kind_type hash for {vec['description']}"
                )
                assert vec["expected_without_both"] in legacy_hashes, (
                    f"Missing without_both hash for {vec['description']}"
                )

    def test_legacy_asset_original_hash(self, vectors: Dict[str, Any]) -> None:
        """Verify original hashes match across SDKs."""
        for vec in vectors["legacy_hash_asset"]:
            result = calculate_hex_hash(vec["payload"])
            assert result == vec["original_hash"], (
                f"Original hash mismatch for {vec['description']}: "
                f"got {result}, expected {vec['original_hash']}"
            )


class TestCrossSdkCanonicalString:
    """The TPV1 canonical string, pinned by value across all four SDKs.

    **Nothing pinned this before 2026-09-10, and that is how a real interop break shipped.**
    The ``hmac_sha256`` group HMACs a hardcoded string that merely *looks* like a canonical
    message, and no consumer routed through the signing path -- so Python and TypeScript
    could upper-case the HTTP method while Java and Go signed it verbatim, leaving a caller
    who issued a lowercase ``get`` unable to authenticate against one of the two families.

    Asserting the SIGNATURE is what pins the MESSAGE: the secret is fixed, so a signature
    match means the exact byte string was signed. Reconstructing the message in the test
    would assert a copy of the implementation instead -- the trap this repo records.

    This section was consumed by the Go suite alone until now; the other three had the code
    fix and no gate, which is the same asymmetry that let the split survive in the first
    place.
    """

    def test_canonical_string_vectors(self, vectors: Dict[str, Any]) -> None:
        import time
        import uuid
        from unittest.mock import patch

        from taurus_protect.crypto.tpv1 import TPV1Auth

        group = vectors["canonical_string"]
        cases = group["cases"]
        assert len(cases) == group["count"], (
            f"canonical_string: {len(cases)} cases, file declares {group['count']}"
        )

        for vec in cases:
            # Driven through the REAL sign_request, which mints its own nonce and
            # timestamp -- so they are patched rather than the message rebuilt.
            auth = TPV1Auth(vec["api_key"], group["secret_hex"])
            try:
                with patch.object(uuid, "uuid4", return_value=vec["nonce"]), patch.object(
                    time, "time", return_value=vec["timestamp"] / 1000
                ):
                    header = auth.sign_request(
                        vec["method"],
                        vec["host"],
                        vec["path"],
                        vec["query"] or None,
                        vec["content_type"] or None,
                        vec["body"] or None,
                    )
            finally:
                auth.close()

            assert f"Signature={vec['expected_signature']}" in header, (
                f"canonical string mismatch for {vec['description']}: "
                f"got {header}, expected a signature over {vec['expected_message']!r}"
            )
