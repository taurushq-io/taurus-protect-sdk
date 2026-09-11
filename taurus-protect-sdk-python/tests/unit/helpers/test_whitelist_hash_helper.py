"""Duplicate-JSON-key rejection on every path that parses a signed payload."""

from __future__ import annotations

import pytest

from taurus_protect.errors import IntegrityError
from taurus_protect.helpers.whitelist_hash_helper import (
    parse_whitelisted_address_from_json,
    parse_whitelisted_asset_identity_from_json,
    resolve_rule_key,
)


class TestResolveRuleKeyRefusesDuplicateKeys:
    """``resolve_rule_key`` is the THIRD place a signed payload is parsed, and the one with
    the widest consequence: ``currency`` and ``network`` decide WHICH governance rule judges
    the row, so ``json.loads``' last-duplicate-wins could steer an address to a tier with a
    weaker quorum.

    Go, Java and TypeScript all guarded this site. This SDK used a bare ``json.loads`` here
    while both of its parse functions went through the strict loader -- found by verifying
    the fix site-by-site across the four SDKs rather than by a failing test, which is the
    drift pattern this repo keeps producing.
    """

    def test_a_duplicated_currency_is_refused(self) -> None:
        payload = '{"address":"0xREAL","currency":"ETH","currency":"XTZ"}'

        with pytest.raises(IntegrityError, match="duplicate key"):
            resolve_rule_key(payload, "ETH", "mainnet")

    def test_a_duplicated_network_is_refused(self) -> None:
        payload = (
            '{"address":"0xREAL","currency":"ETH","network":"mainnet","network":"goerli"}'
        )

        with pytest.raises(IntegrityError, match="duplicate key"):
            resolve_rule_key(payload, "ETH", "mainnet")

    def test_a_clean_payload_still_resolves(self) -> None:
        """The guard must not be a blanket refusal."""
        payload = '{"address":"0xREAL","currency":"ETH","network":"mainnet"}'

        assert resolve_rule_key(payload, "ETH", "mainnet") == ("ETH", "mainnet")


class TestParseFunctionsRefuseDuplicateKeys:
    """The other two sites, so all three are pinned in one place.

    ``json.loads`` keeps the LAST of two duplicate keys and reports no error. Combined with
    the legacy-hash tolerance -- whose strips are not injective -- that let a server append a
    duplicate member immediately before the closing brace, have the strip recover the
    genuinely signed bytes, pass every signature check, and hand the appended value back to
    the caller as verified.
    """

    def test_a_duplicated_address_is_refused(self) -> None:
        with pytest.raises(IntegrityError, match="duplicate key"):
            parse_whitelisted_address_from_json(
                '{"address":"0xREAL","address":"0xATTACKER"}'
            )

    def test_a_duplicated_label_is_refused(self) -> None:
        with pytest.raises(IntegrityError, match="duplicate key"):
            parse_whitelisted_address_from_json(
                '{"address":"0xREAL","label":"treasury","label":"ATTACKER"}'
            )

    def test_a_duplicated_nested_label_is_refused(self) -> None:
        # The strip removes inner linkedInternalAddresses labels wholesale, so an injection
        # inside one of those objects is reachable on any legacy-era row.
        with pytest.raises(IntegrityError, match="duplicate key"):
            parse_whitelisted_address_from_json(
                '{"address":"0xREAL","linkedInternalAddresses":'
                '[{"id":"1","label":"ops","label":"ATTACKER"}]}'
            )

    def test_siblings_in_different_objects_may_share_a_key(self) -> None:
        # Per-object key sets, not one global set: two array elements each carrying "label"
        # is ordinary data, and rejecting it would refuse real payloads.
        parsed = parse_whitelisted_address_from_json(
            '{"address":"0xREAL","linkedWallets":[{"label":"a"},{"label":"b"}]}'
        )

        assert parsed.address == "0xREAL"

    def test_a_duplicated_contract_address_is_refused_on_the_asset_side(self) -> None:
        with pytest.raises(IntegrityError, match="duplicate key"):
            parse_whitelisted_asset_identity_from_json(
                '{"name":"USDC","contractAddress":"0xREAL",'
                '"contractAddress":"0xATTACKER"}'
            )
