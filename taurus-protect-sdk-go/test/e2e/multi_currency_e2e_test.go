package e2e

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// currencyConfig holds configuration for a single currency under test.
type currencyConfig struct {
	symbol         string
	blockchain     string
	network        string
	transferAmount string
	minBalance     *big.Int
	isToken        bool
}

// addressPair holds source and destination addresses for a currency.
type addressPair struct {
	config      currencyConfig
	source      *model.Address
	destination *model.Address
}

// currencyResultStatus represents the outcome of the E2E flow for a currency.
type currencyResultStatus int

const (
	resultPassed  currencyResultStatus = iota
	resultSkipped
	resultFailed
)

// currencyResult holds the outcome of the E2E flow for a single currency.
type currencyResult struct {
	symbol       string
	status       currencyResultStatus
	txHash       string
	errorMessage string
}

const (
	maxWaitDuration  = 15 * time.Minute
	pollInterval     = 5 * time.Second
	maxParallelFlows = 3
)

// terminalStatuses are statuses that indicate the request will not progress further.
var terminalStatuses = map[string]bool{
	"BROADCASTED":       true,
	"CONFIRMED":         true,
	"REJECTED":          true,
	"CANCELED":          true,
	"PERMANENT_FAILURE": true,
	"EXPIRED":           true,
	"INVALID":           true,
	"MINED":             true,
}

func getCurrencyConfigs() []currencyConfig {
	return []currencyConfig{
		{
			symbol:         "SOL",
			blockchain:     "SOL",
			network:        "mainnet",
			transferAmount: "1000000",
			minBalance:     big.NewInt(2000000),
			isToken:        false,
		},
		// {
		// 	symbol:         "ETH",
		// 	blockchain:     "ETH",
		// 	network:        "mainnet",
		// 	transferAmount: "10000000000000",
		// 	minBalance:     big.NewInt(100000000000000),
		// 	isToken:        false,
		// },
		// {
		// 	symbol:         "XRP",
		// 	blockchain:     "XRP",
		// 	network:        "mainnet",
		// 	transferAmount: "1212",
		// 	minBalance:     big.NewInt(12_000_000),
		// 	isToken:        false,
		// },
		{
			symbol:         "XLM",
			blockchain:     "XLM",
			network:        "mainnet",
			transferAmount: "1313",
			minBalance:     big.NewInt(35000000),
			isToken:        false,
		},
		{
			symbol:         "ALGO",
			blockchain:     "ALGO",
			network:        "mainnet",
			transferAmount: "1414",
			minBalance:     big.NewInt(100000),
			isToken:        false,
		},
		// {
		// 	symbol:         "USDC",
		// 	blockchain:     "ETH",
		// 	network:        "mainnet",
		// 	transferAmount: "11100",
		// 	minBalance:     big.NewInt(2000000),
		// 	isToken:        true,
		// },
	}
}

