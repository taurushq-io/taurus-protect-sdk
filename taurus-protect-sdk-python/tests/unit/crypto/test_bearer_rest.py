"""Tests for BearerRESTClient and bearer Credentials wiring."""

from unittest.mock import MagicMock, patch

import pytest
from cryptography.hazmat.primitives import serialization
from cryptography.hazmat.primitives.asymmetric import ec

from taurus_protect import ConfigurationError, Credentials, ProtectClient
from taurus_protect.crypto.bearer_rest import BearerRESTClient


def _super_admin_pem() -> str:
    priv = ec.generate_private_key(ec.SECP256R1())
    return priv.public_key().public_bytes(
        encoding=serialization.Encoding.PEM,
        format=serialization.PublicFormat.SubjectPublicKeyInfo,
    ).decode()


def test_create_bearer_client_uses_bearer_credentials():
    """Bearer credentials authenticate without api_key/api_secret."""
    client = ProtectClient.create(
        host="https://api.example.com",
        credentials=Credentials.bearer_token("tok"),
        super_admin_keys_pem=[_super_admin_pem()],
    )
    assert not client.is_closed
    assert type(client._credentials).__name__ == "_BearerCredentials"


def test_bearer_client_requires_super_admin_keys():
    """SuperAdmin keys are mandatory even under bearer auth."""
    with pytest.raises(ConfigurationError):
        ProtectClient.create(
            host="https://api.example.com",
            credentials=Credentials.bearer_token_provider(lambda: "tok"),
        )


def test_create_hmac_client_still_requires_super_admin_keys():
    """The default (HMAC) path must still mandate SuperAdmin keys."""
    with pytest.raises(ConfigurationError):
        ProtectClient.create(
            host="https://api.example.com",
            api_key="k",
            api_secret="0" * 64,
        )


class TestBearerRESTClient:
    """Tests for BearerRESTClient."""

    @pytest.fixture
    def mock_configuration(self):
        """Create a mock Configuration object."""
        config = MagicMock()
        config.verify_ssl = True
        config.ssl_ca_cert = None
        config.cert_file = None
        config.key_file = None
        config.assert_hostname = None
        config.retries = None
        config.tls_server_name = None
        config.socket_options = None
        config.connection_pool_maxsize = None
        config.proxy = None
        config.proxy_headers = None
        return config

    def test_init(self, mock_configuration):
        client = BearerRESTClient(mock_configuration, lambda: "tok")
        assert client._token_provider() == "tok"

    def test_request_adds_bearer_authorization_header(self, mock_configuration):
        client = BearerRESTClient(mock_configuration, lambda: "session-token-xyz")

        with patch.object(
            BearerRESTClient.__bases__[0], "request", return_value=MagicMock()
        ) as mock_request:
            client.request(
                method="GET",
                url="https://api.example.com/v1/wallets",
                headers={},
            )

            mock_request.assert_called_once()
            call_kwargs = mock_request.call_args
            headers = call_kwargs.kwargs.get("headers") or call_kwargs[1].get("headers")
            assert headers["Authorization"] == "Bearer session-token-xyz"

    def test_request_with_no_headers(self, mock_configuration):
        client = BearerRESTClient(mock_configuration, lambda: "tok")

        with patch.object(
            BearerRESTClient.__bases__[0], "request", return_value=MagicMock()
        ) as mock_request:
            client.request(
                method="GET",
                url="https://api.example.com/v1/wallets",
                headers=None,
            )

            mock_request.assert_called_once()
            call_kwargs = mock_request.call_args
            headers = call_kwargs.kwargs.get("headers") or call_kwargs[1].get("headers")
            assert headers is not None
            assert headers["Authorization"] == "Bearer tok"
