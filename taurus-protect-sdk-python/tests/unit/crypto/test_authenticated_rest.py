"""Tests for AuthenticatedRESTClient."""

from unittest.mock import MagicMock, patch

import pytest

from taurus_protect.crypto.authenticated_rest import AuthenticatedRESTClient
from taurus_protect.crypto.tpv1 import TPV1Auth


class TestAuthenticatedRESTClient:
    """Tests for AuthenticatedRESTClient."""

    @pytest.fixture
    def tpv1_auth(self):
        """Create a TPV1Auth instance for testing."""
        # Use a valid 32-byte hex secret (64 hex chars)
        return TPV1Auth("test-api-key", "0" * 64)

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

    def test_init(self, mock_configuration, tpv1_auth):
        """Test that AuthenticatedRESTClient initializes correctly."""
        client = AuthenticatedRESTClient(mock_configuration, tpv1_auth)
        assert client._tpv1_auth is tpv1_auth

    def test_request_adds_authorization_header(self, mock_configuration, tpv1_auth):
        """Test that request method adds Authorization header."""
        client = AuthenticatedRESTClient(mock_configuration, tpv1_auth)

        # Mock the parent class request method
        with patch.object(
            AuthenticatedRESTClient.__bases__[0], "request", return_value=MagicMock()
        ) as mock_request:
            client.request(
                method="GET",
                url="https://api.example.com/v1/wallets",
                headers={"Content-Type": "application/json"},
            )

            # Verify request was called
            mock_request.assert_called_once()

            # Get the headers that were passed
            call_kwargs = mock_request.call_args
            headers = call_kwargs.kwargs.get("headers") or call_kwargs[1].get("headers")

            # Verify Authorization header was added
            assert "Authorization" in headers
            assert headers["Authorization"].startswith("TPV1-HMAC-SHA256 ApiKey=test-api-key")

    def test_request_with_json_body(self, mock_configuration, tpv1_auth):
        """Test that request properly handles JSON body for signing."""
        client = AuthenticatedRESTClient(mock_configuration, tpv1_auth)

        with patch.object(
            AuthenticatedRESTClient.__bases__[0], "request", return_value=MagicMock()
        ) as mock_request:
            body = {"name": "test-wallet", "currency": "BTC"}
            client.request(
                method="POST",
                url="https://api.example.com/v1/wallets",
                headers={"Content-Type": "application/json"},
                body=body,
            )

            mock_request.assert_called_once()

            # Verify the body was passed through unchanged
            call_kwargs = mock_request.call_args
            passed_body = call_kwargs.kwargs.get("body") or call_kwargs[1].get("body")
            assert passed_body == body

    def test_request_with_string_body(self, mock_configuration, tpv1_auth):
        """Test that request properly handles string body."""
        client = AuthenticatedRESTClient(mock_configuration, tpv1_auth)

        with patch.object(
            AuthenticatedRESTClient.__bases__[0], "request", return_value=MagicMock()
        ) as mock_request:
            body = '{"name": "test-wallet"}'
            client.request(
                method="POST",
                url="https://api.example.com/v1/wallets",
                headers={"Content-Type": "application/json"},
                body=body,
            )

            mock_request.assert_called_once()

    def test_request_with_no_headers(self, mock_configuration, tpv1_auth):
        """Test that request works when headers is None."""
        client = AuthenticatedRESTClient(mock_configuration, tpv1_auth)

        with patch.object(
            AuthenticatedRESTClient.__bases__[0], "request", return_value=MagicMock()
        ) as mock_request:
            client.request(
                method="GET",
                url="https://api.example.com/v1/wallets",
                headers=None,
            )

            mock_request.assert_called_once()

            # Verify headers dict was created and Authorization added
            call_kwargs = mock_request.call_args
            headers = call_kwargs.kwargs.get("headers") or call_kwargs[1].get("headers")
            assert headers is not None
            assert "Authorization" in headers

    def test_request_with_query_string(self, mock_configuration, tpv1_auth):
        """Test that request properly parses query string from URL."""
        client = AuthenticatedRESTClient(mock_configuration, tpv1_auth)

        with patch.object(
            AuthenticatedRESTClient.__bases__[0], "request", return_value=MagicMock()
        ) as mock_request:
            client.request(
                method="GET",
                url="https://api.example.com/v1/wallets?limit=10&offset=0",
                headers={},
            )

            mock_request.assert_called_once()

            # Verify the full URL was passed through
            call_kwargs = mock_request.call_args
            passed_url = call_kwargs.kwargs.get("url") or call_kwargs[1].get("url")
            assert passed_url == "https://api.example.com/v1/wallets?limit=10&offset=0"
