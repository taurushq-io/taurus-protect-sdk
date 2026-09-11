"""Governance rule service for Taurus-PROTECT SDK."""

from __future__ import annotations

import hashlib
import hmac
import threading
from typing import TYPE_CHECKING, Any, List, Optional, Set

from cryptography.hazmat.primitives.asymmetric.ec import (
    EllipticCurvePrivateKey,
    EllipticCurvePublicKey,
)

from taurus_protect._strict_base64 import strict_b64decode
from taurus_protect.errors import APIError, IntegrityError, WhitelistError
from taurus_protect.helpers.signature_verifier import verify_governance_rules
from taurus_protect.mappers.governance_rules import rules_container_from_base64
from taurus_protect.models.governance_rules import (
    ExcludedRuleset,
    DecodedRulesContainer,
    GovernanceRules,
    GovernanceRulesHistoryResult,
    RuleUser,
    RuleUserSignature,
    SuperAdminPublicKey,
)
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.governance_rules_api import (
        GovernanceRulesApi,
    )


_MAX_VERIFIED_RULESETS = 256


def _update_length_prefixed(digest: Any, data: bytes) -> None:
    """
    Write an 8-byte big-endian length, then the bytes, so a sequence of fields cannot be
    re-partitioned into a different sequence that encodes identically.
    """
    digest.update(len(data).to_bytes(8, "big"))
    digest.update(data)


def _ruleset_verification_key(rules: GovernanceRules) -> Optional[str]:
    """
    Identify the exact document verification cleared: the DECODED container bytes -- the
    same bytes ``verify_governance_rules`` checks the signatures over -- plus every
    signature over them. A container whose signature set changed is a memo MISS rather
    than a stale hit; signature order is server-controlled, so it is sorted out of the key.

    Every field is LENGTH-PREFIXED, and that is load-bearing. A plain concatenation is not
    injective: with ``sha256(container || sig || ...)`` the boundary between the container
    and the signature list is not committed, so a response-controlling attacker can shift
    bytes across it and make a MODIFIED container collide with a genuine one's key --
    inheriting its "already verified" status and skipping ECDSA entirely. Interleaving a
    ``b"\\x00"`` separator does NOT fix this; only length prefixes do.

    ``user_id`` is deliberately excluded: verification never reads it (see
    ``verify_governance_rules_signatures``, which consumes only the signature), so a
    server-controlled field that verification ignores must not be able to influence a
    verification-SKIP decision.

    Keying on the decoded bytes rather than the base64 text also removes base64 leniency
    as a lever, and guarantees the key cannot identify something other than what was
    verified.

    Returns:
        The hex digest, or None when the container does not decode -- in which case there
        is nothing worth memoising and the caller must fall through to verification.
    """
    try:
        data = strict_b64decode(rules.rules_container)
    except ValueError:
        # binascii.Error subclasses ValueError; a container we cannot decode has no
        # stable identity to memoise on.
        return None

    digest = hashlib.sha256()
    _update_length_prefixed(digest, data)

    signatures = sorted((sig.signature or "") for sig in (rules.rules_signatures or []))
    digest.update(len(signatures).to_bytes(8, "big"))
    for signature in signatures:
        _update_length_prefixed(digest, signature.encode("utf-8"))
    return digest.hexdigest()


