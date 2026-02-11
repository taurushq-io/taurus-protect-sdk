"""TPV1-HMAC-SHA256 authentication for Taurus-PROTECT API."""

from __future__ import annotations

import hashlib
import hmac
import time
import uuid
from typing import Optional
from urllib.parse import urlparse


def _wipe_bytes(data: bytearray) -> None:
    """
    Securely wipe bytes from memory by overwriting with zeros.

    SECURITY NOTE: Python's garbage collector may retain copies of sensitive data
    in memory due to string interning, object copying, and reference counting.
    This wiping is best-effort and provides defense-in-depth, but cannot guarantee
    complete removal of sensitive data from process memory. For maximum security
    in production environments, consider using hardware security modules (HSMs)
    or dedicated secrets management solutions.
    """
    if data:
        for i in range(len(data)):
            data[i] = 0


class TPV1Auth:
    """
    TPV1-HMAC-SHA256 authentication handler.

    Generates Authorization headers for Taurus-PROTECT API requests.

    Example:
        auth = TPV1Auth("api-key", "hex-encoded-secret")
        header = auth.sign_request("POST", "api.example.com", "/v1/wallets", ...)
        # header = "TPV1-HMAC-SHA256 ApiKey=... Nonce=... Timestamp=... Signature=..."
    """

    def __init__(self, api_key: str, api_secret_hex: str) -> None:
        """
        Initialize TPV1 authentication.

        Args:
            api_key: The API key.
            api_secret_hex: The API secret as a hex-encoded string.

        Raises:
            ValueError: If api_key or api_secret_hex is empty.
        """
        self._closed = False  # Initialize early for __del__ safety
        self._secret = bytearray()

        if not api_key:
            raise ValueError("api_key cannot be empty")
        if not api_secret_hex:
            raise ValueError("api_secret cannot be empty")

        self.api_key = api_key
        # Store as bytearray for secure wiping
        self._secret = bytearray(bytes.fromhex(api_secret_hex))

    def close(self) -> None:
        """
        Securely wipe the API secret from memory.

        Note: Due to Python's garbage collector, string interning, and memory
        management, this is best-effort wiping. See _wipe_bytes() for details.
        """
        if not self._closed:
            _wipe_bytes(self._secret)
            self._closed = True

    def __del__(self) -> None:
        """Ensure secret is wiped on garbage collection."""
        self.close()

    def sign_request(
        self,
        method: str,
        host: str,
        path: str,
        query: Optional[str] = None,
        content_type: Optional[str] = None,
        body: Optional[str] = None,
    ) -> str:
        """
        Generate TPV1 Authorization header for a request.

        Args:
            method: HTTP method (GET, POST, etc.).
            host: Request host (e.g., "api.protect.taurushq.com").
            path: Request path (e.g., "/v1/wallets").
            query: Query string without leading "?" (optional).
            content_type: Content-Type header value (optional).
            body: Request body as string (optional).

        Returns:
            The Authorization header value.

        Raises:
            RuntimeError: If auth has been closed.
        """
        if self._closed:
            raise RuntimeError("TPV1Auth has been closed")

        nonce = str(uuid.uuid4())
        timestamp = int(time.time() * 1000)

        # Build message from non-empty parts
        parts = [
            "TPV1",
            self.api_key,
            nonce,
            str(timestamp),
            method.upper(),
            host,
            path,
        ]

        if query:
            parts.append(query)
        if content_type:
            parts.append(content_type)
        if body:
            parts.append(body)

        message = " ".join(parts)

        # Calculate HMAC-SHA256
        signature = hmac.new(
            bytes(self._secret),
            message.encode("utf-8"),
            hashlib.sha256,
        ).digest()

        import base64

        sig_b64 = base64.b64encode(signature).decode("utf-8")

        return (
            f"TPV1-HMAC-SHA256 ApiKey={self.api_key} "
            f"Nonce={nonce} Timestamp={timestamp} Signature={sig_b64}"
        )

    @staticmethod
    def parse_url(url: str) -> tuple[str, str, Optional[str]]:
        """
        Parse URL into host, path, and query components.

        Args:
            url: Full URL to parse.

        Returns:
            Tuple of (host, path, query).
        """
        parsed = urlparse(url)
        host = parsed.netloc
        path = parsed.path or "/"
        query = parsed.query if parsed.query else None
        return host, path, query


def calculate_base64_hmac(secret: bytes, data: str) -> str:
    """
    Calculate HMAC-SHA256 and return base64-encoded result.

    Args:
        secret: The secret key as bytes.
        data: The data to sign.

    Returns:
        Base64-encoded HMAC-SHA256 signature.
    """
    import base64

    signature = hmac.new(secret, data.encode("utf-8"), hashlib.sha256).digest()
    return base64.b64encode(signature).decode("utf-8")


def verify_base64_hmac(secret: bytes, data: str, expected_hmac: str) -> bool:
    """
    Verify HMAC-SHA256 signature using constant-time comparison.

    Args:
        secret: The secret key as bytes.
        data: The original data.
        expected_hmac: The expected base64-encoded HMAC.

    Returns:
        True if the HMAC matches.
    """
    computed = calculate_base64_hmac(secret, data)
    return hmac.compare_digest(computed, expected_hmac)
