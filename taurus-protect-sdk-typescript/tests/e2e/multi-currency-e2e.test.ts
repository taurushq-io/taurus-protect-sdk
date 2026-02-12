/**
 * Multi-currency end-to-end integration test.
 *
 * Runs the full transfer lifecycle for multiple currencies in parallel.
 *
 * The test proceeds in two phases:
 *   1. Discovery phase (sequential): For each currency, find two existing
 *      funded addresses to use as source and destination.
 *   2. Transfer phase (parallel): Run the transfer flow concurrently for all
 *      currencies that have two funded addresses.
 *
 * Per-currency steps:
 *   Step 0: Find two funded addresses for the currency
 *   Step 1: Create an internal transfer request between them
 *   Step 2: Verify metadata matches the original transfer intent
 *   Step 3: Approve the request
 *   Step 4: Wait for BROADCASTED/CONFIRMED status
 *   Step 5: Verify the transaction by hash
 *
 * Currencies without two funded addresses are skipped. The test passes if at
 * least one currency completes successfully.
 */

import { skipIfNotEnabled, getTestClient, getPrivateKey } from "../testutil";
import { TEAM1_PRIVATE_KEY_PEM } from "../integration/config";
import { ValidationError } from "../../src/errors";
import { RequestStatus } from "../../src/models/request";
import type { Request } from "../../src/models/request";
import type { Wallet } from "../../src/models/wallet";
import type { ProtectClient } from "../../src/client";
import { decodePrivateKeyPem } from "../../src/crypto/keys";
import { getSourceAddress, getDestinationAddress, getAmount } from "../../src/helpers/metadata-utils";
import type {
  TgvalidatordAddress,
  TgvalidatordBalance,
} from "../../src/internal/openapi";

// ── Constants ───────────────────────────────────────────────────────────

/** Maximum time to wait for a single request confirmation (15 minutes). */
const MAX_WAIT_MS = 15 * 60 * 1000;

/** Polling interval between status checks (5 seconds). */
const POLL_INTERVAL_MS = 5000;

/** Maximum concurrent transfer flows to avoid API throttling. */
const MAX_PARALLEL_TRANSFERS = 3;

/** Terminal statuses that indicate the request will not progress further. */
const TERMINAL_STATUSES = new Set<string>([
  RequestStatus.BROADCASTED,
  RequestStatus.CONFIRMED,
  RequestStatus.REJECTED,
  RequestStatus.CANCELED,
  RequestStatus.PERMANENT_FAILURE,
  RequestStatus.EXPIRED,
  "INVALID",
  RequestStatus.MINED,
]);

const SUCCESS_STATUSES = new Set<string>([
  RequestStatus.BROADCASTED,
  RequestStatus.CONFIRMED,
]);

// ── Currency configuration ──────────────────────────────────────────────

interface CurrencyConfig {
  symbol: string;
  blockchain: string;
  network: string;
  transferAmount: bigint;
  minBalance: bigint;
  isToken: boolean;
}

function getCurrencyConfigs(): CurrencyConfig[] {
  return [
    {
      symbol: "SOL", blockchain: "SOL", network: "mainnet",
      transferAmount: 1_000_000n,         // 0.001 SOL (lamports)
      minBalance: 2_000_000n,
      isToken: false,
    },
    // {
    //   symbol: "ETH", blockchain: "ETH", network: "mainnet",
    //   transferAmount: 10_000_000_000_000n, // 0.00001 ETH (wei)
    //   minBalance: 100_000_000_000_000n,
    //   isToken: false,
    // },
    // {
    //   symbol: "XRP", blockchain: "XRP", network: "mainnet",
    //   transferAmount: 1212n,               // 0.001212 XRP (drops)
    //   minBalance: 12_000_000n,             // 12 XRP min balance
    //   isToken: false,
    // },
    {
      symbol: "XLM", blockchain: "XLM", network: "mainnet",
      transferAmount: 1313n,               // 0.0001313 XLM (stroops)
      minBalance: 35_000_000n,             // 3.5 XLM min balance
      isToken: false,
    },
    {
      symbol: "ALGO", blockchain: "ALGO", network: "mainnet",
      transferAmount: 1414n,               // 0.001414 ALGO (microAlgos)
      minBalance: 100_000n,                // 0.1 ALGO min balance
      isToken: false,
    },
    // {
    //   symbol: "USDC", blockchain: "ETH", network: "mainnet",
    //   transferAmount: 11_100n,             // 0.0111 USDC (6 decimals)
    //   minBalance: 2_000_000n,              // 2.0 USDC min balance
    //   isToken: true,
    // },
  ];
}