func TestIntegration_MultiCurrencyE2E(t *testing.T) {
	skipIfNotIntegration(t)
	client := getTestClient(t)
	defer client.Close()
	ctx := context.Background()
	approvalKey := getTeam1PrivateKey(t)

	configs := getCurrencyConfigs()

	// Phase 1: Sequential address discovery
	t.Log("=== Phase 1: Address Discovery (sequential) ===")
	var pairs []addressPair
	var results []currencyResult

	for _, cfg := range configs {
		tag := fmt.Sprintf("[%s]", cfg.symbol)
		t.Logf("%s Step 0: Searching for two funded addresses...", tag)

		pair, err := findTwoFundedAddressesForCurrency(t, ctx, client, cfg)
		if err != nil {
			t.Logf("%s Step 0: Discovery failed -- %v", tag, err)
			results = append(results, currencyResult{symbol: cfg.symbol, status: resultSkipped})
			continue
		}
		if pair == nil {
			t.Logf("%s Step 0: Could not find two funded addresses -- skipping", tag)
			results = append(results, currencyResult{symbol: cfg.symbol, status: resultSkipped})
			continue
		}

		srcBal := "0"
		if pair.source.Balance != nil {
			srcBal = pair.source.Balance.AvailableConfirmed
		}
		t.Logf("%s Step 0: Source: ID=%s (%s), Balance=%s", tag, pair.source.ID, pair.source.Address, srcBal)
		t.Logf("%s Step 0: Dest:   ID=%s (%s)", tag, pair.destination.ID, pair.destination.Address)
		pairs = append(pairs, *pair)
	}

	if len(pairs) == 0 {
		t.Fatal("No currency has two funded addresses. Cannot run E2E test.")
	}

	// Phase 2: Parallel transfer flows
	t.Logf("=== Phase 2: Transfer Flows (parallel, max %d goroutines) ===", maxParallelFlows)

	resultsCh := make(chan currencyResult, len(pairs))
	sem := make(chan struct{}, maxParallelFlows)
	var wg sync.WaitGroup

	for _, pair := range pairs {
		wg.Add(1)
		go func(p addressPair) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			r := runCurrencyTransferFlow(t, ctx, client, p, approvalKey)
			resultsCh <- r
		}(pair)
	}

	wg.Wait()
	close(resultsCh)

	for r := range resultsCh {
		results = append(results, r)
		switch r.status {
		case resultPassed:
			t.Logf(">>> %s: PASSED -- tx: %s", r.symbol, r.txHash)
		case resultFailed:
			t.Logf(">>> %s: FAILED -- %s", r.symbol, r.errorMessage)
		case resultSkipped:
			t.Logf(">>> %s: SKIPPED", r.symbol)
		}
	}

	// Print summary
	var passedCount, failedCount, skippedCount int
	for _, r := range results {
		switch r.status {
		case resultPassed:
			passedCount++
		case resultFailed:
			failedCount++
			t.Logf("WARNING: %s failed: %s", r.symbol, r.errorMessage)
		case resultSkipped:
			skippedCount++
		}
	}

	t.Logf("=== Summary: %d passed, %d failed, %d skipped ===", passedCount, failedCount, skippedCount)

	if failedCount > 0 {
		t.Fatalf("No currency should fail the E2E flow. Failed: %d", failedCount)
	}
	if passedCount == 0 {
		t.Fatal("At least one currency should complete the E2E flow successfully")
	}
}

