"""Whitelisted asset verification utilities.

This module implements the 5-step verification flow for whitelisted assets (contracts):
1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
2. Verify rules container signatures (SuperAdmin signatures)
3. Decode rules container (base64 -> protobuf -> model)
4. Verify hash coverage (metadata.hash in at least one signature.hashes)
5. Verify whitelist signatures (user signatures meet governance thresholds)
"""

from __future__ import annotations

import base64
import binascii
import hmac
import json
from dataclasses import dataclass
from typing import TYPE_CHECKING, Callable, List, Optional

from cryptography.exceptions import InvalidSignature

from taurus_protect._strict_base64 import strict_b64decode
from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.crypto.signing import verify_signature
from taurus_protect.errors import IntegrityError, WhitelistError
from taurus_protect.helpers.constant_time import constant_time_compare
from taurus_protect.helpers.signature_verifier import (
    key_fingerprint,
    verify_governance_rules_signatures,
)
from taurus_protect.helpers.whitelist_hash_helper import (
    compute_asset_legacy_hashes,
    contains_hash,
    resolve_rule_key,
    verify_hash_coverage,
)
from taurus_protect.models.governance_rules import (
    DecodedRulesContainer,
    RuleUserSignature,
    SequentialThresholds,
)
from taurus_protect.models.whitelisted_address import WhitelistedAsset

if TYPE_CHECKING:
    from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey


@dataclass
class AssetVerificationResult:
    """Result of whitelisted asset verification."""

    rules_container: DecodedRulesContainer