// ── Result tracking ─────────────────────────────────────────────────────

enum ResultStatus { PASSED = "PASSED", SKIPPED = "SKIPPED", FAILED = "FAILED" }

interface CurrencyResult {
  symbol: string;
  status: ResultStatus;
  txHash?: string;
  errorMessage?: string;
}

function passedResult(symbol: string, txHash?: string): CurrencyResult {
  return { symbol, status: ResultStatus.PASSED, txHash };
}

function skippedResult(symbol: string): CurrencyResult {
  return { symbol, status: ResultStatus.SKIPPED };
}

function failedResult(symbol: string, errorMessage: string): CurrencyResult {
  return { symbol, status: ResultStatus.FAILED, errorMessage };
}

// ── Address pair ────────────────────────────────────────────────────────

interface AddressPair {
  config: CurrencyConfig;
  source: TgvalidatordAddress;
  destination: TgvalidatordAddress;
}

// ── Helpers ─────────────────────────────────────────────────────────────

function sleep(ms: number): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

function getAvailableConfirmed(balance: TgvalidatordBalance | undefined): bigint {
  if (!balance?.availableConfirmed) return 0n;
  try {
    return BigInt(balance.availableConfirmed);
  } catch {
    return 0n;
  }
}

// ── Address discovery ───────────────────────────────────────────────────

/**
 * Finds two funded addresses for the given currency.
 * The first address (source) must have at least minBalance.
 * The second address (destination) can have any balance.
 */
async function findTwoFundedAddresses(
  client: ProtectClient,
  config: CurrencyConfig,
): Promise<AddressPair | undefined> {
  const tag = `[${config.symbol}]`;
  let candidates: TgvalidatordAddress[];

  // For tokens, build a native balance lookup so we can check gas availability
  const nativeBalanceByAddr = new Map<string, bigint>();

  if (config.isToken) {
    console.log(`${tag}   Step 0: Searching via AssetService...`);
    const tokenReply = await client.assetsApi.walletServiceGetAssetAddresses({
      body: { asset: { currency: config.symbol } },
    });
    candidates = tokenReply.addresses ?? [];
    console.log(`${tag}   Step 0: AssetService returned ${candidates.length} token addresses`);

    // Fetch native currency addresses to check gas balance
    console.log(`${tag}   Step 0: Fetching ${config.blockchain} addresses for gas balance check...`);
    const gasReply = await client.assetsApi.walletServiceGetAssetAddresses({
      body: { asset: { currency: config.blockchain } },
    });
    const nativeAddresses = gasReply.addresses ?? [];
    for (const nAddr of nativeAddresses) {
      const nBal = getAvailableConfirmed(nAddr.balance);
      nativeBalanceByAddr.set((nAddr.address ?? "").toLowerCase(), nBal);
    }
    console.log(`${tag}   Step 0: Found ${nativeAddresses.length} ${config.blockchain} addresses`);

    for (const addr of candidates) {
      const tokenBal = getAvailableConfirmed(addr.balance);
      const gasBal = nativeBalanceByAddr.get((addr.address ?? "").toLowerCase()) ?? 0n;
      const tokenOk = tokenBal >= config.minBalance;
      const gasOk = gasBal > 0n;
      console.log(
        `${tag}     Candidate ID=${addr.id}, Addr=${addr.address}` +
        `, TokenBalance=${tokenBal}${tokenOk ? " [OK]" : " [LOW]"}` +
        `, GasBalance=${gasBal}${gasOk ? " [OK]" : " [EMPTY]"}`,
      );
    }
    // Filter out disabled addresses
    candidates = candidates.filter(addr => {
      if (addr.disabled) {
        console.log(`${tag}   DISABLED Address ID=${addr.id}, Addr=${addr.address} -- skipping`);
        return false;
      }
      return true;
    });
  } else {
    candidates = await findNativeAddresses(client, config);
  }

  if (candidates.length < 2) {
    console.log(`${tag}   Step 0: Not enough candidates (${candidates.length} found, need 2)`);
    return undefined;
  }

  // Find source: needs sufficient balance AND (for tokens) native gas balance
  let source: TgvalidatordAddress | undefined;
  for (const addr of candidates) {
    const balance = getAvailableConfirmed(addr.balance);
    if (balance < config.minBalance) continue;

    if (config.isToken) {
      const gasBal = nativeBalanceByAddr.get((addr.address ?? "").toLowerCase()) ?? 0n;
      if (gasBal <= 0n) {
        console.log(`${tag}   Step 0: Skipping ID=${addr.id} -- has token balance but no gas`);
        continue;
      }
      console.log(`${tag}   Step 0: Selected source ID=${addr.id} -- token=${balance}, gas=${gasBal}`);
    }
    source = addr;
    break;
  }

  if (!source) {
    console.log(
      `${tag}   Step 0: No address has sufficient balance` +
      `${config.isToken ? " with gas" : ""} (min=${config.minBalance})`,
    );
    return undefined;
  }

  // Find destination: any other address for the same currency
  let destination: TgvalidatordAddress | undefined;
  for (const addr of candidates) {
    if (addr.id !== source.id) {
      destination = addr;
      break;
    }
  }

  if (!destination) return undefined;

  return { config, source, destination };
}