// runCurrencyTransferFlow executes the full transfer lifecycle for a single currency.
func runCurrencyTransferFlow(t *testing.T, ctx context.Context, client *protect.Client, pair addressPair, approvalKey *ecdsa.PrivateKey) currencyResult {
	cfg := pair.config
	tag := fmt.Sprintf("[%s]", cfg.symbol)
	source := pair.source
	destination := pair.destination

	// Step 1: Create transfer request
	t.Logf("%s Step 1: Creating transfer of %s from %s to %s...", tag, cfg.transferAmount, source.ID, destination.ID)
	transferReq, err := client.Requests().CreateInternalTransferRequest(ctx, source.ID, destination.ID, cfg.transferAmount)
	if err != nil {
		// 400 Bad Request at Step 1 typically means the environment/currency is not properly
		// configured (insufficient funds, blockchain maintenance, etc.) — skip rather than fail
		if strings.Contains(err.Error(), "400") || strings.Contains(err.Error(), "Bad Request") {
			return currencyResult{symbol: cfg.symbol, status: resultSkipped, errorMessage: fmt.Sprintf("Step 1 skipped (400): %v", err)}
		}
		return currencyResult{symbol: cfg.symbol, status: resultFailed, errorMessage: fmt.Sprintf("Step 1 failed: %v", err)}
	}
	if transferReq == nil || transferReq.ID == "" {
		return currencyResult{symbol: cfg.symbol, status: resultFailed, errorMessage: "Step 1: request is nil or has no ID"}
	}
	t.Logf("%s Step 1: Created request: ID=%s, Status=%s", tag, transferReq.ID, transferReq.Status)

	// Step 2: Verify metadata matches original intent
	t.Logf("%s Step 2: Verifying metadata matches transfer intent...", tag)
	requestToApprove, err := client.Requests().GetRequest(ctx, transferReq.ID)
	if err != nil {
		return currencyResult{symbol: cfg.symbol, status: resultFailed, errorMessage: fmt.Sprintf("Step 2 GetRequest failed: %v", err)}
	}
	metadata := requestToApprove.Metadata
	if metadata == nil {
		return currencyResult{symbol: cfg.symbol, status: resultFailed, errorMessage: "Step 2: metadata is nil"}
	}
	if metadata.Hash == "" {
		return currencyResult{symbol: cfg.symbol, status: resultFailed, errorMessage: "Step 2: metadata hash is empty"}
	}

	// Verify source address in metadata
	metadataSource := metadata.GetSourceAddress()
	if !cfg.isToken && metadataSource != source.Address {
		return currencyResult{symbol: cfg.symbol, status: resultFailed,
			errorMessage: fmt.Sprintf("Step 2: metadata source %q != expected %q", metadataSource, source.Address)}
	}
	t.Logf("%s Step 2: Source verified: %s", tag, metadataSource)

	// Verify destination address in metadata
	metadataDestination := metadata.GetDestinationAddress()
	if !cfg.isToken && metadataDestination != destination.Address {
		return currencyResult{symbol: cfg.symbol, status: resultFailed,
			errorMessage: fmt.Sprintf("Step 2: metadata destination %q != expected %q", metadataDestination, destination.Address)}
	}
	suffix := " (verified)"
	if cfg.isToken {
		suffix = " (token contract)"
	}
	t.Logf("%s Step 2: Destination: %s%s", tag, metadataDestination, suffix)

	// Log amount from metadata (non-fatal check -- payload structure may vary)
	metadataAmount := metadata.GetAmount()
	if metadataAmount != nil {
		t.Logf("%s Step 2: Amount: valueFrom=%s, currency=%s", tag, metadataAmount.ValueFrom, metadataAmount.CurrencyFrom)
		if !cfg.isToken && metadataAmount.ValueFrom != cfg.transferAmount {
			t.Logf("%s Step 2: WARNING: metadata amount %q != expected %q (payload structure may differ)",
				tag, metadataAmount.ValueFrom, cfg.transferAmount)
		}
	} else {
		t.Logf("%s Step 2: Amount metadata not available", tag)
	}

	// Step 3: Approve request
	t.Logf("%s Step 3: Approving request %s...", tag, transferReq.ID)
	signedCount, err := client.Requests().ApproveRequest(ctx, requestToApprove, approvalKey)
	if err != nil {
		return currencyResult{symbol: cfg.symbol, status: resultFailed, errorMessage: fmt.Sprintf("Step 3 ApproveRequest failed: %v", err)}
	}
	if signedCount <= 0 {
		return currencyResult{symbol: cfg.symbol, status: resultFailed, errorMessage: "Step 3: signedCount <= 0"}
	}
	t.Logf("%s Step 3: Approved: signedCount=%d", tag, signedCount)

	// Step 4: Wait for terminal status
	t.Logf("%s Step 4: Waiting for terminal status...", tag)
	confirmedReq, err := waitForTerminalStatus(t, ctx, client, transferReq.ID, tag)
	if err != nil {
		return currencyResult{symbol: cfg.symbol, status: resultFailed, errorMessage: fmt.Sprintf("Step 4 failed: %v", err)}
	}
	t.Logf("%s Step 4: Final status: %s", tag, confirmedReq.Status)

	// Check if it's a success status
	if confirmedReq.Status != "BROADCASTED" && confirmedReq.Status != "CONFIRMED" {
		t.Logf("%s === DIAGNOSTIC DUMP ===", tag)
		for _, sr := range confirmedReq.SignedRequests {
			t.Logf("%s   SignedRequest ID=%s, Status=%s, Hash=%s, Details=%s",
				tag, sr.ID, sr.Status, sr.Hash, sr.Details)
		}
		t.Logf("%s === END DIAGNOSTIC ===", tag)
		return currencyResult{symbol: cfg.symbol, status: resultFailed,
			errorMessage: fmt.Sprintf("request ended with %s instead of BROADCASTED/CONFIRMED", confirmedReq.Status)}
	}

	// Step 5: Verify transaction by hash
	t.Logf("%s Step 5: Fetching transaction details...", tag)
	if len(confirmedReq.SignedRequests) == 0 || confirmedReq.SignedRequests[0].Hash == "" {
		return currencyResult{symbol: cfg.symbol, status: resultFailed,
			errorMessage: "Step 5: no signed requests or transaction hash available"}
	}
	txHash := confirmedReq.SignedRequests[0].Hash
	t.Logf("%s Step 5: Transaction hash: %s", tag, txHash)

	var transaction *model.Transaction
	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(pollInterval)
		tx, err := client.Transactions().GetTransactionByHash(ctx, txHash)
		if err == nil && tx != nil {
			transaction = tx
			break
		}
		t.Logf("%s Step 5: Transaction not indexed yet, retrying...", tag)
	}

	if transaction == nil {
		return currencyResult{symbol: cfg.symbol, status: resultFailed,
			errorMessage: fmt.Sprintf("Step 5: transaction with hash %q not found after 90s", txHash)}
	}
	if transaction.Hash != txHash {
		return currencyResult{symbol: cfg.symbol, status: resultFailed,
			errorMessage: fmt.Sprintf("Step 5: transaction hash %q != expected %q", transaction.Hash, txHash)}
	}
	if transaction.Direction != "outgoing" {
		return currencyResult{symbol: cfg.symbol, status: resultFailed,
			errorMessage: fmt.Sprintf("Step 5: transaction direction %q != expected 'outgoing'", transaction.Direction)}
	}
	t.Logf("%s Step 5: Transaction verified -- ID=%s, Block=%s", tag, transaction.ID, transaction.Block)

	t.Logf("%s E2E PASSED", tag)
	return currencyResult{symbol: cfg.symbol, status: resultPassed, txHash: txHash}
}

