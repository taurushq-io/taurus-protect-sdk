"""Asset v2 models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import Dict, List, Optional

from pydantic import BaseModel, Field

from taurus_protect.models.pagination import CursorPage
from taurus_protect.models.whitelisted_address import ExcludedWhitelistedAddress


class AssetV2(BaseModel):
    """
    An asset definition from the v2 asset API.

    Attributes:
        id: Unique asset identifier.
        label: Asset label.
        asset_type: Asset type.
        status: Asset status.
        blockchain: Blockchain name.
        network: Network name.
        currency_id: Currency identifier.
        name: Asset name.
        symbol: Asset symbol.
        decimals: Number of decimal places, as the server sent it.
        contract_address: Token contract address.
        attributes: Key/value attributes.
        version: Definition version.
        created_at: When the asset was created.
        updated_at: When the asset was last updated.
    """

    id: str = Field(description="Unique asset identifier")
    label: Optional[str] = Field(default=None, description="Asset label")
    asset_type: Optional[str] = Field(default=None, description="Asset type")
    status: Optional[str] = Field(default=None, description="Asset status")
    blockchain: Optional[str] = Field(default=None, description="Blockchain name")
    network: Optional[str] = Field(default=None, description="Network name")
    currency_id: Optional[str] = Field(default=None, description="Currency ID")
    name: Optional[str] = Field(default=None, description="Asset name")
    symbol: Optional[str] = Field(default=None, description="Asset symbol")
    decimals: Optional[str] = Field(default=None, description="Decimal places")
    contract_address: Optional[str] = Field(default=None, description="Contract address")
    attributes: Dict[str, str] = Field(default_factory=dict, description="Attributes")
    version: Optional[str] = Field(default=None, description="Definition version")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}


class AssetAddressV2(BaseModel):
    """
    An address holding an asset, from the v2 asset API.

    The v2 row carries no signature. ``verified`` is True only for an INTERNAL or
    WHITELISTED row whose address the SDK re-read and verified (HSM signature, or the
    whitelist's 6-step chain); its ``address`` then comes from that verified read. A row
    with ``verified`` False is on-chain data (an EXTERNAL holder), never a Taurus-PROTECT
    address.

    Attributes:
        address: The blockchain address; verified only when ``verified`` is True.
        address_id: The internal address ID, when the address is one of the tenant's.
        whitelisted_address_id: The whitelisted address ID, when it is whitelisted.
        address_type: Address type.
        kyc_status: KYC status.
        balance: Balance of the asset at this address.
        verified: Whether ``address`` comes from a verified read.
    """

    address: Optional[str] = Field(default=None, description="Blockchain address")
    address_id: Optional[str] = Field(default=None, description="Internal address ID")
    whitelisted_address_id: Optional[str] = Field(
        default=None, description="Whitelisted address ID"
    )
    address_type: Optional[str] = Field(default=None, description="Address type")
    kyc_status: Optional[str] = Field(default=None, description="KYC status")
    balance: Optional[str] = Field(default=None, description="Asset balance")
    verified: bool = Field(default=False, description="Whether the address was verified")

    model_config = {"frozen": True}


class QueryAssetAddressesResult(BaseModel):
    """
    One page of the addresses holding a v2 asset.

    Attributes:
        addresses: The rows kept: verified INTERNAL/WHITELISTED rows and unverified
            EXTERNAL ones.
        page: The page window; exclusions never move the cursor.
        excluded_unverified: INTERNAL/WHITELISTED rows that could not be verified, with the
            reason. The id is the address or whitelisted address ID, else the address.
    """

    addresses: List[AssetAddressV2] = Field(default_factory=list)
    page: CursorPage = Field(default_factory=CursorPage)
    excluded_unverified: List[ExcludedWhitelistedAddress] = Field(default_factory=list)

    model_config = {"frozen": True}


class AssetOperationV2(BaseModel):
    """
    An operation on an asset (create, mint, burn, ...), from the v2 asset API.

    Attributes:
        id: Unique operation identifier.
        asset_id: The asset operated on.
        type: Operation type.
        status: Operation status.
        initiated_by_address_id: The address that initiated the operation.
        failure_reason: Why the operation failed, if it did.
        blocking_reason: Why the operation is blocked, if it is.
        created_at: When the operation was created.
        updated_at: When the operation was last updated.
    """

    id: str = Field(description="Unique operation identifier")
    asset_id: Optional[str] = Field(default=None, description="Asset ID")
    type: Optional[str] = Field(default=None, description="Operation type")
    status: Optional[str] = Field(default=None, description="Operation status")
    initiated_by_address_id: Optional[str] = Field(
        default=None, description="Initiating address ID"
    )
    failure_reason: Optional[str] = Field(default=None, description="Failure reason")
    blocking_reason: Optional[str] = Field(default=None, description="Blocking reason")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}
