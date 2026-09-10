"""
The rules-container hash label must match validatord's convention exactly.

In normalized list mode a row picks its rules container out of ``rules_containers`` by
``rules_container_hash`` — and both the label and the container arrive in the same
response. Recomputing the label is what stops a server filing container A under
container B's label and steering any row to any *other* validly-signed container: an
older ruleset with a weaker group threshold, say. Both would pass the SuperAdmin check,
so signature verification alone does not catch it.

The convention is one line in validatord's whitelist controller::

    h := base64.StdEncoding.EncodeToString(crypto.Sha256([]byte(e.GetRulesContainer())))

``GetRulesContainer()`` is ALREADY a base64 string there, so the digest is over the
base64 TEXT, not the decoded protobuf, and the output is base64 rather than hex. An
implementation that decodes first and hashes the protobuf, or emits hex, would reject
every container and break every list call — which is why the convention is pinned here
directly rather than only through its use.

Do not confuse it with ``enforcedRulesHash``: that one IS base64(SHA256(raw protobuf)),
but it is a backlink to a ruleset's predecessor rather than its own identity, and it is
never verified server-side. Both are 44-character base64 SHA-256 digests, so the mix-up
is silent.
"""

from taurus_protect.services.whitelisted_address_service import _container_hash_label


def test_label_is_base64_sha256_of_the_base64_text() -> None:
    # echo -n "abc" | sha256sum -> ba7816bf...; base64 of those raw digest bytes:
    assert _container_hash_label("abc") == "ungWv48Bz+pBQUDeXa4iI7ADYaOWF3qctBD/YfIAFa0="


def test_label_is_not_hex() -> None:
    # A hex digest would be 64 chars of [0-9a-f]; the label is 44 chars of base64.
    label = _container_hash_label("abc")
    assert len(label) == 44
    assert label.endswith("=")


def test_all_four_sdks_agree_on_this_vector() -> None:
    # The same input is pinned in Go (TestContainerHashLabel_MatchesValidatordConvention),
    # Java (WhitelistedAddressExclusionTest) and TypeScript, so a drift in any one SDK
    # shows up as a single failing assertion rather than as rejected containers at runtime.
    assert _container_hash_label("abc") == "ungWv48Bz+pBQUDeXa4iI7ADYaOWF3qctBD/YfIAFa0="