/**
 * Collects funded addresses across all wallets holding the native currency.
 * Scans all wallets until at least two funded addresses are found.
 */
async function findNativeAddresses(
  client: ProtectClient,
  config: CurrencyConfig,
): Promise<TgvalidatordAddress[]> {
  const tag = `[${config.symbol}]`;

  const wallets: Wallet[] = await client.assets.getAssetWallets({ currency: config.symbol });
  if (wallets.length === 0) {
    console.log(`${tag}     No wallets found via getAssetWallets`);
    return [];
  }
  console.log(`${tag}     Found ${wallets.length} wallet(s) via getAssetWallets`);

  const fundedAddresses: TgvalidatordAddress[] = [];
  let totalScanned = 0;

  for (const wallet of wallets) {
    console.log(`${tag}     Scanning wallet ID=${wallet.id}, Name=${wallet.name}`);
    const limit = 50;
    let offset = 0;

    while (true) {
      // Use raw API to get addresses WITH balance
      const response = await client.addressesApi.walletServiceGetAddresses({
        walletId: wallet.id,
        limit: String(limit),
        offset: String(offset),
      });

      const addresses = response.result ?? [];
      if (addresses.length === 0) break;

      for (const addr of addresses) {
        totalScanned++;
        if (addr.disabled) {
          console.log(
            `${tag}       DISABLED Address ID=${addr.id}, Addr=${addr.address} -- skipping`,
          );
          continue;
        }
        const bal = getAvailableConfirmed(addr.balance);
        if (bal > 0n) {
          console.log(
            `${tag}       FUNDED Address ID=${addr.id}, Addr=${addr.address}, Balance=${bal}`,
          );
          fundedAddresses.push(addr);
        }
      }

      if (fundedAddresses.length >= 2) {
        console.log(
          `${tag}     Found ${fundedAddresses.length} funded addresses after scanning ${totalScanned} total`,
        );
        return fundedAddresses;
      }
      offset += addresses.length;
    }
  }

  console.log(
    `${tag}     Found ${fundedAddresses.length} funded addresses after scanning ${totalScanned}` +
    ` total across ${wallets.length} wallets`,
  );
  return fundedAddresses;
}

// ── Transfer flow ───────────────────────────────────────────────────────

/**
 * Runs the transfer flow for a single currency using two existing addresses.
 */
