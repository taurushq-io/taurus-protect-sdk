package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestBlockchainFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordBlockchainEntity
		want func(t *testing.T, got interface{})
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
			want: func(t *testing.T, got interface{}) {
				if got != nil {
					t.Errorf("expected nil, got %v", got)
				}
			},
		},
		{
			name: "empty DTO returns blockchain with zero values",
			dto:  &openapi.TgvalidatordBlockchainEntity{},
			want: func(t *testing.T, got interface{}) {
				if got == nil {
					t.Error("expected non-nil blockchain")
				}
			},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordBlockchainEntity {
				symbol := "ETH"
				name := "Ethereum"
				network := "mainnet"
				chainId := "1"
				blackholeAddress := "0x0000000000000000000000000000000000000000"
				confirmations := "12"
				blockHeight := "18000000"
				isLayer2Chain := false
				layer1Network := ""
				return &openapi.TgvalidatordBlockchainEntity{
					Symbol:           &symbol,
					Name:             &name,
					Network:          &network,
					ChainId:          &chainId,
					BlackholeAddress: &blackholeAddress,
					Confirmations:    &confirmations,
					BlockHeight:      &blockHeight,
					IsLayer2Chain:    &isLayer2Chain,
					Layer1Network:    &layer1Network,
				}
			}(),
			want: func(t *testing.T, got interface{}) {
				// Validated in main test body
			},
		},
		{
			name: "layer 2 blockchain",
			dto: func() *openapi.TgvalidatordBlockchainEntity {
				symbol := "ARB"
				name := "Arbitrum"
				network := "mainnet"
				chainId := "42161"
				isLayer2Chain := true
				layer1Network := "ETH"
				return &openapi.TgvalidatordBlockchainEntity{
					Symbol:        &symbol,
					Name:          &name,
					Network:       &network,
					ChainId:       &chainId,
					IsLayer2Chain: &isLayer2Chain,
					Layer1Network: &layer1Network,
				}
			}(),
			want: func(t *testing.T, got interface{}) {
				// Validated in main test body
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BlockchainFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("BlockchainFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("BlockchainFromDTO() returned nil for non-nil input")
			}
			// Verify specific fields
			if tt.dto.Symbol != nil && got.Symbol != *tt.dto.Symbol {
				t.Errorf("Symbol = %v, want %v", got.Symbol, *tt.dto.Symbol)
			}
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Network != nil && got.Network != *tt.dto.Network {
				t.Errorf("Network = %v, want %v", got.Network, *tt.dto.Network)
			}
			if tt.dto.ChainId != nil && got.ChainID != *tt.dto.ChainId {
				t.Errorf("ChainID = %v, want %v", got.ChainID, *tt.dto.ChainId)
			}
			if tt.dto.BlackholeAddress != nil && got.BlackholeAddress != *tt.dto.BlackholeAddress {
				t.Errorf("BlackholeAddress = %v, want %v", got.BlackholeAddress, *tt.dto.BlackholeAddress)
			}
			if tt.dto.Confirmations != nil && got.Confirmations != *tt.dto.Confirmations {
				t.Errorf("Confirmations = %v, want %v", got.Confirmations, *tt.dto.Confirmations)
			}
			if tt.dto.BlockHeight != nil && got.BlockHeight != *tt.dto.BlockHeight {
				t.Errorf("BlockHeight = %v, want %v", got.BlockHeight, *tt.dto.BlockHeight)
			}
			if tt.dto.IsLayer2Chain != nil && got.IsLayer2Chain != *tt.dto.IsLayer2Chain {
				t.Errorf("IsLayer2Chain = %v, want %v", got.IsLayer2Chain, *tt.dto.IsLayer2Chain)
			}
			if tt.dto.Layer1Network != nil && got.Layer1Network != *tt.dto.Layer1Network {
				t.Errorf("Layer1Network = %v, want %v", got.Layer1Network, *tt.dto.Layer1Network)
			}
		})
	}
}

