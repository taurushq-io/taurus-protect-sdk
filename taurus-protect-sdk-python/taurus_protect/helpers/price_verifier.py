"""Price signature verification.

Rate and decimals feed amount conversion, so an unverified price is a wrong number a
caller acts on.
"""

from __future__ import annotations

import json
from typing import TYPE_CHECKING, List, Optional

from taurus_protect.crypto.signing import verify_signature
from taurus_protect.errors import IntegrityError

if TYPE_CHECKING:
    from taurus_protect.models.governance_rules import DecodedRulesContainer
    from taurus_protect.models.statistics import Price


def price_signed_bytes(price: "Price") -> bytes:
    """Return the exact bytes a price signature covers.

    This is validatord's CurrencyPrice JSON projection: the five fields it serialises, in
    its field order, with no spaces. The key order IS the signed byte sequence — do not
    reorder or add fields.
    """
    canonical = {
        "blockchain": price.blockchain or "",
        "currencyFrom": price.currency_from or "",
        "currencyTo": price.currency_to or "",
        "decimals": price.decimals or "",
        "rate": price.rate or "",
    }
    return json.dumps(canonical, separators=(",", ":")).encode("utf-8")


def verify_price(price: "Price", rules_container: "DecodedRulesContainer") -> None:
    """Check a price against the PRICEUPDATER keys in a verified rules container.

    The container decides whether prices must be signed at all. It is SuperAdmin-verified,
    so "this tenant has no price signer" is trustworthy, whereas "this price carries no
    signatures" is not — which is why a stripped signatures list on a tenant that DOES
    have a PRICEUPDATER is an error rather than a skip.

    Raises:
        IntegrityError: If no PRICEUPDATER signature verifies.
    """
    if price is None:
        raise IntegrityError("price cannot be None")
    if rules_container is None:
        raise IntegrityError("rules container required for price signature verification")

    keys = _price_updater_keys(rules_container)
    if not keys:
        # This tenant does not sign prices.
        return

    if not price.signatures:
        raise IntegrityError(
            f"price {price.currency_from}/{price.currency_to} carries no signatures but "
            "the rules container configures a PRICEUPDATER"
        )

    data = price_signed_bytes(price)
    for sig in price.signatures:
        if not sig.signature:
            continue
        for key in keys:
            try:
                if verify_signature(key, data, sig.signature):
                    return
            except Exception:  # noqa: BLE001 - a malformed signature is just a non-match
                continue

    raise IntegrityError(
        f"no PRICEUPDATER signature verifies for price "
        f"{price.currency_from}/{price.currency_to} "
        f"({len(price.signatures)} signature(s) offered)"
    )


def verify_prices(
    prices: List["Price"], rules_container: "DecodedRulesContainer"
) -> None:
    """Verify each price, raising on the first failure."""
    for price in prices or []:
        verify_price(price, rules_container)


def _price_updater_keys(rules_container: "DecodedRulesContainer") -> List[object]:
    """Public keys of every user carrying the PRICEUPDATER role."""
    keys: List[object] = []
    for user in getattr(rules_container, "users", None) or []:
        roles = getattr(user, "roles", None) or []
        if "PRICEUPDATER" not in roles:
            continue
        pem = getattr(user, "public_key_pem", None)
        if not pem:
            continue
        try:
            keys.append(rules_container.get_user_public_key(pem))
        except ValueError:
            continue
    return keys