class WhitelistedAssetVerifier:
    """
    Verifier for whitelisted assets.

    Performs the complete 5-step cryptographic verification to ensure
    asset integrity and proper approval according to governance rules.
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

    # 5-step verification for a whitelisted asset, plus the parse that makes it usable.
    #
    #   envelope ──▶ 1 hash ──▶ 2 SuperAdmin sigs ──▶ 3 decode rules
    #                                                      │
    #                6 parse VERIFIED payload ◀── 5 thresholds ◀── 4 coverage
    #                           │                       │                │
    #                           ▼                       │      returns the hash it
    #                    returned to caller             │      matched, which is
    #                                                   │      what step 5 checks
    #                                                   ▼
    #                                per group: DISTINCT signers, keyed on the
    #                                container's public key, never the entry's userId
    #
    # Steps 1-5 prove the envelope is authentic; step 6 is what stops an unsigned
    # value reaching the caller. Assets use ContractAddressWhitelistingRules and the
    # ASSET legacy hashes (isNFT / kindType) — not the address ones.
    #
    # Steps 1-5 run here; step 6 runs in services/whitelisted_asset_service.py,
    # which maps the asset from the payload this verifier just cleared.
    def verify_whitelisted_asset(
        self,
        asset: WhitelistedAsset,
        rules_container_decoder: Callable[[str], DecodedRulesContainer],
        user_signatures_decoder: Callable[[str], List[RuleUserSignature]],
        dto_blockchain: Optional[str] = None,
        dto_network: Optional[str] = None,
    ) -> AssetVerificationResult:
        """
        Perform the complete 5-step verification of a whitelisted asset.

        Steps:
        1. Verify metadata hash (SHA-256 of payloadAsString == metadata.hash)
        2. Verify rules container signatures (SuperAdmin signatures)
        3. Decode rules container (base64 -> protobuf -> model)
        4. Verify hash coverage (metadata.hash in at least one signature.hashes)
        5. Verify whitelist signatures (user signatures meet governance thresholds)

        Args:
            asset: The whitelisted asset to verify.
            rules_container_decoder: Function to decode base64 rules container.
            user_signatures_decoder: Function to decode base64 user signatures.
            dto_blockchain: Blockchain from DTO envelope (used for rules lookup
                when payload blockchain is missing).
            dto_network: Network from DTO envelope (used for rules lookup
                when payload network is missing).

        Returns:
            Verification result with decoded rules container.

        Raises:
            IntegrityError: If verification fails at any step.
            WhitelistError: If signature thresholds are not met.
        """
        if asset is None:
            raise ValueError("asset cannot be None")
        if asset.metadata is None:
            raise ValueError("metadata cannot be None")

        # Step 1: Verify metadata hash
        self._verify_metadata_hash(asset)

        # Step 2: Verify rules container signatures
        self._verify_rules_container_signatures(asset, user_signatures_decoder)

        # Step 3: Decode rules container
        rules_container = self._decode_rules_container(asset, rules_container_decoder)

        # Step 4: Verify hash coverage
        verified_hash = self._verify_hash_in_signed_hashes(asset)

        # Step 5: Verify whitelist signatures.
        #
        # The DTO values are passed through unchanged. `asset.blockchain` is sourced FROM
        # the payload, so substituting it here fed resolve_rule_key the payload on both
        # sides and its payload-vs-DTO mismatch check could never fire. Go, Java and
        # TypeScript all pass the genuine DTO values.
        self._verify_whitelist_signatures(
            asset, rules_container, verified_hash, dto_blockchain, dto_network
        )

        return AssetVerificationResult(rules_container=rules_container)

    def _verify_metadata_hash(self, asset: WhitelistedAsset) -> None:
        """
        Verify that the computed hash matches the provided hash.

        Step 1 of the verification flow.
        """
        if not asset.metadata or not asset.metadata.payload_as_string:
            raise IntegrityError("payloadAsString is empty")
        if not asset.metadata.hash:
            raise IntegrityError("metadata hash is empty")

        computed_hash = calculate_hex_hash(asset.metadata.payload_as_string)
        if not constant_time_compare(computed_hash, asset.metadata.hash):
            raise IntegrityError("metadata hash verification failed")

    def _verify_rules_container_signatures(
        self,
        asset: WhitelistedAsset,
        user_signatures_decoder: Callable[[str], List[RuleUserSignature]],
    ) -> None:
        """
        Verify SuperAdmin signatures on the rules container.

        Step 2 of the verification flow.
        """
        if not self._super_admin_keys:
            raise IntegrityError("no SuperAdmin keys configured for verification")

        if not asset.rules_container:
            raise IntegrityError("rulesContainer is empty")
        if not asset.rules_signatures:
            raise IntegrityError("rulesSignatures is empty")

        # Decode rules signatures (protobuf UserSignatures)
        try:
            signatures = user_signatures_decoder(asset.rules_signatures)
        except (ValueError, binascii.Error, KeyError) as e:
            raise IntegrityError(f"failed to decode rules signatures: {e}") from e

        # Decode rules container data
        try:
            rules_data = strict_b64decode(asset.rules_container)
        except (binascii.Error, ValueError) as e:
            raise IntegrityError(f"failed to decode rules container: {e}") from e

        try:
            verify_governance_rules_signatures(
                rules_data,
                signatures,
                self._super_admin_keys,
                self._min_valid_signatures,
            )
        except IntegrityError as e:
            raise IntegrityError(f"rules container signature verification failed: {e}") from e

    def _decode_rules_container(
        self,
        asset: WhitelistedAsset,
        rules_container_decoder: Callable[[str], DecodedRulesContainer],
    ) -> DecodedRulesContainer:
        """
        Decode the base64 protobuf rules container.

        Step 3 of the verification flow.
        """
        try:
            return rules_container_decoder(asset.rules_container)
        except (ValueError, KeyError, binascii.Error) as e:
            raise IntegrityError(f"failed to decode rules container: {e}") from e

    def _verify_hash_in_signed_hashes(self, asset: WhitelistedAsset) -> str:
        """
        Verify that the metadata hash is covered by at least one signature.

        Step 4 of the verification flow.
        Returns the hash that was found (may be a legacy hash).

        Uses constant-time comparison to prevent timing side-channel attacks
        that could leak information about the hash value.

        Supports legacy hashes for backward compatibility with assets signed
        before schema changes (e.g., before isNFT or kindType was added).
        """
        if asset.signed_contract_address is None:
            raise IntegrityError("signedContractAddress is nil")

        signatures = asset.signed_contract_address.signatures
        if not signatures:
            raise IntegrityError("no signatures in signedContractAddress")

        metadata_hash = asset.metadata.hash

        # Try the provided hash first using constant-time comparison
        if verify_hash_coverage(metadata_hash, signatures):
            return metadata_hash

        # Try legacy hashes for backward compatibility
        # This handles assets signed before schema changes (e.g., before isNFT or kindType was added)
        legacy_hashes = compute_asset_legacy_hashes(asset.metadata.payload_as_string)
        for legacy_hash in legacy_hashes:
            if verify_hash_coverage(legacy_hash, signatures):
                return legacy_hash

        raise IntegrityError("metadata hash is not covered by any signature")

    def _verify_whitelist_signatures(
        self,
        asset: WhitelistedAsset,
        rules_container: DecodedRulesContainer,
        verified_hash: str,
        lookup_blockchain: Optional[str] = None,
        lookup_network: Optional[str] = None,
    ) -> None:
        """
        Verify user signatures meet governance threshold requirements.

        Step 5 of the verification flow.

        Args:
            asset: The whitelisted asset being verified.
            rules_container: Decoded governance rules container.
            verified_hash: The hash that was verified in step 4.
            lookup_blockchain: Blockchain for rules lookup (defaults to asset.blockchain).
            lookup_network: Network for rules lookup (defaults to asset.network).
        """
        metadata_hash = verified_hash

        # Keyed off the SIGNED payload, not the response. See the address verifier.
        blockchain, network = resolve_rule_key(
            asset.metadata.payload_as_string if asset.metadata else None,
            lookup_blockchain or asset.blockchain,
            lookup_network or asset.network,
        )

        # Find matching contract address whitelisting rules
        whitelist_rules = rules_container.find_contract_address_whitelisting_rules(
            blockchain, network
        )
        if whitelist_rules is None:
            raise WhitelistError(
                f"no contract address whitelisting rules found for blockchain={blockchain} "
                f"network={network}"
            )

        # Contract whitelisting uses parallelThresholds directly
        parallel_thresholds = whitelist_rules.parallel_thresholds
        if not parallel_thresholds:
            raise WhitelistError("no threshold rules defined")

        # Try to verify all paths (OR logic - only one needs to succeed)
        path_failures = self._try_verify_all_paths(
            parallel_thresholds,
            rules_container,
            asset.signed_contract_address.signatures,
            metadata_hash,
        )
        if path_failures:
            raise WhitelistError(
                f"signature verification failed for whitelisted asset (ID: {asset.id}): "
                f"no approval path satisfied the threshold requirements. {'; '.join(path_failures)}"
            )

    def _try_verify_all_paths(
        self,
        parallel_thresholds: List[SequentialThresholds],
        rules_container: DecodedRulesContainer,
        signatures: list,
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
        signatures: list,
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
        signatures: list,
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

        # A populated group with a zero threshold is a malformed container, not a
        # group anyone may satisfy. There is no post-loop threshold check -- the
        # only success exit is inside the loop after an increment -- so a zero
        # here silently means "one signature suffices", turning a 2-of-N group
        # into 1-of-N. Fail closed.
        if min_sigs <= 0:
            return (
                f"group '{group_id}' has {len(group.user_ids)} user(s) but requires "
                "0 signature(s): minimum_signatures must be positive"
            )

        # Build set for faster lookup
        group_user_id_set = set(group.user_ids)

        # Count DISTINCT signers, not signature entries.
        #
        # The entries come from the server-supplied userSignatures blob, so counting
        # them let a duplicated entry from one group member satisfy an N-of-M group.
        # Keyed on the container-resolved public key rather than the server-supplied
        # user_id, so one compromised key shared by two IDs counts once.
        signers: set = set()
        skipped_reasons = []

        for sig_idx, sig in enumerate(signatures):
            if sig.user_signature is None:
                skipped_reasons.append("signature has nil userSig")
                continue

            sig_user_id = sig.user_signature.user_id
            if sig_user_id not in group_user_id_set:
                continue  # Signer not in this group - not an error

            # Check that metadata hash is covered by this signature (constant-time)
            if not contains_hash(sig.hashes, metadata_hash):
                # Carries the hash and the covered list, matching Go, Java and TS:
                # a hash is SHA-256 of a payload the caller already holds, and
                # without it a threshold failure is undebuggable.
                skipped_reasons.append(
                    f"user '{sig_user_id}' signature does not cover metadata hash "
                    f"'{metadata_hash}' (signed hashes={sig.hashes})"
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
                    signers.add(key_fingerprint(public_key))
                    if len(signers) >= min_sigs:
                        return None  # Threshold met
                else:
                    skipped_reasons.append(f"user '{sig_user_id}' signature verification failed")
            except (InvalidSignature, ValueError, binascii.Error) as e:
                skipped_reasons.append(f"user '{sig_user_id}' signature verification error: {e}")

        # Threshold not met
        message = (
            f"group '{group_id}' requires {min_sigs} distinct signer(s) "
            f"but only {len(signers)} valid"
        )
        if skipped_reasons:
            message += f" [{'; '.join(skipped_reasons)}]"
        return message
