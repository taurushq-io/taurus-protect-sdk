"""Multi-currency end-to-end integration test for Taurus-PROTECT Python SDK.

Runs the full transfer lifecycle for multiple currencies in parallel,
mirroring the Java SDK's MultiCurrencyE2EIntegrationTest.

The test proceeds in two phases:
1. Discovery phase (sequential): For each currency, find two funded addresses.
2. Transfer phase (parallel): Run the transfer flow concurrently for all
   currencies that have two funded addresses.

Per-currency steps:
  Step 0: Find two funded addresses for the currency
  Step 1: Create an internal transfer request between them
  Step 2: Verify metadata matches the original transfer intent
  Step 3: Approve the request with ECDSA private key
  Step 4: Wait for terminal status (BROADCASTED/CONFIRMED)

The test passes if at least one currency completes successfully.
"""

from __future__ import annotations

import time
from concurrent.futures import ThreadPoolExecutor, as_completed
from dataclasses import dataclass, field
from enum import Enum
from typing import Dict, List, Optional, Tuple

import pytest

from taurus_protect.client import ProtectClient
from taurus_protect.crypto.keys import decode_private_key_pem
from taurus_protect.models.request import RequestStatus

from tests.testutil import get_private_key

# ---------------------------------------------------------------------------
# Constants
# ---------------------------------------------------------------------------

MAX_WAIT_SECONDS = 15 * 60  # 15 minutes per request
POLL_INTERVAL_SECONDS = 5
MAX_PARALLEL_TRANSFERS = 3

TERMINAL_STATUSES = {
    RequestStatus.BROADCASTED,
    RequestStatus.CONFIRMED,
    RequestStatus.REJECTED,
    RequestStatus.CANCELED,
    RequestStatus.PERMANENT_FAILURE,
    RequestStatus.EXPIRED,
    RequestStatus.MINED,
    # Note: RequestStatus(str, Enum) so enum members match string comparisons
}

SUCCESS_STATUSES = {RequestStatus.BROADCASTED, RequestStatus.CONFIRMED}


# ---------------------------------------------------------------------------
# Data classes
# ---------------------------------------------------------------------------


@dataclass(frozen=True)
class CurrencyConfig:
    """Configuration for a currency under test."""

    symbol: str
    blockchain: str
    network: str
    transfer_amount: str
    min_balance: str
    is_token: bool


class ResultStatus(Enum):
    PASSED = "PASSED"
    SKIPPED = "SKIPPED"
    FAILED = "FAILED"


@dataclass
class CurrencyResult:
    symbol: str
    status: ResultStatus
    tx_hash: Optional[str] = None
    error_message: Optional[str] = None


@dataclass
class AddressPair:
    """A source/destination address pair for a currency."""

    config: CurrencyConfig
    source_id: str
    source_address: str
    source_balance: str
    dest_id: str
    dest_address: str


# ---------------------------------------------------------------------------
# Currency configurations (matching Java SDK)
# ---------------------------------------------------------------------------

CURRENCY_CONFIGS = [
    CurrencyConfig("SOL", "SOL", "mainnet", "1000000", "2000000", False),
    # CurrencyConfig("ETH", "ETH", "mainnet", "10000000000000", "100000000000000", False),
    # CurrencyConfig("XRP", "XRP", "mainnet", "1212", "12000000", False),
    CurrencyConfig("XLM", "XLM", "mainnet", "1313", "35000000", False),
    CurrencyConfig("ALGO", "ALGO", "mainnet", "1414", "100000", False),
    # CurrencyConfig("USDC", "ETH", "mainnet", "11100", "2000000", True),
]


# ---------------------------------------------------------------------------
# Test class
# ---------------------------------------------------------------------------


