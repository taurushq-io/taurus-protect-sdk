"""Whitelisted address verification utilities.

This module implements the 6-step verification flow for whitelisted addresses:
1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
2. Verify rules container signatures (SuperAdmin signatures)
3. Decode rules container (base64 -> protobuf -> model)
4. Verify hash coverage (metadata.hash in at least one signature.hashes)
5. Verify whitelist signatures (user signatures meet governance thresholds)
6. Parse WhitelistedAddress from verified payload
"""

from __future__ import annotations

import base64
import binascii
import hmac
import json
from dataclasses import dataclass
from typing import TYPE_CHECKING, Callable, List, Optional

from cryptography.exceptions import InvalidSignature

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.crypto.signing import verify_signature
from taurus_protect.errors import IntegrityError, WhitelistError
from taurus_protect.helpers.constant_time import constant_time_compare
from taurus_protect.helpers.signature_verifier import is_valid_signature
from taurus_protect.helpers.whitelist_hash_helper import (
    compute_legacy_hashes,
    parse_whitelisted_address_from_json,
)
from taurus_protect.models.governance_rules import (
    RULE_SOURCE_TYPE_INTERNAL_WALLET,
    AddressWhitelistingLine,
    AddressWhitelistingRules,
    DecodedRulesContainer,
    RuleUserSignature,
    SequentialThresholds,
)
from taurus_protect.models.whitelisted_address import (
    SignedWhitelistedAddressEnvelope,
    WhitelistedAddress,
    WhitelistSignatureEntry,
)

if TYPE_CHECKING:
    from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey


@dataclass
class AddressVerificationResult:
    """Result of 6-step whitelisted address verification."""

    rules_container: DecodedRulesContainer
    verified_hash: str
    verified_whitelisted_address: WhitelistedAddress


