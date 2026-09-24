"""Earn models for Taurus-PROTECT SDK."""

from __future__ import annotations

from typing import Optional

from pydantic import BaseModel, Field


class EarnReward(BaseModel):
    """
    A reward earned by an address.

    Attributes:
        id: Unique reward identifier.
        recipient_address_id: The internal address ID receiving the reward.
        recipient_address: The receiving blockchain address.
        reward_type: Reward type (e.g. ``RewardTypeMerklToken``).
        amount: Total reward amount.
        claimed: Amount already claimed.
        pending: Amount still pending.
        token_address: Contract address of the reward token.
        token_symbol: Symbol of the reward token.
        token_asset_id: Asset ID of the reward token.
    """

    id: str = Field(description="Unique reward identifier")
    recipient_address_id: Optional[str] = Field(default=None, description="Recipient address ID")
    recipient_address: Optional[str] = Field(default=None, description="Recipient address")
    reward_type: Optional[str] = Field(default=None, description="Reward type")
    amount: Optional[str] = Field(default=None, description="Total reward amount")
    claimed: Optional[str] = Field(default=None, description="Claimed amount")
    pending: Optional[str] = Field(default=None, description="Pending amount")
    token_address: Optional[str] = Field(default=None, description="Reward token address")
    token_symbol: Optional[str] = Field(default=None, description="Reward token symbol")
    token_asset_id: Optional[str] = Field(default=None, description="Reward token asset ID")

    model_config = {"frozen": True}
