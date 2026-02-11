"""Pledge service for Taurus-PROTECT SDK Taurus Network."""

from __future__ import annotations

import json
from typing import TYPE_CHECKING, Any, List, Optional, Tuple

from cryptography.hazmat.primitives.asymmetric.ec import EllipticCurvePrivateKey

from taurus_protect.crypto.hashing import calculate_hex_hash
from taurus_protect.crypto.signing import sign_data
from taurus_protect.mappers.taurus_network.pledge import (
    pledge_action_from_dto,
    pledge_actions_from_dto,
    pledge_from_dto,
    pledge_withdrawal_from_dto,
    pledge_withdrawals_from_dto,
    pledges_from_dto,
)
from taurus_protect.models.pagination import Pagination
from taurus_protect.models.taurus_network.pledge import (
    AddPledgeCollateralRequest,
    CreatePledgeRequest,
    InitiateWithdrawPledgeRequest,
    ListPledgeActionsOptions,
    ListPledgesOptions,
    ListPledgeWithdrawalsOptions,
    Pledge,
    PledgeAction,
    PledgeWithdrawal,
    RejectPledgeActionsRequest,
    RejectPledgeRequest,
    UpdatePledgeRequest,
    WithdrawPledgeRequest,
)
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    pass  # For OpenAPI types when available


