"""Price signature verification against the PRICEUPDATER role.

Rate and decimals feed amount conversion, so an unverified price is a wrong number a
caller acts on. Whether prices must be signed is decided by the SuperAdmin-verified
rules container.
"""

from __future__ import annotations

import pytest
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect.crypto.signing import sign_data
from taurus_protect.errors import IntegrityError
from taurus_protect.helpers.price_verifier import (
    price_signed_bytes,
    verify_price,
    verify_prices,
)
from taurus_protect.models.governance_rules import DecodedRulesContainer, RuleUser
from taurus_protect.models.statistics import Price, PriceSignature


def _key_pair():
    private_key = ec.generate_private_key(ec.SECP256R1())
    return private_key, private_key.public_key()


def _pem(public_key) -> str:
    return public_key.public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    ).decode("utf-8")


def _price(**overrides) -> Price:
    fields = {
        "blockchain": "ETH",
        "currency_from": "ETH",
        "currency_to": "USD",
        "decimals": "18",
        "rate": "2500.00",
    }
    fields.update(overrides)
    return Price(**fields)


def _container(roles, public_key) -> DecodedRulesContainer:
    return DecodedRulesContainer(
        users=[RuleUser(id="price@bank.com", public_key_pem=_pem(public_key), roles=roles)],
        groups=[],
    )


def _signed(price: Price, private_key) -> Price:
    signature = sign_data(private_key, price_signed_bytes(price))
    return price.model_copy(
        update={"signatures": [PriceSignature(user_id="price@bank.com", signature=signature)]}
    )


# The canonical form is validatord's CurrencyPrice JSON projection. Pinned by value: a
# reordered or extended dict silently changes every signature this SDK will accept.
def test_signed_bytes_is_the_canonical_projection():
    assert price_signed_bytes(_price()) == (
        b'{"blockchain":"ETH","currencyFrom":"ETH","currencyTo":"USD",'
        b'"decimals":"18","rate":"2500.00"}'
    )


def test_accepts_a_price_updater_signature():
    priv, pub = _key_pair()
    verify_price(_signed(_price(), priv), _container(["PRICEUPDATER"], pub))


def test_rejects_a_signature_from_a_non_price_updater():
    signer_priv, _ = _key_pair()
    _, updater_pub = _key_pair()

    with pytest.raises(IntegrityError):
        verify_price(
            _signed(_price(), signer_priv), _container(["PRICEUPDATER"], updater_pub)
        )


def test_rejects_a_tampered_rate():
    priv, pub = _key_pair()
    price = _signed(_price(), priv)
    # The rate is what a caller converts with.
    tampered = price.model_copy(update={"rate": "1.00"})

    with pytest.raises(IntegrityError):
        verify_price(tampered, _container(["PRICEUPDATER"], pub))


def test_stripped_signatures_are_rejected_when_a_price_updater_exists():
    _, pub = _key_pair()
    with pytest.raises(IntegrityError):
        verify_price(_price(), _container(["PRICEUPDATER"], pub))


def test_passes_through_when_no_price_updater_is_configured():
    """A tenant with no price signer legitimately serves unsigned prices."""
    _, pub = _key_pair()
    verify_price(_price(), _container(["REQUESTAPPROVER"], pub))


def test_requires_a_rules_container():
    with pytest.raises(IntegrityError):
        verify_price(_price(), None)


def test_verify_prices_reports_the_first_failure():
    priv, pub = _key_pair()
    good = _signed(_price(), priv)
    bad = _price(currency_to="EUR")

    with pytest.raises(IntegrityError):
        verify_prices([good, bad], _container(["PRICEUPDATER"], pub))
