"""Integration tests for BlockchainService.

These tests verify blockchain operations against a live Taurus-PROTECT API.
"""

from __future__ import annotations

import pytest

from taurus_protect.client import ProtectClient


@pytest.mark.integration
def test_list_blockchains(client: ProtectClient) -> None:
    """Test listing blockchains."""
    blockchains = client.blockchains.list()

    print(f"Found {len(blockchains)} blockchains")

    for blockchain in blockchains:
        print(
            f"Blockchain: Symbol={blockchain.name}, Name={blockchain.display_name}, Network={blockchain.network}"
        )
