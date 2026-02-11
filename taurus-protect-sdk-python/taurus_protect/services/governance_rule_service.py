"""Governance rule service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePublicKey

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.errors import APIError, IntegrityError
from taurus_protect.helpers.signature_verifier import verify_governance_rules
from taurus_protect.models.governance_rules import (
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
            return self.verify_governance_rules(rules)
        except IntegrityError:
            raise
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

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
            return self.verify_governance_rules(rules)
        except IntegrityError:
            raise
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

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
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

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

            rules_list: List[GovernanceRules] = []
            if reply.result:
                for dto in reply.result:
                    rules_list.append(self._map_rules_from_dto(dto))

            # Extract cursor (may be bytes or str from OpenAPI)
            raw_cursor = getattr(reply, "cursor", None)
            next_cursor: Optional[str] = None
            if raw_cursor:
                next_cursor = raw_cursor.decode("utf-8") if isinstance(raw_cursor, bytes) else raw_cursor

            return GovernanceRulesHistoryResult(
                rules=rules_list,
                cursor=next_cursor,
                total_items=getattr(reply, "total_items", None) or None,
            )
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

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
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

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

        # Verify signatures
        verify_governance_rules(
            rules,
            self._min_valid_signatures,
            self._super_admin_keys,
        )

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