async function runTransferFlow(
  client: ProtectClient,
  pair: AddressPair,
  approvalKey: ReturnType<typeof decodePrivateKeyPem>,
): Promise<CurrencyResult> {
  const config = pair.config;
  const tag = `[${config.symbol}]`;

  try {
    const source = pair.source;
    const destination = pair.destination;

    // Step 1: Create transfer request
    console.log(
      `${tag} Step 1: Creating transfer of ${config.transferAmount}` +
      ` from address ${source.id} to address ${destination.id}...`,
    );
    let transferRequest: Request;
    try {
      transferRequest = await client.requests.createInternalTransferRequest({
        fromAddressId: parseInt(source.id ?? "0", 10),
        toAddressId: parseInt(destination.id ?? "0", 10),
        amount: config.transferAmount.toString(),
      });
    } catch (stepErr: unknown) {
      // 400 Bad Request at Step 1 typically means the environment/currency is
      // not properly configured (insufficient funds, blockchain maintenance,
      // business-rule limits, etc.) -- skip rather than fail.
      if (stepErr instanceof ValidationError) {
        const msg = stepErr.message;
        console.log(`${tag} Step 1: Skipped (400 Bad Request): ${msg}`);
        return skippedResult(config.symbol);
      }
      throw stepErr;
    }
    expect(transferRequest).toBeDefined();
    expect(transferRequest.id).toBeGreaterThan(0);
    console.log(
      `${tag} Step 1: Created request: ID=${transferRequest.id}, Status=${transferRequest.status}`,
    );

    // Step 2: Verify metadata matches original intent
    console.log(`${tag} Step 2: Verifying metadata matches transfer intent...`);
    const requestToApprove: Request = await client.requests.get(transferRequest.id);
    const metadata = requestToApprove.metadata;
    expect(metadata).toBeDefined();
    expect(metadata!.hash).toBeDefined();

    // Verify source address in metadata matches what we sent
    if (metadata!.payloadAsString) {
      const metadataSource = getSourceAddress(metadata!.payloadAsString);
      expect(metadataSource).toBe(source.address);
      console.log(`${tag} Step 2: Source verified: ${metadataSource}`);

      // Verify destination address in metadata
      // For token transfers, the metadata destination is the token contract address,
      // not the recipient address
      const metadataDestination = getDestinationAddress(metadata!.payloadAsString);
      if (!config.isToken) {
        expect(metadataDestination).toBe(destination.address);
      }
      console.log(
        `${tag} Step 2: Destination: ${metadataDestination}` +
        `${config.isToken ? " (token contract)" : " (verified)"}`,
      );

      // Verify amount in metadata
      // For token transfers, the native value is 0 (token amount is in contract call data)
      const metadataAmount = getAmount(metadata!.payloadAsString);
      expect(metadataAmount).toBeDefined();
      if (!config.isToken) {
        expect(metadataAmount!.valueFrom).toBe(config.transferAmount.toString());
      }
      console.log(
        `${tag} Step 2: Amount: valueFrom=${metadataAmount!.valueFrom}` +
        `, currency=${metadataAmount!.currencyFrom}` +
        `${config.isToken ? " (token transfer -- native value is 0)" : " (verified)"}`,
      );
    }

    // Step 3: Approve request
    console.log(`${tag} Step 3: Approving request ${transferRequest.id}...`);
    const signedCount = await client.requests.approveRequest(requestToApprove, approvalKey);
    expect(signedCount).toBeGreaterThan(0);
    console.log(`${tag} Step 3: Approved: signedCount=${signedCount}`);

    // Step 4: Wait for terminal status
    console.log(`${tag} Step 4: Waiting for terminal status...`);
    const confirmedRequest = await waitForTerminalStatus(client, transferRequest.id, tag);
    console.log(`${tag} Step 4: Final status: ${confirmedRequest.status}`);

    // Dump diagnostics if not a success status
    if (!SUCCESS_STATUSES.has(confirmedRequest.status as string)) {
      console.log(`${tag} === DIAGNOSTIC DUMP ===`);
      if (confirmedRequest.signedRequests.length > 0) {
        for (const sr of confirmedRequest.signedRequests) {
          console.log(
            `${tag}   SignedRequest ID=${sr.id}, Status=${sr.status}` +
            `, Hash=${sr.hash}, Details=${sr.details}`,
          );
        }
      }
      // Trails still via raw API (Request model lacks trails field)
      const rawResponse = await client.requestsApi.requestServiceGetRequest({
        id: String(transferRequest.id),
      });
      const rawRequest = rawResponse.result;
      if (rawRequest?.trails) {
        console.log(`${tag}   Trails:`);
        for (const trail of rawRequest.trails) {
          console.log(
            `${tag}     action=${trail.action}, status=${trail.requestStatus}` +
            `, comment=${trail.comment}, date=${trail.date}`,
          );
        }
      }
      console.log(`${tag} === END DIAGNOSTIC ===`);
      throw new Error(
        `${config.symbol} request ended with ${confirmedRequest.status} instead of BROADCASTED/CONFIRMED`,
      );
    }

    // Step 5: Verify transaction by hash
    console.log(`${tag} Step 5: Fetching transaction details...`);
    const txHash = confirmedRequest.signedRequests.length > 0
      ? confirmedRequest.signedRequests[0].hash
      : undefined;
    expect(txHash).toBeDefined();
    console.log(`${tag} Step 5: Transaction hash: ${txHash}`);

    let transaction: any = undefined;
    const txLookupDeadline = Date.now() + 90_000;
    while (!transaction && Date.now() < txLookupDeadline) {
      await sleep(POLL_INTERVAL_MS);
      try {
        transaction = await client.transactions.getByHash(txHash!);
      } catch {
        console.log(`${tag} Step 5: Transaction not indexed yet, retrying...`);
      }
    }
    expect(transaction).toBeDefined();
    expect(transaction!.hash).toBe(txHash);
    expect(transaction!.direction).toBe("outgoing");
    expect(transaction!.requestId).toBe(String(transferRequest.id));
    console.log(
      `${tag} Step 5: Transaction verified -- ID=${transaction!.id}, Block=${transaction!.block}`,
    );

    console.log(`${tag} E2E PASSED`);
    return passedResult(config.symbol, txHash);

  } catch (e) {
    const msg = e instanceof Error ? e.message : String(e);
    console.log(`${tag} E2E FAILED: ${msg}`);
    return failedResult(config.symbol, msg);
  }
}

