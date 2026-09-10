"""Contract whitelisting service for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import TYPE_CHECKING, Any, List, Optional

from taurus_protect._internal.openapi.exceptions import ApiException
from taurus_protect.services._base import BaseService

if TYPE_CHECKING:
    from taurus_protect._internal.openapi.api.contract_whitelisting_api import (
        ContractWhitelistingApi,
    )


class ContractWhitelistingService(BaseService):
    """
    Service for managing whitelisted smart contracts.

    Whitelisted contracts are pre-approved smart contracts that
    can be interacted with through transaction requests.
    """

    def __init__(
        self,
        api_client: Any,
        contract_whitelisting_api: "ContractWhitelistingApi",
    ) -> None:
        super().__init__(api_client)
        self._api = contract_whitelisting_api

    def create(
        self,
        address: str,
        name: str,
        blockchain: str,
        network: Optional[str] = None,
        abi: Optional[str] = None,
    ) -> int:
        """Create a whitelisted contract request."""
        self._validate_required(address, "address")
        self._validate_required(name, "name")
        self._validate_required(blockchain, "blockchain")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_create_whitelisted_contract_request import (
                TgvalidatordCreateWhitelistedContractRequest,
            )

            body = TgvalidatordCreateWhitelistedContractRequest(
                address=address,
                name=name,
                blockchain=blockchain,
                network=network,
                abi=abi,
            )

            reply = self._api.whitelist_service_create_whitelisted_contract(body=body)
            return int(reply.result) if reply.result else 0
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def delete(self, contract_id: int) -> None:
        """Delete a whitelisted contract."""
        if contract_id <= 0:
            raise ValueError("contract_id must be positive")

        try:
            self._api.whitelist_service_delete_whitelisted_contract(str(contract_id))
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def approve_whitelisted_contracts(
        self,
        contract_ids: List[str],
        signature: str,
        comment: Optional[str] = None,
    ) -> None:
        """
        Approve whitelisted contracts with a signature.

        Requires a cryptographic signature computed over the metadata hashes
        of the contracts being approved.

        .. deprecated::
            The signature is an opaque blob over hashes nothing verified, so the caller
            cannot know what they signed. Use ``WhitelistedAssetService.approve``, which
            re-reads and verifies the rows first.

        Args:
            contract_ids: List of contract IDs to approve.
            signature: The approval signature (base64-encoded).
            comment: Optional approval comment.

        Raises:
            ValueError: If contract_ids is empty or signature is empty.
            APIError: If the API request fails.
        """
        if not contract_ids:
            raise ValueError("contract_ids cannot be empty")
        self._validate_required(signature, "signature")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_approve_whitelisted_contract_address_request import (
                TgvalidatordApproveWhitelistedContractAddressRequest,
            )

            body = TgvalidatordApproveWhitelistedContractAddressRequest(
                ids=contract_ids,
                signature=signature,
                comment=comment or "",
            )

            self._api.whitelist_service_approve_whitelisted_contract(body=body)
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def create_attribute(
        self,
        contract_id: str,
        key: str,
        value: str,
    ) -> None:
        """
        Create an attribute on a whitelisted contract.

        Args:
            contract_id: The whitelisted contract ID.
            key: The attribute key.
            value: The attribute value.

        Raises:
            ValueError: If any required parameter is empty.
            APIError: If the API request fails.
        """
        self._validate_required(contract_id, "contract_id")
        self._validate_required(key, "key")

        try:
            from taurus_protect._internal.openapi.models.tgvalidatord_create_whitelisted_contract_address_attribute_request import (
                TgvalidatordCreateWhitelistedContractAddressAttributeRequest,
            )
            from taurus_protect._internal.openapi.models.whitelist_service_create_whitelisted_contract_attributes_body import (
                WhitelistServiceCreateWhitelistedContractAttributesBody,
            )

            attr_request = TgvalidatordCreateWhitelistedContractAddressAttributeRequest(
                key=key,
                value=value,
            )

            body = WhitelistServiceCreateWhitelistedContractAttributesBody(
                attributes=[attr_request],
            )

            self._api.whitelist_service_create_whitelisted_contract_attributes(
                whitelisted_contract_address_id=contract_id,
                body=body,
            )
        except Exception as e:
            if isinstance(e, ApiException):
                raise self._handle_error(e)
            raise

    def get_attribute(
        self,
        contract_id: str,
        key: str,
    ) -> Optional[str]:
        """
        Get an attribute value from a whitelisted contract.

        Args:
            contract_id: The whitelisted contract ID.
            key: The attribute ID to retrieve.

        Returns:
            The attribute value, or None if not found.

        Raises:
            ValueError: If any required parameter is empty.
            APIError: If the API request fails.
        """
        self._validate_required(contract_id, "contract_id")
        self._validate_required(key, "key")

        try:
            reply = self._api.whitelist_service_get_whitelisted_contract_attribute(
                whitelisted_contract_address_id=contract_id,
                id=key,
            )
            result = reply.result
            if result is None:
                return None
            return getattr(result, "value", None)
        except Exception as e:
            if "ApiException" in type(e).__name__:
                # Check for 404 not found
                status = getattr(e, "status", 500)
                if status == 404:
                    return None
                raise self._handle_error(e)
            raise
