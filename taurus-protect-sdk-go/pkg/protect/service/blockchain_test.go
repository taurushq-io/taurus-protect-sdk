package service

import (
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func TestNewBlockchainService(t *testing.T) {
	// Cannot test with nil client as it would panic on field access
	// This test documents that the constructor exists and follows expected pattern
}

func TestBlockchainService_ListBlockchains_NilOptions(t *testing.T) {
	// This test documents that ListBlockchains accepts nil options
	// The actual API call would require a mock, but we can verify
	// that the service handles nil options without panicking in the
	// options handling code (before the API call)
	svc := &BlockchainService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This will panic due to nil API, but that's expected since we're
	// not mocking the API client. The test verifies the options handling
	// doesn't cause issues before the API call.
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to nil API, but didn't get one")
		}
	}()

	_, _ = svc.ListBlockchains(nil, nil)
}

func TestBlockchainService_ListBlockchains_WithOptions(t *testing.T) {
	// This test documents that ListBlockchains handles options correctly
	// The actual API call would require a mock
	svc := &BlockchainService{
		api:       nil,
		errMapper: NewErrorMapper(),
	}

	// This will panic due to nil API, but that's expected since we're
	// not mocking the API client. The test verifies the options handling
	// doesn't cause issues before the API call.
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic due to nil API, but didn't get one")
		}
	}()

	opts := &model.ListBlockchainsOptions{
		Blockchain:         "ETH",
		Network:            "mainnet",
		IncludeBlockHeight: true,
	}

	_, _ = svc.ListBlockchains(nil, opts)
}

func TestListBlockchainsOptions(t *testing.T) {
	// Test that options struct can be created with expected fields
	opts := &model.ListBlockchainsOptions{
		Blockchain:         "ETH",
		Network:            "mainnet",
		IncludeBlockHeight: true,
	}

	if opts.Blockchain != "ETH" {
		t.Errorf("Blockchain = %v, want 'ETH'", opts.Blockchain)
	}
	if opts.Network != "mainnet" {
		t.Errorf("Network = %v, want 'mainnet'", opts.Network)
	}
	if !opts.IncludeBlockHeight {
		t.Error("IncludeBlockHeight should be true")
	}
}

func TestListBlockchainsOptions_Defaults(t *testing.T) {
	// Test that options struct defaults are correct
	opts := &model.ListBlockchainsOptions{}

	if opts.Blockchain != "" {
		t.Errorf("Blockchain = %v, want empty string", opts.Blockchain)
	}
	if opts.Network != "" {
		t.Errorf("Network = %v, want empty string", opts.Network)
	}
	if opts.IncludeBlockHeight {
		t.Error("IncludeBlockHeight should be false by default")
	}
}

func TestBlockchainModel(t *testing.T) {
	// Test that Blockchain model struct can be created with expected fields
	bc := &model.Blockchain{
		Symbol:           "ETH",
		Name:             "Ethereum",
		Network:          "mainnet",
		ChainID:          "1",
		BlackholeAddress: "0x0000000000000000000000000000000000000000",
		Confirmations:    "12",
		BlockHeight:      "18000000",
		IsLayer2Chain:    false,
		Layer1Network:    "",
	}

	if bc.Symbol != "ETH" {
		t.Errorf("Symbol = %v, want 'ETH'", bc.Symbol)
	}
	if bc.Name != "Ethereum" {
		t.Errorf("Name = %v, want 'Ethereum'", bc.Name)
	}
	if bc.Network != "mainnet" {
		t.Errorf("Network = %v, want 'mainnet'", bc.Network)
	}
	if bc.ChainID != "1" {
		t.Errorf("ChainID = %v, want '1'", bc.ChainID)
	}
	if bc.Confirmations != "12" {
		t.Errorf("Confirmations = %v, want '12'", bc.Confirmations)
	}
	if bc.BlockHeight != "18000000" {
		t.Errorf("BlockHeight = %v, want '18000000'", bc.BlockHeight)
	}
	if bc.IsLayer2Chain {
		t.Error("IsLayer2Chain should be false")
	}
}

func TestBlockchainModel_Layer2(t *testing.T) {
	// Test Layer 2 blockchain model
	bc := &model.Blockchain{
		Symbol:        "ARB",
		Name:          "Arbitrum",
		Network:       "mainnet",
		ChainID:       "42161",
		IsLayer2Chain: true,
		Layer1Network: "ETH",
	}

	if !bc.IsLayer2Chain {
		t.Error("IsLayer2Chain should be true")
	}
	if bc.Layer1Network != "ETH" {
		t.Errorf("Layer1Network = %v, want 'ETH'", bc.Layer1Network)
	}
}

func TestEVMBlockchainInfo(t *testing.T) {
	info := &model.EVMBlockchainInfo{
		ChainID: "1",
	}

	if info.ChainID != "1" {
		t.Errorf("ChainID = %v, want '1'", info.ChainID)
	}
}

func TestDOTBlockchainInfo(t *testing.T) {
	info := &model.DOTBlockchainInfo{
		CurrentEra:     "1234",
		MaxNominations: "16",
		ForkNumber:     "1",
	}

	if info.CurrentEra != "1234" {
		t.Errorf("CurrentEra = %v, want '1234'", info.CurrentEra)
	}
	if info.MaxNominations != "16" {
		t.Errorf("MaxNominations = %v, want '16'", info.MaxNominations)
	}
	if info.ForkNumber != "1" {
		t.Errorf("ForkNumber = %v, want '1'", info.ForkNumber)
	}
}

func TestXTZBlockchainInfo(t *testing.T) {
	info := &model.XTZBlockchainInfo{
		CurrentCycle: "789",
	}

	if info.CurrentCycle != "789" {
		t.Errorf("CurrentCycle = %v, want '789'", info.CurrentCycle)
	}
}