class WhitelistedAddressVerifier:
    """
    Verifier for whitelisted addresses.

    Performs the complete 6-step cryptographic verification to ensure
    address integrity and proper approval according to governance rules:
    1. Verify metadata hash
    2. Verify rules container signatures
    3. Decode rules container
    4. Verify hash coverage
    5. Verify whitelist signatures
    6. Parse WhitelistedAddress from verified payload
    """

    def __init__(
        self,
        super_admin_keys: List["EllipticCurvePublicKey"],
        min_valid_signatures: int = 1,
    ) -> None:
        """
        Initialize the verifier.

        Args:
            super_admin_keys: List of SuperAdmin public keys for verification.
            min_valid_signatures: Minimum number of valid signatures required.
        """
        self._super_admin_keys = super_admin_keys
        self._min_valid_signatures = min_valid_signatures

    def verify_whitelisted_address(
        self,
        envelope: SignedWhitelistedAddressEnvelope,
        rules_container_decoder: Callable[[str], DecodedRulesContainer],
        user_signatures_decoder: Callable[[str], List[RuleUserSignature]],
        cached_rules_container: Optional[DecodedRulesContainer] = None,
    ) -> AddressVerificationResult:
        """
        Perform the complete 6-step verification of a whitelisted address.

        Args:
            envelope: The signed envelope to verify.
            rules_container_decoder: Function to decode base64 rules container.
            user_signatures_decoder: Function to decode base64 user signatures.
            cached_rules_container: Pre-verified and decoded rules container.
                When provided, steps 2-3 are skipped (already done during cache building).

        Returns:
            Verification result with decoded rules container and verified hash.

        Raises:
            IntegrityError: If verification fails at any step.
            WhitelistError: If signature thresholds are not met.
        """
        if envelope is None:
            raise ValueError("envelope cannot be None")
        if envelope.metadata is None:
            raise ValueError("metadata cannot be None")

        # Step 1: Verify metadata hash
        self._verify_metadata_hash(envelope)

        if cached_rules_container is not None:
            # Steps 2-3 already done during cache building
            rules_container = cached_rules_container
        else:
            # Step 2: Verify rules container signatures
            self._verify_rules_container_signatures(envelope, user_signatures_decoder)

            # Step 3: Decode rules container
            rules_container = self._decode_rules_container(envelope, rules_container_decoder)

        # Step 4: Verify hash coverage
        verified_hash = self._verify_hash_in_signed_hashes(envelope)

        # Step 5: Verify whitelist signatures
        self._verify_whitelist_signatures(envelope, rules_container, verified_hash)

        # Step 6: Parse WhitelistedAddress from verified payload
        verified_whitelisted_address = parse_whitelisted_address_from_json(
            envelope.metadata.payload_as_string
        )

        return AddressVerificationResult(
            rules_container=rules_container,
            verified_hash=verified_hash,
            verified_whitelisted_address=verified_whitelisted_address,
        )

    def _verify_metadata_hash(self, envelope: SignedWhitelistedAddressEnvelope) -> None:
        """
        Verify that the computed hash matches the provided hash.

        Step 1 of the verification flow.
        """
        if not envelope.metadata or not envelope.metadata.payload_as_string:
            raise IntegrityError("payloadAsString is empty")
        if not envelope.metadata.hash:
            raise IntegrityError("metadata hash is empty")

        computed_hash = calculate_hex_hash(envelope.metadata.payload_as_string)
        if not constant_time_compare(computed_hash, envelope.metadata.hash):
            raise IntegrityError("metadata hash verification failed")

    def _verify_rules_container_signatures(
        self,
        envelope: SignedWhitelistedAddressEnvelope,
        user_signatures_decoder: Callable[[str], List[RuleUserSignature]],
    ) -> None:
        """
        Verify SuperAdmin signatures on the rules container.

        Step 2 of the verification flow.
        """
        if not self._super_admin_keys:
            raise IntegrityError("no SuperAdmin keys configured for verification")

        if not envelope.rules_container:
            raise IntegrityError("rulesContainer is empty")
        if not envelope.rules_signatures:
            raise IntegrityError("rulesSignatures is empty")

        # Decode rules signatures (protobuf UserSignatures)
        try:
            signatures = user_signatures_decoder(envelope.rules_signatures)
        except (ValueError, binascii.Error, KeyError) as e:
            raise IntegrityError(f"failed to decode rules signatures: {e}") from e

        # Decode rules container data
        try:
            rules_data = base64.b64decode(envelope.rules_container)
        except (binascii.Error, ValueError) as e:
            raise IntegrityError(f"failed to decode rules container: {e}") from e

        # Verify signatures
        valid_count = 0
        for sig in signatures:
            if sig.signature and is_valid_signature(
                rules_data, sig.signature, self._super_admin_keys
            ):
                valid_count += 1

        if valid_count < self._min_valid_signatures:
            raise IntegrityError(
                f"rules container signature verification failed: only {valid_count} valid signatures, "
                f"minimum {self._min_valid_signatures} required"
            )

    def _decode_rules_container(
        self,
        envelope: SignedWhitelistedAddressEnvelope,
        rules_container_decoder: Callable[[str], DecodedRulesContainer],
    ) -> DecodedRulesContainer:
        """
        Decode the base64 protobuf rules container.

        Step 3 of the verification flow.
        """
        try:
            return rules_container_decoder(envelope.rules_container)
        except (ValueError, KeyError, binascii.Error) as e:
            raise IntegrityError(f"failed to decode rules container: {e}") from e

    def _verify_hash_in_signed_hashes(
        self,
        envelope: SignedWhitelistedAddressEnvelope,
    ) -> str:
        """
        Verify that the metadata hash is covered by at least one signature.

        Step 4 of the verification flow.
        Returns the hash that was found (may be a legacy hash).
        """
        if envelope.signed_address is None:
            raise IntegrityError("signedAddress is nil")

        signatures = envelope.signed_address.signatures
        if not signatures:
            raise IntegrityError("no signatures in signedAddress")

        metadata_hash = envelope.metadata.hash

        # Try the provided hash first using constant-time comparison
        if _verify_hash_coverage(metadata_hash, signatures):
            return metadata_hash

        # Try legacy hashes for backward compatibility
        legacy_hashes = compute_legacy_hashes(envelope.metadata.payload_as_string)
        for legacy_hash in legacy_hashes:
            if _verify_hash_coverage(legacy_hash, signatures):
                return legacy_hash

        raise IntegrityError("metadata hash is not covered by any signature")

    def _verify_whitelist_signatures(
        self,
        envelope: SignedWhitelistedAddressEnvelope,
        rules_container: DecodedRulesContainer,
        metadata_hash: str,
    ) -> None:
        """
        Verify user signatures meet governance threshold requirements.

        Step 5 of the verification flow.
        """
        blockchain = envelope.blockchain
        network = envelope.network

        # Find matching address whitelisting rules
        whitelist_rules = rules_container.find_address_whitelisting_rules(
            blockchain, network
        )
        if whitelist_rules is None:
            raise WhitelistError(
                f"no address whitelisting rules found for blockchain={blockchain} "
                f"network={network}"
            )

        # Determine which thresholds to use based on rule lines
        parallel_thresholds = self._get_applicable_thresholds(whitelist_rules, envelope)
        if not parallel_thresholds:
            raise WhitelistError("no threshold rules defined")

        # Try to verify all paths (OR logic - only one needs to succeed)
        path_failures = self._try_verify_all_paths(
            parallel_thresholds,
            rules_container,
            envelope.signed_address.signatures,
            metadata_hash,
        )
        if path_failures:
            envelope_id = ""
            if envelope.metadata and envelope.metadata.payload_as_string:
                try:
                    payload = json.loads(envelope.metadata.payload_as_string)
                    envelope_id = str(payload.get("id", ""))
                except (json.JSONDecodeError, ValueError):
                    pass
            raise WhitelistError(
                f"signature verification failed for whitelisted address (ID: {envelope_id}): "
                f"no approval path satisfied the threshold requirements. {'; '.join(path_failures)}"
            )

    def _get_applicable_thresholds(
        self,
        rules: AddressWhitelistingRules,
        envelope: SignedWhitelistedAddressEnvelope,
    ) -> List[SequentialThresholds]:
        """
        Determine which thresholds to use based on rule lines.

        Checks rule lines only when: NO linked addresses AND exactly 1 linked wallet.
        """
        # Parse linked internal addresses from payload to check count
        has_linked_addresses = False
        if envelope.metadata and envelope.metadata.payload_as_string:
            try:
                payload = json.loads(envelope.metadata.payload_as_string)
                linked_addrs = payload.get("linkedInternalAddresses", [])
                has_linked_addresses = len(linked_addrs) > 0
            except (json.JSONDecodeError, ValueError):
                pass

        wallet_count = len(envelope.linked_wallets)
        should_check_rule_lines = not has_linked_addresses and wallet_count == 1

        if should_check_rule_lines and rules.lines:
            wallet_path = envelope.linked_wallets[0].path

            # Find matching line by wallet path
            for line in rules.lines:
                if self._matches_wallet_path(line, wallet_path):
                    return line.parallel_thresholds

        # Fallback to default thresholds
        return rules.parallel_thresholds

    @staticmethod
    def _matches_wallet_path(line: AddressWhitelistingLine, wallet_path: Optional[str]) -> bool:
        """Check if a rule line matches the given wallet path."""
        if not line.cells:
            return False

        source = line.cells[0]
        if source.type != RULE_SOURCE_TYPE_INTERNAL_WALLET:
            return False

        if source.internal_wallet is None:
            return False

        return wallet_path is not None and wallet_path == source.internal_wallet.path

    def _try_verify_all_paths(
        self,
        parallel_thresholds: List[SequentialThresholds],
        rules_container: DecodedRulesContainer,
        signatures: List[WhitelistSignatureEntry],
        metadata_hash: str,
    ) -> List[str]:
        """
        Try to verify all parallel threshold paths.

        Returns empty list if verification passed, or list of failure messages.
        """
        # Pre-compute JSON serialization of each signature's hashes array once,
        # so it can be reused across all group threshold checks.
        precomputed_hashes_json: List[Optional[bytes]] = []
        for sig in signatures:
            if sig.hashes is not None:
                precomputed_hashes_json.append(
                    json.dumps(sig.hashes, separators=(",", ":")).encode("utf-8")
                )
            else:
                precomputed_hashes_json.append(None)

        path_failures = []

        for i, seq_threshold in enumerate(parallel_thresholds):
            err = self._verify_sequential_thresholds(
                seq_threshold, rules_container, signatures, metadata_hash,
                precomputed_hashes_json,
            )
            if err is None:
                return []  # Verification passed
            path_failures.append(f"Path {i + 1}: {err}")

        return path_failures

    def _verify_sequential_thresholds(
        self,
        seq_threshold: SequentialThresholds,
        rules_container: DecodedRulesContainer,
        signatures: List[WhitelistSignatureEntry],
        metadata_hash: str,
        precomputed_hashes_json: List[Optional[bytes]],
    ) -> Optional[str]:
        """
        Verify all group thresholds in a sequential threshold path.

        Returns None if successful, or error message on failure.
        """
        if seq_threshold is None or not seq_threshold.thresholds:
            return "no group thresholds defined"

        # ALL group thresholds must be satisfied (AND logic)
        for group_threshold in seq_threshold.thresholds:
            err = self._verify_group_threshold(
                group_threshold, rules_container, signatures, metadata_hash,
                precomputed_hashes_json,
            )
            if err:
                return err

        return None

    def _verify_group_threshold(
        self,
        group_threshold,
        rules_container: DecodedRulesContainer,
        signatures: List[WhitelistSignatureEntry],
        metadata_hash: str,
        precomputed_hashes_json: List[Optional[bytes]],
    ) -> Optional[str]:
        """
        Verify that a group threshold is met.

        Returns None if successful, or error message on failure.
        """
        group_id = group_threshold.group_id
        min_sigs = group_threshold.get_min_signatures()

        group = rules_container.find_group_by_id(group_id)
        if group is None:
            return f"group '{group_id}' not found in rules container"

        if not group.user_ids:
            if min_sigs > 0:
                return f"group '{group_id}' has no users but requires {min_sigs} signature(s)"
            return None  # min_signatures == 0, so empty group is OK

        # Build set for faster lookup
        group_user_id_set = set(group.user_ids)

        # Count valid signatures from users in this group
        valid_count = 0
        skipped_reasons = []

        for sig_idx, sig in enumerate(signatures):
            if sig.user_signature is None:
                skipped_reasons.append("signature has nil userSig")
                continue

            sig_user_id = sig.user_signature.user_id
            if sig_user_id not in group_user_id_set:
                continue  # Signer not in this group - not an error

            # Check that metadata hash is covered by this signature
            if not _contains_hash(sig.hashes, metadata_hash):
                skipped_reasons.append(
                    f"user '{sig_user_id}' signature does not cover metadata hash"
                )
                continue

            user = rules_container.find_user_by_id(sig_user_id)
            if user is None:
                skipped_reasons.append(f"user '{sig_user_id}' not found in rules container")
                continue
            if not user.public_key_pem:
                skipped_reasons.append(f"user '{sig_user_id}' has no public key")
                continue

            # Use cached public key from rules container
            try:
                public_key = rules_container.get_user_public_key(user.public_key_pem)
            except ValueError as e:
                skipped_reasons.append(f"failed to decode public key for user '{sig_user_id}': {e}")
                continue

            # Verify signature against pre-computed JSON-encoded hashes
            try:
                hashes_data = precomputed_hashes_json[sig_idx]
                if hashes_data is None:
                    skipped_reasons.append(f"user '{sig_user_id}' has no hashes")
                    continue
                if verify_signature(
                    public_key,
                    hashes_data,
                    sig.user_signature.signature,
                ):
                    valid_count += 1
                    if valid_count >= min_sigs:
                        return None  # Threshold met
                else:
                    skipped_reasons.append(f"user '{sig_user_id}' signature verification failed")
            except (InvalidSignature, ValueError, binascii.Error) as e:
                skipped_reasons.append(f"user '{sig_user_id}' signature verification error: {e}")

        # Threshold not met
        message = (
            f"group '{group_id}' requires {min_sigs} signature(s) but only {valid_count} valid"
        )
        if skipped_reasons:
            message += f" [{'; '.join(skipped_reasons)}]"
        return message


def _verify_hash_coverage(
    metadata_hash: str, signatures: List[WhitelistSignatureEntry]
) -> bool:
    """
    Check if the metadata hash is covered by at least one signature.

    Uses constant-time comparison to prevent timing side-channel attacks.
    """
    found = False
    for sig in signatures:
        for h in sig.hashes:
            if hmac.compare_digest(metadata_hash, h):
                found = True
                # Continue checking to maintain constant time
    return found


def _contains_hash(hashes: List[str], target: str) -> bool:
    """Check if a hash is in the list using constant-time comparison."""
    for h in hashes:
        if hmac.compare_digest(h, target):
            return True
    return False
