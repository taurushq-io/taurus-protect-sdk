package protect

import (
	"context"
	"crypto/ecdsa"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/cache"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/service"
)

// Client is the main entry point for the Taurus-PROTECT SDK.
// It provides access to all API services and handles authentication.
//
// Use NewClient to create a new instance with the functional options pattern:
//
//	client, err := protect.NewClient(
//	    "https://api.taurus.example.com",
//	    protect.WithCredentials(apiKey, apiSecret),
//	    protect.WithSuperAdminKeysPEM(pemKeys),
//	    protect.WithMinValidSignatures(2),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
// Client implements io.Closer and should be closed when no longer needed
// to securely wipe credentials from memory.
type Client struct {
	baseURL            string
	httpClient         *http.Client
	auth               *crypto.TPV1Auth
	superAdminKeys     []*ecdsa.PublicKey
	minValidSignatures int
	rulesCache         *cache.RulesContainerCache

	// OpenAPI client for service creation
	apiClient *openapi.APIClient

	// Service instances (lazily initialized)
	mu                   sync.RWMutex
	wallets              *service.WalletService
	addresses            *service.AddressService
	requests             *service.RequestService
	transactions         *service.TransactionService
	governanceRules      *service.GovernanceRuleService
	balances             *service.BalanceService
	currencies           *service.CurrencyService
	whitelistedAddresses *service.WhitelistedAddressService
	whitelistedAssets    *service.WhitelistedAssetService
	fees                 *service.FeeService
	audits               *service.AuditService
	changes              *service.ChangeService
	prices               *service.PriceService
	airGap               *service.AirGapService
	staking              *service.StakingService
	whitelistedContracts *service.WhitelistedContractService
	businessRules        *service.BusinessRuleService
	reservations         *service.ReservationService
	users                *service.UserService
	groups               *service.GroupService
	visibilityGroups     *service.VisibilityGroupService
	config               *service.ConfigService
	webhooks             *service.WebhookService
	webhookCalls         *service.WebhookCallService
	tags                 *service.TagService
	assets               *service.AssetService
	actions              *service.ActionService
	blockchains          *service.BlockchainService
	exchanges            *service.ExchangeService
	fiat                 *service.FiatService
	feePayers            *service.FeePayerService
	health               *service.HealthService
	jobs                 *service.JobService
	scores               *service.ScoreService
	statistics           *service.StatisticsService
	tokenMetadata        *service.TokenMetadataService
	userDevices             *service.UserDeviceService
	multiFactorSignature    *service.MultiFactorSignatureService
	// Taurus Network namespace client
	taurusNetwork *TaurusNetworkClient
}

// NewClient creates a new Taurus-PROTECT API client.
//
// The host parameter should be the base URL of the Taurus-PROTECT API
// (e.g., "https://api.taurus.example.com").
//
// At minimum, WithCredentials, WithSuperAdminKeysPEM (or WithSuperAdminKeys),
// and WithMinValidSignatures must be provided:
//
//	client, err := protect.NewClient(
//	    "https://api.taurus.example.com",
//	    protect.WithCredentials(apiKey, apiSecret),
//	    protect.WithSuperAdminKeysPEM(pemKeys),
//	    protect.WithMinValidSignatures(2),
//	)
func NewClient(host string, opts ...Option) (*Client, error) {
	// Set defaults
	config := &clientConfig{
		host:          strings.TrimSuffix(host, "/"),
		rulesCacheTTL: DefaultRulesCacheTTL,
		httpTimeout:   DefaultHTTPTimeout,
	}

	// Apply options
	for _, opt := range opts {
		if err := opt(config); err != nil {
			return nil, err
		}
	}

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	// Create TPV1 authentication
	auth, err := crypto.NewTPV1Auth(config.apiKey, config.apiSecret)
	if err != nil {
		return nil, err
	}

	// Create HTTP client with TPV1 transport
	baseClient := config.httpClient
	if baseClient == nil {
		baseClient = &http.Client{
			Timeout: config.httpTimeout,
		}
	}
	httpClient := newHTTPClient(auth, baseClient)

	// Create OpenAPI client configuration
	apiConfig := openapi.NewConfiguration()
	apiConfig.Servers = []openapi.ServerConfiguration{
		{URL: config.host},
	}
	apiConfig.HTTPClient = httpClient
	apiClient := openapi.NewAPIClient(apiConfig)

	// Create the client
	client := &Client{
		baseURL:            config.host,
		httpClient:         httpClient,
		auth:               auth,
		superAdminKeys:     config.superAdminKeys,
		minValidSignatures: config.minValidSignatures,
		apiClient:          apiClient,
	}

	// Initialize rules cache with a fetcher that uses the GovernanceRuleService
	// The fetcher closure captures 'client' to lazily access the service
	client.rulesCache = cache.NewRulesContainerCache(config.rulesCacheTTL, func(ctx context.Context) (*model.DecodedRulesContainer, error) {
		govService := client.GovernanceRules()
		rules, err := govService.GetRules(ctx)
		if err != nil {
			return nil, err
		}
		return govService.GetDecodedRulesContainer(rules)
	})

	return client, nil
}