@pytest.mark.e2e
class TestMultiCurrencyE2E:
    """Multi-currency end-to-end transfer test."""

    def test_multi_currency_transfer_e2e(self, client: ProtectClient) -> None:
        """Run E2E transfer flows for multiple currencies in parallel."""
        pem_bytes = get_private_key(1)
        if pem_bytes is None:
            pytest.skip("No private key configured for identity 1")
        private_key = decode_private_key_pem(pem_bytes.decode("utf-8"))

        # Phase 1: Sequential address discovery
        print("\n=== Phase 1: Address Discovery (sequential) ===")
        address_pairs: List[AddressPair] = []
        results: List[CurrencyResult] = []

        for config in CURRENCY_CONFIGS:
            tag = f"[{config.symbol}] "
            print(f"{tag}Step 0: Searching for two funded addresses...")
            try:
                pair = _find_two_funded_addresses(client, config)
                if pair is None:
                    print(f"{tag}Step 0: Could not find two funded addresses -- skipping")
                    results.append(CurrencyResult(config.symbol, ResultStatus.SKIPPED))
                else:
                    print(
                        f"{tag}Step 0: Source: ID={pair.source_id}"
                        f" ({pair.source_address}), Balance={pair.source_balance}"
                    )
                    print(f"{tag}Step 0: Dest:   ID={pair.dest_id} ({pair.dest_address})")
                    address_pairs.append(pair)
            except Exception as exc:
                print(f"{tag}Step 0: Discovery failed -- {exc}")
                results.append(CurrencyResult(config.symbol, ResultStatus.SKIPPED))

        if not address_pairs:
            print("No currency has two funded addresses. Cannot run E2E test.")
            pytest.fail("No currency has two funded addresses available for testing")

        # Phase 2: Parallel transfer flows
        num_workers = min(len(address_pairs), MAX_PARALLEL_TRANSFERS)
        print(f"\n=== Phase 2: Transfer Flows (parallel, max {num_workers} threads) ===")

        with ThreadPoolExecutor(max_workers=num_workers) as executor:
            future_to_pair = {
                executor.submit(_run_transfer_flow, client, pair, private_key): pair
                for pair in address_pairs
            }

            for future in as_completed(future_to_pair):
                pair = future_to_pair[future]
                try:
                    result = future.result()
                    results.append(result)
                    extra = ""
                    if result.tx_hash:
                        extra = f" -- tx: {result.tx_hash}"
                    if result.error_message:
                        extra += f" -- {result.error_message}"
                    print(f">>> {result.symbol}: {result.status.value}{extra}")
                except Exception as exc:
                    result = CurrencyResult("UNKNOWN", ResultStatus.FAILED, error_message=str(exc))
                    results.append(result)
                    print(f">>> UNKNOWN: FAILED -- {exc}")

        # Print summary
        _print_summary(results)

        # Assertions
        passed = sum(1 for r in results if r.status == ResultStatus.PASSED)
        failed = sum(1 for r in results if r.status == ResultStatus.FAILED)
        skipped = len(results) - passed - failed

        print(f"Summary: {passed} passed, {failed} failed, {skipped} skipped")

        for r in results:
            if r.status == ResultStatus.FAILED:
                print(f"WARNING: {r.symbol} failed: {r.error_message}")

        # At least one currency must pass. Individual currencies may fail due to
        # environment-specific issues (business rule limits, insufficient balance, etc.)
        assert passed > 0, (
            f"At least one currency should complete the E2E flow successfully. "
            f"Passed: {passed}, Failed: {failed}, Skipped: {skipped}"
        )


# ---------------------------------------------------------------------------
# Transfer flow (runs per-currency in a thread)
# ---------------------------------------------------------------------------


