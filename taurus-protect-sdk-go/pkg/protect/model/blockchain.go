package model

import "time"

// Blockchain represents a blockchain entity with its networks and configuration.
type Blockchain struct {
	// Symbol is the blockchain symbol (e.g., BTC, ETH, FTM, Cosmos).
	Symbol string `json:"symbol"`
	// Name is the blockchain name.
	Name string `json:"name"`
	// Network is the network or environment (e.g., mainnet, testnet).
	Network string `json:"network"`
	// ChainID is the chain identifier for the blockchain.
	ChainID string `json:"chain_id"`
	// BlackholeAddress is the destination address used for cancel requests.
	BlackholeAddress string `json:"blackhole_address,omitempty"`
	// Confirmations is the number of blocks needed for a transaction to be considered confirmed.
	Confirmations string `json:"confirmations,omitempty"`
	// BlockHeight is the current block height of the blockchain.
	BlockHeight string `json:"block_height,omitempty"`
	// IsLayer2Chain indicates if this is a Layer 2 blockchain.
	IsLayer2Chain bool `json:"is_layer2_chain"`
	// Layer1Network is the underlying layer 1 blockchain network, only relevant when IsLayer2Chain is true.
	Layer1Network string `json:"layer1_network,omitempty"`
	// BaseCurrency is the native currency of the blockchain.
	BaseCurrency *Currency `json:"base_currency,omitempty"`
	// EVMInfo contains EVM-specific blockchain information.
	EVMInfo *EVMBlockchainInfo `json:"evm_info,omitempty"`
	// DOTInfo contains Polkadot-specific blockchain information.
	DOTInfo *DOTBlockchainInfo `json:"dot_info,omitempty"`
	// XTZInfo contains Tezos-specific blockchain information.
	XTZInfo *XTZBlockchainInfo `json:"xtz_info,omitempty"`
}

// EVMBlockchainInfo contains EVM-specific blockchain information.
type EVMBlockchainInfo struct {
	// ChainID is the EVM chain identifier.
	ChainID string `json:"chain_id,omitempty"`
}

// DOTBlockchainInfo contains Polkadot-specific blockchain information.
type DOTBlockchainInfo struct {
	// CurrentEra is the current era in the Polkadot network.
	CurrentEra string `json:"current_era,omitempty"`
	// MaxNominations is the maximum number of nominations allowed.
	MaxNominations string `json:"max_nominations,omitempty"`
	// ForkNumber is set to 1 for AssetHub.
	ForkNumber string `json:"fork_number,omitempty"`
	// ForkMigratedAt is empty when DOT hasn't migrated to AssetHub.
	ForkMigratedAt *time.Time `json:"fork_migrated_at,omitempty"`
}

// XTZBlockchainInfo contains Tezos-specific blockchain information.
type XTZBlockchainInfo struct {
	// CurrentCycle is the current cycle in the Tezos network.
	CurrentCycle string `json:"current_cycle,omitempty"`
}

// ListBlockchainsOptions contains options for listing blockchains.
type ListBlockchainsOptions struct {
	// Blockchain filters by blockchain symbol (e.g., "ETH", "BTC").
	Blockchain string
	// Network filters by network or environment (e.g., "mainnet", "testnet").
	Network string
	// IncludeBlockHeight includes the current block height in the response.
	// Requires Blockchain and Network to be specified.
	IncludeBlockHeight bool
}