// ── Status polling ──────────────────────────────────────────────────────

async function waitForTerminalStatus(
  client: ProtectClient,
  requestId: number,
  tag: string,
): Promise<Request> {
  const startTime = Date.now();

  while (true) {
    const request = await client.requests.get(requestId);
    const status = request.status as string;
    const elapsed = Math.floor((Date.now() - startTime) / 1000);
    console.log(`${tag}   [${elapsed}s] Request ${requestId} -- Status: ${status}`);

    if (TERMINAL_STATUSES.has(status)) {
      return request;
    }

    if (Date.now() - startTime >= MAX_WAIT_MS) {
      throw new Error(
        `Request ${requestId} did not reach terminal status within ` +
        `${MAX_WAIT_MS / 1000} seconds. Last status: ${status}`,
      );
    }

    await sleep(POLL_INTERVAL_MS);
  }
}

// ── Summary output ──────────────────────────────────────────────────────

function printSummary(results: CurrencyResult[]): void {
  console.log("");
  console.log("=".repeat(65));
  console.log("                 Multi-Currency E2E Results");
  console.log("=".repeat(65));

  for (const result of results) {
    switch (result.status) {
      case ResultStatus.PASSED:
        console.log(` ${result.symbol.padEnd(6)}: PASSED`);
        if (result.txHash) {
          console.log(`   tx : ${result.txHash}`);
        }
        break;
      case ResultStatus.SKIPPED:
        console.log(` ${result.symbol.padEnd(6)}: SKIPPED (fewer than 2 addresses found)`);
        break;
      case ResultStatus.FAILED: {
        let msg = result.errorMessage ?? "";
        if (msg.length > 50) msg = msg.substring(0, 50) + "...";
        console.log(` ${result.symbol.padEnd(6)}: FAILED -- ${msg}`);
        break;
      }
    }
  }

  console.log("=".repeat(65));
  console.log("");
}