func TestBlockchainFromDTO_WithBaseCurrency(t *testing.T) {
	currencyID := "eth-mainnet"
	currencySymbol := "ETH"
	currencyName := "ethereum"
	blockchain := "ETH"
	network := "mainnet"

	dto := &openapi.TgvalidatordBlockchainEntity{
		Symbol:  &blockchain,
		Network: &network,
		BaseCurrency: &openapi.TgvalidatordCurrency{
			Id:         &currencyID,
			Symbol:     &currencySymbol,
			Name:       &currencyName,
			Blockchain: &blockchain,
			Network:    &network,
		},
	}

	got := BlockchainFromDTO(dto)

	if got.BaseCurrency == nil {
		t.Fatal("BaseCurrency should not be nil")
	}
	if got.BaseCurrency.ID != currencyID {
		t.Errorf("BaseCurrency.ID = %v, want %v", got.BaseCurrency.ID, currencyID)
	}
	if got.BaseCurrency.Symbol != currencySymbol {
		t.Errorf("BaseCurrency.Symbol = %v, want %v", got.BaseCurrency.Symbol, currencySymbol)
	}
}

func TestBlockchainFromDTO_WithEVMInfo(t *testing.T) {
	symbol := "ETH"
	evmChainId := "1"

	dto := &openapi.TgvalidatordBlockchainEntity{
		Symbol: &symbol,
		EthInfo: &openapi.TgvalidatordEVMBlockchainInfo{
			ChainId: &evmChainId,
		},
	}

	got := BlockchainFromDTO(dto)

	if got.EVMInfo == nil {
		t.Fatal("EVMInfo should not be nil")
	}
	if got.EVMInfo.ChainID != evmChainId {
		t.Errorf("EVMInfo.ChainID = %v, want %v", got.EVMInfo.ChainID, evmChainId)
	}
}

func TestBlockchainFromDTO_WithDOTInfo(t *testing.T) {
	symbol := "DOT"
	currentEra := "1234"
	maxNominations := "16"
	forkNumber := "1"
	forkMigratedAt := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	dto := &openapi.TgvalidatordBlockchainEntity{
		Symbol: &symbol,
		DotInfo: &openapi.TgvalidatordDOTBlockchainInfo{
			CurrentEra:     &currentEra,
			MaxNominations: &maxNominations,
			ForkNumber:     &forkNumber,
			ForkMigratedAt: &forkMigratedAt,
		},
	}

	got := BlockchainFromDTO(dto)

	if got.DOTInfo == nil {
		t.Fatal("DOTInfo should not be nil")
	}
	if got.DOTInfo.CurrentEra != currentEra {
		t.Errorf("DOTInfo.CurrentEra = %v, want %v", got.DOTInfo.CurrentEra, currentEra)
	}
	if got.DOTInfo.MaxNominations != maxNominations {
		t.Errorf("DOTInfo.MaxNominations = %v, want %v", got.DOTInfo.MaxNominations, maxNominations)
	}
	if got.DOTInfo.ForkNumber != forkNumber {
		t.Errorf("DOTInfo.ForkNumber = %v, want %v", got.DOTInfo.ForkNumber, forkNumber)
	}
	if got.DOTInfo.ForkMigratedAt == nil || !got.DOTInfo.ForkMigratedAt.Equal(forkMigratedAt) {
		t.Errorf("DOTInfo.ForkMigratedAt = %v, want %v", got.DOTInfo.ForkMigratedAt, forkMigratedAt)
	}
}

func TestBlockchainFromDTO_WithXTZInfo(t *testing.T) {
	symbol := "XTZ"
	currentCycle := "789"

	dto := &openapi.TgvalidatordBlockchainEntity{
		Symbol: &symbol,
		XtzInfo: &openapi.TgvalidatordXTZBlockchainInfo{
			CurrentCycle: &currentCycle,
		},
	}

	got := BlockchainFromDTO(dto)

	if got.XTZInfo == nil {
		t.Fatal("XTZInfo should not be nil")
	}
	if got.XTZInfo.CurrentCycle != currentCycle {
		t.Errorf("XTZInfo.CurrentCycle = %v, want %v", got.XTZInfo.CurrentCycle, currentCycle)
	}
}

func TestBlockchainsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordBlockchainEntity
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordBlockchainEntity{},
			want: 0,
		},
		{
			name: "converts multiple blockchains",
			dtos: func() []openapi.TgvalidatordBlockchainEntity {
				symbol1 := "ETH"
				symbol2 := "BTC"
				return []openapi.TgvalidatordBlockchainEntity{
					{Symbol: &symbol1},
					{Symbol: &symbol2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BlockchainsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("BlockchainsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("BlockchainsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestBlockchainsFromDTO_PreservesOrder(t *testing.T) {
	symbol1 := "ETH"
	symbol2 := "BTC"
	symbol3 := "DOT"
	dtos := []openapi.TgvalidatordBlockchainEntity{
		{Symbol: &symbol1},
		{Symbol: &symbol2},
		{Symbol: &symbol3},
	}

	got := BlockchainsFromDTO(dtos)

	if len(got) != 3 {
		t.Fatalf("BlockchainsFromDTO() length = %v, want 3", len(got))
	}
	if got[0].Symbol != "ETH" {
		t.Errorf("got[0].Symbol = %v, want 'ETH'", got[0].Symbol)
	}
	if got[1].Symbol != "BTC" {
		t.Errorf("got[1].Symbol = %v, want 'BTC'", got[1].Symbol)
	}
	if got[2].Symbol != "DOT" {
		t.Errorf("got[2].Symbol = %v, want 'DOT'", got[2].Symbol)
	}
}

func TestEVMBlockchainInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordEVMBlockchainInfo
		want *string
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
			want: nil,
		},
		{
			name: "empty DTO returns info with zero values",
			dto:  &openapi.TgvalidatordEVMBlockchainInfo{},
			want: func() *string { s := ""; return &s }(),
		},
		{
			name: "complete DTO maps chainId",
			dto: func() *openapi.TgvalidatordEVMBlockchainInfo {
				chainId := "1"
				return &openapi.TgvalidatordEVMBlockchainInfo{
					ChainId: &chainId,
				}
			}(),
			want: func() *string { s := "1"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EVMBlockchainInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("EVMBlockchainInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("EVMBlockchainInfoFromDTO() returned nil for non-nil input")
			}
			if tt.want != nil && got.ChainID != *tt.want {
				t.Errorf("ChainID = %v, want %v", got.ChainID, *tt.want)
			}
		})
	}
}

func TestDOTBlockchainInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordDOTBlockchainInfo
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns info with zero values",
			dto:  &openapi.TgvalidatordDOTBlockchainInfo{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordDOTBlockchainInfo {
				currentEra := "1234"
				maxNominations := "16"
				forkNumber := "1"
				forkMigratedAt := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
				return &openapi.TgvalidatordDOTBlockchainInfo{
					CurrentEra:     &currentEra,
					MaxNominations: &maxNominations,
					ForkNumber:     &forkNumber,
					ForkMigratedAt: &forkMigratedAt,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DOTBlockchainInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("DOTBlockchainInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("DOTBlockchainInfoFromDTO() returned nil for non-nil input")
			}
			// Verify specific fields
			if tt.dto.CurrentEra != nil && got.CurrentEra != *tt.dto.CurrentEra {
				t.Errorf("CurrentEra = %v, want %v", got.CurrentEra, *tt.dto.CurrentEra)
			}
			if tt.dto.MaxNominations != nil && got.MaxNominations != *tt.dto.MaxNominations {
				t.Errorf("MaxNominations = %v, want %v", got.MaxNominations, *tt.dto.MaxNominations)
			}
			if tt.dto.ForkNumber != nil && got.ForkNumber != *tt.dto.ForkNumber {
				t.Errorf("ForkNumber = %v, want %v", got.ForkNumber, *tt.dto.ForkNumber)
			}
			if tt.dto.ForkMigratedAt != nil {
				if got.ForkMigratedAt == nil || !got.ForkMigratedAt.Equal(*tt.dto.ForkMigratedAt) {
					t.Errorf("ForkMigratedAt = %v, want %v", got.ForkMigratedAt, *tt.dto.ForkMigratedAt)
				}
			}
		})
	}
}

func TestXTZBlockchainInfoFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordXTZBlockchainInfo
		want *string
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
			want: nil,
		},
		{
			name: "empty DTO returns info with zero values",
			dto:  &openapi.TgvalidatordXTZBlockchainInfo{},
			want: func() *string { s := ""; return &s }(),
		},
		{
			name: "complete DTO maps currentCycle",
			dto: func() *openapi.TgvalidatordXTZBlockchainInfo {
				currentCycle := "789"
				return &openapi.TgvalidatordXTZBlockchainInfo{
					CurrentCycle: &currentCycle,
				}
			}(),
			want: func() *string { s := "789"; return &s }(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := XTZBlockchainInfoFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("XTZBlockchainInfoFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("XTZBlockchainInfoFromDTO() returned nil for non-nil input")
			}
			if tt.want != nil && got.CurrentCycle != *tt.want {
				t.Errorf("CurrentCycle = %v, want %v", got.CurrentCycle, *tt.want)
			}
		})
	}
}
