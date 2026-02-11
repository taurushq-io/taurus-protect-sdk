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
9. [Pagination Patterns](#pagination-patterns)
10. [Error Handling](#error-handling)
11. [Context and Cancellation](#context-and-cancellation)

---

## Client Setup

### Basic Initialization

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

func main() {
    // Configuration
    host := "https://api.taurus-protect.com"
    apiKey := "your-api-key-uuid"
    apiSecret := "your-api-secret-hex"

    // Create client
    client, err := protect.NewClient(
        host,
        protect.WithCredentials(apiKey, apiSecret),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer client.Close()

    // Use client
    ctx := context.Background()
    wallets, _, err := client.Wallets().ListWallets(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list wallets: %v", err)
    }

    fmt.Printf("Found %d wallets\n", len(wallets))
}
```

### With SuperAdmin Keys (Recommended)

```go
package main

import (
    "time"

    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

func createClient() (*protect.Client, error) {
    host := "https://api.taurus-protect.com"
    apiKey := "your-api-key-uuid"
    apiSecret := "your-api-secret-hex"

    // SuperAdmin public keys in PEM format
    superAdminKeys := []string{
        `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----`,
        `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE...
-----END PUBLIC KEY-----`,
    }

    return protect.NewClient(
        host,
        protect.WithCredentials(apiKey, apiSecret),
        protect.WithSuperAdminKeysPEM(superAdminKeys),
        protect.WithMinValidSignatures(2),
        protect.WithRulesCacheTTL(10 * time.Minute),
        protect.WithHTTPTimeout(30 * time.Second),
    )
}
```

### Environment-Based Configuration

```go
package main

import (
    "os"
    "strconv"
    "strings"
    "time"

    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

func createClientFromEnv() (*protect.Client, error) {
    host := os.Getenv("TAURUS_API_HOST")
    apiKey := os.Getenv("TAURUS_API_KEY")
    apiSecret := os.Getenv("TAURUS_API_SECRET")

    opts := []protect.ClientOption{
        protect.WithCredentials(apiKey, apiSecret),
    }

    // Optional SuperAdmin keys
    if keysEnv := os.Getenv("TAURUS_SUPERADMIN_KEYS"); keysEnv != "" {
        keys := strings.Split(keysEnv, "|||") // Use delimiter for multiple PEM keys
        opts = append(opts, protect.WithSuperAdminKeysPEM(keys))

        minSigs := 2
        if v := os.Getenv("TAURUS_MIN_SIGNATURES"); v != "" {
            minSigs, _ = strconv.Atoi(v)
        }
        opts = append(opts, protect.WithMinValidSignatures(minSigs))
    }

    // Optional cache TTL
    if ttlEnv := os.Getenv("TAURUS_CACHE_TTL_SECONDS"); ttlEnv != "" {
        if ttlSecs, err := strconv.Atoi(ttlEnv); err == nil {
            opts = append(opts, protect.WithRulesCacheTTL(time.Duration(ttlSecs)*time.Second))
        }
    }

    return protect.NewClient(host, opts...)
}
```

---

## Wallet Management

### Create a Wallet

```go
func createWalletExample(ctx context.Context, client *protect.Client) error {
    // Basic wallet creation
    wallet, err := client.Wallets().CreateWallet(ctx, &model.CreateWalletRequest{
        Name:     "Treasury",
        Currency: "ETH",
    })
    if err != nil {
        return fmt.Errorf("create wallet: %w", err)
    }

    fmt.Printf("Created wallet: %s (ID: %s)\n", wallet.Name, wallet.ID)

    // Wallet with metadata
    walletWithMeta, err := client.Wallets().CreateWallet(ctx, &model.CreateWalletRequest{
        Name:       "Operations Wallet",
        Currency:   "ETH",
        Comment:    "Main operations wallet",
        CustomerID: "CUST-12345",
    })
    if err != nil {
        return fmt.Errorf("create wallet with metadata: %w", err)
    }

    fmt.Printf("Created wallet with customer ID: %s\n", walletWithMeta.CustomerID)
    return nil
}
```

### List and Search Wallets

```go
func listWalletsExample(ctx context.Context, client *protect.Client) error {
    // List all wallets with pagination
    var allWallets []*model.Wallet
    opts := &model.ListWalletsOptions{
        Limit:  50,
        Offset: 0,
    }

    for {
        wallets, pagination, err := client.Wallets().ListWallets(ctx, opts)
        if err != nil {
            return fmt.Errorf("list wallets: %w", err)
        }

        allWallets = append(allWallets, wallets...)

        if !pagination.HasMore {
            break
        }
        opts.Offset += opts.Limit
    }

    fmt.Printf("Total wallets: %d\n", len(allWallets))

    // Search by currency
    ethWallets, _, err := client.Wallets().ListWallets(ctx, &model.ListWalletsOptions{
        Limit:    10,
        Currency: "ETH",
    })
    if err != nil {
        return fmt.Errorf("search wallets: %w", err)
    }

    fmt.Printf("Found %d ETH wallets\n", len(ethWallets))
    return nil
}
```

### Get Wallet Details

```go
func getWalletExample(ctx context.Context, client *protect.Client, walletID string) error {
    wallet, err := client.Wallets().GetWallet(ctx, walletID)
    if err != nil {
        return fmt.Errorf("get wallet: %w", err)
    }

    fmt.Printf("Wallet: %s\n", wallet.Name)
    fmt.Printf("  Blockchain: %s/%s\n", wallet.Blockchain, wallet.Network)
    fmt.Printf("  Addresses: %d\n", wallet.AddressesCount)

    if wallet.Balance != nil {
        fmt.Printf("  Available: %s\n", wallet.Balance.AvailableConfirmed)
        fmt.Printf("  Reserved: %s\n", wallet.Balance.ReservedConfirmed)
    }

    return nil
}
```

---

## Address Management

### Create Addresses

```go
func createAddressExample(ctx context.Context, client *protect.Client, walletID string) error {
    address, err := client.Addresses().CreateAddress(ctx, walletID, &model.CreateAddressRequest{
        Label:      "Customer Deposit",
        Comment:    "Auto-generated for customer",
        CustomerID: "USER-789",
    })
    if err != nil {
        return fmt.Errorf("create address: %w", err)
    }

    fmt.Printf("Address: %s\n", address.Address)
    fmt.Printf("ID: %s\n", address.ID)
    fmt.Printf("Status: %s\n", address.Status)
    return nil
}
```

### Get Address with Balance

```go
func getAddressExample(ctx context.Context, client *protect.Client, addressID string) error {
    address, err := client.Addresses().GetAddress(ctx, addressID)
    if err != nil {
        return fmt.Errorf("get address: %w", err)
    }

    fmt.Printf("Address: %s\n", address.Address)
    fmt.Printf("Label: %s\n", address.Label)
    fmt.Printf("Status: %s\n", address.Status)

    if address.Balance != nil {
        fmt.Printf("Available Confirmed: %s\n", address.Balance.AvailableConfirmed)
        fmt.Printf("Total Confirmed: %s\n", address.Balance.TotalConfirmed)
    }

    return nil
}
```

### List Addresses in Wallet

```go
func listAddressesExample(ctx context.Context, client *protect.Client, walletID string) error {
    addresses, pagination, err := client.Addresses().ListAddresses(ctx, walletID, &model.ListAddressesOptions{
        Limit:  100,
        Offset: 0,
    })
    if err != nil {
        return fmt.Errorf("list addresses: %w", err)
    }

    fmt.Printf("Found %d addresses (total: %d)\n", len(addresses), pagination.TotalItems)

    for _, addr := range addresses {
        fmt.Printf("  %s - %s\n", addr.Address, addr.Label)
    }

    return nil
}
```

---

## Request Approval Workflow

### Complete Approval Flow

```go
import (
    "context"
    "fmt"

    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func completeApprovalFlow(ctx context.Context, client *protect.Client, userPrivateKeyPEM string) error {
    // Step 1: Create a transfer request
    request, err := client.Requests().CreateOutgoingRequest(ctx, &model.CreateOutgoingRequest{
        FromAddressID:          "addr-123",
        ToWhitelistedAddressID: "wla-456",
        Amount:                 "1000000000000000000", // 1 ETH in wei
    })
    if err != nil {
        return fmt.Errorf("create request: %w", err)
    }

    fmt.Printf("Created request ID: %s\n", request.ID)
    fmt.Printf("Status: %s\n", request.Status)

    // Step 2: Retrieve request to get metadata (with hash verification)
    retrieved, err := client.Requests().GetRequest(ctx, request.ID)
    if err != nil {
        return fmt.Errorf("get request: %w", err)
    }

    // Step 3: Review metadata before signing
    fmt.Println("=== Review Request Metadata ===")
    fmt.Printf("Request ID: %s\n", retrieved.ID)
    fmt.Printf("Currency: %s\n", retrieved.Currency)
    if retrieved.Metadata != nil {
        fmt.Printf("Hash: %s\n", retrieved.Metadata.Hash)
        fmt.Printf("Payload: %v\n", retrieved.Metadata.Payload)
    }

    // Step 4: Load private key
    privateKey, err := crypto.DecodePrivateKeyPEM(userPrivateKeyPEM)
    if err != nil {
        return fmt.Errorf("decode private key: %w", err)
    }

    // Step 5: Approve request - SDK handles signing internally
    signedCount, err := client.Requests().ApproveRequest(ctx, retrieved, privateKey)
    if err != nil {
        return fmt.Errorf("approve request: %w", err)
    }

    fmt.Printf("Request approved successfully (%d signed)\n", signedCount)

    // Step 6: Check updated status
    updated, err := client.Requests().GetRequest(ctx, request.ID)
    if err != nil {
        return fmt.Errorf("get updated request: %w", err)
    }

    fmt.Printf("New status: %s\n", updated.Status)
    return nil
}
```

### List Requests Pending Approval

```go
func listPendingRequestsExample(ctx context.Context, client *protect.Client) error {
    requests, _, err := client.Requests().ListRequests(ctx, &model.ListRequestsOptions{
        Limit:  100,
        Status: "APPROVING",
    })
    if err != nil {
        return fmt.Errorf("list requests: %w", err)
    }

    fmt.Printf("Found %d pending requests\n", len(requests))

    for _, req := range requests {
        fmt.Printf("  %s: %s %s (needs approval from: %v)\n",
            req.ID, req.Currency, req.Type, req.NeedsApprovalFrom)
    }

    return nil
}
```

### Reject Request

```go
func rejectRequestExample(ctx context.Context, client *protect.Client, requestID string) error {
    err := client.Requests().RejectRequest(ctx, requestID, "Amount exceeds daily limit")
    if err != nil {
        return fmt.Errorf("reject request: %w", err)
    }

    fmt.Println("Request rejected")
    return nil
}
```

---

## Transaction Queries

### Query Transactions

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func queryTransactionsExample(ctx context.Context, client *protect.Client) error {
    // Query with filters
    from := time.Now().AddDate(0, 0, -7) // Last 7 days
    to := time.Now()

    txs, pagination, err := client.Transactions().ListTransactions(ctx, &model.ListTransactionsOptions{
        Limit:     100,
        From:      from,
        To:        to,
        Currency:  "ETH",
        Direction: "outgoing",
    })
    if err != nil {
        return fmt.Errorf("list transactions: %w", err)
    }

    fmt.Printf("Found %d transactions (total: %d)\n", len(txs), pagination.TotalItems)

    for _, tx := range txs {
        fmt.Printf("  %s: %s %s (%s)\n",
            tx.Hash[:16]+"...",
            tx.Amount,
            tx.Currency,
            tx.Direction,
        )
    }

    return nil
}
```

---

## Balance Queries

### Get Balances

```go
func getBalancesExample(ctx context.Context, client *protect.Client) error {
    // Get all asset balances for the tenant
    result, err := client.Balances().GetBalances(ctx, &model.GetBalancesOptions{
        Limit: 100,
    })
    if err != nil {
        return fmt.Errorf("get balances: %w", err)
    }

    fmt.Printf("Total assets: %d\n", result.Total)
    for _, assetBalance := range result.Balances {
        if assetBalance.Asset != nil && assetBalance.Balance != nil {
            fmt.Printf("  %s: Available=%s, Reserved=%s\n",
                assetBalance.Asset.Currency,
                assetBalance.Balance.AvailableConfirmed,
                assetBalance.Balance.ReservedConfirmed,
            )
        }
    }

    // Filter by currency
    ethResult, err := client.Balances().GetBalances(ctx, &model.GetBalancesOptions{
        Currency: "ETH",
    })
    if err != nil {
        return fmt.Errorf("get ETH balances: %w", err)
    }

    for _, balance := range ethResult.Balances {
        fmt.Printf("ETH Balance: %s\n", balance.Balance.TotalConfirmed)
    }

    // Cursor-based pagination
    if result.NextCursor != "" {
        nextPage, err := client.Balances().GetBalances(ctx, &model.GetBalancesOptions{
            Cursor: result.NextCursor,
        })
        if err != nil {
            return fmt.Errorf("get next page: %w", err)
        }
        fmt.Printf("Next page has %d balances\n", len(nextPage.Balances))
    }

    return nil
}
```

---

## Whitelisted Addresses

### Get Whitelisted Address

```go
func getWhitelistedAddressExample(ctx context.Context, client *protect.Client, id string) error {
    addr, err := client.WhitelistedAddresses().GetWhitelistedAddress(ctx, id)
    if err != nil {
        if protect.IsIntegrityError(err) {
            return fmt.Errorf("integrity verification failed: %w", err)
        }
        return fmt.Errorf("get whitelisted address: %w", err)
    }

    fmt.Printf("Blockchain: %s/%s\n", addr.Blockchain, addr.Network)
    fmt.Printf("Address: %s\n", addr.Address)
    fmt.Printf("Name: %s\n", addr.Name)
    fmt.Printf("Status: %s\n", addr.Status)

    if addr.Approvers != nil {
        fmt.Printf("Approvers: %d parallel groups\n", len(addr.Approvers.Parallel))
    }

    return nil
}
```

### List Whitelisted Addresses

```go
func listWhitelistedAddressesExample(ctx context.Context, client *protect.Client) error {
    // List all
    addresses, pagination, err := client.WhitelistedAddresses().ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{
        Limit:  100,
        Offset: 0,
    })
    if err != nil {
        return fmt.Errorf("list whitelisted addresses: %w", err)
    }

    fmt.Printf("Found %d whitelisted addresses\n", len(addresses))

    // Filter by blockchain
    ethAddresses, _, err := client.WhitelistedAddresses().ListWhitelistedAddresses(ctx, &model.ListWhitelistedAddressesOptions{
        Limit:      100,
        Blockchain: "ETH",
        Network:    "mainnet",
    })
    if err != nil {
        return fmt.Errorf("list ETH whitelisted addresses: %w", err)
    }

    for _, addr := range ethAddresses {
        fmt.Printf("  %s: %s (%s)\n", addr.Blockchain, addr.Address, addr.Name)
    }

    return nil
}
```

---

## Webhook Management

### Create and Manage Webhooks

```go
func webhookExample(ctx context.Context, client *protect.Client) error {
    // Create a webhook
    webhook, err := client.Webhooks().CreateWebhook(ctx, &model.CreateWebhookRequest{
        URL:    "https://example.com/webhook/transactions",
        Type:   "TRANSACTION",
        Secret: "my-webhook-secret-key",
    })
    if err != nil {
        return fmt.Errorf("create webhook: %w", err)
    }

    fmt.Printf("Created webhook: %s\n", webhook.ID)

    // List all webhooks
    webhooks, _, err := client.Webhooks().ListWebhooks(ctx, nil)
    if err != nil {
        return fmt.Errorf("list webhooks: %w", err)
    }

    for _, wh := range webhooks {
        fmt.Printf("Webhook: %s (%s) - Status: %s\n", wh.ID, wh.Type, wh.Status)
    }

    // Disable a webhook
    disabled, err := client.Webhooks().UpdateWebhookStatus(ctx, webhook.ID, "DISABLED")
    if err != nil {
        return fmt.Errorf("disable webhook: %w", err)
    }

    fmt.Printf("Webhook status: %s\n", disabled.Status)

    // Delete the webhook
    err = client.Webhooks().DeleteWebhook(ctx, webhook.ID)
    if err != nil {
        return fmt.Errorf("delete webhook: %w", err)
    }

    fmt.Println("Webhook deleted")
    return nil
}
```

---

## Pagination Patterns

### Standard Pagination Loop

```go
func paginationExample(ctx context.Context, client *protect.Client) error {
    var allWallets []*model.Wallet
    opts := &model.ListWalletsOptions{
        Limit:  50,
        Offset: 0,
    }

    pageNum := 0
    for {
        wallets, pagination, err := client.Wallets().ListWallets(ctx, opts)
        if err != nil {
            return fmt.Errorf("list wallets page %d: %w", pageNum, err)
        }

        pageNum++
        allWallets = append(allWallets, wallets...)
        fmt.Printf("Page %d: %d items (total: %d)\n", pageNum, len(wallets), pagination.TotalItems)

        if !pagination.HasMore {
            break
        }
        opts.Offset += opts.Limit
    }

    fmt.Printf("Total items retrieved: %d\n", len(allWallets))
    return nil
}
```

### Concurrent Pagination (Advanced)

```go
import (
    "context"
    "sync"

    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func concurrentPaginationExample(ctx context.Context, client *protect.Client) ([]*model.Wallet, error) {
    // First, get total count
    _, pagination, err := client.Wallets().ListWallets(ctx, &model.ListWalletsOptions{
        Limit:  1,
        Offset: 0,
    })
    if err != nil {
        return nil, err
    }

    totalItems := pagination.TotalItems
    pageSize := int64(100)
    numPages := (totalItems + pageSize - 1) / pageSize

    results := make([][]*model.Wallet, numPages)
    var wg sync.WaitGroup
    errChan := make(chan error, numPages)

    for i := int64(0); i < numPages; i++ {
        wg.Add(1)
        go func(pageIdx int64) {
            defer wg.Done()

            wallets, _, err := client.Wallets().ListWallets(ctx, &model.ListWalletsOptions{
                Limit:  pageSize,
                Offset: pageIdx * pageSize,
            })
            if err != nil {
                errChan <- err
                return
            }
            results[pageIdx] = wallets
        }(i)
    }

    wg.Wait()
    close(errChan)

    if err := <-errChan; err != nil {
        return nil, err
    }

    // Flatten results
    var allWallets []*model.Wallet
    for _, page := range results {
        allWallets = append(allWallets, page...)
    }

    return allWallets, nil
}
```

---

## Error Handling

### Basic Error Handling

```go
import (
    "errors"
    "fmt"

    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
)

func errorHandlingExample(ctx context.Context, client *protect.Client) {
    wallet, err := client.Wallets().GetWallet(ctx, "nonexistent-id")
    if err != nil {
        // Check for specific error types
        if apiErr, ok := protect.IsAPIError(err); ok {
            fmt.Printf("API Error:\n")
            fmt.Printf("  HTTP Status: %d\n", apiErr.Code)
            fmt.Printf("  Error Code: %s\n", apiErr.ErrorCode)
            fmt.Printf("  Message: %s\n", apiErr.Message)

            // Handle specific codes
            switch {
            case errors.Is(err, protect.ErrNotFound):
                fmt.Println("Wallet not found")
            case errors.Is(err, protect.ErrAuthentication):
                fmt.Println("Authentication failed - check credentials")
            case errors.Is(err, protect.ErrAuthorization):
                fmt.Println("Permission denied")
            case errors.Is(err, protect.ErrRateLimit):
                fmt.Printf("Rate limited - retry after: %v\n", apiErr.RetryAfter)
            }

            if apiErr.IsRetryable() {
                fmt.Println("This error is retryable")
            }
            return
        }

        if protect.IsIntegrityError(err) {
            fmt.Println("Integrity verification failed - potential tampering!")
            return
        }

        if protect.IsWhitelistError(err) {
            fmt.Println("Whitelist verification failed")
            return
        }

        fmt.Printf("Unknown error: %v\n", err)
        return
    }

    fmt.Printf("Wallet: %s\n", wallet.Name)
}
```

### Retry Logic

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
    "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func withRetry[T any](
    ctx context.Context,
    maxRetries int,
    operation func() (T, error),
) (T, error) {
    var result T
    var lastErr error

    for attempt := 0; attempt < maxRetries; attempt++ {
        result, lastErr = operation()
        if lastErr == nil {
            return result, nil
        }

        if apiErr, ok := protect.IsAPIError(lastErr); ok {
            if !apiErr.IsRetryable() {
                return result, lastErr
            }

            // Use suggested retry delay or exponential backoff
            delay := apiErr.SuggestedRetryDelay()
            if delay == 0 {
                delay = time.Duration(1<<attempt) * 100 * time.Millisecond
            }

            select {
            case <-ctx.Done():
                return result, ctx.Err()
            case <-time.After(delay):
                continue
            }
        }

        return result, lastErr
    }

    return result, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// Usage
func retryExample(ctx context.Context, client *protect.Client, walletID string) (*model.Wallet, error) {
    return withRetry(ctx, 3, func() (*model.Wallet, error) {
        return client.Wallets().GetWallet(ctx, walletID)
    })
}
```

---

## Context and Cancellation

### Timeout Context

```go
func timeoutExample(client *protect.Client, walletID string) error {
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    wallet, err := client.Wallets().GetWallet(ctx, walletID)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return fmt.Errorf("request timed out")
        }
        return err
    }

    fmt.Printf("Wallet: %s\n", wallet.Name)
    return nil
}
```

### Cancellation

```go
import (
    "context"
    "os"
    "os/signal"
    "syscall"
)

func cancellationExample(client *protect.Client) error {
    // Create context that cancels on SIGINT/SIGTERM
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Handle signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        fmt.Println("\nReceived shutdown signal, cancelling...")
        cancel()
    }()

    // Long-running operation
    var allWallets []*model.Wallet
    opts := &model.ListWalletsOptions{Limit: 50}

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        wallets, pagination, err := client.Wallets().ListWallets(ctx, opts)
        if err != nil {
            return err
        }

        allWallets = append(allWallets, wallets...)
        if !pagination.HasMore {
            break
        }
        opts.Offset += opts.Limit
    }

    fmt.Printf("Retrieved %d wallets\n", len(allWallets))
    return nil
}
```

---

## Related Documentation

- [SDK Overview](SDK_OVERVIEW.md) - Architecture and modules
- [Authentication](AUTHENTICATION.md) - Security and signing
- [Services Reference](SERVICES.md) - Complete API documentation
- [Whitelisted Address Verification](WHITELISTED_ADDRESS_VERIFICATION.md) - Verification details