// ── Concurrency limiter ─────────────────────────────────────────────────

/**
 * Run async tasks with limited concurrency (mirrors Java's fixed thread pool).
 */
async function runWithConcurrency<T>(
  tasks: (() => Promise<T>)[],
  maxConcurrency: number,
): Promise<T[]> {
  const results: T[] = [];
  const executing = new Set<Promise<void>>();

  for (const task of tasks) {
    const p = task().then((result) => {
      results.push(result);
    });

    const wrapped = p.then(
      () => { executing.delete(wrapped); },
      () => { executing.delete(wrapped); },
    );
    executing.add(wrapped);

    if (executing.size >= maxConcurrency) {
      await Promise.race(executing);
    }
  }

  await Promise.allSettled(executing);
  return results;
}

// ── Test ─────────────────────────────────────────────────────────────────

describe("Integration: Multi-Currency E2E", () => {
  it("should complete transfer lifecycle for multiple currencies", async () => {
    try {
      skipIfNotEnabled();
    } catch (e) {
      if (e instanceof Error && e.message === "SKIP_TEST") return;
      throw e;
    }

    const client = getTestClient();
    try {
      const approvalKey = decodePrivateKeyPem(TEAM1_PRIVATE_KEY_PEM);
      const configs = getCurrencyConfigs();

      // Phase 1: Sequential address discovery to avoid API overload
      console.log("=== Phase 1: Address Discovery (sequential) ===");
      const addressPairs: AddressPair[] = [];
      const results: CurrencyResult[] = [];

      for (const config of configs) {
        const tag = `[${config.symbol}]`;
        console.log(`${tag} Step 0: Searching for two funded addresses...`);
        try {
          const pair = await findTwoFundedAddresses(client, config);
          if (!pair) {
            console.log(`${tag} Step 0: Could not find two funded addresses -- skipping`);
            results.push(skippedResult(config.symbol));
          } else {
            console.log(
              `${tag} Step 0: Source: ID=${pair.source.id}` +
              ` (${pair.source.address})` +
              `, Balance=${getAvailableConfirmed(pair.source.balance)}`,
            );
            console.log(
              `${tag} Step 0: Dest:   ID=${pair.destination.id}` +
              ` (${pair.destination.address})`,
            );
            addressPairs.push(pair);
          }
        } catch (e) {
          const msg = e instanceof Error ? e.message : String(e);
          console.log(`${tag} Step 0: Discovery failed -- ${msg}`);
          results.push(skippedResult(config.symbol));
        }
      }

      if (addressPairs.length === 0) {
        console.log("No currency has two funded addresses. Cannot run E2E test.");
        expect(addressPairs.length).toBeGreaterThan(0);
        return;
      }

      // Phase 2: Parallel transfer flows
      console.log(
        `\n=== Phase 2: Transfer Flows (parallel, max ${MAX_PARALLEL_TRANSFERS} threads) ===`,
      );

      const transferTasks = addressPairs.map(
        (pair) => () => runTransferFlow(client, pair, approvalKey),
      );

      const transferResults = await runWithConcurrency(transferTasks, MAX_PARALLEL_TRANSFERS);
      results.push(...transferResults);

      // Print summary
      printSummary(results);

      // Assertions: at least one currency must pass.
      let passedCount = 0;
      let failedCount = 0;
      for (const result of results) {
        if (result.status === ResultStatus.PASSED) {
          passedCount++;
        } else if (result.status === ResultStatus.FAILED) {
          failedCount++;
          console.log(`WARNING: ${result.symbol} failed: ${result.errorMessage}`);
        }
      }

      const skippedCount = results.length - passedCount - failedCount;
      console.log(
        `Summary: ${passedCount} passed, ${failedCount} failed, ${skippedCount} skipped`,
      );

      expect(failedCount).toBe(0);
      expect(passedCount).toBeGreaterThan(0);
    } finally {
      client.close();
    }
  }, 25 * 60 * 1000); // 25 minute timeout
});