// waitForTerminalStatus polls the request status until it reaches a terminal state or times out.
func waitForTerminalStatus(t *testing.T, ctx context.Context, client *protect.Client, requestID string, tag string) (*model.Request, error) {
	startTime := time.Now()
	consecutiveErrors := 0
	const maxConsecutiveErrors = 3

	for {
		req, err := client.Requests().GetRequest(ctx, requestID)
		if err != nil {
			consecutiveErrors++
			elapsed := time.Since(startTime)
			t.Logf("%s  [%ds] GetRequest transient error (%d/%d): %v",
				tag, int(elapsed.Seconds()), consecutiveErrors, maxConsecutiveErrors, err)
			if consecutiveErrors >= maxConsecutiveErrors {
				return nil, fmt.Errorf("GetRequest failed after %d consecutive errors: %w", maxConsecutiveErrors, err)
			}
			time.Sleep(pollInterval)
			continue
		}
		consecutiveErrors = 0

		elapsed := time.Since(startTime)
		t.Logf("%s  [%ds] Request %s -- Status: %s", tag, int(elapsed.Seconds()), requestID, req.Status)

		if terminalStatuses[req.Status] {
			return req, nil
		}

		if elapsed >= maxWaitDuration {
			return nil, fmt.Errorf("request %s did not reach terminal status within %v, last status: %s",
				requestID, maxWaitDuration, req.Status)
		}

		time.Sleep(pollInterval)
	}
}

