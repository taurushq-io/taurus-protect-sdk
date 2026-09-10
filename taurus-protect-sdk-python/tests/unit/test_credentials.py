"""Tests for the Credentials authentication mechanisms."""

import pytest

from taurus_protect import ConfigurationError, Credentials


def test_api_key_rejects_empty():
    with pytest.raises(ConfigurationError):
        Credentials.api_key("", "abcdef")
    with pytest.raises(ConfigurationError):
        Credentials.api_key("k", "")


def test_api_key_rejects_non_hex_secret():
    with pytest.raises(ConfigurationError):
        Credentials.api_key("k", "not-hex")


def test_bearer_token_rejects_empty():
    with pytest.raises(ConfigurationError):
        Credentials.bearer_token("")


def test_bearer_token_provider_rejects_none():
    with pytest.raises(ConfigurationError):
        Credentials.bearer_token_provider(None)


def test_api_key_close_wipes_secret():
    creds = Credentials.api_key("k", "00" * 32)
    creds.close()  # wipes the TPV1 secret; must not raise


def test_bearer_close_is_noop():
    Credentials.bearer_token("tok").close()  # no-op; must not raise
