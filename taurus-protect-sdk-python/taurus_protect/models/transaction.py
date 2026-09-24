"""Transaction models for Taurus-PROTECT SDK."""

from __future__ import annotations

from datetime import datetime
from typing import List, Optional

from pydantic import BaseModel, Field


class Transaction(BaseModel):
    """
    Confirmed blockchain transaction.

    Attributes:
        id: Unique transaction identifier.
        request_id: Associated request ID.
        wallet_id: Wallet ID.
        address_id: Address ID.
        currency: Currency symbol.
        blockchain: Blockchain name.
        tx_hash: Blockchain transaction hash.
        block_height: Block number.
        block_hash: Block hash.
        amount: Transaction amount.
        fee: Transaction fee.
        direction: Transaction direction (in/out).
        status: Transaction status.
        confirmations: Number of confirmations.
        created_at: When the transaction was recorded.
        confirmed_at: When the transaction was confirmed.
    """

    id: str = Field(description="Unique transaction identifier")
    request_id: Optional[str] = Field(default=None, description="Associated request ID")
    wallet_id: Optional[str] = Field(default=None, description="Wallet ID")
    address_id: Optional[str] = Field(default=None, description="Address ID")
    currency: Optional[str] = Field(default=None, description="Currency symbol")
    blockchain: Optional[str] = Field(default=None, description="Blockchain name")
    tx_hash: Optional[str] = Field(default=None, description="Blockchain transaction hash")
    block_height: Optional[int] = Field(default=None, description="Block number")
    block_hash: Optional[str] = Field(default=None, description="Block hash")
    amount: Optional[str] = Field(default=None, description="Transaction amount")
    fee: Optional[str] = Field(default=None, description="Transaction fee")
    direction: Optional[str] = Field(default=None, description="Direction (in/out)")
    status: Optional[str] = Field(default=None, description="Transaction status")
    confirmations: int = Field(default=0, description="Number of confirmations")
    created_at: Optional[datetime] = Field(default=None, description="Record timestamp")
    confirmed_at: Optional[datetime] = Field(default=None, description="Confirmation timestamp")

    model_config = {"frozen": True}


class TransactionExport(BaseModel):
    """
    One transaction export.

    The export cannot page: the server always starts at the first matching row. When
    ``total_items`` exceeds the rows exported, raise the limit or narrow the filters.

    Attributes:
        content: The exported text, in the requested format.
        total_items: The server's count of matching transactions.
    """

    content: str = Field(default="", description="Exported text")
    total_items: int = Field(default=0, description="Matching transactions on the server")

    model_config = {"frozen": True}
