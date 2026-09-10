"""Per-group step-5 threshold: distinct signers, not signature entries.

The signature entries arrive in the server-supplied userSignatures blob, so counting them
let a duplicated entry from one group member satisfy an N-of-M group — promoting an
under-approved whitelist entry to approved.
"""

from __future__ import annotations

import json

import pytest
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect.crypto.signing import sign_data
from taurus_protect.helpers.whitelisted_address_verifier import WhitelistedAddressVerifier
from taurus_protect.models.governance_rules import (
    DecodedRulesContainer,
    GroupThreshold,
    RuleGroup,
    RuleUser,
)
from taurus_protect.models.whitelisted_address import (
    WhitelistSignatureEntry,
    WhitelistUserSignature,
)

METADATA_HASH = "abc123"


def _key_pair():
    private_key = ec.generate_private_key(ec.SECP256R1())
    return private_key, private_key.public_key()


def _pem(public_key) -> str:
    from cryptography.hazmat.primitives import serialization

    return public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    ).decode("utf-8")


def _signed_entry(user_id: str, private_key) -> WhitelistSignatureEntry:
    """A genuinely valid approval entry, as it arrives on the wire."""
    hashes = [METADATA_HASH]
    hashes_json = json.dumps(hashes, separators=(",", ":"))
    signature = sign_data(private_key, hashes_json.encode("utf-8"))
    return WhitelistSignatureEntry(
        hashes=hashes,
        user_signature=WhitelistUserSignature(user_id=user_id, signature=signature),
    )


def _container(user_pems: dict) -> DecodedRulesContainer:
    return DecodedRulesContainer(
        users=[
            RuleUser(id=user_id, public_key_pem=pem, roles=["USER"])
            for user_id, pem in user_pems.items()
        ],
        groups=[RuleGroup(id="approvers", user_ids=list(user_pems.keys()))],
    )


def _walk(container, signatures, min_sigs):
    """Drives the group walk directly; returns the failure message, or None on success."""
    verifier = WhitelistedAddressVerifier([], 0)
    threshold = GroupThreshold(group_id="approvers", minimum_signatures=min_sigs)
    precomputed = {
        idx: json.dumps(sig.hashes, separators=(",", ":")).encode("utf-8")
        for idx, sig in enumerate(signatures)
    }
    return verifier._verify_group_threshold(
        threshold, container, signatures, METADATA_HASH, precomputed
    )


def test_duplicate_entry_from_one_member_does_not_meet_two_of_n():
    priv1, pub1 = _key_pair()
    priv2, pub2 = _key_pair()
    container = _container({"user1@bank.com": _pem(pub1), "user2@bank.com": _pem(pub2)})

    entry = _signed_entry("user1@bank.com", priv1)
    assert _walk(container, [entry, entry], 2) is not None


def test_re_signed_entry_from_one_member_does_not_meet_two_of_n():
    priv1, pub1 = _key_pair()
    priv2, pub2 = _key_pair()
    container = _container({"user1@bank.com": _pem(pub1), "user2@bank.com": _pem(pub2)})

    signatures = [
        _signed_entry("user1@bank.com", priv1),
        _signed_entry("user1@bank.com", priv1),
    ]
    # ECDSA is randomized, so re-signing yields a different signature for the same key.
    assert (
        signatures[0].user_signature.signature != signatures[1].user_signature.signature
    )
    assert _walk(container, signatures, 2) is not None


def test_two_distinct_members_meet_two_of_n():
    priv1, pub1 = _key_pair()
    priv2, pub2 = _key_pair()
    container = _container({"user1@bank.com": _pem(pub1), "user2@bank.com": _pem(pub2)})

    signatures = [
        _signed_entry("user1@bank.com", priv1),
        _signed_entry("user2@bank.com", priv2),
    ]
    assert _walk(container, signatures, 2) is None


def test_one_member_meets_one_of_n():
    priv1, pub1 = _key_pair()
    _, pub2 = _key_pair()
    container = _container({"user1@bank.com": _pem(pub1), "user2@bank.com": _pem(pub2)})

    assert _walk(container, [_signed_entry("user1@bank.com", priv1)], 1) is None


def test_two_user_ids_sharing_one_key_count_once():
    """One compromised secret must not satisfy 2-of-N, whatever it is called."""
    priv, pub = _key_pair()
    container = _container({"user1@bank.com": _pem(pub), "user2@bank.com": _pem(pub)})

    signatures = [
        _signed_entry("user1@bank.com", priv),
        _signed_entry("user2@bank.com", priv),
    ]
    assert _walk(container, signatures, 2) is not None


def test_signer_outside_group_contributes_nothing():
    priv1, pub1 = _key_pair()
    _, pub2 = _key_pair()
    priv_out, pub_out = _key_pair()

    container = _container({"user1@bank.com": _pem(pub1), "user2@bank.com": _pem(pub2)})
    container.users.append(
        RuleUser(id="outsider@bank.com", public_key_pem=_pem(pub_out), roles=["USER"])
    )

    signatures = [
        _signed_entry("user1@bank.com", priv1),
        _signed_entry("outsider@bank.com", priv_out),
    ]
    assert _walk(container, signatures, 2) is not None
