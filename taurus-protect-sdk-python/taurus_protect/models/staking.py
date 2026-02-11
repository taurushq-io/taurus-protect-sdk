"""Staking and fee models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from decimal import Decimal
from typing import Optional

from pydantic import BaseModel, Field


class Validator(BaseModel):
    """
    Blockchain validator information.

    Represents a validator node that can be used for staking operations.

    Attributes:
        id: Unique validator identifier.
        name: Human-readable validator name.
        blockchain: Blockchain type (e.g., "ETH", "SOL", "ADA").
        network: Network identifier (e.g., "mainnet", "testnet").
        address: Validator's blockchain address.
        commission: Commission rate as a percentage.
        total_stake: Total amount staked with this validator.
        active: Whether the validator is currently active.
        status: Validator status (e.g., "active", "inactive", "jailed").
    """

    id: str = Field(description="Unique validator identifier")
    name: str = Field(default="", description="Human-readable validator name")
    blockchain: str = Field(description="Blockchain type")
    network: str = Field(description="Network identifier")
    address: str = Field(default="", description="Validator's blockchain address")
    commission: Optional[Decimal] = Field(default=None, description="Commission rate as percentage")
    total_stake: Optional[Decimal] = Field(default=None, description="Total amount staked")
    active: bool = Field(default=True, description="Whether validator is active")
    status: str = Field(default="active", description="Validator status")

    model_config = {"frozen": True}


class StakingInfo(BaseModel):
    """
    Staking information for an address.

    Contains details about staking positions and rewards for a specific address.

    Attributes:
        address_id: The address ID this staking info belongs to.
        blockchain: Blockchain type.
        network: Network identifier.
        validator_id: ID of the validator being staked with.
        validator_address: Address of the validator.
        staked_amount: Total amount currently staked.
        rewards: Accumulated staking rewards.
        status: Staking status (e.g., "staked", "unstaking", "withdrawn").
        staked_at: When the stake was initiated.
        unbonding_at: When unbonding will complete (if unstaking).
    """

    address_id: str = Field(description="Address ID this staking info belongs to")
    blockchain: str = Field(default="", description="Blockchain type")
    network: str = Field(default="", description="Network identifier")
    validator_id: Optional[str] = Field(default=None, description="Validator ID")
    validator_address: Optional[str] = Field(default=None, description="Validator address")
    staked_amount: Optional[Decimal] = Field(default=None, description="Total staked amount")
    rewards: Optional[Decimal] = Field(default=None, description="Accumulated rewards")
    status: str = Field(default="", description="Staking status")
    staked_at: Optional[datetime] = Field(default=None, description="When stake was initiated")
    unbonding_at: Optional[datetime] = Field(default=None, description="When unbonding completes")

    model_config = {"frozen": True}


class FeeEstimate(BaseModel):
    """
    Transaction fee estimate.

    Provides estimated fees for a potential transaction.

    Attributes:
        currency: Currency symbol (e.g., "ETH", "BTC").
        blockchain: Blockchain type.
        network: Network identifier.
        amount: Transaction amount.
        fee_low: Low priority fee estimate.
        fee_medium: Medium priority fee estimate.
        fee_high: High priority fee estimate.
        gas_limit: Estimated gas limit (for EVM chains).
        gas_price: Current gas price (for EVM chains).
    """

    currency: str = Field(description="Currency symbol")
    blockchain: str = Field(default="", description="Blockchain type")
    network: str = Field(default="", description="Network identifier")
    amount: Optional[Decimal] = Field(default=None, description="Transaction amount")
    fee_low: Optional[Decimal] = Field(default=None, description="Low priority fee")
    fee_medium: Optional[Decimal] = Field(default=None, description="Medium priority fee")
    fee_high: Optional[Decimal] = Field(default=None, description="High priority fee")
    gas_limit: Optional[int] = Field(default=None, description="Estimated gas limit")
    gas_price: Optional[Decimal] = Field(default=None, description="Current gas price")

    model_config = {"frozen": True}


class FeePayer(BaseModel):
    """
    Fee payer account.

    A fee payer is an account used to pay transaction fees on behalf of other accounts.
    This is commonly used for gas station networks or meta-transactions.

    Attributes:
        id: Unique fee payer identifier.
        blockchain: Blockchain type.
        network: Network identifier.
        address: Fee payer's blockchain address.
        balance: Current balance of the fee payer.
        status: Fee payer status (e.g., "active", "disabled").
        created_at: When the fee payer was created.
        updated_at: When the fee payer was last updated.
    """

    id: str = Field(description="Unique fee payer identifier")
    blockchain: str = Field(default="", description="Blockchain type")
    network: str = Field(default="", description="Network identifier")
    address: str = Field(default="", description="Fee payer's blockchain address")
    balance: Optional[Decimal] = Field(default=None, description="Current balance")
    status: str = Field(default="active", description="Fee payer status")
    created_at: Optional[datetime] = Field(default=None, description="Creation timestamp")
    updated_at: Optional[datetime] = Field(default=None, description="Last update timestamp")

    model_config = {"frozen": True}


class UnsignedPayload(BaseModel):
    """
    Unsigned transaction payload for air-gap signing.

    Used in air-gap (cold storage) signing workflows where the transaction
    is prepared online but signed offline.

    Attributes:
        request_id: The request ID this payload belongs to.
        payload: The raw unsigned payload bytes as hex.
        hash: Hash of the payload to be signed.
        blockchain: Blockchain type.
        network: Network identifier.
    """

    request_id: str = Field(description="Request ID this payload belongs to")
    payload: str = Field(description="Raw unsigned payload as hex")
    hash: str = Field(default="", description="Hash of payload to be signed")
    blockchain: str = Field(default="", description="Blockchain type")
    network: str = Field(default="", description="Network identifier")

    model_config = {"frozen": True}
