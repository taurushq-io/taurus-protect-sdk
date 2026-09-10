"""
Payload accessors must refuse metadata that verification has not cleared.

Reading the payload used to succeed regardless of ``hash_verified``, so a caller
could not tell a verified request from one that merely came back from the API.
"""

import json

import pytest

from taurus_protect.errors import RequestMetadataError, UnverifiedMetadataError
from taurus_protect.models.request import RequestMetadata

PAYLOAD = json.dumps(
    [
        {"key": "currency", "value": "XLM"},
        {"key": "source", "value": {"payload": {"address": "SRC"}}},
        {"key": "destination", "value": {"payload": {"address": "DST"}}},
        {"key": "amount", "value": {"valueFrom": "2", "valueTo": "0.1", "rate": "0.05"}},
    ]
)


class TestUnverifiedAccess:
    def test_source_address_refuses_unverified_payload(self):
        metadata = RequestMetadata(payload_as_string=PAYLOAD, hash="deadbeef")
        with pytest.raises(UnverifiedMetadataError):
            metadata.get_source_address()

    def test_destination_address_refuses_unverified_payload(self):
        metadata = RequestMetadata(payload_as_string=PAYLOAD, hash="deadbeef")
        with pytest.raises(UnverifiedMetadataError):
            metadata.get_destination_address()

    def test_amount_refuses_unverified_payload(self):
        metadata = RequestMetadata(payload_as_string=PAYLOAD, hash="deadbeef")
        with pytest.raises(UnverifiedMetadataError):
            metadata.get_amount()

    def test_unverified_is_a_request_metadata_error(self):
        """Existing ``except RequestMetadataError`` blocks must keep working."""
        assert issubclass(UnverifiedMetadataError, RequestMetadataError)


class TestVerifiedAccess:
    def test_accessors_serve_data_once_verified(self):
        metadata = RequestMetadata(
            payload_as_string=PAYLOAD, hash="deadbeef", hash_verified=True
        )
        assert metadata.get_source_address() == "SRC"
        assert metadata.get_destination_address() == "DST"
        assert metadata.get_amount().value_from == "2"

    def test_absent_field_is_none_not_an_error(self):
        """'not in the payload' and 'not verified' must stay distinguishable."""
        metadata = RequestMetadata(
            payload_as_string=json.dumps([{"key": "currency", "value": "XLM"}]),
            hash="deadbeef",
            hash_verified=True,
        )
        assert metadata.get_source_address() is None

    def test_metadata_without_a_payload_is_not_an_error(self):
        """A request in an early status has nothing to verify and nothing to read."""
        metadata = RequestMetadata()
        assert metadata.get_source_address() is None