// Close releases resources and securely wipes credentials from memory.
// It is safe to call Close multiple times.
func (c *Client) Close() error {
	if c.auth != nil {
		c.auth.Close()
		c.auth = nil // Invalidate auth reference to prevent use after close
	}
	return nil
}

// BaseURL returns the base URL of the API.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// HTTPClient returns the underlying HTTP client.
// This can be used for advanced use cases, but most users should use
// the service methods instead.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}

// SuperAdminKeys returns the configured SuperAdmin public keys.
func (c *Client) SuperAdminKeys() []*ecdsa.PublicKey {
	return c.superAdminKeys
}

// MinValidSignatures returns the minimum number of valid signatures required.
func (c *Client) MinValidSignatures() int {
	return c.minValidSignatures
}

// RulesCache returns the rules container cache used for address signature verification.
// This is always non-nil as address signature verification is mandatory.
func (c *Client) RulesCache() *cache.RulesContainerCache {
	return c.rulesCache
}

// Wallets returns the wallet service for managing cryptocurrency wallets.
func (c *Client) Wallets() *service.WalletService {
	c.mu.RLock()
	if c.wallets != nil {
		c.mu.RUnlock()
		return c.wallets
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.wallets == nil {
		c.wallets = service.NewWalletService(c.apiClient)
	}
	return c.wallets
}

// Addresses returns the address service for managing blockchain addresses.
// Address signature verification is mandatory on all get and list operations.
func (c *Client) Addresses() *service.AddressService {
	c.mu.RLock()
	if c.addresses != nil {
		c.mu.RUnlock()
		return c.addresses
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.addresses == nil {
		c.addresses = service.NewAddressService(c.apiClient, c.rulesCache)
	}
	return c.addresses
}

// Requests returns the request service for managing transaction requests.
func (c *Client) Requests() *service.RequestService {
	c.mu.RLock()
	if c.requests != nil {
		c.mu.RUnlock()
		return c.requests
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.requests == nil {
		c.requests = service.NewRequestService(c.apiClient)
	}
	return c.requests
}

// Transactions returns the transaction service for querying transactions.
func (c *Client) Transactions() *service.TransactionService {
	c.mu.RLock()
	if c.transactions != nil {
		c.mu.RUnlock()
		return c.transactions
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.transactions == nil {
		c.transactions = service.NewTransactionService(c.apiClient)
	}
	return c.transactions
}

// GovernanceRules returns the governance rule service for managing governance rules.
func (c *Client) GovernanceRules() *service.GovernanceRuleService {
	c.mu.RLock()
	if c.governanceRules != nil {
		c.mu.RUnlock()
		return c.governanceRules
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.governanceRules == nil {
		c.governanceRules = service.NewGovernanceRuleService(c.apiClient)
	}
	return c.governanceRules
}

// Balances returns the balance service for querying asset balances.
func (c *Client) Balances() *service.BalanceService {
	c.mu.RLock()
	if c.balances != nil {
		c.mu.RUnlock()
		return c.balances
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.balances == nil {
		c.balances = service.NewBalanceService(c.apiClient)
	}
	return c.balances
}

// Currencies returns the currency service for querying available currencies.
func (c *Client) Currencies() *service.CurrencyService {
	c.mu.RLock()
	if c.currencies != nil {
		c.mu.RUnlock()
		return c.currencies
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.currencies == nil {
		c.currencies = service.NewCurrencyService(c.apiClient)
	}
	return c.currencies
}

// WhitelistedAddresses returns the whitelisted address service for managing external addresses.
// The service always verifies the integrity of all retrieved addresses using the
// 6-step cryptographic verification flow with the configured SuperAdmin keys.
func (c *Client) WhitelistedAddresses() *service.WhitelistedAddressService {
	c.mu.RLock()
	if c.whitelistedAddresses != nil {
		c.mu.RUnlock()
		return c.whitelistedAddresses
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.whitelistedAddresses == nil {
		c.whitelistedAddresses = service.NewWhitelistedAddressServiceWithVerification(
			c.apiClient,
			&service.WhitelistedAddressServiceConfig{
				SuperAdminKeys:     c.superAdminKeys,
				MinValidSignatures: c.minValidSignatures,
			},
		)
	}
	return c.whitelistedAddresses
}

// WhitelistedAssets returns the whitelisted asset service for managing whitelisted contracts/tokens.
// The service always verifies the integrity of all retrieved assets using the
// cryptographic verification flow with the configured SuperAdmin keys.
func (c *Client) WhitelistedAssets() *service.WhitelistedAssetService {
	c.mu.RLock()
	if c.whitelistedAssets != nil {
		c.mu.RUnlock()
		return c.whitelistedAssets
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.whitelistedAssets == nil {
		c.whitelistedAssets = service.NewWhitelistedAssetServiceWithVerification(
			c.apiClient,
			&service.WhitelistedAssetServiceConfig{
				SuperAdminKeys:     c.superAdminKeys,
				MinValidSignatures: c.minValidSignatures,
			},
		)
	}
	return c.whitelistedAssets
}

// Fees returns the fee service for querying blockchain fee estimates.
func (c *Client) Fees() *service.FeeService {
	c.mu.RLock()
	if c.fees != nil {
		c.mu.RUnlock()
		return c.fees
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.fees == nil {
		c.fees = service.NewFeeService(c.apiClient)
	}
	return c.fees
}

// Audits returns the audit service for querying audit trails.
func (c *Client) Audits() *service.AuditService {
	c.mu.RLock()
	if c.audits != nil {
		c.mu.RUnlock()
		return c.audits
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.audits == nil {
		c.audits = service.NewAuditService(c.apiClient)
	}
	return c.audits
}

// Changes returns the change service for managing configuration changes.
func (c *Client) Changes() *service.ChangeService {
	c.mu.RLock()
	if c.changes != nil {
		c.mu.RUnlock()
		return c.changes
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.changes == nil {
		c.changes = service.NewChangeService(c.apiClient)
	}
	return c.changes
}

// Prices returns the price service for querying currency prices and conversions.
func (c *Client) Prices() *service.PriceService {
	c.mu.RLock()
	if c.prices != nil {
		c.mu.RUnlock()
		return c.prices
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.prices == nil {
		c.prices = service.NewPriceService(c.apiClient)
	}
	return c.prices
}

// AirGap returns the air-gap service for cold HSM integration operations.
func (c *Client) AirGap() *service.AirGapService {
	c.mu.RLock()
	if c.airGap != nil {
		c.mu.RUnlock()
		return c.airGap
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.airGap == nil {
		c.airGap = service.NewAirGapService(c.apiClient)
	}
	return c.airGap
}

// Staking returns the staking service for blockchain staking operations.
func (c *Client) Staking() *service.StakingService {
	c.mu.RLock()
	if c.staking != nil {
		c.mu.RUnlock()
		return c.staking
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.staking == nil {
		c.staking = service.NewStakingService(c.apiClient)
	}
	return c.staking
}

// WhitelistedContracts returns the whitelisted contract service for smart contract whitelisting.
func (c *Client) WhitelistedContracts() *service.WhitelistedContractService {
	c.mu.RLock()
	if c.whitelistedContracts != nil {
		c.mu.RUnlock()
		return c.whitelistedContracts
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.whitelistedContracts == nil {
		c.whitelistedContracts = service.NewWhitelistedContractService(c.apiClient)
	}
	return c.whitelistedContracts
}

// BusinessRules returns the business rule service for querying business rules.
func (c *Client) BusinessRules() *service.BusinessRuleService {
	c.mu.RLock()
	if c.businessRules != nil {
		c.mu.RUnlock()
		return c.businessRules
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.businessRules == nil {
		c.businessRules = service.NewBusinessRuleService(c.apiClient)
	}
	return c.businessRules
}

// Reservations returns the reservation service for managing UTXO reservations.
func (c *Client) Reservations() *service.ReservationService {
	c.mu.RLock()
	if c.reservations != nil {
		c.mu.RUnlock()
		return c.reservations
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.reservations == nil {
		c.reservations = service.NewReservationService(c.apiClient)
	}
	return c.reservations
}

// Users returns the user service for managing users.
func (c *Client) Users() *service.UserService {
	c.mu.RLock()
	if c.users != nil {
		c.mu.RUnlock()
		return c.users
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.users == nil {
		c.users = service.NewUserService(c.apiClient)
	}
	return c.users
}

// Groups returns the group service for managing user groups.
func (c *Client) Groups() *service.GroupService {
	c.mu.RLock()
	if c.groups != nil {
		c.mu.RUnlock()
		return c.groups
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.groups == nil {
		c.groups = service.NewGroupService(c.apiClient)
	}
	return c.groups
}

// VisibilityGroups returns the visibility group service for managing restricted visibility groups.
func (c *Client) VisibilityGroups() *service.VisibilityGroupService {
	c.mu.RLock()
	if c.visibilityGroups != nil {
		c.mu.RUnlock()
		return c.visibilityGroups
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.visibilityGroups == nil {
		c.visibilityGroups = service.NewVisibilityGroupService(c.apiClient)
	}
	return c.visibilityGroups
}

// Config returns the config service for querying tenant configuration.
func (c *Client) Config() *service.ConfigService {
	c.mu.RLock()
	if c.config != nil {
		c.mu.RUnlock()
		return c.config
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.config == nil {
		c.config = service.NewConfigService(c.apiClient)
	}
	return c.config
}

// Webhooks returns the webhook service for managing webhook configurations.
func (c *Client) Webhooks() *service.WebhookService {
	c.mu.RLock()
	if c.webhooks != nil {
		c.mu.RUnlock()
		return c.webhooks
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.webhooks == nil {
		c.webhooks = service.NewWebhookService(c.apiClient)
	}
	return c.webhooks
}

// WebhookCalls returns the webhook call service for querying webhook call history.
func (c *Client) WebhookCalls() *service.WebhookCallService {
	c.mu.RLock()
	if c.webhookCalls != nil {
		c.mu.RUnlock()
		return c.webhookCalls
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.webhookCalls == nil {
		c.webhookCalls = service.NewWebhookCallService(c.apiClient)
	}
	return c.webhookCalls
}

// Tags returns the tag service for managing tags.
func (c *Client) Tags() *service.TagService {
	c.mu.RLock()
	if c.tags != nil {
		c.mu.RUnlock()
		return c.tags
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.tags == nil {
		c.tags = service.NewTagService(c.apiClient)
	}
	return c.tags
}

// Assets returns the asset service for querying asset balances at address and wallet level.
func (c *Client) Assets() *service.AssetService {
	c.mu.RLock()
	if c.assets != nil {
		c.mu.RUnlock()
		return c.assets
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.assets == nil {
		c.assets = service.NewAssetService(c.apiClient)
	}
	return c.assets
}

// Actions returns the action service for managing automated actions.
func (c *Client) Actions() *service.ActionService {
	c.mu.RLock()
	if c.actions != nil {
		c.mu.RUnlock()
		return c.actions
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.actions == nil {
		c.actions = service.NewActionService(c.apiClient)
	}
	return c.actions
}

// Blockchains returns the blockchain service for querying blockchain metadata.
func (c *Client) Blockchains() *service.BlockchainService {
	c.mu.RLock()
	if c.blockchains != nil {
		c.mu.RUnlock()
		return c.blockchains
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.blockchains == nil {
		c.blockchains = service.NewBlockchainService(c.apiClient)
	}
	return c.blockchains
}

// Exchanges returns the exchange service for managing exchange accounts.
func (c *Client) Exchanges() *service.ExchangeService {
	c.mu.RLock()
	if c.exchanges != nil {
		c.mu.RUnlock()
		return c.exchanges
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.exchanges == nil {
		c.exchanges = service.NewExchangeService(c.apiClient)
	}
	return c.exchanges
}

// Fiat returns the fiat service for managing fiat provider accounts.
func (c *Client) Fiat() *service.FiatService {
	c.mu.RLock()
	if c.fiat != nil {
		c.mu.RUnlock()
		return c.fiat
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.fiat == nil {
		c.fiat = service.NewFiatService(c.apiClient)
	}
	return c.fiat
}

// FeePayers returns the fee payer service for managing fee payers.
func (c *Client) FeePayers() *service.FeePayerService {
	c.mu.RLock()
	if c.feePayers != nil {
		c.mu.RUnlock()
		return c.feePayers
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.feePayers == nil {
		c.feePayers = service.NewFeePayerService(c.apiClient)
	}
	return c.feePayers
}

// Health returns the health service for querying system health checks.
func (c *Client) Health() *service.HealthService {
	c.mu.RLock()
	if c.health != nil {
		c.mu.RUnlock()
		return c.health
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.health == nil {
		c.health = service.NewHealthService(c.apiClient)
	}
	return c.health
}

// Jobs returns the job service for managing background jobs.
func (c *Client) Jobs() *service.JobService {
	c.mu.RLock()
	if c.jobs != nil {
		c.mu.RUnlock()
		return c.jobs
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.jobs == nil {
		c.jobs = service.NewJobService(c.apiClient)
	}
	return c.jobs
}

// Scores returns the score service for managing address risk scores.
func (c *Client) Scores() *service.ScoreService {
	c.mu.RLock()
	if c.scores != nil {
		c.mu.RUnlock()
		return c.scores
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.scores == nil {
		c.scores = service.NewScoreService(c.apiClient)
	}
	return c.scores
}

// Statistics returns the statistics service for querying portfolio and tag statistics.
func (c *Client) Statistics() *service.StatisticsService {
	c.mu.RLock()
	if c.statistics != nil {
		c.mu.RUnlock()
		return c.statistics
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.statistics == nil {
		c.statistics = service.NewStatisticsService(c.apiClient)
	}
	return c.statistics
}

// TokenMetadata returns the token metadata service for querying NFT and token metadata.
func (c *Client) TokenMetadata() *service.TokenMetadataService {
	c.mu.RLock()
	if c.tokenMetadata != nil {
		c.mu.RUnlock()
		return c.tokenMetadata
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.tokenMetadata == nil {
		c.tokenMetadata = service.NewTokenMetadataService(c.apiClient)
	}
	return c.tokenMetadata
}

// UserDevices returns the user device service for managing device pairing.
func (c *Client) UserDevices() *service.UserDeviceService {
	c.mu.RLock()
	if c.userDevices != nil {
		c.mu.RUnlock()
		return c.userDevices
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.userDevices == nil {
		c.userDevices = service.NewUserDeviceService(c.apiClient)
	}
	return c.userDevices
}

// MultiFactorSignature returns the multi-factor signature service for MFA approval workflows.
func (c *Client) MultiFactorSignature() *service.MultiFactorSignatureService {
	c.mu.RLock()
	if c.multiFactorSignature != nil {
		c.mu.RUnlock()
		return c.multiFactorSignature
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.multiFactorSignature == nil {
		c.multiFactorSignature = service.NewMultiFactorSignatureService(c.apiClient)
	}
	return c.multiFactorSignature
}

// TaurusNetwork returns the Taurus Network namespace client providing access to all
// Taurus Network services: Participants, Pledges, Lending, Settlements, and Sharing.
//
// Example usage:
//
//	client.TaurusNetwork().Participants().GetMyParticipant(ctx)
//	client.TaurusNetwork().Pledges().Create(ctx, req)
//	client.TaurusNetwork().Lending().GetAgreement(ctx, agreementID)
//	client.TaurusNetwork().Settlements().Get(ctx, settlementID)
//	client.TaurusNetwork().Sharing().ShareAddress(ctx, req)
func (c *Client) TaurusNetwork() *TaurusNetworkClient {
	c.mu.RLock()
	if c.taurusNetwork != nil {
		c.mu.RUnlock()
		return c.taurusNetwork
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()
	if c.taurusNetwork == nil {
		c.taurusNetwork = &TaurusNetworkClient{
			participants: service.NewTaurusNetworkParticipantService(c.apiClient),
			pledges:      service.NewTaurusNetworkPledgeService(c.apiClient),
			lending:      service.NewTaurusNetworkLendingService(c.apiClient),
			settlements:  service.NewTaurusNetworkSettlementService(c.apiClient),
			sharing:      service.NewTaurusNetworkSharingService(c.apiClient),
		}
	}
	return c.taurusNetwork
}

// Ensure Client implements io.Closer
var _ io.Closer = (*Client)(nil)
