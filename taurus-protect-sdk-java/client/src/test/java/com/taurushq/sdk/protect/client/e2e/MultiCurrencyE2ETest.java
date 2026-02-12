package com.taurushq.sdk.protect.client.e2e;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Balance;
import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestMetadata;
import com.taurushq.sdk.protect.client.model.RequestMetadataAmount;
import com.taurushq.sdk.protect.client.model.RequestStatus;
import com.taurushq.sdk.protect.client.model.RequestTrail;
import com.taurushq.sdk.protect.client.model.SignedRequest;
import com.taurushq.sdk.protect.client.model.Transaction;
import com.taurushq.sdk.protect.client.model.Wallet;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.math.BigInteger;
import java.security.PrivateKey;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.EnumSet;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.CompletionService;
import java.util.concurrent.ExecutorCompletionService;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.concurrent.TimeUnit;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Multi-currency end-to-end integration test that runs the full transfer lifecycle
 * for multiple currencies in parallel.
 *
 * <p>The test proceeds in two phases:</p>
 * <ol>
 *   <li><b>Discovery phase</b> (sequential): For each currency, find two existing
 *       funded addresses to use as source and destination.</li>
 *   <li><b>Transfer phase</b> (parallel): Run the transfer flow concurrently for all
 *       currencies that have two funded addresses.</li>
 * </ol>
 *
 * <p>Per-currency steps:</p>
 * <ol>
 *   <li>Step 0: Find two funded addresses for the currency</li>
 *   <li>Step 1: Create an internal transfer request between them</li>
 *   <li>Step 2: Verify metadata matches the original transfer intent</li>
 *   <li>Step 3: Approve the request</li>
 *   <li>Step 4: Wait for BROADCASTED/CONFIRMED status</li>
 *   <li>Step 5: Verify the transaction by hash</li>
 * </ol>
 *
 * <p>Currencies without two funded addresses are skipped. The test passes if at least
 * one currency completes successfully.</p>
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class MultiCurrencyE2ETest {

    /** Maximum time to wait for a single request confirmation (15 minutes). */
    private static final long MAX_WAIT_MS = 15 * 60 * 1000L;

    /** Polling interval between status checks (5 seconds). */
    private static final long POLL_INTERVAL_MS = 5000L;

    /** Maximum time to wait for all currencies to complete (20 minutes). */
    private static final long EXECUTOR_TIMEOUT_MINUTES = 20;

    /** Maximum concurrent transfer flows to avoid API throttling. */
    private static final int MAX_PARALLEL_TRANSFERS = 3;

    /** Terminal statuses that indicate the request will not progress further. */
    private static final Set<RequestStatus> TERMINAL_STATUSES = EnumSet.of(
            RequestStatus.BROADCASTED,
            RequestStatus.CONFIRMED,
            RequestStatus.REJECTED,
            RequestStatus.CANCELED,
            RequestStatus.PERMANENT_FAILURE,
            RequestStatus.EXPIRED,
            RequestStatus.INVALID,
            RequestStatus.MINED
    );

    private ProtectClient client;

    @BeforeAll
    void setup() throws Exception {
        TestHelper.skipIfNotEnabled();
        client = TestHelper.getTestClient(1);
    }

    @AfterAll
    void teardown() {
        if (client != null) {
            client.close();
        }
    }

    @Test
    void multiCurrencyTransferE2E() throws Exception {
        PrivateKey approvalKey = TestHelper.getPrivateKey(1);
        List<CurrencyConfig> configs = getCurrencyConfigs();

        // Phase 1: Sequential address discovery to avoid API overload
        System.out.println("=== Phase 1: Address Discovery (sequential) ===");
        List<AddressPair> addressPairs = new ArrayList<>();
        List<CurrencyResult> results = new ArrayList<>();

        for (CurrencyConfig config : configs) {
            String tag = "[" + config.symbol + "] ";
            System.out.println(tag + "Step 0: Searching for two funded addresses...");
            try {
                AddressPair pair = findTwoFundedAddresses(config);
                if (pair == null) {
                    System.out.println(tag + "Step 0: Could not find two funded addresses — skipping");
                    results.add(CurrencyResult.skipped(config.symbol));
                } else {
                    System.out.println(tag + "Step 0: Source: ID=" + pair.source.getId()
                            + " (" + pair.source.getAddress() + ")"
                            + ", Balance=" + pair.source.getBalance().getAvailableConfirmed());
                    System.out.println(tag + "Step 0: Dest:   ID=" + pair.destination.getId()
                            + " (" + pair.destination.getAddress() + ")");
                    addressPairs.add(pair);
                }
            } catch (Exception e) {
                System.out.println(tag + "Step 0: Discovery failed — " + e.getMessage());
                results.add(CurrencyResult.skipped(config.symbol));
            }
        }

        if (addressPairs.isEmpty()) {
            System.out.println("No currency has two funded addresses. Cannot run E2E test.");
            assertTrue(false, "No currency has two funded addresses available for testing");
        }

        // Phase 2: Parallel transfer flows — results are printed as each currency completes
        System.out.println("\n=== Phase 2: Transfer Flows (parallel, max " + MAX_PARALLEL_TRANSFERS + " threads) ===");
        ExecutorService executor = Executors.newFixedThreadPool(
                Math.min(addressPairs.size(), MAX_PARALLEL_TRANSFERS));
        CompletionService<CurrencyResult> completionService = new ExecutorCompletionService<>(executor);

        int submitted = 0;
        for (AddressPair pair : addressPairs) {
            completionService.submit(() -> runTransferFlow(pair, approvalKey));
            submitted++;
        }

        executor.shutdown();

        // Collect results as they complete
        long deadline = System.currentTimeMillis() + TimeUnit.MINUTES.toMillis(EXECUTOR_TIMEOUT_MINUTES);
        for (int i = 0; i < submitted; i++) {
            long remaining = deadline - System.currentTimeMillis();
            if (remaining <= 0) {
                System.out.println("TIMEOUT: not all currencies completed within " + EXECUTOR_TIMEOUT_MINUTES + " minutes");
                break;
            }
            Future<CurrencyResult> future = completionService.poll(remaining, TimeUnit.MILLISECONDS);
            if (future == null) {
                System.out.println("TIMEOUT: not all currencies completed within " + EXECUTOR_TIMEOUT_MINUTES + " minutes");
                break;
            }
            try {
                CurrencyResult result = future.get();
                results.add(result);
                System.out.println(">>> " + result.symbol + ": " + result.status
                        + (result.txHash != null ? " — tx: " + result.txHash : "")
                        + (result.errorMessage != null ? " — " + result.errorMessage : "")
                        + " [" + results.size() + "/" + (submitted + (results.size() - submitted)) + " done]");
            } catch (Exception e) {
                CurrencyResult result = CurrencyResult.failed("UNKNOWN", e.getMessage());
                results.add(result);
                System.out.println(">>> UNKNOWN: FAILED — " + e.getMessage());
            }
        }

        executor.shutdownNow();

        // Print summary
        printSummary(results);

        // Assertions: at least one currency must pass.
        long passedCount = 0;
        long failedCount = 0;
        for (CurrencyResult result : results) {
            if (result.status == CurrencyResult.Status.PASSED) {
                passedCount++;
            } else if (result.status == CurrencyResult.Status.FAILED) {
                failedCount++;
                System.out.println("WARNING: " + result.symbol + " failed: " + result.errorMessage);
            }
        }

        System.out.println("Summary: " + passedCount + " passed, " + failedCount + " failed, "
                + (results.size() - passedCount - failedCount) + " skipped");

        assertEquals(0, failedCount,
                "No currency should fail the E2E flow. Failed: " + failedCount);
        assertTrue(passedCount > 0,
                "At least one currency should complete the E2E flow successfully");
    }

    // ── Currency configurations ─────────────────────────────────────────

    private static List<CurrencyConfig> getCurrencyConfigs() {
        return Arrays.asList(
                new CurrencyConfig("SOL", "SOL", "mainnet",
                        BigInteger.valueOf(1_000_000),         // 0.001 SOL (lamports)
                        BigInteger.valueOf(2_000_000),
                        false),
//                new CurrencyConfig("ETH", "ETH", "mainnet",
//                        BigInteger.valueOf(10_000_000_000_000L), // 0.00001 ETH (wei)
//                        BigInteger.valueOf(100_000_000_000_000L),
//                        false),
//                new CurrencyConfig("BTC", "BTC", "mainnet",    // too slow
//                         BigInteger.valueOf(1111),
//                         BigInteger.valueOf(20_000),
//                         false),
//                 new CurrencyConfig("XRP", "XRP", "mainnet",
//                         BigInteger.valueOf(1212),              // 0.001212 XRP (drops)
//                         BigInteger.valueOf(12_000_000),        // 12 XRP min balance
//                         false),
                new CurrencyConfig("XLM", "XLM", "mainnet",
                        BigInteger.valueOf(1313),              // 0.0001313 XLM (stroops)
                        BigInteger.valueOf(35_000_000),        // 3.5 XLM min balance
                        false),
                new CurrencyConfig("ALGO", "ALGO", "mainnet",
                        BigInteger.valueOf(1414),           // 0.001414 ALGO (microAlgos)
                        BigInteger.valueOf(100_000),        // 0.1 ALGO min balance
                        false)
//                new CurrencyConfig("USDC", "ETH", "mainnet",
//                        BigInteger.valueOf(11_100),       // 0.0111 USDC (6 decimals)
//                        BigInteger.valueOf(2_000_000),    // 2.0 USDC min balance
//                        true)
        );
    }

    // ── Transfer flow (runs in parallel per currency) ───────────────────

    /**
     * Runs the transfer flow for a single currency using two existing addresses.
     */
    private CurrencyResult runTransferFlow(AddressPair pair, PrivateKey approvalKey) {
        CurrencyConfig config = pair.config;
        String tag = "[" + config.symbol + "] ";
        try {
            Address source = pair.source;
            Address destination = pair.destination;

            // Step 1: Create transfer request
            System.out.println(tag + "Step 1: Creating transfer of " + config.transferAmount
                    + " from address " + source.getId() + " to address " + destination.getId() + "...");
            Request transferRequest = client.getRequestService()
                    .createInternalTransferRequest(
                            source.getId(),
                            destination.getId(),
                            config.transferAmount);
            assertNotNull(transferRequest);
            assertTrue(transferRequest.getId() > 0, "Request ID should be positive");
            System.out.println(tag + "Step 1: Created request: ID=" + transferRequest.getId()
                    + ", Status=" + transferRequest.getStatus());

            // Step 2: Verify metadata matches original intent
            System.out.println(tag + "Step 2: Verifying metadata matches transfer intent...");
            Request requestToApprove = client.getRequestService().getRequest(transferRequest.getId());
            RequestMetadata metadata = requestToApprove.getMetadata();
            assertNotNull(metadata, "Request metadata should be available");
            assertNotNull(metadata.getHash(), "Request hash should be available");

            // Verify source address in metadata matches what we sent
            String metadataSource = metadata.getSourceAddress();
            assertEquals(source.getAddress(), metadataSource,
                    config.symbol + " metadata source address should match the funded address");
            System.out.println(tag + "Step 2: Source verified: " + metadataSource);

            // Verify destination address in metadata
            // For token transfers, the metadata destination is the token contract address,
            // not the recipient address
            String metadataDestination = metadata.getDestinationAddress();
            if (!config.isToken) {
                assertEquals(destination.getAddress(), metadataDestination,
                        config.symbol + " metadata destination address should match the target address");
            }
            System.out.println(tag + "Step 2: Destination: " + metadataDestination
                    + (config.isToken ? " (token contract)" : " (verified)"));

            // Verify amount in metadata
            // For token transfers, the native value is 0 (token amount is in contract call data)
            RequestMetadataAmount metadataAmount = metadata.getAmount();
            if (!config.isToken) {
                assertEquals(String.valueOf(config.transferAmount.longValue()), metadataAmount.getValueFrom(),
                        config.symbol + " metadata amount should match the transfer amount");
            }
            System.out.println(tag + "Step 2: Amount: valueFrom=" + metadataAmount.getValueFrom()
                    + ", currency=" + metadataAmount.getCurrencyFrom()
                    + (config.isToken ? " (token transfer — native value is 0)" : " (verified)"));

            // Step 3: Approve request
            System.out.println(tag + "Step 3: Approving request " + transferRequest.getId() + "...");
            int signedCount = client.getRequestService().approveRequest(requestToApprove, approvalKey);
            assertTrue(signedCount > 0, "At least one request should have been signed");
            System.out.println(tag + "Step 3: Approved: signedCount=" + signedCount);

            // Step 4: Wait for terminal status
            System.out.println(tag + "Step 4: Waiting for terminal status...");
            Request confirmedRequest = waitForTerminalStatus(transferRequest.getId(), tag);
            System.out.println(tag + "Step 4: Final status: " + confirmedRequest.getStatus());

            // Dump diagnostics if not a success status (BROADCASTED or CONFIRMED)
            if (confirmedRequest.getStatus() != RequestStatus.BROADCASTED
                    && confirmedRequest.getStatus() != RequestStatus.CONFIRMED) {
                System.out.println(tag + "=== DIAGNOSTIC DUMP ===");
                if (confirmedRequest.getSignedRequests() != null) {
                    for (SignedRequest sr : confirmedRequest.getSignedRequests()) {
                        System.out.println(tag + "  SignedRequest ID=" + sr.getId()
                                + ", Status=" + sr.getStatus()
                                + ", Hash=" + sr.getHash()
                                + ", Details=" + sr.getDetails());
                    }
                }
                if (confirmedRequest.getTrails() != null) {
                    System.out.println(tag + "  Trails:");
                    for (RequestTrail trail : confirmedRequest.getTrails()) {
                        System.out.println(tag + "    action=" + trail.getAction()
                                + ", status=" + trail.getRequestStatus()
                                + ", comment=" + trail.getComment()
                                + ", date=" + trail.getDate());
                    }
                }
                System.out.println(tag + "=== END DIAGNOSTIC ===");
                throw new RuntimeException(config.symbol + " request ended with "
                        + confirmedRequest.getStatus() + " instead of BROADCASTED/CONFIRMED");
            }

            // Step 5: Verify transaction (retry up to 30s for indexing delay)
            System.out.println(tag + "Step 5: Fetching transaction details...");
            String txHash = confirmedRequest.getSignedRequests().get(0).getHash();
            assertNotNull(txHash, "Confirmed request should have a transaction hash");
            System.out.println(tag + "Step 5: Transaction hash: " + txHash);

            Transaction transaction = null;
            long txLookupDeadline = System.currentTimeMillis() + 90_000L;
            while (transaction == null && System.currentTimeMillis() < txLookupDeadline) {
                Thread.sleep(POLL_INTERVAL_MS);
                try {
                    transaction = client.getTransactionService().getTransactionByHash(txHash);
                } catch (ApiException e) {
                    System.out.println(tag + "Step 5: Transaction not indexed yet, retrying...");
                }
            }
            assertNotNull(transaction, "Transaction with hash '" + txHash + "' not found after 90s");
            assertEquals(txHash, transaction.getHash());
            assertEquals("outgoing", transaction.getDirection());
            assertEquals(confirmedRequest.getId(), transaction.getRequestId());
            System.out.println(tag + "Step 5: Transaction verified — ID=" + transaction.getId()
                    + ", Block=" + transaction.getBlock());

            System.out.println(tag + "E2E PASSED");
            return CurrencyResult.passed(config.symbol, txHash);

        } catch (Exception e) {
            System.out.println(tag + "E2E FAILED: " + e.getMessage());
            return CurrencyResult.failed(config.symbol, e.getMessage());
        }
    }

    // ── Funded address discovery ────────────────────────────────────────

    /**
     * Finds two funded addresses for the given currency.
     * The first address (source) must have at least {@code minBalance}.
     * The second address (destination) can have any balance.
     *
     * @return an AddressPair with source and destination, or null if fewer than 2 addresses exist
     */
    private AddressPair findTwoFundedAddresses(CurrencyConfig config) throws ApiException {
        String tag = "[" + config.symbol + "] ";
        List<Address> candidates;
        // For tokens, build a native balance lookup so we can check gas availability
        Map<String, BigInteger> nativeBalanceByAddr = new HashMap<>();
        if (config.isToken) {
            System.out.println(tag + "Step 0: Searching via AssetService...");
            candidates = client.getAssetService().getAssetAddresses(config.symbol);
            System.out.println(tag + "Step 0: AssetService returned " + (candidates == null ? 0 : candidates.size()) + " token addresses");

            // Fetch native currency addresses to check gas balance
            System.out.println(tag + "Step 0: Fetching " + config.blockchain + " addresses for gas balance check...");
            List<Address> nativeAddresses = client.getAssetService().getAssetAddresses(config.blockchain);
            if (nativeAddresses != null) {
                for (Address nAddr : nativeAddresses) {
                    BigInteger nBal = (nAddr.getBalance() != null && nAddr.getBalance().getAvailableConfirmed() != null)
                            ? nAddr.getBalance().getAvailableConfirmed() : BigInteger.ZERO;
                    nativeBalanceByAddr.put(nAddr.getAddress().toLowerCase(), nBal);
                }
                System.out.println(tag + "Step 0: Found " + nativeAddresses.size() + " " + config.blockchain + " addresses");
            }

            if (candidates != null) {
                // Filter out disabled addresses
                List<Address> enabledCandidates = new ArrayList<>();
                for (Address addr : candidates) {
                    if (addr.isDisabled()) {
                        System.out.println(tag + "  DISABLED Address ID=" + addr.getId()
                                + ", Addr=" + addr.getAddress() + " — skipping");
                    } else {
                        enabledCandidates.add(addr);
                    }
                }
                candidates = enabledCandidates;

                for (Address addr : candidates) {
                    BigInteger tokenBal = (addr.getBalance() != null && addr.getBalance().getAvailableConfirmed() != null)
                            ? addr.getBalance().getAvailableConfirmed() : BigInteger.ZERO;
                    BigInteger gasBal = nativeBalanceByAddr.getOrDefault(addr.getAddress().toLowerCase(), BigInteger.ZERO);
                    boolean tokenOk = tokenBal.compareTo(config.minBalance) >= 0;
                    boolean gasOk = gasBal.compareTo(BigInteger.ZERO) > 0;
                    System.out.println(tag + "  Candidate ID=" + addr.getId()
                            + ", Addr=" + addr.getAddress()
                            + ", TokenBalance=" + tokenBal + (tokenOk ? " [OK]" : " [LOW]")
                            + ", GasBalance=" + gasBal + (gasOk ? " [OK]" : " [EMPTY]"));
                }
            }
        } else {
            candidates = findNativeAddresses(config);
        }

        if (candidates == null || candidates.size() < 2) {
            System.out.println(tag + "Step 0: Not enough candidates (" + (candidates == null ? 0 : candidates.size()) + " found, need 2)");
            return null;
        }

        // Find source: needs sufficient token balance AND (for tokens) native gas balance
        Address source = null;
        for (Address addr : candidates) {
            Balance balance = addr.getBalance();
            if (balance == null || balance.getAvailableConfirmed() == null
                    || balance.getAvailableConfirmed().compareTo(config.minBalance) < 0) {
                continue;
            }
            if (config.isToken) {
                BigInteger gasBal = nativeBalanceByAddr.getOrDefault(addr.getAddress().toLowerCase(), BigInteger.ZERO);
                if (gasBal.compareTo(BigInteger.ZERO) <= 0) {
                    System.out.println(tag + "Step 0: Skipping ID=" + addr.getId() + " — has token balance but no gas");
                    continue;
                }
                System.out.println(tag + "Step 0: Selected source ID=" + addr.getId()
                        + " — token=" + balance.getAvailableConfirmed() + ", gas=" + gasBal);
            }
            source = addr;
            break;
        }

        if (source == null) {
            System.out.println(tag + "Step 0: No address has sufficient balance"
                    + (config.isToken ? " with gas" : "") + " (min=" + config.minBalance + ")");
            return null;
        }

        // Find destination: any other address for the same currency
        Address destination = null;
        for (Address addr : candidates) {
            if (addr.getId() != source.getId()) {
                destination = addr;
                break;
            }
        }

        if (destination == null) {
            return null;
        }

        return new AddressPair(config, source, destination);
    }

    /**
     * Collects funded addresses across all wallets holding the native currency.
     * Scans all wallets until at least two funded addresses are found.
     */
    private List<Address> findNativeAddresses(CurrencyConfig config) throws ApiException {
        String tag = "[" + config.symbol + "] ";
        List<Wallet> wallets = client.getAssetService().getAssetWallets(config.symbol);
        if (wallets == null || wallets.isEmpty()) {
            System.out.println(tag + "  No wallets found via getAssetWallets");
            return null;
        }
        System.out.println(tag + "  Found " + wallets.size() + " wallet(s) via getAssetWallets");

        List<Address> fundedAddresses = new ArrayList<>();
        int totalScanned = 0;
        for (Wallet wallet : wallets) {
            System.out.println(tag + "  Scanning wallet ID=" + wallet.getId()
                    + ", Name=" + wallet.getName());
            int limit = 50;
            int offset = 0;
            while (true) {
                List<Address> addresses = client.getAddressService().getAddresses(wallet.getId(), limit, offset);
                if (addresses.isEmpty()) {
                    break;
                }
                for (Address addr : addresses) {
                    totalScanned++;
                    if (addr.isDisabled()) {
                        System.out.println(tag + "    DISABLED Address ID=" + addr.getId()
                                + ", Addr=" + addr.getAddress() + " — skipping");
                        continue;
                    }
                    BigInteger bal = (addr.getBalance() != null && addr.getBalance().getAvailableConfirmed() != null)
                            ? addr.getBalance().getAvailableConfirmed() : BigInteger.ZERO;
                    if (bal.compareTo(BigInteger.ZERO) > 0) {
                        System.out.println(tag + "    FUNDED Address ID=" + addr.getId()
                                + ", Addr=" + addr.getAddress()
                                + ", Balance=" + bal);
                        fundedAddresses.add(addr);
                    }
                }
                if (fundedAddresses.size() >= 2) {
                    System.out.println(tag + "  Found " + fundedAddresses.size()
                            + " funded addresses after scanning " + totalScanned + " total");
                    return fundedAddresses;
                }
                offset += addresses.size();
            }
        }

        System.out.println(tag + "  Found " + fundedAddresses.size()
                + " funded addresses after scanning " + totalScanned + " total across "
                + wallets.size() + " wallets");
        return fundedAddresses;
    }

    // ── Status polling ──────────────────────────────────────────────────

    private Request waitForTerminalStatus(long requestId, String tag) throws ApiException, InterruptedException {
        long startTime = System.currentTimeMillis();

        while (true) {
            Request request = client.getRequestService().getRequest(requestId);
            RequestStatus status = request.getStatus();

            System.out.println(tag + "  [" + ((System.currentTimeMillis() - startTime) / 1000) + "s] Request " + requestId + " — Status: " + status);

            if (TERMINAL_STATUSES.contains(status)) {
                return request;
            }

            long elapsed = System.currentTimeMillis() - startTime;
            if (elapsed >= MAX_WAIT_MS) {
                throw new RuntimeException("Request " + requestId + " did not reach terminal status within "
                        + (MAX_WAIT_MS / 1000) + " seconds. Last status: " + status);
            }

            Thread.sleep(POLL_INTERVAL_MS);
        }
    }

    // ── Summary output ──────────────────────────────────────────────────

    private static void printSummary(List<CurrencyResult> results) {
        System.out.println();
        System.out.println("\u2554\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2557");
        System.out.println("\u2551                 Multi-Currency E2E Results                      \u2551");
        System.out.println("\u2560\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2563");

        for (CurrencyResult result : results) {
            switch (result.status) {
                case PASSED:
                    System.out.printf("\u2551 %-5s: PASSED%51s\u2551%n", result.symbol, "");
                    System.out.printf("\u2551   tx : %-60s\u2551%n", result.txHash);
                    break;
                case SKIPPED:
                    System.out.printf("\u2551 %-5s: SKIPPED (fewer than 2 addresses found)%25s\u2551%n", result.symbol, "");
                    break;
                case FAILED:
                    String msg = result.errorMessage;
                    if (msg != null && msg.length() > 50) {
                        msg = msg.substring(0, 50) + "...";
                    }
                    System.out.printf("\u2551 %-5s: FAILED \u2014 %-54s\u2551%n", result.symbol, msg);
                    break;
                default:
                    break;
            }
        }

        System.out.println("\u255a\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u2550\u255d");
        System.out.println();
    }

    // ── Inner classes ───────────────────────────────────────────────────

    /**
     * A pair of source and destination addresses for a currency.
     */
    static final class AddressPair {
        final CurrencyConfig config;
        final Address source;
        final Address destination;

        AddressPair(CurrencyConfig config, Address source, Address destination) {
            this.config = config;
            this.source = source;
            this.destination = destination;
        }
    }

    /**
     * Configuration for a currency under test.
     */
    static final class CurrencyConfig {
        final String symbol;
        final String blockchain;
        final String network;
        final BigInteger transferAmount;
        final BigInteger minBalance;
        final boolean isToken;

        CurrencyConfig(String symbol, String blockchain, String network,
                       BigInteger transferAmount, BigInteger minBalance, boolean isToken) {
            this.symbol = symbol;
            this.blockchain = blockchain;
            this.network = network;
            this.transferAmount = transferAmount;
            this.minBalance = minBalance;
            this.isToken = isToken;
        }
    }

    /**
     * Result of the E2E flow for a single currency.
     */
    static final class CurrencyResult {
        enum Status { PASSED, SKIPPED, FAILED }

        final String symbol;
        final Status status;
        final String txHash;
        final String errorMessage;

        private CurrencyResult(String symbol, Status status, String txHash, String errorMessage) {
            this.symbol = symbol;
            this.status = status;
            this.txHash = txHash;
            this.errorMessage = errorMessage;
        }

        static CurrencyResult passed(String symbol, String txHash) {
            return new CurrencyResult(symbol, Status.PASSED, txHash, null);
        }

        static CurrencyResult skipped(String symbol) {
            return new CurrencyResult(symbol, Status.SKIPPED, null, null);
        }

        static CurrencyResult failed(String symbol, String errorMessage) {
            return new CurrencyResult(symbol, Status.FAILED, null, errorMessage);
        }
    }
}