def _run_transfer_flow(
    client: ProtectClient,
    pair: AddressPair,
    private_key: object,
) -> CurrencyResult:
    """Run the full transfer flow for a single currency."""
    config = pair.config
    tag = f"[{config.symbol}] "

    try:
        # Step 1: Create internal transfer request
        print(
            f"{tag}Step 1: Creating transfer of {config.transfer_amount}"
            f" from address {pair.source_id} to address {pair.dest_id}..."
        )
        transfer_request = client.requests.create_internal_transfer(
            from_address_id=int(pair.source_id),
            to_address_id=int(pair.dest_id),
            amount=config.transfer_amount,
        )
        assert transfer_request is not None
        assert int(transfer_request.id) > 0, "Request ID should be positive"
        print(
            f"{tag}Step 1: Created request: ID={transfer_request.id}"
            f", Status={transfer_request.status.value}"
        )

        # Step 2: Verify metadata matches original intent
        print(f"{tag}Step 2: Verifying metadata matches transfer intent...")
        request_to_approve = client.requests.get(int(transfer_request.id))
        metadata = request_to_approve.metadata
        assert metadata is not None, "Request metadata should be available"
        assert metadata.hash, "Request hash should be available"

        # Verify source address from verified payload
        metadata_source = metadata.get_source_address()
        assert metadata_source == pair.source_address, (
            f"{config.symbol} metadata source address should match the funded address"
        )
        print(f"{tag}Step 2: Source verified: {metadata_source}")

        # Verify destination address
        # For token transfers, the metadata destination is the token contract address
        metadata_destination = metadata.get_destination_address()
        if not config.is_token:
            assert metadata_destination == pair.dest_address, (
                f"{config.symbol} metadata destination address should match the target address"
            )
        token_note = " (token contract)" if config.is_token else " (verified)"
        print(f"{tag}Step 2: Destination: {metadata_destination}{token_note}")

        # Verify amount (for non-token transfers)
        metadata_amount = metadata.get_amount()
        metadata_amount_value = metadata_amount.value_from if metadata_amount else None
        if not config.is_token and metadata_amount_value is not None:
            assert int(metadata_amount_value) == int(config.transfer_amount), (
                f"{config.symbol} metadata amount should match the transfer amount"
            )
        amount_note = " (token transfer -- native value is 0)" if config.is_token else " (verified)"
        print(f"{tag}Step 2: Amount: {metadata_amount_value}{amount_note}")

        # Step 3: Approve request
        print(f"{tag}Step 3: Approving request {transfer_request.id}...")
        signed_count = client.requests.approve_request(request_to_approve, private_key)
        assert signed_count > 0, "At least one request should have been signed"
        print(f"{tag}Step 3: Approved: signedCount={signed_count}")

        # Step 4: Wait for terminal status
        print(f"{tag}Step 4: Waiting for terminal status...")
        confirmed_request = _wait_for_terminal_status(client, int(transfer_request.id), tag)
        print(f"{tag}Step 4: Final status: {confirmed_request.status.value}")

        # Dump diagnostics if not a success status
        if confirmed_request.status not in SUCCESS_STATUSES:
            print(f"{tag}=== DIAGNOSTIC DUMP ===")
            if confirmed_request.signed_requests:
                for sr in confirmed_request.signed_requests:
                    print(
                        f"{tag}  SignedRequest ID={sr.id}"
                        f", Status={sr.status}"
                        f", Hash={sr.hash}"
                        f", Details={sr.details}"
                    )
            if confirmed_request.trails:
                print(f"{tag}  Trails:")
                for trail in confirmed_request.trails:
                    print(
                        f"{tag}    action={trail.action}"
                        f", comment={trail.comment}"
                        f", timestamp={trail.timestamp}"
                    )
            print(f"{tag}=== END DIAGNOSTIC ===")
            raise RuntimeError(
                f"{config.symbol} request ended with {confirmed_request.status.value}"
                f" instead of BROADCASTED/CONFIRMED"
            )

        # Step 5: Verify transaction by hash
        tx_hash = None
        if confirmed_request.signed_requests and confirmed_request.signed_requests[0].hash:
            tx_hash = confirmed_request.signed_requests[0].hash
            print(f"{tag}Step 5: Transaction hash: {tx_hash}")

            transaction = None
            tx_deadline = time.time() + 90
            while transaction is None and time.time() < tx_deadline:
                time.sleep(POLL_INTERVAL_SECONDS)
                try:
                    transaction = client.transactions.get_by_hash(tx_hash)
                except Exception:
                    print(f"{tag}Step 5: Transaction not indexed yet, retrying...")

            assert transaction is not None, f"{tag} transaction should be found by hash"
            assert transaction.tx_hash == tx_hash, f"{tag} transaction hash should match"
            assert transaction.direction == "outgoing", f"{tag} transaction direction should be outgoing"
            assert transaction.request_id == str(confirmed_request.id), f"{tag} transaction request_id should match"
            print(f"{tag}Step 5: Transaction verified -- ID={transaction.id}, Block={transaction.block_height}")
        else:
            print(f"{tag}Step 5: No signed requests or transaction hash available, skipping")

        print(f"{tag}E2E PASSED")
        return CurrencyResult(config.symbol, ResultStatus.PASSED, tx_hash=tx_hash)

    except Exception as exc:
        print(f"{tag}E2E FAILED: {exc}")
        return CurrencyResult(config.symbol, ResultStatus.FAILED, error_message=str(exc))


# ---------------------------------------------------------------------------
# Address discovery
# ---------------------------------------------------------------------------