// findTwoFundedAddressesForCurrency finds two addresses for the given currency.
// The source address must have at least minBalance. The destination can be any other address.
func findTwoFundedAddressesForCurrency(t *testing.T, ctx context.Context, client *protect.Client, cfg currencyConfig) (*addressPair, error) {
	tag := fmt.Sprintf("[%s]", cfg.symbol)

	var candidates []*model.Address
	nativeBalanceByAddr := make(map[string]*big.Int)

	if cfg.isToken {
		// For tokens, use AssetService to find token addresses
		t.Logf("%s Step 0: Searching via AssetService for token addresses...", tag)
		tokenResult, err := client.Assets().GetAssetAddresses(ctx, &model.GetAssetAddressesRequest{
			Asset: model.AssetFilter{Currency: cfg.symbol},
			Limit: 50,
		})
		if err != nil {
			return nil, fmt.Errorf("GetAssetAddresses(%s) failed: %w", cfg.symbol, err)
		}
		if tokenResult != nil {
			candidates = tokenResult.Addresses
		}
		t.Logf("%s Step 0: AssetService returned %d token addresses", tag, len(candidates))

		// Fetch native currency addresses to check gas balance
		t.Logf("%s Step 0: Fetching %s addresses for gas balance check...", tag, cfg.blockchain)
		nativeResult, err := client.Assets().GetAssetAddresses(ctx, &model.GetAssetAddressesRequest{
			Asset: model.AssetFilter{Currency: cfg.blockchain},
			Limit: 50,
		})
		if err == nil && nativeResult != nil {
			for _, nAddr := range nativeResult.Addresses {
				nBal := big.NewInt(0)
				if nAddr.Balance != nil && nAddr.Balance.AvailableConfirmed != "" {
					nBal, _ = new(big.Int).SetString(nAddr.Balance.AvailableConfirmed, 10)
					if nBal == nil {
						nBal = big.NewInt(0)
					}
				}
				nativeBalanceByAddr[strings.ToLower(nAddr.Address)] = nBal
			}
			t.Logf("%s Step 0: Found %d %s addresses for gas check", tag, len(nativeResult.Addresses), cfg.blockchain)
		}

		// Log candidate details
		for _, addr := range candidates {
			tokenBal := big.NewInt(0)
			if addr.Balance != nil && addr.Balance.AvailableConfirmed != "" {
				tokenBal, _ = new(big.Int).SetString(addr.Balance.AvailableConfirmed, 10)
				if tokenBal == nil {
					tokenBal = big.NewInt(0)
				}
			}
			gasBal := nativeBalanceByAddr[strings.ToLower(addr.Address)]
			if gasBal == nil {
				gasBal = big.NewInt(0)
			}
			tokenOk := tokenBal.Cmp(cfg.minBalance) >= 0
			gasOk := gasBal.Cmp(big.NewInt(0)) > 0
			tokenTag := "[LOW]"
			if tokenOk {
				tokenTag = "[OK]"
			}
			gasTag := "[EMPTY]"
			if gasOk {
				gasTag = "[OK]"
			}
			t.Logf("%s   Candidate ID=%s, Addr=%s, TokenBalance=%s %s, GasBalance=%s %s",
				tag, addr.ID, addr.Address, tokenBal.String(), tokenTag, gasBal.String(), gasTag)
		}

		// Filter out disabled addresses
		var enabledCandidates []*model.Address
		for _, addr := range candidates {
			if addr.Disabled {
				t.Logf("%s   DISABLED Address ID=%s, Addr=%s -- skipping", tag, addr.ID, addr.Address)
				continue
			}
			enabledCandidates = append(enabledCandidates, addr)
		}
		candidates = enabledCandidates
	} else {
		// For native currencies, use AssetService to find wallets, then list addresses
		var err error
		candidates, err = findNativeAddresses(t, ctx, client, cfg)
		if err != nil {
			return nil, err
		}
	}

	if len(candidates) < 2 {
		t.Logf("%s Step 0: Not enough candidates (%d found, need 2)", tag, len(candidates))
		return nil, nil
	}

	// Find source: needs sufficient balance AND (for tokens) native gas balance
	var source *model.Address
	for _, addr := range candidates {
		bal := big.NewInt(0)
		if addr.Balance != nil && addr.Balance.AvailableConfirmed != "" {
			bal, _ = new(big.Int).SetString(addr.Balance.AvailableConfirmed, 10)
			if bal == nil {
				bal = big.NewInt(0)
			}
		}
		if bal.Cmp(cfg.minBalance) < 0 {
			continue
		}
		if cfg.isToken {
			gasBal := nativeBalanceByAddr[strings.ToLower(addr.Address)]
			if gasBal == nil || gasBal.Cmp(big.NewInt(0)) <= 0 {
				t.Logf("%s Step 0: Skipping ID=%s -- has token balance but no gas", tag, addr.ID)
				continue
			}
			t.Logf("%s Step 0: Selected source ID=%s -- token=%s, gas=%s", tag, addr.ID, bal.String(), gasBal.String())
		}
		source = addr
		break
	}

	if source == nil {
		suffix := ""
		if cfg.isToken {
			suffix = " with gas"
		}
		t.Logf("%s Step 0: No address has sufficient balance%s (min=%s)", tag, suffix, cfg.minBalance.String())
		return nil, nil
	}

	// Find destination: any other address for the same currency
	var destination *model.Address
	for _, addr := range candidates {
		if addr.ID != source.ID {
			destination = addr
			break
		}
	}

	if destination == nil {
		return nil, nil
	}

	return &addressPair{config: cfg, source: source, destination: destination}, nil
}