class GovernanceRuleService(BaseService):
    """
    Service for managing governance rules.

    Governance rules define the approval workflows and policies for
    transaction requests and address whitelisting.
    """

    def __init__(
        self,
        api_client: Any,
        governance_rules_api: "GovernanceRulesApi",
        super_admin_keys: List[EllipticCurvePublicKey],
        min_valid_signatures: int,
    ) -> None:
        """
        Initialize the governance rule service.

        Args:
            api_client: The OpenAPI client instance.
            governance_rules_api: The governance rules API instance.
            super_admin_keys: List of SuperAdmin public keys for verification.
            min_valid_signatures: Minimum number of valid signatures required.
        """
        super().__init__(api_client)
        self._api = governance_rules_api
        self._super_admin_keys = super_admin_keys
        self._min_valid_signatures = min_valid_signatures
        # Rulesets whose SuperAdmin signatures already checked out, keyed by
        # _ruleset_verification_key. See _verified_ruleset.
        self._verified: Set[str] = set()
        self._verified_lock = threading.RLock()

    def get_rules(self) -> Optional[GovernanceRules]:
        """
        Get the currently enforced governance rules.

        Returns:
            The governance rules, or None if not available.

        Raises:
            APIError: If the API call fails.
            IntegrityError: If signature verification fails.
        """
        try:
            reply = self._api.rule_service_get_rules()
            result = reply.result
            if result is None:
                return None

            rules = self._map_rules_from_dto(result)
            return self._verified_ruleset(rules)
        except IntegrityError:
            raise
        except Exception as e:
            # Funnel on the SDK error taxonomy, not on the raw generated ApiException.
            # Keying on ApiException left every other failure -- a urllib3 transport
            # error, say -- propagating unmapped, so a caller could not treat it as an
            # APIError at all. And the pass-through list must name IntegrityError and
            # WhitelistError explicitly: both are plain Exceptions, so _handle_error
            # would map them to ServerError(500), whose is_retryable() is True. That
            # inverts the documented "security error, DO NOT retry" contract.
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_rules_by_id(self, rules_id: str) -> Optional[GovernanceRules]:
        """
        Get a governance ruleset by its ID.

        Args:
            rules_id: The ruleset ID.

        Returns:
            The governance rules, or None if not found.

        Raises:
            APIError: If the API call fails.
            IntegrityError: If signature verification fails.
            ValueError: If rules_id is empty.
        """
        self._validate_required(rules_id, "rules_id")

        try:
            reply = self._api.rule_service_get_rules_by_id(rules_id)
            result = reply.result
            if result is None:
                return None

            rules = self._map_rules_from_dto(result)
            return self._verified_ruleset(rules)
        except IntegrityError:
            raise
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_rules_proposal(self) -> Optional[GovernanceRules]:
        """
        Get the proposed governance rules.

        Requires SuperAdmin or SuperAdminReadOnly role.

        Returns:
            The proposed governance rules, or None if not available.

        Raises:
            APIError: If the API call fails.
        """
        try:
            reply = self._api.rule_service_get_rules_proposal()
            result = reply.result
            if result is None:
                return None

            # Proposal rules are not verified
            return self._map_rules_from_dto(result)
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def update_rules_proposal(self, container: DecodedRulesContainer) -> None:
        """Submit a rules container as a governance proposal (SuperAdmin only).

        The container is encoded to the wire format internally; the
        server-controlled ``enforcedRulesHash`` and ``timestamp`` are stripped
        (the server recomputes them and rejects submissions asserting stale
        values). The endpoint returns no body — callers that need the persisted
        proposal should call :meth:`get_rules_proposal` (rules reads are cached
        server-side, so an immediate read-back may be stale).

        Args:
            container: The rules container to submit.

        Raises:
            APIError: If the API call fails.
            ValueError: If the container is None.
        """
        if container is None:
            raise ValueError("rules container cannot be None")
        from taurus_protect._internal.openapi.models.tgvalidatord_update_rules_proposal_request import (
            TgvalidatordUpdateRulesProposalRequest,
        )
        from taurus_protect.mappers.rules_container_encode import rules_container_to_base64

        encoded = rules_container_to_base64(container)
        try:
            self._api.rule_service_update_rules_proposal(
                body=TgvalidatordUpdateRulesProposalRequest(rules_container=encoded)
            )
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def proposal_container_hash(self, rules: GovernanceRules) -> str:
        """Return the canonical SHA-256 hex digest of a ruleset's decoded container.

        This is the value to pass as ``expected_container_hash`` to
        :meth:`approve_rules_proposal`, and it is what pins the approval to reviewed
        content.

        It digests the DECODED bytes, not the base64 text: base64 has multiple encodings
        of the same bytes, so hashing the text would let a re-encoded but byte-identical
        container read as a mismatch. Those decoded bytes are also exactly what gets
        signed.

        Args:
            rules: The ruleset whose container to digest.

        Returns:
            The hex digest.

        Raises:
            ValueError: If the ruleset carries no container.
            IntegrityError: If the container is not valid base64.
        """
        if rules is None or not rules.rules_container:
            raise ValueError("ruleset carries no rules container")
        try:
            data = strict_b64decode(rules.rules_container)
        except ValueError as e:
            raise IntegrityError(f"failed to decode rules container: {e}") from e
        return hashlib.sha256(data).hexdigest()

    def decode_proposal_for_review(self, rules: GovernanceRules) -> DecodedRulesContainer:
        """Decode a PENDING proposal so a SuperAdmin can inspect it before approving.

        UNVERIFIED BY DESIGN. A pending proposal legitimately carries 0..N signatures --
        signing IS the approval step -- so there is no threshold to check yet, and
        :meth:`get_decoded_rules_container` (which verifies unconditionally) always fails
        on one. That is why this is a separate, explicitly-named entry point rather than a
        flag: the absence of verification has to be visible at the call site.

        Pair it with :meth:`proposal_container_hash` and pass that digest to
        :meth:`approve_rules_proposal`, so the bytes signed are the bytes reviewed.

        Args:
            rules: The pending proposal, from :meth:`get_rules_proposal`.

        Returns:
            The decoded rules container, NOT verified.

        Raises:
            ValueError: If the ruleset carries no container.
            IntegrityError: If decoding fails.
        """
        if rules is None or not rules.rules_container:
            raise ValueError("ruleset carries no rules container")
        return rules_container_from_base64(rules.rules_container)

    def approve_rules_proposal(
        self,
        private_key: EllipticCurvePrivateKey,
        comment: str,
        expected_container_hash: str,
    ) -> None:
        """Sign the pending proposal's rules container and submit the approval.

        ``expected_container_hash`` PINS the content being approved: it must be the
        :meth:`proposal_container_hash` of the proposal the caller actually reviewed. The
        proposal is re-fetched here, and if its container does not match that digest the
        call aborts WITHOUT signing.

        That pin is the security control. Without it a server able to shape responses
        could serve the benign proposal to the review call and a different container to
        the re-fetch inside this method, obtaining a GENUINE SuperAdmin signature over
        bytes of its choosing -- and such a container then verifies clean everywhere,
        including in independent and air-gapped verifiers that never trusted that server.
        Nothing else binds the two calls: the wire format carries no proposal id, hash or
        version, and a pending proposal has no signatures to verify against, so
        re-verification cannot substitute for pinning.

        The signature is computed over the decoded bytes of the pending proposal's rules
        container (SHA-256 + P-256 ECDSA, base64 raw r||s).

        Args:
            private_key: SuperAdmin P-256 private key.
            comment: Approval comment.
            expected_container_hash: :meth:`proposal_container_hash` of the reviewed
                proposal.

        Raises:
            ValueError: If no proposal is pending, the key is None, or the pin is missing.
            IntegrityError: If the pending container differs from the pin.
            APIError: If the API call fails.
        """
        if private_key is None:
            raise ValueError("private_key cannot be None")
        if not expected_container_hash:
            raise ValueError(
                "expected_container_hash is required: it pins the approval to the "
                "container you reviewed (see proposal_container_hash)"
            )

        from taurus_protect._internal.openapi.models.tgvalidatord_approve_rules_proposal_request import (
            TgvalidatordApproveRulesProposalRequest,
        )
        from taurus_protect.crypto.signing import sign_data

        proposal = self.get_rules_proposal()
        if proposal is None or not proposal.rules_container:
            raise ValueError("no pending rules proposal to approve")

        actual_hash = self.proposal_container_hash(proposal)
        # Constant-time, matching how every other hash comparison in this SDK is done.
        if not hmac.compare_digest(actual_hash, expected_container_hash):
            raise IntegrityError(
                "refusing to sign: the pending rules proposal changed since it was "
                f"reviewed (reviewed {expected_container_hash}, pending {actual_hash})"
            )

        signature = sign_data(private_key, strict_b64decode(proposal.rules_container))
        try:
            self._api.rule_service_approve_rules_proposal(
                body=TgvalidatordApproveRulesProposalRequest(signature=signature, comment=comment)
            )
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def reject_rules_proposal(self, comment: str) -> None:
        """Reject the pending rules proposal with a comment (SuperAdmin only).

        Args:
            comment: Rejection comment.

        Raises:
            APIError: If the API call fails.
        """
        from taurus_protect._internal.openapi.models.tgvalidatord_reject_rules_proposal_request import (
            TgvalidatordRejectRulesProposalRequest,
        )

        try:
            self._api.rule_service_reject_rules_proposal(
                body=TgvalidatordRejectRulesProposalRequest(comment=comment)
            )
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_rules_history(
        self, page_size: int = 50, cursor: Optional[str] = None
    ) -> GovernanceRulesHistoryResult:
        """
        Get the history of governance rules with cursor-based pagination.

        Args:
            page_size: Maximum number of rules to return per page (default 50).
            cursor: Opaque pagination cursor from a previous call (None for first page).

        Returns:
            A GovernanceRulesHistoryResult with rules, cursor, and total_items.

        Raises:
            APIError: If the API call fails.
            ValueError: If page_size is not positive.
        """
        if page_size <= 0:
            raise ValueError("page_size must be positive")

        try:
            reply = self._api.rule_service_get_rules_history(
                limit=str(page_size),
                cursor=cursor,
            )

            # Every historical ruleset is a past ENFORCED document, so each is verified.
            # The memo makes this one ECDSA pass per distinct container, not one per row.
            #
            #   entries -> _verified_ruleset -+- ok    -> kept
            #                                 +- error -> excluded + named in the result
            #
            # LENIENT, unlike the whitelist lists, and it does NOT raise when nothing
            # survives: a SuperAdmin key rotation makes every pre-rotation ruleset
            # unverifiable, so aborting would deny the whole audit trail from then on.
            rules_list: List[GovernanceRules] = []
            excluded: List[ExcludedRuleset] = []
            if reply.result:
                for dto in reply.result:
                    entry = self._map_rules_from_dto(dto)
                    try:
                        rules_list.append(self._verified_ruleset(entry))
                    except IntegrityError as exc:
                        excluded.append(
                            ExcludedRuleset(creation_date=entry.creation_date, reason=exc.message)
                        )

            # Extract cursor (may be bytes or str from OpenAPI)
            raw_cursor = getattr(reply, "cursor", None)
            next_cursor: Optional[str] = None
            if raw_cursor:
                next_cursor = (
                    raw_cursor.decode("utf-8") if isinstance(raw_cursor, bytes) else raw_cursor
                )

            # Reduced by the exclusions: the server counts rows it returned, the caller
            # receives only those that verified.
            # total_items is a STRING on this model (Go carries an int64), so the
            # reduced value is re-stringified rather than changing the public type.
            total = getattr(reply, "total_items", None) or None
            if total is not None and excluded:
                try:
                    total = str(max(0, int(total) - len(excluded)))
                except ValueError:
                    # A non-numeric total is the server's to explain; pass it through
                    # rather than inventing a count.
                    pass

            return GovernanceRulesHistoryResult(
                rules=rules_list,
                cursor=next_cursor,
                total_items=total,
                excluded_unverified=excluded,
            )
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def get_public_keys(self) -> List[SuperAdminPublicKey]:
        """
        Get the list of SuperAdmin public keys.

        Returns the public keys of all SuperAdmin users configured for
        the tenant. These keys are used for signing governance rules.

        Returns:
            List of SuperAdmin public keys.

        Raises:
            APIError: If the API call fails.
        """
        try:
            reply = self._api.rule_service_get_public_keys()

            result: List[SuperAdminPublicKey] = []
            if reply.public_keys:
                for dto in reply.public_keys:
                    result.append(
                        SuperAdminPublicKey(
                            user_id=getattr(dto, "user_id", None),
                            public_key=getattr(dto, "public_key", None),
                        )
                    )

            return result
        except Exception as e:
            # See above: funnel on the SDK taxonomy, and never let IntegrityError or
            # WhitelistError be remapped to a retryable ServerError(500).
            if isinstance(e, (APIError, IntegrityError, WhitelistError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def _verified_ruleset(self, rules: GovernanceRules) -> GovernanceRules:
        """
        Verify a ruleset's SuperAdmin signatures, memoising the outcome so the same
        document is not re-verified on every read.

        get_rules / get_rules_by_id / get_rules_history -> _verified_ruleset
                                                             |- memo hit -> return
                                                             +- miss -> ECDSA -> memoise

        The container itself is already cached for the address/asset/price paths by
        ``RulesContainerCache``, which fetches through ``get_decoded_rules_container``.
        This memo covers only the reads that return the RAW ruleset, which that cache
        does not serve. Do not add a third cache of the same bytes.

        Only SUCCESSES are memoised: a failure must resurface its error on every call.

        Raises:
            IntegrityError: If verification fails.
        """
        key = _ruleset_verification_key(rules)
        if key is not None:
            with self._verified_lock:
                if key in self._verified:
                    return rules

        verify_governance_rules(rules, self._min_valid_signatures, self._super_admin_keys)

        # An undecodable container has no stable identity to memoise on; verification
        # above has already had the final say.
        if key is None:
            return rules

        with self._verified_lock:
            # Bounded so walking a long history cannot grow it without limit.
            # Containers change rarely, so a plain reset beats LRU bookkeeping.
            if len(self._verified) >= _MAX_VERIFIED_RULESETS:
                self._verified.clear()
            self._verified.add(key)
        return rules

    def verify_governance_rules(self, rules: GovernanceRules) -> GovernanceRules:
        """
        Verify that governance rules have enough valid SuperAdmin signatures.

        Args:
            rules: The governance rules to verify.

        Returns:
            The verified rules.

        Raises:
            IntegrityError: If verification fails.
        """
        verify_governance_rules(
            rules,
            self._min_valid_signatures,
            self._super_admin_keys,
        )
        return rules

    def get_decoded_rules_container(self, rules: GovernanceRules) -> DecodedRulesContainer:
        """
        Get the decoded rules container from governance rules.

        Verifies signatures and decodes the rules container.

        Args:
            rules: The governance rules.

        Returns:
            The decoded rules container.

        Raises:
            IntegrityError: If signature verification fails or decoding fails.
            ValueError: If rules is None.
        """
        if rules is None:
            raise ValueError("rules cannot be None")

        # Verify signatures before decoding. Through the memo, which is also the ONE
        # verification site on this path: get_rules() verified and then this verified
        # again, so removing either left a test asserting verification still green.
        self._verified_ruleset(rules)

        # Decode the rules container
        if rules.rules_container is None:
            raise IntegrityError("Rules container is None")

        return self._decode_rules_container(rules.rules_container)

    def _decode_rules_container(self, rules_container_b64: str) -> DecodedRulesContainer:
        """
        Decode a base64-encoded rules container.

        Args:
            rules_container_b64: Base64-encoded rules container.

        Returns:
            Decoded rules container.

        Raises:
            IntegrityError: If decoding fails.
        """
        from taurus_protect.mappers.governance_rules import rules_container_from_base64

        return rules_container_from_base64(rules_container_b64)

    @property
    def super_admin_keys(self) -> List[EllipticCurvePublicKey]:
        """Get the configured SuperAdmin public keys."""
        return list(self._super_admin_keys)

    @property
    def min_valid_signatures(self) -> int:
        """Get the configured minimum valid signatures."""
        return self._min_valid_signatures

    @staticmethod
    def _map_rules_from_dto(dto: Any) -> GovernanceRules:
        """Map OpenAPI DTO to domain model."""
        signatures: List[RuleUserSignature] = []
        if dto.rules_signatures:
            for sig_dto in dto.rules_signatures:
                signatures.append(
                    RuleUserSignature(
                        user_id=getattr(sig_dto, "user_id", None),
                        signature=getattr(sig_dto, "signature", None),
                    )
                )

        return GovernanceRules(
            rules_container=dto.rules_container,
            rules_signatures=signatures,
            locked=dto.locked if hasattr(dto, "locked") else False,
            creation_date=dto.creation_date if hasattr(dto, "creation_date") else None,
            update_date=dto.update_date if hasattr(dto, "update_date") else None,
        )
