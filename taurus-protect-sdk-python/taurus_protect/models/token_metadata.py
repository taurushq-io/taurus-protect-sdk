"""Token metadata domain models for Taurus-PROTECT SDK."""

from __future__ import annotations

from dataclasses import dataclass
from typing import Optional


@dataclass
class TokenMetadata:
    """
    Represents token metadata for ERC tokens.

    Contains information about a token including its name, description,
    and optional media data.
    """

    name: str = ""
    description: str = ""
    decimals: str = ""
    uri: str = ""
    data_type: str = ""
    base64_data: str = ""


@dataclass
class FATokenMetadata:
    """
    Represents token metadata for FA tokens (Tezos).

    Contains information about a Tezos FA1.2 or FA2 token.
    """

    name: str = ""
    symbol: str = ""
    decimals: str = ""
    description: str = ""
    uri: str = ""
    data_type: str = ""
    base64_data: str = ""


@dataclass
class CryptoPunkMetadata:
    """
    Represents metadata for CryptoPunk tokens.

    Contains information about a CryptoPunk NFT including attributes.
    """

    punk_index: str = ""
    image_url: str = ""
    attributes: Optional[dict] = None