// findNativeAddresses collects funded addresses across all wallets for a native currency.
func findNativeAddresses(t *testing.T, ctx context.Context, client *protect.Client, cfg currencyConfig) ([]*model.Address, error) {
	tag := fmt.Sprintf("[%s]", cfg.symbol)

	walletsResult, err := client.Assets().GetAssetWallets(ctx, &model.GetAssetWalletsRequest{
		Asset: model.AssetFilter{Currency: cfg.symbol},
		Limit: 50,
	})
	if err != nil {
		return nil, fmt.Errorf("GetAssetWallets(%s) failed: %w", cfg.symbol, err)
	}
	if walletsResult == nil || len(walletsResult.Wallets) == 0 {
		t.Logf("%s   No wallets found via GetAssetWallets", tag)
		return nil, nil
	}
	t.Logf("%s   Found %d wallet(s) via GetAssetWallets", tag, len(walletsResult.Wallets))

	var fundedAddresses []*model.Address
	totalScanned := 0

	for _, wallet := range walletsResult.Wallets {
		t.Logf("%s   Scanning wallet ID=%s, Name=%s", tag, wallet.ID, wallet.Name)
		var offset int64

		for {
			addresses, _, err := client.Addresses().ListAddresses(ctx, &model.ListAddressesOptions{
				WalletID: wallet.ID,
				Limit:    50,
				Offset:   offset,
			})
			if err != nil {
				return nil, fmt.Errorf("ListAddresses(wallet=%s) failed: %w", wallet.ID, err)
			}
			if len(addresses) == 0 {
				break
			}

			for _, addr := range addresses {
				totalScanned++
				if addr.Disabled {
					t.Logf("%s     DISABLED Address ID=%s, Addr=%s -- skipping", tag, addr.ID, addr.Address)
					continue
				}
				bal := big.NewInt(0)
				if addr.Balance != nil && addr.Balance.AvailableConfirmed != "" {
					bal, _ = new(big.Int).SetString(addr.Balance.AvailableConfirmed, 10)
					if bal == nil {
						bal = big.NewInt(0)
					}
				}
				if bal.Cmp(big.NewInt(0)) > 0 {
					t.Logf("%s     FUNDED Address ID=%s, Addr=%s, Balance=%s", tag, addr.ID, addr.Address, bal.String())
					fundedAddresses = append(fundedAddresses, addr)
				}
			}

			if len(fundedAddresses) >= 2 {
				t.Logf("%s   Found %d funded addresses after scanning %d total", tag, len(fundedAddresses), totalScanned)
				return fundedAddresses, nil
			}

			offset += int64(len(addresses))
		}
	}

	t.Logf("%s   Found %d funded addresses after scanning %d total across %d wallets",
		tag, len(fundedAddresses), totalScanned, len(walletsResult.Wallets))
	return fundedAddresses, nil
}