class PledgeService(BaseService):
    """
    Service for managing Taurus Network pledges.

    Provides operations for creating, updating, and managing pledges between
    Taurus Network participants. Pledges represent reserved funds from one
    participant (pledgor) to another (pledgee).

    **Security Features:**
    - Pledge action approval uses ECDSA signatures for cryptographic verification
    - Actions are sorted by ID before signing to ensure consistent signatures

    Example:
        >>> # Create a pledge
        >>> request = CreatePledgeRequest(
        ...     shared_address_id="123",
        ...     currency_id="ETH",
        ...     amount="1000000000000000000",
        ...     pledge_type="PLEDGEE_WITHDRAWALS_RIGHTS",
        ... )
        >>> pledge, action = client.taurus_network.pledges.create_pledge(request)
        >>>
        >>> # Approve pledge actions
        >>> actions, _ = client.taurus_network.pledges.list_pledge_actions_for_approval()
        >>> count = client.taurus_network.pledges.approve_pledge_actions(actions, private_key)
    """

    def __init__(self, api_client: Any, pledge_api: Any) -> None:
        """
        Initialize pledge service.

        Args:
            api_client: The OpenAPI client instance.
            pledge_api: The TaurusNetworkPledgeApi service from OpenAPI client.
        """
        super().__init__(api_client)
        self._pledge_api = pledge_api

    def get_pledge(self, pledge_id: str) -> Pledge:
        """
        Get a pledge by ID.

        Args:
            pledge_id: The pledge ID to retrieve.

        Returns:
            The pledge.

        Raises:
            ValueError: If pledge_id is invalid.
            NotFoundError: If pledge not found.
            APIError: If API request fails.
        """
        self._validate_required(pledge_id, "pledge_id")

        try:
            resp = self._pledge_api.taurus_network_service_get_pledge(pledge_id)

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Pledge {pledge_id} not found")

            pledge = pledge_from_dto(result)
            if pledge is None:
                from taurus_protect.errors import NotFoundError

                raise NotFoundError(f"Pledge {pledge_id} not found")

            return pledge
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_pledges(
        self,
        opts: Optional[ListPledgesOptions] = None,
    ) -> Tuple[List[Pledge], Optional[Pagination]]:
        """
        List pledges with optional filtering.

        Args:
            opts: Optional filtering and pagination options.

        Returns:
            Tuple of (pledges list, pagination info).

        Raises:
            APIError: If API request fails.
        """
        options = opts or ListPledgesOptions()

        try:
            resp = self._pledge_api.taurus_network_service_get_pledges(
                statuses=options.statuses if options.statuses else None,
                currency_id=options.currency_id if options.currency_id else None,
                owner_participant_id=options.participant_id if options.participant_id else None,
                sort_order=options.direction if options.direction else None,
                cursor_page_size=str(options.limit) if options.limit > 0 else None,
            )

            result = getattr(resp, "result", None)
            pledges = pledges_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=getattr(resp, "offset", None),
                limit=options.limit,
            )

            return pledges, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def create_pledge(
        self,
        req: CreatePledgeRequest,
    ) -> Tuple[Pledge, PledgeAction]:
        """
        Create a new pledge.

        Creates a pledge of funds from an internal address to a Taurus Network
        participant. The funds will be reserved until unpledged or withdrawn.
        Returns both the pledge and the action that needs approval.

        Args:
            req: Pledge creation parameters.

        Returns:
            Tuple of (created pledge, pledge action requiring approval).

        Raises:
            ValueError: If required fields are missing.
            APIError: If API request fails.
        """
        if req is None:
            raise ValueError("request cannot be None")
        self._validate_required(req.shared_address_id, "shared_address_id")
        self._validate_required(req.currency_id, "currency_id")
        self._validate_required(req.amount, "amount")

        try:
            # Build request body
            body: dict[str, Any] = {
                "sharedAddressID": req.shared_address_id,
                "currencyID": req.currency_id,
                "amount": req.amount,
            }

            if req.pledge_type:
                body["pledgeType"] = req.pledge_type
            if req.external_reference_id:
                body["externalReferenceId"] = req.external_reference_id
            if req.reconciliation_note:
                body["reconciliationNote"] = req.reconciliation_note
            if req.pledge_duration_setup:
                body["pledgeDurationSetup"] = {}
                if req.pledge_duration_setup.start_date:
                    body["pledgeDurationSetup"][
                        "startDate"
                    ] = req.pledge_duration_setup.start_date.isoformat()
                if req.pledge_duration_setup.end_date:
                    body["pledgeDurationSetup"][
                        "endDate"
                    ] = req.pledge_duration_setup.end_date.isoformat()
            if req.key_value_attributes:
                body["keyValueAttributes"] = [
                    {"key": attr.key, "value": attr.value} for attr in req.key_value_attributes
                ]

            resp = self._pledge_api.taurus_network_service_create_pledge(body=body)

            pledge_result = getattr(resp, "result", None)
            if pledge_result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create pledge: no result returned")

            pledge = pledge_from_dto(pledge_result)
            if pledge is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create pledge: invalid response")

            action_result = getattr(resp, "action", None)
            action = pledge_action_from_dto(action_result)
            if action is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to create pledge: no action returned")

            return pledge, action
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def update_pledge(
        self,
        pledge_id: str,
        req: UpdatePledgeRequest,
    ) -> Pledge:
        """
        Update a pledge's default destination.

        Args:
            pledge_id: The pledge ID to update.
            req: Update parameters.

        Returns:
            The updated pledge.

        Raises:
            ValueError: If pledge_id is invalid.
            APIError: If API request fails.
        """
        self._validate_required(pledge_id, "pledge_id")
        if req is None:
            raise ValueError("request cannot be None")

        try:
            body: dict[str, Any] = {}
            if req.default_destination_shared_address_id:
                body["defaultDestinationSharedAddressID"] = (
                    req.default_destination_shared_address_id
                )
            if req.default_destination_internal_address_id:
                body["defaultDestinationInternalAddressID"] = (
                    req.default_destination_internal_address_id
                )

            resp = self._pledge_api.taurus_network_service_update_pledge(
                pledge_id=pledge_id,
                body=body,
            )

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to update pledge: no result returned")

            pledge = pledge_from_dto(result)
            if pledge is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to update pledge: invalid response")

            return pledge
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def add_pledge_collateral(
        self,
        pledge_id: str,
        req: AddPledgeCollateralRequest,
    ) -> Tuple[Pledge, PledgeAction]:
        """
        Add collateral to an existing pledge.

        Increases the pledged amount. Only the pledgor can call this.
        Returns the updated pledge and a new action requiring approval.

        Args:
            pledge_id: The pledge ID to add collateral to.
            req: Collateral addition parameters.

        Returns:
            Tuple of (updated pledge, pledge action requiring approval).

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(pledge_id, "pledge_id")
        if req is None:
            raise ValueError("request cannot be None")
        self._validate_required(req.amount, "amount")

        try:
            body = {"amount": req.amount}

            resp = self._pledge_api.taurus_network_service_add_pledge_collateral(
                pledge_id=pledge_id,
                body=body,
            )

            pledge_result = getattr(resp, "result", None)
            if pledge_result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to add collateral: no result returned")

            pledge = pledge_from_dto(pledge_result)
            if pledge is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to add collateral: invalid response")

            action_result = getattr(resp, "action", None)
            action = pledge_action_from_dto(action_result)
            if action is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to add collateral: no action returned")

            return pledge, action
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def withdraw_pledge(
        self,
        pledge_id: str,
        req: WithdrawPledgeRequest,
    ) -> Tuple[PledgeWithdrawal, PledgeAction]:
        """
        Withdraw from a pledge (pledgee operation).

        Allows the pledgee (target participant) to withdraw funds from the pledge.
        Returns the withdrawal record and an action requiring approval (unless
        auto-approved based on pledge type).

        Args:
            pledge_id: The pledge ID to withdraw from.
            req: Withdrawal parameters.

        Returns:
            Tuple of (withdrawal record, pledge action).

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(pledge_id, "pledge_id")
        if req is None:
            raise ValueError("request cannot be None")
        self._validate_required(req.amount, "amount")

        try:
            body: dict[str, Any] = {"amount": req.amount}
            if req.destination_shared_address_id:
                body["destinationSharedAddressID"] = req.destination_shared_address_id
            if req.destination_internal_address_id:
                body["destinationInternalAddressID"] = req.destination_internal_address_id
            if req.external_reference_id:
                body["externalReferenceID"] = req.external_reference_id

            resp = self._pledge_api.taurus_network_service_withdraw_pledge(
                pledge_id=pledge_id,
                body=body,
            )

            withdrawal_result = getattr(resp, "result", None)
            if withdrawal_result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to withdraw: no result returned")

            withdrawal = pledge_withdrawal_from_dto(withdrawal_result)
            if withdrawal is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to withdraw: invalid response")

            action_result = getattr(resp, "action", None)
            action = pledge_action_from_dto(action_result)
            if action is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to withdraw: no action returned")

            return withdrawal, action
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def initiate_withdraw_pledge(
        self,
        pledge_id: str,
        req: InitiateWithdrawPledgeRequest,
    ) -> Tuple[PledgeWithdrawal, PledgeAction]:
        """
        Initiate withdrawal from a pledge (pledgor operation).

        Allows the pledgor (owner participant) to initiate a withdrawal.
        This is typically used when the pledgee does not have withdrawal rights.
        Returns the withdrawal record and an action requiring approval.

        Args:
            pledge_id: The pledge ID to withdraw from.
            req: Withdrawal parameters.

        Returns:
            Tuple of (withdrawal record, pledge action).

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(pledge_id, "pledge_id")
        if req is None:
            raise ValueError("request cannot be None")
        self._validate_required(req.amount, "amount")

        try:
            body: dict[str, Any] = {"amount": req.amount}
            if req.destination_shared_address_id:
                body["destinationSharedAddressID"] = req.destination_shared_address_id

            resp = self._pledge_api.taurus_network_service_initiate_withdraw_pledge(
                pledge_id=pledge_id,
                body=body,
            )

            withdrawal_result = getattr(resp, "result", None)
            if withdrawal_result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to initiate withdrawal: no result returned")

            withdrawal = pledge_withdrawal_from_dto(withdrawal_result)
            if withdrawal is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to initiate withdrawal: invalid response")

            action_result = getattr(resp, "action", None)
            action = pledge_action_from_dto(action_result)
            if action is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to initiate withdrawal: no action returned")

            return withdrawal, action
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def unpledge(
        self,
        pledge_id: str,
    ) -> Tuple[Pledge, PledgeAction]:
        """
        Unpledge all funds from a pledge.

        Releases all pledged funds back to the pledgor. Only the pledgor
        can call this. Returns the updated pledge and an action requiring approval.

        Args:
            pledge_id: The pledge ID to unpledge.

        Returns:
            Tuple of (updated pledge, pledge action requiring approval).

        Raises:
            ValueError: If pledge_id is invalid.
            APIError: If API request fails.
        """
        self._validate_required(pledge_id, "pledge_id")

        try:
            resp = self._pledge_api.taurus_network_service_unpledge(
                pledge_id=pledge_id,
            )

            pledge_result = getattr(resp, "result", None)
            if pledge_result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to unpledge: no result returned")

            pledge = pledge_from_dto(pledge_result)
            if pledge is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to unpledge: invalid response")

            action_result = getattr(resp, "action", None)
            action = pledge_action_from_dto(action_result)
            if action is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to unpledge: no action returned")

            return pledge, action
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def reject_pledge(
        self,
        pledge_id: str,
        req: RejectPledgeRequest,
    ) -> Pledge:
        """
        Reject a pledge.

        Rejects a pending pledge. Only the pledgee can reject an incoming pledge.

        Args:
            pledge_id: The pledge ID to reject.
            req: Rejection parameters with comment.

        Returns:
            The rejected pledge.

        Raises:
            ValueError: If arguments are invalid.
            APIError: If API request fails.
        """
        self._validate_required(pledge_id, "pledge_id")
        if req is None:
            raise ValueError("request cannot be None")
        self._validate_required(req.comment, "comment")

        try:
            body = {"comment": req.comment}

            resp = self._pledge_api.taurus_network_service_reject_pledge(
                pledge_id=pledge_id,
                body=body,
            )

            result = getattr(resp, "result", None)
            if result is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to reject pledge: no result returned")

            pledge = pledge_from_dto(result)
            if pledge is None:
                from taurus_protect.errors import APIError

                raise APIError(500, "Failed to reject pledge: invalid response")

            return pledge
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_pledge_actions(
        self,
        opts: Optional[ListPledgeActionsOptions] = None,
    ) -> Tuple[List[PledgeAction], Optional[Pagination]]:
        """
        List all pledge actions with optional filtering.

        Args:
            opts: Optional filtering and pagination options.

        Returns:
            Tuple of (pledge actions list, pagination info).

        Raises:
            APIError: If API request fails.
        """
        options = opts or ListPledgeActionsOptions()

        try:
            resp = self._pledge_api.taurus_network_service_get_pledge_actions(
                statuses=options.statuses,
                action_types=options.action_types,
                pledge_id=options.pledge_id,
                limit=str(options.limit) if options.limit > 0 else None,
                offset=str(options.offset) if options.offset > 0 else None,
            )

            result = getattr(resp, "result", None)
            actions = pledge_actions_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=getattr(resp, "offset", None),
                limit=options.limit,
            )

            return actions, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def list_pledge_actions_for_approval(
        self,
        opts: Optional[ListPledgeActionsOptions] = None,
    ) -> Tuple[List[PledgeAction], Optional[Pagination]]:
        """
        List pledge actions pending approval.

        Returns only actions that require approval from the current user.

        Args:
            opts: Optional filtering and pagination options.

        Returns:
            Tuple of (pledge actions list, pagination info).

        Raises:
            APIError: If API request fails.
        """
        options = opts or ListPledgeActionsOptions()

        try:
            resp = self._pledge_api.taurus_network_service_get_pledge_actions_for_approval(
                action_types=options.action_types,
                pledge_id=options.pledge_id,
                limit=str(options.limit) if options.limit > 0 else None,
                offset=str(options.offset) if options.offset > 0 else None,
            )

            result = getattr(resp, "result", None)
            actions = pledge_actions_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=getattr(resp, "offset", None),
                limit=options.limit,
            )

            return actions, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e

    def approve_pledge_actions(
        self,
        actions: List[PledgeAction],
        private_key: EllipticCurvePrivateKey,
        comment: str = "approving via taurus-protect-sdk-python",
    ) -> int:
        """
        Approve multiple pledge actions with ECDSA signature.

        The actions are sorted by ID, and a signature is computed over
        the concatenated hashes of their metadata using the provided private key.

        **Signature Format:**
        base64(ecdsa_sign(sha256([hex(sha256(action1_metadata)), hex(sha256(action2_metadata)), ...])))

        Args:
            actions: List of pledge actions to approve.
            private_key: ECDSA private key for signing.
            comment: Optional approval comment.

        Returns:
            Number of actions successfully approved.

        Raises:
            ValueError: If actions list is empty or private_key is None.
            APIError: If API request fails.
        """
        if not actions:
            raise ValueError("actions list cannot be empty")
        if private_key is None:
            raise ValueError("private_key cannot be None")

        # Validate all actions have metadata with hash
        for action in actions:
            if action.metadata is None:
                raise ValueError("action metadata cannot be None")
            if not action.metadata.hash:
                raise ValueError("action metadata hash cannot be empty")

        try:
            # Sort actions by ID (string sort, as IDs might be UUIDs)
            sorted_actions = sorted(actions, key=lambda a: a.id)

            # Build concatenated hash string - array of hex hashes
            hashes = [a.metadata.hash for a in sorted_actions]
            to_sign = json.dumps(hashes)

            # Sign with ECDSA
            signature = sign_data(private_key, to_sign.encode("utf-8"))

            # Build API request
            body = {
                "ids": [a.id for a in sorted_actions],
                "comment": comment,
                "signature": signature,
            }

            resp = self._pledge_api.taurus_network_service_approve_pledge_actions(body=body)

            approved_count = getattr(resp, "approved_count", None)
            if approved_count is not None:
                return int(approved_count)

            # If no count returned, assume all were approved
            return len(actions)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def reject_pledge_actions(
        self,
        req: RejectPledgeActionsRequest,
    ) -> int:
        """
        Reject multiple pledge actions.

        Args:
            req: Rejection request with action IDs and comment.

        Returns:
            Number of actions rejected.

        Raises:
            ValueError: If request is invalid.
            APIError: If API request fails.
        """
        if req is None:
            raise ValueError("request cannot be None")
        if not req.ids:
            raise ValueError("ids list cannot be empty")
        self._validate_required(req.comment, "comment")

        try:
            body = {
                "ids": req.ids,
                "comment": req.comment,
            }

            resp = self._pledge_api.taurus_network_service_reject_pledge_actions(body=body)

            rejected_count = getattr(resp, "rejected_count", None)
            if rejected_count is not None:
                return int(rejected_count)

            # If no count returned, assume all were rejected
            return len(req.ids)
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, (APIError, ValueError)):
                raise
            raise self._handle_error(e) from e

    def list_pledge_withdrawals(
        self,
        opts: Optional[ListPledgeWithdrawalsOptions] = None,
    ) -> Tuple[List[PledgeWithdrawal], Optional[Pagination]]:
        """
        List pledge withdrawals with optional filtering.

        Args:
            opts: Optional filtering and pagination options.

        Returns:
            Tuple of (withdrawals list, pagination info).

        Raises:
            APIError: If API request fails.
        """
        options = opts or ListPledgeWithdrawalsOptions()

        try:
            resp = self._pledge_api.taurus_network_service_get_pledges_withdrawals(
                statuses=options.statuses,
                pledge_id=options.pledge_id,
                limit=str(options.limit) if options.limit > 0 else None,
                offset=str(options.offset) if options.offset > 0 else None,
            )

            result = getattr(resp, "result", None)
            withdrawals = pledge_withdrawals_from_dto(result) if result else []

            pagination = self._extract_pagination(
                total_items=getattr(resp, "total_items", None),
                offset=getattr(resp, "offset", None),
                limit=options.limit,
            )

            return withdrawals, pagination
        except Exception as e:
            from taurus_protect.errors import APIError

            if isinstance(e, APIError):
                raise
            raise self._handle_error(e) from e
