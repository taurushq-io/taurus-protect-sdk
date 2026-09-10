"""The asset flow must reject a DTO whose chain disagrees with the signed payload.

`asset.blockchain` is sourced FROM the payload, so passing it as the DTO side fed
`resolve_rule_key` the payload on both sides and its mismatch check could never fire.
Nothing binds the DTO to the signatures, and which chain is chosen selects which
governance rules judge the asset.
"""

from __future__ import annotations

import pytest

from taurus_protect.errors import IntegrityError
from taurus_protect.helpers.whitelisted_asset_verifier import (
    AssetVerificationResult,
    WhitelistedAssetVerifier,
)

from tests.unit.helpers.test_whitelisted_asset_verifier import (
    _build_full_asset_envelope,
    superadmin_keys,  # noqa: F401 — pytest fixture
    user1_keys,  # noqa: F401 — pytest fixture
)


def test_dto_chain_disagreeing_with_payload_is_rejected(superadmin_keys, user1_keys):  # noqa: F811
    sa_priv, sa_pub = superadmin_keys
    u_priv, u_pub = user1_keys
    asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)
    verifier = WhitelistedAssetVerifier([sa_pub], min_valid_signatures=1)

    # The envelope's payload says ETH; the surrounding response claims TRX.
    with pytest.raises(IntegrityError):
        verifier.verify_whitelisted_asset(
            asset, rc_dec, us_dec, dto_blockchain="TRX", dto_network="mainnet"
        )


def test_dto_chain_agreeing_with_payload_still_verifies(superadmin_keys, user1_keys):  # noqa: F811
    """Guards against over-tightening: agreement must not itself be an error."""
    sa_priv, sa_pub = superadmin_keys
    u_priv, u_pub = user1_keys
    asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)
    verifier = WhitelistedAssetVerifier([sa_pub], min_valid_signatures=1)

    result = verifier.verify_whitelisted_asset(
        asset, rc_dec, us_dec, dto_blockchain="ETH", dto_network="mainnet"
    )
    assert isinstance(result, AssetVerificationResult)


def test_absent_dto_chain_falls_back_to_the_payload(superadmin_keys, user1_keys):  # noqa: F811
    """A response carrying no chain must not be treated as a mismatch."""
    sa_priv, sa_pub = superadmin_keys
    u_priv, u_pub = user1_keys
    asset, rc_dec, us_dec = _build_full_asset_envelope(u_priv, u_pub, sa_priv, sa_pub)
    verifier = WhitelistedAssetVerifier([sa_pub], min_valid_signatures=1)

    result = verifier.verify_whitelisted_asset(asset, rc_dec, us_dec)
    assert isinstance(result, AssetVerificationResult)