def _find_two_funded_addresses(
    client: ProtectClient,
    config: CurrencyConfig,
) -> Optional[AddressPair]:
    """Find two funded addresses for the given currency."""
    tag = f"[{config.symbol}] "

    if config.is_token:
        return _find_token_addresses(client, config, tag)
    else:
        return _find_native_addresses(client, config, tag)


def _find_native_addresses(
    client: ProtectClient,
    config: CurrencyConfig,
    tag: str,
) -> Optional[AddressPair]:
    """Find two funded native-currency addresses by scanning wallets."""
    wallets_dtos, _ = client.assets.get_wallets(config.symbol)
    if not wallets_dtos:
        print(f"{tag}  No wallets found via get_wallets")
        return None
    print(f"{tag}  Found {len(wallets_dtos)} wallet(s) via get_wallets")

    funded: List[Tuple[str, str, str]] = []  # (id, address, balance)
    total_scanned = 0

    for wallet_dto in wallets_dtos:
        wallet_id = str(getattr(wallet_dto, "wallet_id", None) or getattr(wallet_dto, "id", ""))
        wallet_name = getattr(wallet_dto, "name", "")
        if not wallet_id:
            continue
        print(f"{tag}  Scanning wallet ID={wallet_id}, Name={wallet_name}")

        limit = 50
        offset = 0
        while True:
            addresses, _ = client.addresses.list(int(wallet_id), limit=limit, offset=offset)
            if not addresses:
                break
            for addr in addresses:
                total_scanned += 1
                if addr.disabled:
                    print(
                        f"{tag}    DISABLED Address ID={addr.id}"
                        f", Addr={addr.address} -- skipping"
                    )
                    continue
                bal_str = "0"
                if addr.balance is not None:
                    bal_str = addr.balance.available_confirmed or "0"
                bal_int = int(bal_str)
                if bal_int > 0:
                    print(
                        f"{tag}    FUNDED Address ID={addr.id}"
                        f", Addr={addr.address}, Balance={bal_str}"
                    )
                    funded.append((addr.id, addr.address, bal_str))
            if len(funded) >= 2:
                print(
                    f"{tag}  Found {len(funded)} funded addresses"
                    f" after scanning {total_scanned} total"
                )
                break
            offset += len(addresses)

        if len(funded) >= 2:
            break

    print(
        f"{tag}  Found {len(funded)} funded addresses"
        f" after scanning {total_scanned} total across {len(wallets_dtos)} wallets"
    )

    if len(funded) < 2:
        return None

    # Source: needs sufficient balance
    source = None
    for addr_id, addr_str, bal_str in funded:
        if int(bal_str) >= int(config.min_balance):
            source = (addr_id, addr_str, bal_str)
            break

    if source is None:
        print(f"{tag}Step 0: No address has sufficient balance (min={config.min_balance})")
        return None

    # Destination: any other address
    dest = None
    for addr_id, addr_str, bal_str in funded:
        if addr_id != source[0]:
            dest = (addr_id, addr_str, bal_str)
            break

    if dest is None:
        return None

    return AddressPair(
        config=config,
        source_id=source[0],
        source_address=source[1],
        source_balance=source[2],
        dest_id=dest[0],
        dest_address=dest[1],
    )


