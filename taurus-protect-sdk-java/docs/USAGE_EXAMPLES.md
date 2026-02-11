# Usage Examples

This document provides complete code examples for common SDK operations.

## Table of Contents

1. [Client Setup](#client-setup)
2. [Wallet Management](#wallet-management)
3. [Address Management](#address-management)
4. [Request Approval Workflow](#request-approval-workflow)
5. [Transaction Queries](#transaction-queries)
6. [Balance Queries](#balance-queries)
7. [Whitelisted Addresses](#whitelisted-addresses)
8. [Webhook Management](#webhook-management)
9. [Staking Information](#staking-information)
10. [Contract Whitelisting](#contract-whitelisting)
11. [Pagination Patterns](#pagination-patterns)
12. [Error Handling](#error-handling)

---

## Client Setup

### Basic Initialization

```java
import com.taurushq.sdk.protect.client.ProtectClient;
import java.util.Arrays;

public class ProtectClientSetup {

    public static ProtectClient createClient() {
        // Configuration
        String host = "https://api.taurus-protect.com";
        String apiKey = "your-api-key-uuid";
        String apiSecret = "your-api-secret-hex";

        // SuperAdmin public keys (PEM format)
        String superAdmin1 = "-----BEGIN PUBLIC KEY-----\n" +
            "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n" +
            "-----END PUBLIC KEY-----";
        String superAdmin2 = "-----BEGIN PUBLIC KEY-----\n" +
            "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...\n" +
            "-----END PUBLIC KEY-----";

        // Create client
        return ProtectClient.createFromPem(
            host,
            apiKey,
            apiSecret,
            Arrays.asList(superAdmin1, superAdmin2),
            2  // Minimum valid signatures
        );
    }
}
```

### Environment-Based Configuration

```java
import com.taurushq.sdk.protect.client.ProtectClient;
import java.util.Arrays;
import java.util.List;
import java.util.stream.Collectors;

public class ProtectClientConfig {

    public static ProtectClient createFromEnvironment() {
        String host = System.getenv("TAURUS_API_HOST");
        String apiKey = System.getenv("TAURUS_API_KEY");
        String apiSecret = System.getenv("TAURUS_API_SECRET");

        // SuperAdmin keys as comma-separated PEM strings
        String keysEnv = System.getenv("TAURUS_SUPERADMIN_KEYS");
        List<String> superAdminKeys = Arrays.asList(keysEnv.split(","));

        int minSignatures = Integer.parseInt(
            System.getenv().getOrDefault("TAURUS_MIN_SIGNATURES", "2")
        );

        return ProtectClient.createFromPem(
            host, apiKey, apiSecret, superAdminKeys, minSignatures
        );
    }
}
```

---

## Wallet Management

### Create a Wallet

```java
import com.taurushq.sdk.protect.client.model.Wallet;

public void createWalletExample(ProtectClient client) throws Exception {
    // Basic wallet creation
    Wallet wallet = client.getWalletService().createWallet(
        "ETH",       // blockchain
        "mainnet",   // network
        "Treasury",  // name
        false        // omnibus flag
    );

    System.out.println("Created wallet: " + wallet.getId());
    System.out.println("Path: " + wallet.getAccountPath());

    // Wallet with metadata
    Wallet walletWithMeta = client.getWalletService().createWallet(
        "ETH",
        "mainnet",
        "Operations Wallet",
        false,
        "Used for daily operations",  // comment
        "CUST-12345"                   // customerId
    );
}
```

### List and Search Wallets

```java
import com.taurushq.sdk.protect.client.model.Wallet;
import java.util.ArrayList;
import java.util.List;

public void listWalletsExample(ProtectClient client) throws Exception {
    // List all wallets with pagination
    int limit = 50;
    int offset = 0;
    List<Wallet> allWallets = new ArrayList<>();

    List<Wallet> page;
    do {
        page = client.getWalletService().getWallets(limit, offset);
        allWallets.addAll(page);
        offset += limit;
    } while (page.size() == limit);

    System.out.println("Total wallets: " + allWallets.size());

    // Search by name
    List<Wallet> treasuryWallets = client.getWalletService()
        .getWalletsByName("Treasury", 10, 0);
}
```

### Get Wallet Balance History

```java
import com.taurushq.sdk.protect.client.model.BalanceHistoryPoint;
import java.util.List;

public void walletBalanceHistoryExample(ProtectClient client, long walletId) throws Exception {
    // Get 24-hour history with hourly intervals
    List<BalanceHistoryPoint> history = client.getWalletService()
        .getWalletBalanceHistory(walletId, 1);  // 1-hour intervals

    for (BalanceHistoryPoint point : history) {
        System.out.printf("%s: %s%n",
            point.getTimestamp(),
            point.getAvailableConfirmed()
        );
    }
}
```

---

## Address Management

### Create Addresses

```java
import com.taurushq.sdk.protect.client.model.Address;

public void createAddressExample(ProtectClient client, long walletId) throws Exception {
    Address address = client.getAddressService().createAddress(
        walletId,
        "Customer Deposit",    // label
        "Auto-generated",      // comment
        "USER-789"            // customerId
    );

    System.out.println("Address: " + address.getAddress());
    System.out.println("ID: " + address.getId());
}
```

### Get Address with Verification

```java
import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.Balance;

public void getAddressExample(ProtectClient client, long addressId) throws Exception {
    // This performs signature verification automatically
    Address address = client.getAddressService().getAddress(addressId);

    Balance balance = address.getBalance();
    System.out.println("Address: " + address.getAddress());
    System.out.println("Available: " + balance.getAvailableConfirmed());
    System.out.println("Pending: " + balance.getPendingConfirmed());
}
```

### Add Custom Attributes

```java
public void addAttributeExample(ProtectClient client, long addressId) throws Exception {
    client.getAddressService().createAddressAttribute(
        addressId,
        "department",
        "Finance"
    );

    client.getAddressService().createAddressAttribute(
        addressId,
        "costCenter",
        "CC-100"
    );
}
```

---

## Request Approval Workflow

### Complete Approval Flow

```java
import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestMetadata;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import java.math.BigInteger;
import java.security.PrivateKey;

public class RequestApprovalExample {

    public void completeApprovalFlow(ProtectClient client) throws Exception {
        // Step 1: Create a transfer request
        Request request = client.getRequestService().createExternalTransferRequest(
            1001L,                              // fromAddressId
            2001L,                              // toWhitelistedAddressId
            new BigInteger("1000000000000000000")  // 1 ETH in wei
        );

        System.out.println("Created request ID: " + request.getId());
        System.out.println("Status: " + request.getStatus());

        // Step 2: Retrieve and verify request
        Request retrieved = client.getRequestService().getRequest(request.getId());

        // Step 3: Review metadata before signing
        RequestMetadata metadata = retrieved.getMetadata();
        System.out.println("=== Review Request Metadata ===");
        System.out.println("Request ID: " + metadata.getRequestId());
        System.out.println("Currency: " + retrieved.getCurrency());
        System.out.println("Source: " + metadata.getSourceAddress());
        System.out.println("Destination: " + metadata.getDestinationAddress());
        System.out.println("Amount: " + metadata.getAmount());
        System.out.println("Hash: " + metadata.getHash());

        // Step 4: Approve with private key (after human review)
        String privateKeyPem = loadPrivateKey();  // From secure storage
        PrivateKey privateKey = CryptoTPV1.decodePrivateKey(privateKeyPem);

        int signatures = client.getRequestService().approveRequest(
            retrieved,
            privateKey
        );

        System.out.println("Signatures performed: " + signatures);

        // Step 5: Check updated status
        Request updated = client.getRequestService().getRequest(request.getId());
        System.out.println("New status: " + updated.getStatus());
    }

    private String loadPrivateKey() {
        // Load from secure storage (HSM, vault, etc.)
        return System.getenv("USER_PRIVATE_KEY");
    }
}
```

### Batch Approval

```java
import com.taurushq.sdk.protect.client.model.Request;
import com.taurushq.sdk.protect.client.model.RequestResult;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
import java.util.List;
import java.util.stream.Collectors;

public void batchApprovalExample(ProtectClient client, PrivateKey privateKey) throws Exception {
    // Get all requests pending my approval
    ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 100);
    RequestResult result = client.getRequestService().getRequestsForApproval(cursor);

    List<Request> toApprove = result.getRequests().stream()
        .filter(r -> r.getCurrency().equals("ETH"))  // Filter by currency
        .collect(Collectors.toList());

    if (!toApprove.isEmpty()) {
        // Batch approve
        int sigs = client.getRequestService().approveRequests(toApprove, privateKey);
        System.out.println("Approved " + toApprove.size() + " requests with " + sigs + " signatures");
    }
}
```

### Reject Requests

```java
public void rejectRequestExample(ProtectClient client, long requestId) throws Exception {
    client.getRequestService().rejectRequest(
        requestId,
        "Amount exceeds daily limit"
    );

    // Or batch reject
    client.getRequestService().rejectRequests(
        Arrays.asList(101L, 102L, 103L),
        "Suspicious activity detected"
    );
}
```

---

## Transaction Queries

### Query Transactions

```java
import com.taurushq.sdk.protect.client.model.Transaction;
import java.time.OffsetDateTime;
import java.util.List;

public void queryTransactionsExample(ProtectClient client) throws Exception {
    // Get by ID
    Transaction tx = client.getTransactionService().getTransactionById(12345L);
    System.out.println("Hash: " + tx.getHash());
    System.out.println("Status: " + tx.getStatus());

    // Get by hash
    Transaction txByHash = client.getTransactionService().getTransactionByHash(
        "0xe41578e07623a4a3646cf6393a512d975adbc4d6446849148d8c742069dfb34f"
    );

    // Query with filters
    OffsetDateTime from = OffsetDateTime.now().minusDays(7);
    OffsetDateTime to = OffsetDateTime.now();

    List<Transaction> recentOutgoing = client.getTransactionService().getTransactions(
        from,
        to,
        "ETH",       // currency
        "outgoing",  // direction
        100,         // limit
        0            // offset
    );

    for (Transaction t : recentOutgoing) {
        System.out.printf("%s: %s -> %s (%s)%n",
            t.getHash().substring(0, 10),
            t.getFrom(),
            t.getTo(),
            t.getValue()
        );
    }
}
```

### Export to CSV

```java
public void exportTransactionsExample(ProtectClient client) throws Exception {
    OffsetDateTime from = OffsetDateTime.now().minusMonths(1);
    OffsetDateTime to = OffsetDateTime.now();

    String csv = client.getTransactionService().exportTransactions(
        from, to, null, null, 1000, 0
    );

    // Save to file
    java.nio.file.Files.writeString(
        java.nio.file.Path.of("transactions.csv"),
        csv
    );
}
```

---

## Balance Queries

### Get All Balances

```java
import com.taurushq.sdk.protect.client.model.BalanceResult;
import com.taurushq.sdk.protect.client.model.AssetBalance;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
import java.util.ArrayList;
import java.util.List;

public void getAllBalancesExample(ProtectClient client) throws Exception {
    List<AssetBalance> allBalances = new ArrayList<>();
    ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 100);

    BalanceResult result;
    do {
        result = client.getBalanceService().getBalances(cursor);
        allBalances.addAll(result.getBalances());
        cursor = result.nextCursor(100);
    } while (result.hasNext());

    // Summarize by currency
    for (AssetBalance balance : allBalances) {
        System.out.printf("%s: Available=%s, Pending=%s%n",
            balance.getAsset().getSymbol(),
            balance.getBalance().getAvailableConfirmed(),
            balance.getBalance().getPendingConfirmed()
        );
    }
}
```

### NFT Balances

```java
import com.taurushq.sdk.protect.client.model.NFTCollectionBalanceResult;

public void getNFTBalancesExample(ProtectClient client) throws Exception {
    ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 50);

    NFTCollectionBalanceResult result = client.getBalanceService()
        .getNFTCollectionBalances("ETH", "mainnet", cursor);

    result.getBalances().forEach(nft -> {
        System.out.printf("Collection: %s, Count: %d%n",
            nft.getCollectionName(),
            nft.getCount()
        );
    });
}
```

---

## Whitelisted Addresses

### Get Verified Whitelisted Address

```java
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddressEnvelope;

public void getWhitelistedAddressExample(ProtectClient client) throws Exception {
    // Simple: get verified address directly
    WhitelistedAddress addr = client.getWhitelistedAddressService()
        .getWhitelistedAddress(1001L);

    System.out.println("Blockchain: " + addr.getBlockchain());
    System.out.println("Address: " + addr.getAddress());
    System.out.println("Label: " + addr.getLabel());

    // Advanced: get full envelope with approval details
    SignedWhitelistedAddressEnvelope envelope = client.getWhitelistedAddressService()
        .getWhitelistedAddressEnvelope(1001L);

    WhitelistedAddress verified = envelope.getWhitelistedAddress();
    System.out.println("Approvers: " + envelope.getApprovers());
    System.out.println("Trails: " + envelope.getTrails().size());
}
```

### List Whitelisted Addresses

```java
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddressEnvelope;
import java.util.List;

public void listWhitelistedAddressesExample(ProtectClient client) throws Exception {
    // List all
    List<SignedWhitelistedAddressEnvelope> addresses =
        client.getWhitelistedAddressService().getWhitelistedAddresses(100, 0);

    // Filter by blockchain
    List<SignedWhitelistedAddressEnvelope> ethAddresses =
        client.getWhitelistedAddressService().getWhitelistedAddresses(100, 0, "ETH");

    // Filter by blockchain and network
    List<SignedWhitelistedAddressEnvelope> mainnetAddresses =
        client.getWhitelistedAddressService().getWhitelistedAddresses(100, 0, "ETH", "mainnet");

    for (SignedWhitelistedAddressEnvelope env : mainnetAddresses) {
        WhitelistedAddress wla = env.getWhitelistedAddress();
        System.out.printf("%s: %s (%s)%n",
            wla.getBlockchain(),
            wla.getAddress(),
            wla.getLabel()
        );
    }
}
```

---

## Webhook Management

### Create and Configure Webhooks

```java
import com.taurushq.sdk.protect.client.model.Webhook;
import com.taurushq.sdk.protect.client.model.WebhookResult;
import com.taurushq.sdk.protect.client.model.WebhookStatus;

public void webhookExample(ProtectClient client) throws Exception {
    // Create a webhook for transaction notifications
    Webhook webhook = client.getWebhookService().createWebhook(
        "https://example.com/webhook/transactions",
        "TRANSACTION",
        "my-webhook-secret-key"
    );
    System.out.println("Created webhook: " + webhook.getId());

    // List all webhooks
    WebhookResult result = client.getWebhookService().getWebhooks(null, null, null);
    for (Webhook wh : result.getWebhooks()) {
        System.out.printf("Webhook: %s (%s) - Status: %s%n",
            wh.getId(), wh.getType(), wh.getStatus());
    }

    // Disable a webhook temporarily
    Webhook disabled = client.getWebhookService()
        .updateWebhookStatus(webhook.getId(), WebhookStatus.DISABLED);
    System.out.println("Webhook status: " + disabled.getStatus());

    // Re-enable the webhook
    client.getWebhookService().updateWebhookStatus(webhook.getId(), WebhookStatus.ENABLED);

    // Delete the webhook
    client.getWebhookService().deleteWebhook(webhook.getId());
}
```

---

## Staking Information

### Query Staking Across Blockchains

```java
import com.taurushq.sdk.protect.client.model.ADAStakePoolInfo;
import com.taurushq.sdk.protect.client.model.ETHValidatorInfo;
import com.taurushq.sdk.protect.client.model.StakeAccountResult;
import com.taurushq.sdk.protect.client.model.XTZStakingRewards;
import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.List;

public void stakingExample(ProtectClient client) throws Exception {
    // Get Cardano stake pool info
    ADAStakePoolInfo poolInfo = client.getStakingService()
        .getADAStakePoolInfo("mainnet", "pool1abc123...");
    System.out.println("Pool pledge: " + poolInfo.getPledge());
    System.out.println("Pool margin: " + poolInfo.getMargin());
    System.out.println("Active stake: " + poolInfo.getActiveStake());

    // Get Ethereum validator info
    List<ETHValidatorInfo> validators = client.getStakingService()
        .getETHValidatorsInfo("mainnet", Arrays.asList("validator1", "validator2"));
    for (ETHValidatorInfo v : validators) {
        System.out.printf("Validator %s: balance=%s, status=%s%n",
            v.getPublicKey(), v.getBalance(), v.getStatus());
    }

    // Get Tezos staking rewards
    XTZStakingRewards rewards = client.getStakingService()
        .getXTZStakingRewards(
            "mainnet",
            "address-123",
            OffsetDateTime.now().minusDays(30),
            OffsetDateTime.now()
        );
    System.out.println("Total rewards: " + rewards.getTotalRewards());

    // List stake accounts with pagination
    StakeAccountResult stakeAccounts = client.getStakingService()
        .getStakeAccounts(null, null, null, null);
    stakeAccounts.getStakeAccounts().forEach(account -> {
        System.out.printf("Account: %s, Type: %s%n",
            account.getAccountAddress(), account.getAccountType());
    });
}
```

---

## Contract Whitelisting

### Manage Whitelisted Contracts

```java
import com.taurushq.sdk.protect.client.model.SignedWhitelistedContractAddressEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistedContractAddressResult;
import com.taurushq.sdk.protect.client.model.Attribute;
import java.util.List;

public void contractWhitelistingExample(ProtectClient client) throws Exception {
    // Create a whitelisted ERC20 token
    String id = client.getContractWhitelistingService().createWhitelistedContract(
        "ETH",                                          // blockchain
        "mainnet",                                      // network
        "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", // USDC contract
        "USDC",                                         // symbol
        "USD Coin",                                     // name
        6,                                              // decimals
        "erc20",                                        // kind
        null                                            // tokenId (null for fungible)
    );
    System.out.println("Created whitelist entry: " + id);

    // Create a whitelisted NFT collection (ERC721)
    String nftId = client.getContractWhitelistingService().createWhitelistedContract(
        "ETH",
        "mainnet",
        "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d",
        "BAYC",
        "Bored Ape Yacht Club",
        0,           // 0 decimals for NFTs
        "erc721",
        null
    );

    // Get a whitelisted contract
    SignedWhitelistedContractAddressEnvelope contract =
        client.getContractWhitelistingService().getWhitelistedContract(id);
    System.out.println("Blockchain: " + contract.getBlockchain());
    System.out.println("Status: " + contract.getStatus());

    // List whitelisted contracts with filtering
    WhitelistedContractAddressResult result = client.getContractWhitelistingService()
        .getWhitelistedContracts("ETH", "mainnet", null, false, 50, 0);

    System.out.println("Total contracts: " + result.getTotalItems());
    for (SignedWhitelistedContractAddressEnvelope c : result.getContracts()) {
        System.out.printf("  %s: %s%n", c.getBlockchain(), c.getId());
    }

    // Add an attribute to a contract
    List<Attribute> attrs = client.getContractWhitelistingService().createAttribute(
        id,
        "coingecko_id",
        "usd-coin",
        "text/plain",
        "metadata",
        null
    );
    System.out.println("Added attribute: " + attrs.get(0).getKey());

    // List contracts pending approval
    WhitelistedContractAddressResult pending = client.getContractWhitelistingService()
        .getWhitelistedContractsForApproval(null, 50, 0);
    System.out.println("Contracts pending approval: " + pending.getTotalItems());

    // Update contract metadata
    client.getContractWhitelistingService().updateWhitelistedContract(
        id, "USDC.e", "USD Coin (Bridged)", 6
    );

    // Delete a whitelisted contract
    String deleteId = client.getContractWhitelistingService()
        .deleteWhitelistedContract(id, "No longer supported");
}
```

---

## Pagination Patterns

### Cursor-Based Pagination

```java
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;

public void cursorPaginationExample(ProtectClient client) throws Exception {
    // Initialize cursor for first page
    ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 100);

    int totalItems = 0;
    int pageNumber = 0;

    BalanceResult result;
    do {
        result = client.getBalanceService().getBalances(cursor);

        pageNumber++;
        totalItems += result.getBalances().size();
        System.out.printf("Page %d: %d items%n", pageNumber, result.getBalances().size());

        // Get cursor for next page
        cursor = result.nextCursor(100);

    } while (result.hasNext());

    System.out.println("Total items: " + totalItems);
}
```

### Offset-Based Pagination

```java
public void offsetPaginationExample(ProtectClient client) throws Exception {
    int pageSize = 50;
    int offset = 0;
    int totalWallets = 0;

    List<Wallet> page;
    do {
        page = client.getWalletService().getWallets(pageSize, offset);
        totalWallets += page.size();

        // Process page
        for (Wallet wallet : page) {
            processWallet(wallet);
        }

        offset += pageSize;
    } while (page.size() == pageSize);

    System.out.println("Processed " + totalWallets + " wallets");
}
```

---

## Error Handling

### Basic Exception Handling

```java
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.WhitelistException;

public void errorHandlingExample(ProtectClient client) {
    try {
        Wallet wallet = client.getWalletService().getWallet(999999L);
    } catch (IntegrityException e) {
        // Cryptographic verification failed
        System.err.println("Integrity check failed: " + e.getMessage());
        // This indicates potential tampering - investigate!
    } catch (WhitelistException e) {
        // Whitelist-specific verification failed
        System.err.println("Whitelist verification failed: " + e.getMessage());
    } catch (ApiException e) {
        // General API error
        System.err.println("API Error:");
        System.err.println("  HTTP Status: " + e.getCode());
        System.err.println("  Error Code: " + e.getErrorCode());
        System.err.println("  Message: " + e.getMessage());
        System.err.println("  Error: " + e.getError());

        // Handle specific error codes
        if (e.getCode() == 404) {
            System.err.println("Resource not found");
        } else if (e.getCode() == 401) {
            System.err.println("Authentication failed - check credentials");
        } else if (e.getCode() == 403) {
            System.err.println("Permission denied");
        }
    }
}
```

### Retry Logic

```java
import java.util.concurrent.TimeUnit;

public <T> T withRetry(Callable<T> operation, int maxRetries) throws Exception {
    int attempt = 0;
    while (true) {
        try {
            return operation.call();
        } catch (ApiException e) {
            attempt++;
            if (attempt >= maxRetries || !isRetryable(e)) {
                throw e;
            }
            // Exponential backoff
            TimeUnit.MILLISECONDS.sleep((long) Math.pow(2, attempt) * 100);
        }
    }
}

private boolean isRetryable(ApiException e) {
    // Retry on server errors and rate limiting
    return e.getCode() >= 500 || e.getCode() == 429;
}

// Usage
Wallet wallet = withRetry(() -> client.getWalletService().getWallet(walletId), 3);
```

---

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Authentication](AUTHENTICATION.md) - Security and signing
- [Services Reference](SERVICES.md) - Complete API documentation
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Verification details