def _find_token_addresses(
    client: ProtectClient,
    config: CurrencyConfig,
    tag: str,
) -> Optional[AddressPair]:
    """Find two funded token addresses, checking gas balance on the native chain."""
    print(f"{tag}Step 0: Searching via AssetService (token)...")
    token_dtos, _ = client.assets.get_addresses(config.symbol)
    print(
        f"{tag}Step 0: AssetService returned"
        f" {len(token_dtos) if token_dtos else 0} token addresses"
    )

    if not token_dtos or len(token_dtos) < 2:
        print(
            f"{tag}Step 0: Not enough token addresses"
            f" ({len(token_dtos) if token_dtos else 0} found, need 2)"
        )
        return None

    # Fetch native currency addresses to check gas balance
    print(
        f"{tag}Step 0: Fetching {config.blockchain}"
        f" addresses for gas balance check..."
    )
    native_dtos, _ = client.assets.get_addresses(config.blockchain)
    native_balance_by_addr: Dict[str, int] = {}
    if native_dtos:
        for n_dto in native_dtos:
            n_addr = str(getattr(n_dto, "address", "") or "").lower()
            n_bal_obj = getattr(n_dto, "balance", None)
            n_bal = 0
            if n_bal_obj is not None:
                n_bal = int(getattr(n_bal_obj, "available_confirmed", "0") or "0")
            if n_addr:
                native_balance_by_addr[n_addr] = n_bal
        print(f"{tag}Step 0: Found {len(native_dtos)} {config.blockchain} addresses")

    # Build candidate list with diagnostic output
    candidates: List[Tuple[str, str, int, int]] = []  # (id, address, token_bal, gas_bal)
    for dto in token_dtos:
        dto_id = str(getattr(dto, "id", ""))
        dto_addr = str(getattr(dto, "address", "") or "").lower()

        # Filter out disabled addresses
        if getattr(dto, "disabled", False):
            print(f"{tag}  DISABLED Address ID={dto_id}, Addr={dto_addr} -- skipping")
            continue

        bal_obj = getattr(dto, "balance", None)
        token_bal = 0
        if bal_obj is not None:
            token_bal = int(getattr(bal_obj, "available_confirmed", "0") or "0")
        gas_bal = native_balance_by_addr.get(dto_addr, 0)

        token_ok = token_bal >= int(config.min_balance)
        gas_ok = gas_bal > 0
        print(
            f"{tag}  Candidate ID={dto_id}, Addr={dto_addr}"
            f", TokenBalance={token_bal} {'[OK]' if token_ok else '[LOW]'}"
            f", GasBalance={gas_bal} {'[OK]' if gas_ok else '[EMPTY]'}"
        )
        candidates.append((dto_id, dto_addr, token_bal, gas_bal))

    if len(candidates) < 2:
        print(f"{tag}Step 0: Not enough candidates ({len(candidates)} found, need 2)")
        return None

    # Source: needs sufficient token balance AND gas
    source = None
    for cid, caddr, tbal, gbal in candidates:
        if tbal >= int(config.min_balance) and gbal > 0:
            print(f"{tag}Step 0: Selected source ID={cid} -- token={tbal}, gas={gbal}")
            source = (cid, caddr, str(tbal))
            break

    if source is None:
        print(f"{tag}Step 0: No address has sufficient balance with gas (min={config.min_balance})")
        return None

    # Destination: any other address
    dest = None
    for cid, caddr, tbal, gbal in candidates:
        if cid != source[0]:
            dest = (cid, caddr, str(tbal))
            break

    if dest is None:
        return None

    return AddressPair(
        config=config,
        source_id=source[0],
        source_address=source[1],
        source_balance=source[2],
        dest_id=dest[0],
        dest_address=dest[1],
    )


# ---------------------------------------------------------------------------
# Status polling
# ---------------------------------------------------------------------------


def _wait_for_terminal_status(
    client: ProtectClient,
    request_id: int,
    tag: str,
) -> object:
    """Poll until the request reaches a terminal status or timeout."""
    start = time.monotonic()

    while True:
        request = client.requests.get(request_id)
        status = request.status
        elapsed = int(time.monotonic() - start)
        print(f"{tag}  [{elapsed}s] Request {request_id} -- Status: {status.value}")

        if status in TERMINAL_STATUSES:
            return request

        if time.monotonic() - start >= MAX_WAIT_SECONDS:
            raise RuntimeError(
                f"Request {request_id} did not reach terminal status within"
                f" {MAX_WAIT_SECONDS} seconds. Last status: {status.value}"
            )

        time.sleep(POLL_INTERVAL_SECONDS)


# ---------------------------------------------------------------------------
# Summary output
# ---------------------------------------------------------------------------


def _print_summary(results: List[CurrencyResult]) -> None:
    """Print a formatted summary table."""
    print()
    print("=" * 67)
    print("            Multi-Currency E2E Results")
    print("=" * 67)

    for result in results:
        if result.status == ResultStatus.PASSED:
            print(f"  {result.symbol:<5}: PASSED")
            if result.tx_hash:
                print(f"    tx : {result.tx_hash}")
        elif result.status == ResultStatus.SKIPPED:
            print(f"  {result.symbol:<5}: SKIPPED (fewer than 2 addresses found)")
        elif result.status == ResultStatus.FAILED:
            msg = result.error_message or ""
            if len(msg) > 60:
                msg = msg[:60] + "..."
            print(f"  {result.symbol:<5}: FAILED -- {msg}")

    print("=" * 67)
    print()
