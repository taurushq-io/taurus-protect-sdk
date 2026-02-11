/**
 * Main client for the Taurus-PROTECT API.
 *
 * This is the primary entry point for using the SDK. The client provides
 * lazy-initialized access to all API classes.
 *
 * @example
 * ```typescript
 * const client = ProtectClient.create({
 *   host: 'https://your-protect-instance.example.com',
 *   apiKey: 'your-api-key',
 *   apiSecret: 'your-hex-encoded-secret',
 * });
 *
 * // Access low-level APIs
 * const wallets = await client.walletsApi.walletServiceGetWalletsV2();
 * const addresses = await client.addressesApi.addressServiceGetAddresses();
 *
 * // Clean up when done
 * client.close();
 * ```
 */

import { Configuration, type Middleware } from "./internal/openapi/runtime";
import {
  ActionsApi,
  AddressWhitelistingApi,
  AddressesApi,
  AirGapApi,
  AssetsApi,
  AuditApi,
  AuthenticationApi,
  AuthenticationHMACApi,
  AuthenticationOIDCApi,
  AuthenticationSAMLApi,
  BalancesApi,
  BlockchainApi,
  BusinessRulesApi,
  ChangesApi,
  ConfigApi,
  ContractWhitelistingApi,
  CurrenciesApi,
  ExchangeApi,
  FeeApi,
  FeePayersApi,
  FiatApi,
  GovernanceRulesApi,
  GroupsApi,
  HealthApi,
  JobsApi,
  MultiFactorSignatureApi,
  PricesApi,
  RequestsApi,
  RequestsADAApi,
  RequestsALGOApi,
  RequestsContractsApi,
  RequestsCosmosApi,
  RequestsDOTApi,
  RequestsFTMApi,
  RequestsHederaApi,
  RequestsICPApi,
  RequestsMinaApi,
  RequestsNEARApi,
  RequestsSOLApi,
  RequestsXLMApi,
  RequestsXTZApi,
  ReservationsApi,
  RestrictedVisibilityGroupsApi,
  SCIMApi,
  ScoresApi,
  StakingApi,
  StatisticsApi,
  StewardApi,
  TagsApi,
  TaurusNetworkLendingApi,
  TaurusNetworkParticipantApi,
  TaurusNetworkPledgeApi,
  TaurusNetworkSettlementApi,
  TaurusNetworkSharedAddressAssetApi,
  TokenMetadataApi,
  TransactionsApi,
  UserDeviceApi,
  UsersApi,
  WalletsApi,
  WebhookCallsApi,
  WebhooksApi,
} from "./internal/openapi/apis";
import { ConfigurationError, ServerError } from "./errors";
import { createTPV1Middleware } from "./transport";

// Import service classes
import { AddressService } from "./services/address-service";
import { AirGapService } from "./services/air-gap-service";
import { AssetService } from "./services/asset-service";
import { AuditService } from "./services/audit-service";
import { BalanceService } from "./services/balance-service";
import { ConfigService } from "./services/config-service";
import { CurrencyService } from "./services/currency-service";
import { ExchangeService } from "./services/exchange-service";
import { FeePayerService } from "./services/fee-payer-service";
import { FeeService } from "./services/fee-service";
import { GovernanceRuleService } from "./services/governance-rule-service";
import { GroupService } from "./services/group-service";
import { HealthService } from "./services/health-service";
import { JobService } from "./services/job-service";
import { PriceService } from "./services/price-service";
import { RequestService } from "./services/request-service";
import { StakingService } from "./services/staking-service";
import { StatisticsService } from "./services/statistics-service";
import { TagService } from "./services/tag-service";
import { TokenMetadataService } from "./services/token-metadata-service";
import { TransactionService } from "./services/transaction-service";
import { UserService } from "./services/user-service";
import { VisibilityGroupService } from "./services/visibility-group-service";
import { WalletService } from "./services/wallet-service";
import { WebhookCallService } from "./services/webhook-call-service";
import { WebhookService } from "./services/webhook-service";
import { WhitelistedAddressService } from "./services/whitelisted-address-service";
import { WhitelistedAssetService } from "./services/whitelisted-asset-service";

import { ActionService } from "./services/action-service";
import { BlockchainService } from "./services/blockchain-service";
import { BusinessRuleService } from "./services/business-rule-service";
import { ChangeService } from "./services/change-service";
import { ContractWhitelistingService } from "./services/contract-whitelisting-service";
import { FiatService } from "./services/fiat-service";
import { MultiFactorSignatureService } from "./services/multi-factor-signature-service";
import { ReservationService } from "./services/reservation-service";
import { ScoreService } from "./services/score-service";
import { UserDeviceService } from "./services/user-device-service";

import { ParticipantService } from "./services/taurus-network/participant-service";
import { PledgeService } from "./services/taurus-network/pledge-service";
import { LendingService } from "./services/taurus-network/lending-service";
import { SettlementService } from "./services/taurus-network/settlement-service";
import { SharingService } from "./services/taurus-network/sharing-service";

// Import cache
import { RulesContainerCache } from "./cache";

// Import crypto utilities
import { decodePublicKeysPem } from "./crypto/keys";

// Import decoders for whitelisted address/asset verification
import {
  rulesContainerFromBase64,
  userSignaturesFromBase64,
} from "./mappers/governance-rules";

/**
 * Configuration options for creating a ProtectClient.
 */
export interface ProtectClientConfig {
  /**
   * The base URL of the Taurus-PROTECT API.
   * Example: 'https://your-protect-instance.example.com'
   */
  host: string;

  /**
   * The API key for authentication.
   */
  apiKey: string;

  /**
   * The API secret in hexadecimal format.
   */
  apiSecret: string;

  /**
   * SuperAdmin public keys in PEM format for signature verification.
   * Required for integrity verification of whitelisted addresses and assets.
   * At least one key must be provided.
   */
  superAdminKeysPem: string[];

  /**
   * Minimum number of valid SuperAdmin signatures required.
   * Default: 1
   */
  minValidSignatures?: number;

  /**
   * Time-to-live for the rules container cache in milliseconds.
   * Default: 300000 (5 minutes)
   */
  rulesCacheTtlMs?: number;

  /**
   * Request timeout in milliseconds.
   * Default: 30000 (30 seconds)
   */
  timeout?: number;

  /**
   * Additional middleware to apply to requests.
   */
  middleware?: Middleware[];
}

/**
 * TaurusNetwork namespace providing access to Taurus Network APIs.
 *
 * These APIs are for managing lending, pledges, settlements, participants,
 * and shared address assets in the Taurus Network.
 */
export class TaurusNetworkNamespace {
  private readonly apiConfiguration: Configuration;
  private readonly ensureOpen: () => void;

  // Lazy-initialized API instances
  private _lendingApi?: TaurusNetworkLendingApi;
  private _participantApi?: TaurusNetworkParticipantApi;
  private _pledgeApi?: TaurusNetworkPledgeApi;
  private _settlementApi?: TaurusNetworkSettlementApi;
  private _sharedAddressAssetApi?: TaurusNetworkSharedAddressAssetApi;

  // High-level service instances
  private _participantService?: ParticipantService;
  private _pledgeService?: PledgeService;
  private _lendingService?: LendingService;
  private _settlementService?: SettlementService;
  private _sharingService?: SharingService;

  constructor(apiConfiguration: Configuration, ensureOpen: () => void) {
    this.apiConfiguration = apiConfiguration;
    this.ensureOpen = ensureOpen;
  }

  /**
   * Low-level Taurus Network Lending API access.
   */
  get lendingApi(): TaurusNetworkLendingApi {
    this.ensureOpen();
    if (!this._lendingApi) {
      this._lendingApi = new TaurusNetworkLendingApi(this.apiConfiguration);
    }
    return this._lendingApi;
  }

  /**
   * Low-level Taurus Network Participant API access.
   */
  get participantApi(): TaurusNetworkParticipantApi {
    this.ensureOpen();
    if (!this._participantApi) {
      this._participantApi = new TaurusNetworkParticipantApi(
        this.apiConfiguration
      );
    }
    return this._participantApi;
  }

  /**
   * Low-level Taurus Network Pledge API access.
   */
  get pledgeApi(): TaurusNetworkPledgeApi {
    this.ensureOpen();
    if (!this._pledgeApi) {
      this._pledgeApi = new TaurusNetworkPledgeApi(this.apiConfiguration);
    }
    return this._pledgeApi;
  }

  /**
   * Low-level Taurus Network Settlement API access.
   */
  get settlementApi(): TaurusNetworkSettlementApi {
    this.ensureOpen();
    if (!this._settlementApi) {
      this._settlementApi = new TaurusNetworkSettlementApi(
        this.apiConfiguration
      );
    }
    return this._settlementApi;
  }

  /**
   * Low-level Taurus Network Shared Address Asset API access.
   */
  get sharedAddressAssetApi(): TaurusNetworkSharedAddressAssetApi {
    this.ensureOpen();
    if (!this._sharedAddressAssetApi) {
      this._sharedAddressAssetApi = new TaurusNetworkSharedAddressAssetApi(
        this.apiConfiguration
      );
    }
    return this._sharedAddressAssetApi;
  }

  // ===== High-Level TaurusNetwork Service Accessors =====

  /**
   * High-level Participant service for Taurus Network participant management.
   */
  get participants(): ParticipantService {
    this.ensureOpen();
    if (!this._participantService) {
      this._participantService = new ParticipantService(this.participantApi);
    }
    return this._participantService;
  }

  /**
   * High-level Pledge service for Taurus Network pledge lifecycle.
   */
  get pledges(): PledgeService {
    this.ensureOpen();
    if (!this._pledgeService) {
      this._pledgeService = new PledgeService(this.pledgeApi);
    }
    return this._pledgeService;
  }

  /**
   * High-level Lending service for Taurus Network offers and agreements.
   */
  get lending(): LendingService {
    this.ensureOpen();
    if (!this._lendingService) {
      this._lendingService = new LendingService(this.lendingApi);
    }
    return this._lendingService;
  }

  /**
   * High-level Settlement service for Taurus Network settlement operations.
   */
  get settlements(): SettlementService {
    this.ensureOpen();
    if (!this._settlementService) {
      this._settlementService = new SettlementService(this.settlementApi);
    }
    return this._settlementService;
  }

  /**
   * High-level Sharing service for Taurus Network address/asset sharing.
   */
  get sharing(): SharingService {
    this.ensureOpen();
    if (!this._sharingService) {
      this._sharingService = new SharingService(this.sharedAddressAssetApi);
    }
    return this._sharingService;
  }

  /**
   * Clears all cached API instances.
   * @internal
   */
  clear(): void {
    this._lendingApi = undefined;
    this._participantApi = undefined;
    this._pledgeApi = undefined;
    this._settlementApi = undefined;
    this._sharedAddressAssetApi = undefined;
    this._participantService = undefined;
    this._pledgeService = undefined;
    this._lendingService = undefined;
    this._settlementService = undefined;
    this._sharingService = undefined;
  }
}

/**
 * Main client for interacting with the Taurus-PROTECT API.
 *
 * Creates a client using the static `create` method. The client uses
 * TPV1-HMAC-SHA256 authentication for all API requests.
 *
 * Services are lazy-initialized on first access for efficiency.
 */
export class ProtectClient {
  private readonly config: ProtectClientConfig;
  private readonly apiConfiguration: Configuration;
  private closed = false;

  // Lazy-initialized API instances
  private _actionsApi?: ActionsApi;
  private _addressWhitelistingApi?: AddressWhitelistingApi;
  private _addressesApi?: AddressesApi;
  private _airGapApi?: AirGapApi;
  private _assetsApi?: AssetsApi;
  private _auditApi?: AuditApi;
  private _authenticationApi?: AuthenticationApi;
  private _authenticationHMACApi?: AuthenticationHMACApi;
  private _authenticationOIDCApi?: AuthenticationOIDCApi;
  private _authenticationSAMLApi?: AuthenticationSAMLApi;
  private _balancesApi?: BalancesApi;
  private _blockchainApi?: BlockchainApi;
  private _businessRulesApi?: BusinessRulesApi;
  private _changesApi?: ChangesApi;
  private _configApi?: ConfigApi;
  private _contractWhitelistingApi?: ContractWhitelistingApi;
  private _currenciesApi?: CurrenciesApi;
  private _exchangeApi?: ExchangeApi;
  private _feeApi?: FeeApi;
  private _feePayersApi?: FeePayersApi;
  private _fiatApi?: FiatApi;
  private _governanceRulesApi?: GovernanceRulesApi;
  private _groupsApi?: GroupsApi;
  private _healthApi?: HealthApi;
  private _jobsApi?: JobsApi;
  private _multiFactorSignatureApi?: MultiFactorSignatureApi;
  private _pricesApi?: PricesApi;
  private _requestsApi?: RequestsApi;
  private _requestsADAApi?: RequestsADAApi;
  private _requestsALGOApi?: RequestsALGOApi;
  private _requestsContractsApi?: RequestsContractsApi;
  private _requestsCosmosApi?: RequestsCosmosApi;
  private _requestsDOTApi?: RequestsDOTApi;
  private _requestsFTMApi?: RequestsFTMApi;
  private _requestsHederaApi?: RequestsHederaApi;
  private _requestsICPApi?: RequestsICPApi;
  private _requestsMinaApi?: RequestsMinaApi;
  private _requestsNEARApi?: RequestsNEARApi;
  private _requestsSOLApi?: RequestsSOLApi;
  private _requestsXLMApi?: RequestsXLMApi;
  private _requestsXTZApi?: RequestsXTZApi;
  private _reservationsApi?: ReservationsApi;
  private _restrictedVisibilityGroupsApi?: RestrictedVisibilityGroupsApi;
  private _scimApi?: SCIMApi;
  private _scoresApi?: ScoresApi;
  private _stakingApi?: StakingApi;
  private _statisticsApi?: StatisticsApi;
  private _stewardApi?: StewardApi;
  private _tagsApi?: TagsApi;
  private _tokenMetadataApi?: TokenMetadataApi;
  private _transactionsApi?: TransactionsApi;
  private _userDeviceApi?: UserDeviceApi;
  private _usersApi?: UsersApi;
  private _walletsApi?: WalletsApi;
  private _webhookCallsApi?: WebhookCallsApi;
  private _webhooksApi?: WebhooksApi;

  // TaurusNetwork namespace
  private _taurusNetwork?: TaurusNetworkNamespace;

  // High-level service instances
  private _addressService?: AddressService;
  private _airGapService?: AirGapService;
  private _assetService?: AssetService;
  private _auditService?: AuditService;
  private _balanceService?: BalanceService;
  private _configService?: ConfigService;
  private _currencyService?: CurrencyService;
  private _exchangeService?: ExchangeService;
  private _feeService?: FeeService;
  private _feePayerService?: FeePayerService;
  private _governanceRuleService?: GovernanceRuleService;
  private _groupService?: GroupService;
  private _healthService?: HealthService;
  private _jobService?: JobService;
  private _priceService?: PriceService;
  private _requestService?: RequestService;
  private _stakingService?: StakingService;
  private _statisticsService?: StatisticsService;
  private _tagService?: TagService;
  private _tokenMetadataService?: TokenMetadataService;
  private _transactionService?: TransactionService;
  private _userService?: UserService;
  private _visibilityGroupService?: VisibilityGroupService;
  private _walletService?: WalletService;
  private _webhookService?: WebhookService;
  private _webhookCallService?: WebhookCallService;
  private _whitelistedAddressService?: WhitelistedAddressService;
  private _whitelistedAssetService?: WhitelistedAssetService;

  private _actionService?: ActionService;
  private _blockchainService?: BlockchainService;
  private _businessRuleService?: BusinessRuleService;
  private _changeService?: ChangeService;
  private _contractWhitelistingService?: ContractWhitelistingService;
  private _fiatService?: FiatService;
  private _multiFactorSignatureService?: MultiFactorSignatureService;
  private _reservationService?: ReservationService;
  private _scoreService?: ScoreService;
  private _userDeviceService?: UserDeviceService;

  // Rules container cache for signature verification
  private _rulesCache?: RulesContainerCache;

  /**
   * Private constructor - use ProtectClient.create() instead.
   */
  private constructor(config: ProtectClientConfig) {
    this.config = config;

    // Build middleware chain
    const middleware: Middleware[] = [
      createTPV1Middleware(config.apiKey, config.apiSecret),
      ...(config.middleware ?? []),
    ];

    // Create OpenAPI configuration
    this.apiConfiguration = new Configuration({
      basePath: config.host.replace(/\/+$/, ""), // Remove trailing slashes
      middleware,
    });
  }

  /**
   * Creates a new ProtectClient with the specified configuration.
   *
   * @param config - Client configuration
   * @returns A new ProtectClient instance
   * @throws ConfigurationError if configuration is invalid
   *
   * @example
   * ```typescript
   * const client = ProtectClient.create({
   *   host: 'https://protect.example.com',
   *   apiKey: 'your-api-key',
   *   apiSecret: 'your-hex-secret',
   * });
   * ```
   */
  static create(config: ProtectClientConfig): ProtectClient {
    // Validate configuration
    if (!config.host) {
      throw new ConfigurationError("host is required");
    }
    if (!config.apiKey) {
      throw new ConfigurationError("apiKey is required");
    }
    if (!config.apiSecret) {
      throw new ConfigurationError("apiSecret is required");
    }

    // Validate host URL
    try {
      new URL(config.host);
    } catch {
      throw new ConfigurationError(`Invalid host URL: ${config.host}`);
    }

    // Validate API secret is valid hex
    if (!/^[0-9a-fA-F]+$/.test(config.apiSecret)) {
      throw new ConfigurationError(
        "apiSecret must be a valid hexadecimal string"
      );
    }

    // Validate superAdminKeysPem is non-empty
    if (!config.superAdminKeysPem || config.superAdminKeysPem.length === 0) {
      throw new ConfigurationError(
        "superAdminKeysPem is required: at least one SuperAdmin public key must be provided for integrity verification"
      );
    }

    // Reject minValidSignatures <= 0 whenever explicitly set.
    if (config.minValidSignatures !== undefined && config.minValidSignatures <= 0) {
      throw new ConfigurationError(
        "minValidSignatures must be greater than 0 when specified"
      );
    }

    return new ProtectClient(config);
  }

  /**
   * Closes the client and releases any resources.
   *
   * After calling close(), the client should not be used.
   *
   * **SECURITY NOTE - Memory Wiping Limitations in JavaScript:**
   * Due to JavaScript's garbage collector, copies of sensitive data (like the
   * API secret) may still exist in memory even after this method is called.
   * The V8 engine may create intermediate copies, cache values, or defer garbage
   * collection indefinitely. This cleanup is best-effort and provides defense-in-depth,
   * but cannot guarantee complete removal of sensitive data from memory.
   * See TPV1Auth.close() for more details on memory security limitations.
   */
  close(): void {
    if (this.closed) {
      return;
    }
    this.closed = true;

    // SECURITY: Best-effort secret wiping
    // Attempt to overwrite the API secret in the config
    // Note: JavaScript/Node.js GC may retain copies - this is defense-in-depth
    try {
      if (this.config.apiSecret) {
        // Overwrite the string reference with zeros
        // Note: JavaScript strings are immutable, so this only replaces the reference
        // The original string may still exist in memory until GC runs
        const secretLength = this.config.apiSecret.length;
        (this.config as { apiSecret: string }).apiSecret = '0'.repeat(secretLength);
      }
    } catch {
      // Ignore - secret wiping is best-effort in JavaScript
    }

    // Clear cached API instances using array iteration for maintainability
    const apiFields = [
      '_actionsApi',
      '_addressWhitelistingApi',
      '_addressesApi',
      '_airGapApi',
      '_assetsApi',
      '_auditApi',
      '_authenticationApi',
      '_authenticationHMACApi',
      '_authenticationOIDCApi',
      '_authenticationSAMLApi',
      '_balancesApi',
      '_blockchainApi',
      '_businessRulesApi',
      '_changesApi',
      '_configApi',
      '_contractWhitelistingApi',
      '_currenciesApi',
      '_exchangeApi',
      '_feeApi',
      '_feePayersApi',
      '_fiatApi',
      '_governanceRulesApi',
      '_groupsApi',
      '_healthApi',
      '_jobsApi',
      '_multiFactorSignatureApi',
      '_pricesApi',
      '_requestsApi',
      '_requestsADAApi',
      '_requestsALGOApi',
      '_requestsContractsApi',
      '_requestsCosmosApi',
      '_requestsDOTApi',
      '_requestsFTMApi',
      '_requestsHederaApi',
      '_requestsICPApi',
      '_requestsMinaApi',
      '_requestsNEARApi',
      '_requestsSOLApi',
      '_requestsXLMApi',
      '_requestsXTZApi',
      '_reservationsApi',
      '_restrictedVisibilityGroupsApi',
      '_scimApi',
      '_scoresApi',
      '_stakingApi',
      '_statisticsApi',
      '_stewardApi',
      '_tagsApi',
      '_tokenMetadataApi',
      '_transactionsApi',
      '_userDeviceApi',
      '_usersApi',
      '_walletsApi',
      '_webhookCallsApi',
      '_webhooksApi',
    ] as const;
    for (const field of apiFields) {
      (this as Record<string, unknown>)[field] = undefined;
    }

    // Clear TaurusNetwork namespace
    this._taurusNetwork?.clear();
    this._taurusNetwork = undefined;

    // Clear high-level service instances using array iteration
    const serviceFields = [
      '_actionService',
      '_addressService',
      '_airGapService',
      '_assetService',
      '_auditService',
      '_balanceService',
      '_blockchainService',
      '_businessRuleService',
      '_changeService',
      '_configService',
      '_contractWhitelistingService',
      '_currencyService',
      '_feeService',
      '_exchangeService',
      '_feePayerService',
      '_fiatService',
      '_governanceRuleService',
      '_groupService',
      '_healthService',
      '_jobService',
      '_multiFactorSignatureService',
      '_priceService',
      '_requestService',
      '_reservationService',
      '_scoreService',
      '_stakingService',
      '_statisticsService',
      '_tagService',
      '_tokenMetadataService',
      '_transactionService',
      '_userDeviceService',
      '_userService',
      '_visibilityGroupService',
      '_walletService',
      '_webhookService',
      '_webhookCallService',
      '_whitelistedAddressService',
      '_whitelistedAssetService',
    ] as const;
    for (const field of serviceFields) {
      (this as Record<string, unknown>)[field] = undefined;
    }

    // Clear rules cache
    this._rulesCache?.clear();
    this._rulesCache = undefined;
  }

  /**
   * Checks if the client has been closed.
   */
  get isClosed(): boolean {
    return this.closed;
  }

  private ensureOpen(): void {
    if (this.closed) {
      throw new ConfigurationError("Client has been closed");
    }
  }

  // ===== Configuration Accessors =====

  /**
   * Returns the configured host URL.
   */
  get host(): string {
    return this.config.host;
  }

  /**
   * Returns the SuperAdmin public keys in PEM format.
   */
  get superAdminKeysPem(): string[] {
    return this.config.superAdminKeysPem;
  }

  /**
   * Returns the minimum valid signatures threshold.
   */
  get minValidSignatures(): number {
    return this.config.minValidSignatures ?? 1;
  }

  /**
   * Returns the rules cache TTL in milliseconds.
   */
  get rulesCacheTtlMs(): number {
    return this.config.rulesCacheTtlMs ?? 300000;
  }

  // ===== TaurusNetwork Namespace =====

  /**
   * TaurusNetwork namespace for accessing Taurus Network APIs.
   *
   * @example
   * ```typescript
   * const participants = await client.taurusNetwork.participantApi.getAllParticipants();
   * const pledges = await client.taurusNetwork.pledgeApi.getAllPledges();
   * ```
   */
  get taurusNetwork(): TaurusNetworkNamespace {
    this.ensureOpen();
    if (!this._taurusNetwork) {
      this._taurusNetwork = new TaurusNetworkNamespace(
        this.apiConfiguration,
        () => this.ensureOpen()
      );
    }
    return this._taurusNetwork;
  }

  // ===== Core API Accessors =====

  /**
   * Low-level Actions API access.
   */
  get actionsApi(): ActionsApi {
    this.ensureOpen();
    if (!this._actionsApi) {
      this._actionsApi = new ActionsApi(this.apiConfiguration);
    }
    return this._actionsApi;
  }

  /**
   * Low-level Address Whitelisting API access.
   */
  get addressWhitelistingApi(): AddressWhitelistingApi {
    this.ensureOpen();
    if (!this._addressWhitelistingApi) {
      this._addressWhitelistingApi = new AddressWhitelistingApi(
        this.apiConfiguration
      );
    }
    return this._addressWhitelistingApi;
  }

  /**
   * Low-level Addresses API access.
   */
  get addressesApi(): AddressesApi {
    this.ensureOpen();
    if (!this._addressesApi) {
      this._addressesApi = new AddressesApi(this.apiConfiguration);
    }
    return this._addressesApi;
  }

  /**
   * Low-level Air Gap API access.
   */
  get airGapApi(): AirGapApi {
    this.ensureOpen();
    if (!this._airGapApi) {
      this._airGapApi = new AirGapApi(this.apiConfiguration);
    }
    return this._airGapApi;
  }

  /**
   * Low-level Assets API access.
   */
  get assetsApi(): AssetsApi {
    this.ensureOpen();
    if (!this._assetsApi) {
      this._assetsApi = new AssetsApi(this.apiConfiguration);
    }
    return this._assetsApi;
  }

  /**
   * Low-level Audit API access.
   */
  get auditApi(): AuditApi {
    this.ensureOpen();
    if (!this._auditApi) {
      this._auditApi = new AuditApi(this.apiConfiguration);
    }
    return this._auditApi;
  }

  /**
   * Low-level Authentication API access.
   */
  get authenticationApi(): AuthenticationApi {
    this.ensureOpen();
    if (!this._authenticationApi) {
      this._authenticationApi = new AuthenticationApi(this.apiConfiguration);
    }
    return this._authenticationApi;
  }

  /**
   * Low-level Authentication HMAC API access.
   */
  get authenticationHMACApi(): AuthenticationHMACApi {
    this.ensureOpen();
    if (!this._authenticationHMACApi) {
      this._authenticationHMACApi = new AuthenticationHMACApi(
        this.apiConfiguration
      );
    }
    return this._authenticationHMACApi;
  }

  /**
   * Low-level Authentication OIDC API access.
   */
  get authenticationOIDCApi(): AuthenticationOIDCApi {
    this.ensureOpen();
    if (!this._authenticationOIDCApi) {
      this._authenticationOIDCApi = new AuthenticationOIDCApi(
        this.apiConfiguration
      );
    }
    return this._authenticationOIDCApi;
  }

  /**
   * Low-level Authentication SAML API access.
   */
  get authenticationSAMLApi(): AuthenticationSAMLApi {
    this.ensureOpen();
    if (!this._authenticationSAMLApi) {
      this._authenticationSAMLApi = new AuthenticationSAMLApi(
        this.apiConfiguration
      );
    }
    return this._authenticationSAMLApi;
  }

  /**
   * Low-level Balances API access.
   */
  get balancesApi(): BalancesApi {
    this.ensureOpen();
    if (!this._balancesApi) {
      this._balancesApi = new BalancesApi(this.apiConfiguration);
    }
    return this._balancesApi;
  }

  /**
   * Low-level Blockchain API access.
   */
  get blockchainApi(): BlockchainApi {
    this.ensureOpen();
    if (!this._blockchainApi) {
      this._blockchainApi = new BlockchainApi(this.apiConfiguration);
    }
    return this._blockchainApi;
  }

  /**
   * Low-level Business Rules API access.
   */
  get businessRulesApi(): BusinessRulesApi {
    this.ensureOpen();
    if (!this._businessRulesApi) {
      this._businessRulesApi = new BusinessRulesApi(this.apiConfiguration);
    }
    return this._businessRulesApi;
  }

  /**
   * Low-level Changes API access.
   */
  get changesApi(): ChangesApi {
    this.ensureOpen();
    if (!this._changesApi) {
      this._changesApi = new ChangesApi(this.apiConfiguration);
    }
    return this._changesApi;
  }

  /**
   * Low-level Config API access.
   */
  get configApi(): ConfigApi {
    this.ensureOpen();
    if (!this._configApi) {
      this._configApi = new ConfigApi(this.apiConfiguration);
    }
    return this._configApi;
  }

  /**
   * Low-level Contract Whitelisting API access.
   */
  get contractWhitelistingApi(): ContractWhitelistingApi {
    this.ensureOpen();
    if (!this._contractWhitelistingApi) {
      this._contractWhitelistingApi = new ContractWhitelistingApi(
        this.apiConfiguration
      );
    }
    return this._contractWhitelistingApi;
  }

  /**
   * Low-level Currencies API access.
   */
  get currenciesApi(): CurrenciesApi {
    this.ensureOpen();
    if (!this._currenciesApi) {
      this._currenciesApi = new CurrenciesApi(this.apiConfiguration);
    }
    return this._currenciesApi;
  }

  /**
   * Low-level Exchange API access.
   */
  get exchangeApi(): ExchangeApi {
    this.ensureOpen();
    if (!this._exchangeApi) {
      this._exchangeApi = new ExchangeApi(this.apiConfiguration);
    }
    return this._exchangeApi;
  }

  /**
   * Low-level Fee API access.
   */
  get feeApi(): FeeApi {
    this.ensureOpen();
    if (!this._feeApi) {
      this._feeApi = new FeeApi(this.apiConfiguration);
    }
    return this._feeApi;
  }

  /**
   * Low-level Fee Payers API access.
   */
  get feePayersApi(): FeePayersApi {
    this.ensureOpen();
    if (!this._feePayersApi) {
      this._feePayersApi = new FeePayersApi(this.apiConfiguration);
    }
    return this._feePayersApi;
  }

  /**
   * Low-level Fiat API access.
   */
  get fiatApi(): FiatApi {
    this.ensureOpen();
    if (!this._fiatApi) {
      this._fiatApi = new FiatApi(this.apiConfiguration);
    }
    return this._fiatApi;
  }

  /**
   * Low-level Governance Rules API access.
   */
  get governanceRulesApi(): GovernanceRulesApi {
    this.ensureOpen();
    if (!this._governanceRulesApi) {
      this._governanceRulesApi = new GovernanceRulesApi(this.apiConfiguration);
    }
    return this._governanceRulesApi;
  }

  /**
   * Low-level Groups API access.
   */
  get groupsApi(): GroupsApi {
    this.ensureOpen();
    if (!this._groupsApi) {
      this._groupsApi = new GroupsApi(this.apiConfiguration);
    }
    return this._groupsApi;
  }

  /**
   * Low-level Health API access.
   */
  get healthApi(): HealthApi {
    this.ensureOpen();
    if (!this._healthApi) {
      this._healthApi = new HealthApi(this.apiConfiguration);
    }
    return this._healthApi;
  }

  /**
   * Low-level Jobs API access.
   */
  get jobsApi(): JobsApi {
    this.ensureOpen();
    if (!this._jobsApi) {
      this._jobsApi = new JobsApi(this.apiConfiguration);
    }
    return this._jobsApi;
  }

  /**
   * Low-level Multi Factor Signature API access.
   */
  get multiFactorSignatureApi(): MultiFactorSignatureApi {
    this.ensureOpen();
    if (!this._multiFactorSignatureApi) {
      this._multiFactorSignatureApi = new MultiFactorSignatureApi(
        this.apiConfiguration
      );
    }
    return this._multiFactorSignatureApi;
  }

  /**
   * Low-level Prices API access.
   */
  get pricesApi(): PricesApi {
    this.ensureOpen();
    if (!this._pricesApi) {
      this._pricesApi = new PricesApi(this.apiConfiguration);
    }
    return this._pricesApi;
  }

  /**
   * Low-level Requests API access.
   */
  get requestsApi(): RequestsApi {
    this.ensureOpen();
    if (!this._requestsApi) {
      this._requestsApi = new RequestsApi(this.apiConfiguration);
    }
    return this._requestsApi;
  }

  /**
   * Low-level Requests ADA API access.
   */
  get requestsADAApi(): RequestsADAApi {
    this.ensureOpen();
    if (!this._requestsADAApi) {
      this._requestsADAApi = new RequestsADAApi(this.apiConfiguration);
    }
    return this._requestsADAApi;
  }

  /**
   * Low-level Requests ALGO API access.
   */
  get requestsALGOApi(): RequestsALGOApi {
    this.ensureOpen();
    if (!this._requestsALGOApi) {
      this._requestsALGOApi = new RequestsALGOApi(this.apiConfiguration);
    }
    return this._requestsALGOApi;
  }

  /**
   * Low-level Requests Contracts API access.
   */
  get requestsContractsApi(): RequestsContractsApi {
    this.ensureOpen();
    if (!this._requestsContractsApi) {
      this._requestsContractsApi = new RequestsContractsApi(
        this.apiConfiguration
      );
    }
    return this._requestsContractsApi;
  }

  /**
   * Low-level Requests Cosmos API access.
   */
  get requestsCosmosApi(): RequestsCosmosApi {
    this.ensureOpen();
    if (!this._requestsCosmosApi) {
      this._requestsCosmosApi = new RequestsCosmosApi(this.apiConfiguration);
    }
    return this._requestsCosmosApi;
  }

  /**
   * Low-level Requests DOT API access.
   */
  get requestsDOTApi(): RequestsDOTApi {
    this.ensureOpen();
    if (!this._requestsDOTApi) {
      this._requestsDOTApi = new RequestsDOTApi(this.apiConfiguration);
    }
    return this._requestsDOTApi;
  }

  /**
   * Low-level Requests FTM API access.
   */
  get requestsFTMApi(): RequestsFTMApi {
    this.ensureOpen();
    if (!this._requestsFTMApi) {
      this._requestsFTMApi = new RequestsFTMApi(this.apiConfiguration);
    }
    return this._requestsFTMApi;
  }

  /**
   * Low-level Requests Hedera API access.
   */
  get requestsHederaApi(): RequestsHederaApi {
    this.ensureOpen();
    if (!this._requestsHederaApi) {
      this._requestsHederaApi = new RequestsHederaApi(this.apiConfiguration);
    }
    return this._requestsHederaApi;
  }

  /**
   * Low-level Requests ICP API access.
   */
  get requestsICPApi(): RequestsICPApi {
    this.ensureOpen();
    if (!this._requestsICPApi) {
      this._requestsICPApi = new RequestsICPApi(this.apiConfiguration);
    }
    return this._requestsICPApi;
  }

  /**
   * Low-level Requests Mina API access.
   */
  get requestsMinaApi(): RequestsMinaApi {
    this.ensureOpen();
    if (!this._requestsMinaApi) {
      this._requestsMinaApi = new RequestsMinaApi(this.apiConfiguration);
    }
    return this._requestsMinaApi;
  }

  /**
   * Low-level Requests NEAR API access.
   */
  get requestsNEARApi(): RequestsNEARApi {
    this.ensureOpen();
    if (!this._requestsNEARApi) {
      this._requestsNEARApi = new RequestsNEARApi(this.apiConfiguration);
    }
    return this._requestsNEARApi;
  }

  /**
   * Low-level Requests SOL API access.
   */
  get requestsSOLApi(): RequestsSOLApi {
    this.ensureOpen();
    if (!this._requestsSOLApi) {
      this._requestsSOLApi = new RequestsSOLApi(this.apiConfiguration);
    }
    return this._requestsSOLApi;
  }

  /**
   * Low-level Requests XLM API access.
   */
  get requestsXLMApi(): RequestsXLMApi {
    this.ensureOpen();
    if (!this._requestsXLMApi) {
      this._requestsXLMApi = new RequestsXLMApi(this.apiConfiguration);
    }
    return this._requestsXLMApi;
  }

  /**
   * Low-level Requests XTZ API access.
   */
  get requestsXTZApi(): RequestsXTZApi {
    this.ensureOpen();
    if (!this._requestsXTZApi) {
      this._requestsXTZApi = new RequestsXTZApi(this.apiConfiguration);
    }
    return this._requestsXTZApi;
  }

  /**
   * Low-level Reservations API access.
   */
  get reservationsApi(): ReservationsApi {
    this.ensureOpen();
    if (!this._reservationsApi) {
      this._reservationsApi = new ReservationsApi(this.apiConfiguration);
    }
    return this._reservationsApi;
  }

  /**
   * Low-level Restricted Visibility Groups API access.
   */
  get restrictedVisibilityGroupsApi(): RestrictedVisibilityGroupsApi {
    this.ensureOpen();
    if (!this._restrictedVisibilityGroupsApi) {
      this._restrictedVisibilityGroupsApi = new RestrictedVisibilityGroupsApi(
        this.apiConfiguration
      );
    }
    return this._restrictedVisibilityGroupsApi;
  }

  /**
   * Low-level SCIM API access.
   */
  get scimApi(): SCIMApi {
    this.ensureOpen();
    if (!this._scimApi) {
      this._scimApi = new SCIMApi(this.apiConfiguration);
    }
    return this._scimApi;
  }

  /**
   * Low-level Scores API access.
   */
  get scoresApi(): ScoresApi {
    this.ensureOpen();
    if (!this._scoresApi) {
      this._scoresApi = new ScoresApi(this.apiConfiguration);
    }
    return this._scoresApi;
  }

  /**
   * Low-level Staking API access.
   */
  get stakingApi(): StakingApi {
    this.ensureOpen();
    if (!this._stakingApi) {
      this._stakingApi = new StakingApi(this.apiConfiguration);
    }
    return this._stakingApi;
  }

  /**
   * Low-level Statistics API access.
   */
  get statisticsApi(): StatisticsApi {
    this.ensureOpen();
    if (!this._statisticsApi) {
      this._statisticsApi = new StatisticsApi(this.apiConfiguration);
    }
    return this._statisticsApi;
  }

  /**
   * Low-level Steward API access.
   */
  get stewardApi(): StewardApi {
    this.ensureOpen();
    if (!this._stewardApi) {
      this._stewardApi = new StewardApi(this.apiConfiguration);
    }
    return this._stewardApi;
  }

  /**
   * Low-level Tags API access.
   */
  get tagsApi(): TagsApi {
    this.ensureOpen();
    if (!this._tagsApi) {
      this._tagsApi = new TagsApi(this.apiConfiguration);
    }
    return this._tagsApi;
  }

  /**
   * Low-level Token Metadata API access.
   */
  get tokenMetadataApi(): TokenMetadataApi {
    this.ensureOpen();
    if (!this._tokenMetadataApi) {
      this._tokenMetadataApi = new TokenMetadataApi(this.apiConfiguration);
    }
    return this._tokenMetadataApi;
  }

  /**
   * Low-level Transactions API access.
   */
  get transactionsApi(): TransactionsApi {
    this.ensureOpen();
    if (!this._transactionsApi) {
      this._transactionsApi = new TransactionsApi(this.apiConfiguration);
    }
    return this._transactionsApi;
  }

  /**
   * Low-level User Device API access.
   */
  get userDeviceApi(): UserDeviceApi {
    this.ensureOpen();
    if (!this._userDeviceApi) {
      this._userDeviceApi = new UserDeviceApi(this.apiConfiguration);
    }
    return this._userDeviceApi;
  }

  /**
   * Low-level Users API access.
   */
  get usersApi(): UsersApi {
    this.ensureOpen();
    if (!this._usersApi) {
      this._usersApi = new UsersApi(this.apiConfiguration);
    }
    return this._usersApi;
  }

  /**
   * Low-level Wallets API access.
   */
  get walletsApi(): WalletsApi {
    this.ensureOpen();
    if (!this._walletsApi) {
      this._walletsApi = new WalletsApi(this.apiConfiguration);
    }
    return this._walletsApi;
  }

  /**
   * Low-level Webhook Calls API access.
   */
  get webhookCallsApi(): WebhookCallsApi {
    this.ensureOpen();
    if (!this._webhookCallsApi) {
      this._webhookCallsApi = new WebhookCallsApi(this.apiConfiguration);
    }
    return this._webhookCallsApi;
  }

  /**
   * Low-level Webhooks API access.
   */
  get webhooksApi(): WebhooksApi {
    this.ensureOpen();
    if (!this._webhooksApi) {
      this._webhooksApi = new WebhooksApi(this.apiConfiguration);
    }
    return this._webhooksApi;
  }

  // ===== High-Level Service Accessors =====

  /**
   * Gets the rules container cache for address signature verification.
   *
   * The cache is lazily initialized on first access and shared across services.
   * It requires GovernanceRulesApi to fetch the rules container.
   *
   * @internal
   */
  private getRulesCache(): RulesContainerCache {
    if (!this._rulesCache) {
      const governanceApi = this.governanceRulesApi;
      this._rulesCache = new RulesContainerCache(
        async () => {
          // Fetch rules container from governance rules API
          const response = await governanceApi.ruleServiceGetRules();
          const rulesBase64 = response.result?.rulesContainer;
          if (!rulesBase64) {
            throw new ServerError("No rules container returned from API");
          }
          // Decode the rules container from base64
          // Import the mapper function dynamically to avoid circular deps
          const { rulesContainerFromBase64 } = await import(
            "./mappers/governance-rules"
          );
          return rulesContainerFromBase64(rulesBase64);
        },
        this.rulesCacheTtlMs
      );
    }
    return this._rulesCache;
  }

  /**
   * High-level Wallet service for wallet management operations.
   *
   * Provides methods to list, get, create wallets, and manage wallet attributes.
   *
   * @example
   * ```typescript
   * // List wallets
   * const result = await client.wallets.list({ limit: 50 });
   * for (const wallet of result.items) {
   *   console.log(`${wallet.name}: ${wallet.currency}`);
   * }
   *
   * // Get a wallet by ID
   * const wallet = await client.wallets.get(123);
   * ```
   */
  get wallets(): WalletService {
    this.ensureOpen();
    if (!this._walletService) {
      this._walletService = new WalletService(this.walletsApi);
    }
    return this._walletService;
  }

  /**
   * High-level Address service for address management operations.
   *
   * Provides methods to list, get, create addresses, and manage address attributes.
   * All addresses are automatically verified using the HSM public key from the
   * rules container. Requires SuperAdmin keys to be configured so the rules
   * container can be verified before trusting its contents.
   *
   * @example
   * ```typescript
   * // Get an address with signature verification
   * const address = await client.addresses.get(456);
   * console.log(`Address: ${address.address}`);
   *
   * // List addresses for a wallet
   * const { items, pagination } = await client.addresses.list(123);
   * ```
   */
  get addresses(): AddressService {
    this.ensureOpen();
    if (!this._addressService) {
      // RulesContainerCache is required — address signature verification is mandatory.
      // The rules container must be verified by SuperAdmin keys before trusting the
      // HSM public key it contains, so SuperAdmin keys must be configured.
      if (!this.config.superAdminKeysPem || this.config.superAdminKeysPem.length === 0) {
        throw new ConfigurationError(
          "superAdminKeysPem must be configured to use AddressService — " +
          "the rules container must be verified by SuperAdmin keys before address signature verification"
        );
      }
      this._addressService = new AddressService(this.addressesApi, this.getRulesCache());
    }
    return this._addressService;
  }

  /**
   * High-level AirGap service for air-gapped HSM signing operations.
   *
   * Provides methods to export requests/addresses for cold HSM signing
   * and import signed responses from the cold HSM.
   *
   * @example
   * ```typescript
   * // Export requests for cold HSM signing
   * const payload = await client.airGap.getOutgoingAirGap({
   *   requestIds: ['request-1', 'request-2'],
   * });
   *
   * // Submit signed responses from cold HSM
   * await client.airGap.submitIncomingAirGap({
   *   payload: signedPayloadBase64,
   * });
   * ```
   */
  get airGap(): AirGapService {
    this.ensureOpen();
    if (!this._airGapService) {
      this._airGapService = new AirGapService(this.airGapApi);
    }
    return this._airGapService;
  }

  /**
   * High-level Asset service for querying asset balances.
   *
   * Provides methods to retrieve addresses and wallets that hold a specific asset.
   * This is useful for portfolio management, compliance reporting, and understanding
   * asset distribution across the organization.
   *
   * @example
   * ```typescript
   * // Get all addresses holding ETH
   * const ethAddresses = await client.assets.getAssetAddresses({ currency: 'ETH' });
   *
   * // Get all wallets holding a specific token
   * const usdcWallets = await client.assets.getAssetWallets({ currency: 'USDC' });
   *
   * // Get addresses for a specific asset filtered by wallet
   * const addresses = await client.assets.getAssetAddresses({
   *   currency: 'BTC',
   *   walletId: 'wallet-123',
   * });
   * ```
   */
  get assets(): AssetService {
    this.ensureOpen();
    if (!this._assetService) {
      this._assetService = new AssetService(this.assetsApi);
    }
    return this._assetService;
  }

  /**
   * High-level Request service for transaction request operations.
   *
   * Provides methods to list, get, create, approve, and reject requests.
   * Request hash verification is performed automatically on get() operations.
   *
   * @example
   * ```typescript
   * // Get a request with hash verification
   * const request = await client.requests.get(789);
   *
   * // List requests pending approval
   * const { requests } = await client.requests.listForApproval({ limit: 50 });
   *
   * // Approve requests with ECDSA signature
   * const count = await client.requests.approveRequests(requests, privateKey);
   * ```
   */
  get requests(): RequestService {
    this.ensureOpen();
    if (!this._requestService) {
      this._requestService = new RequestService(this.requestsApi);
    }
    return this._requestService;
  }

  /**
   * High-level Transaction service for transaction operations.
   *
   * Provides methods to list and get transactions.
   */
  get transactions(): TransactionService {
    this.ensureOpen();
    if (!this._transactionService) {
      this._transactionService = new TransactionService(this.transactionsApi);
    }
    return this._transactionService;
  }

  /**
   * High-level Balance service for balance operations.
   *
   * Provides methods to list balances across wallets.
   */
  get balances(): BalanceService {
    this.ensureOpen();
    if (!this._balanceService) {
      this._balanceService = new BalanceService(this.balancesApi);
    }
    return this._balanceService;
  }

  /**
   * High-level Config service for tenant configuration operations.
   *
   * Provides methods to retrieve tenant configuration settings,
   * including security requirements, feature flags, and system parameters.
   *
   * @example
   * ```typescript
   * // Get tenant configuration
   * const tenantConfig = await client.configService.getTenantConfig();
   * console.log(`Base currency: ${tenantConfig.baseCurrency}`);
   * console.log(`MFA required: ${tenantConfig.mfaMandatory}`);
   * ```
   */
  get configService(): ConfigService {
    this.ensureOpen();
    if (!this._configService) {
      this._configService = new ConfigService(this.configApi);
    }
    return this._configService;
  }

  /**
   * High-level Currency service for currency operations.
   *
   * Provides methods to list and get currencies.
   */
  get currencies(): CurrencyService {
    this.ensureOpen();
    if (!this._currencyService) {
      this._currencyService = new CurrencyService(this.currenciesApi);
    }
    return this._currencyService;
  }

  /**
   * High-level Exchange service for exchange account operations.
   *
   * Provides methods to list and get exchange accounts, retrieve counterparty
   * summaries, and calculate withdrawal fees.
   *
   * @example
   * ```typescript
   * // List exchange accounts
   * const result = await client.exchanges.list();
   * for (const exchange of result.items) {
   *   console.log(`${exchange.exchange}: ${exchange.totalBalance}`);
   * }
   *
   * // Get a specific exchange account
   * const exchange = await client.exchanges.get('exchange-123');
   *
   * // Get counterparty summaries
   * const counterparties = await client.exchanges.getCounterparties();
   * ```
   */
  get exchanges(): ExchangeService {
    this.ensureOpen();
    if (!this._exchangeService) {
      this._exchangeService = new ExchangeService(this.exchangeApi);
    }
    return this._exchangeService;
  }

  /**
   * High-level Fee service for retrieving network fee information.
   *
   * Provides methods to get current network fees for various blockchains.
   *
   * @example
   * ```typescript
   * // Get current network fees (v2 - recommended)
   * const fees = await client.fees.getFeesV2();
   * for (const fee of fees) {
   *   console.log(`${fee.currencyInfo?.symbol}: ${fee.value} ${fee.denom}`);
   * }
   * ```
   */
  get fees(): FeeService {
    this.ensureOpen();
    if (!this._feeService) {
      this._feeService = new FeeService(this.feeApi);
    }
    return this._feeService;
  }

  /**
   * High-level Fee Payer service for managing fee payers.
   *
   * Fee payers are accounts used to pay transaction fees on behalf of other
   * addresses, commonly used for sponsored transactions on EVM-compatible
   * blockchains like Ethereum.
   *
   * @example
   * ```typescript
   * // List all fee payers
   * const feePayers = await client.feePayers.list();
   *
   * // List fee payers for a specific blockchain
   * const ethFeePayers = await client.feePayers.list({
   *   blockchain: 'ETH',
   *   network: 'mainnet',
   * });
   *
   * // Get a specific fee payer
   * const feePayer = await client.feePayers.get('fp-123');
   * ```
   */
  get feePayers(): FeePayerService {
    this.ensureOpen();
    if (!this._feePayerService) {
      this._feePayerService = new FeePayerService(this.feePayersApi);
    }
    return this._feePayerService;
  }

  /**
   * High-level Health service for health check operations.
   */
  get health(): HealthService {
    this.ensureOpen();
    if (!this._healthService) {
      this._healthService = new HealthService(this.healthApi);
    }
    return this._healthService;
  }

  /**
   * High-level Job service for monitoring jobs.
   *
   * Jobs are background tasks that process various operations such as
   * transaction monitoring, balance updates, and other async operations.
   *
   * @example
   * ```typescript
   * // List all jobs
   * const jobs = await client.jobs.list();
   * for (const job of jobs) {
   *   console.log(`${job.name}: ${job.statistics?.successes} successes`);
   * }
   *
   * // Get a specific job
   * const job = await client.jobs.get('balance-sync');
   *
   * // Get job execution status
   * const status = await client.jobs.getStatus('balance-sync', 'exec-123');
   * ```
   */
  get jobs(): JobService {
    this.ensureOpen();
    if (!this._jobService) {
      this._jobService = new JobService(this.jobsApi);
    }
    return this._jobService;
  }

  /**
   * High-level Price service for retrieving cryptocurrency prices and conversions.
   *
   * Provides methods to list current prices, get price history, and convert
   * amounts between currencies.
   *
   * @example
   * ```typescript
   * // Get all current prices
   * const prices = await client.prices.list();
   * for (const price of prices) {
   *   console.log(`${price.currencyFrom}/${price.currencyTo}: ${price.rate}`);
   * }
   *
   * // Get price history for a currency pair
   * const history = await client.prices.getHistory({
   *   base: 'ETH',
   *   quote: 'USD',
   *   limit: 100,
   * });
   *
   * // Convert an amount to target currencies
   * const converted = await client.prices.convert({
   *   currency: 'ETH',
   *   amount: '1000000000000000000',
   *   targetCurrencyIds: ['USD', 'BTC'],
   * });
   * ```
   */
  get prices(): PriceService {
    this.ensureOpen();
    if (!this._priceService) {
      this._priceService = new PriceService(this.pricesApi);
    }
    return this._priceService;
  }

  /**
   * High-level User service for user management operations.
   */
  get users(): UserService {
    this.ensureOpen();
    if (!this._userService) {
      this._userService = new UserService(this.usersApi);
    }
    return this._userService;
  }

  /**
   * High-level Group service for group management operations.
   */
  get groups(): GroupService {
    this.ensureOpen();
    if (!this._groupService) {
      this._groupService = new GroupService(this.groupsApi);
    }
    return this._groupService;
  }

  /**
   * High-level Visibility Group service for visibility group management operations.
   *
   * Visibility groups are used to control data access. Users can only see
   * wallets, addresses, and other entities that belong to their assigned
   * visibility groups.
   *
   * @example
   * ```typescript
   * // Get all visibility groups
   * const groups = await client.visibilityGroups.list();
   *
   * // Get users in a specific visibility group
   * const users = await client.visibilityGroups.getUsersByVisibilityGroup('vg-123');
   * ```
   */
  get visibilityGroups(): VisibilityGroupService {
    this.ensureOpen();
    if (!this._visibilityGroupService) {
      this._visibilityGroupService = new VisibilityGroupService(
        this.restrictedVisibilityGroupsApi
      );
    }
    return this._visibilityGroupService;
  }

  /**
   * High-level Tag service for tag management operations.
   */
  get tags(): TagService {
    this.ensureOpen();
    if (!this._tagService) {
      this._tagService = new TagService(this.tagsApi);
    }
    return this._tagService;
  }

  /**
   * High-level Statistics service for portfolio statistics.
   *
   * Provides methods to retrieve aggregated portfolio statistics including
   * total balance, address counts, and wallet counts.
   *
   * @example
   * ```typescript
   * const stats = await client.statistics.getPortfolioStatistics();
   * console.log(`Total wallets: ${stats.walletsCount}`);
   * console.log(`Total addresses: ${stats.addressesCount}`);
   * console.log(`Total balance (base currency): ${stats.totalBalanceBaseCurrency}`);
   * ```
   */
  get statistics(): StatisticsService {
    this.ensureOpen();
    if (!this._statisticsService) {
      this._statisticsService = new StatisticsService(this.statisticsApi);
    }
    return this._statisticsService;
  }

  /**
   * High-level Token Metadata service for retrieving token information.
   */
  get tokenMetadata(): TokenMetadataService {
    this.ensureOpen();
    if (!this._tokenMetadataService) {
      this._tokenMetadataService = new TokenMetadataService(
        this.tokenMetadataApi
      );
    }
    return this._tokenMetadataService;
  }

  /**
   * High-level Webhook service for webhook management operations.
   */
  get webhooks(): WebhookService {
    this.ensureOpen();
    if (!this._webhookService) {
      this._webhookService = new WebhookService(this.webhooksApi);
    }
    return this._webhookService;
  }

  /**
   * High-level Audit service for audit log operations.
   */
  get audits(): AuditService {
    this.ensureOpen();
    if (!this._auditService) {
      this._auditService = new AuditService(this.auditApi);
    }
    return this._auditService;
  }

  /**
   * High-level GovernanceRule service for governance rules operations.
   *
   * Provides methods to get and verify governance rules.
   * When superAdminKeysPem is configured, signature verification is performed.
   */
  get governanceRules(): GovernanceRuleService {
    this.ensureOpen();
    if (!this._governanceRuleService) {
      // Parse SuperAdmin PEM keys if configured
      const superAdminKeys =
        this.config.superAdminKeysPem && this.config.superAdminKeysPem.length > 0
          ? decodePublicKeysPem(this.config.superAdminKeysPem)
          : [];

      this._governanceRuleService = new GovernanceRuleService(
        this.governanceRulesApi,
        {
          superAdminKeys,
          minValidSignatures: this.config.minValidSignatures ?? 1,
        }
      );
    }
    return this._governanceRuleService;
  }

  /**
   * High-level WhitelistedAddress service for whitelisted address operations.
   *
   * Provides methods to list and get whitelisted addresses with full
   * integrity verification using SuperAdmin keys.
   */
  get whitelistedAddresses(): WhitelistedAddressService {
    this.ensureOpen();
    if (!this._whitelistedAddressService) {
      this._whitelistedAddressService = WhitelistedAddressService.withVerification(
        this.addressWhitelistingApi,
        {
          superAdminKeysPem: this.config.superAdminKeysPem,
          minValidSignatures: this.config.minValidSignatures ?? 1,
          rulesContainerDecoder: rulesContainerFromBase64,
          userSignaturesDecoder: userSignaturesFromBase64,
        }
      );
    }
    return this._whitelistedAddressService;
  }

  /**
   * High-level WhitelistedAsset service for whitelisted asset/contract operations.
   *
   * Provides methods to list and get whitelisted assets with full
   * integrity verification using SuperAdmin keys.
   */
  get whitelistedAssets(): WhitelistedAssetService {
    this.ensureOpen();
    if (!this._whitelistedAssetService) {
      this._whitelistedAssetService = WhitelistedAssetService.withVerification(
        this.contractWhitelistingApi,
        {
          superAdminKeysPem: this.config.superAdminKeysPem,
          minValidSignatures: this.config.minValidSignatures ?? 1,
          rulesContainerDecoder: rulesContainerFromBase64,
          userSignaturesDecoder: userSignaturesFromBase64,
        }
      );
    }
    return this._whitelistedAssetService;
  }

  /**
   * High-level Action service for action operations.
   */
  get actions(): ActionService {
    this.ensureOpen();
    if (!this._actionService) {
      this._actionService = new ActionService(this.actionsApi);
    }
    return this._actionService;
  }

  /**
   * High-level Blockchain service for blockchain information.
   */
  get blockchains(): BlockchainService {
    this.ensureOpen();
    if (!this._blockchainService) {
      this._blockchainService = new BlockchainService(this.blockchainApi);
    }
    return this._blockchainService;
  }

  /**
   * High-level Business Rule service for business rule operations.
   */
  get businessRules(): BusinessRuleService {
    this.ensureOpen();
    if (!this._businessRuleService) {
      this._businessRuleService = new BusinessRuleService(this.businessRulesApi);
    }
    return this._businessRuleService;
  }

  /**
   * High-level Change service for change log operations.
   */
  get changes(): ChangeService {
    this.ensureOpen();
    if (!this._changeService) {
      this._changeService = new ChangeService(this.changesApi);
    }
    return this._changeService;
  }

  /**
   * High-level Contract Whitelisting service for whitelisted contract operations.
   */
  get contractWhitelisting(): ContractWhitelistingService {
    this.ensureOpen();
    if (!this._contractWhitelistingService) {
      this._contractWhitelistingService = new ContractWhitelistingService(
        this.contractWhitelistingApi
      );
    }
    return this._contractWhitelistingService;
  }

  /**
   * High-level Fiat service for fiat account operations.
   */
  get fiatAccounts(): FiatService {
    this.ensureOpen();
    if (!this._fiatService) {
      this._fiatService = new FiatService(this.fiatApi);
    }
    return this._fiatService;
  }

  /**
   * High-level Multi-Factor Signature service for MFA operations.
   */
  get multiFactorSignature(): MultiFactorSignatureService {
    this.ensureOpen();
    if (!this._multiFactorSignatureService) {
      this._multiFactorSignatureService = new MultiFactorSignatureService(
        this.multiFactorSignatureApi
      );
    }
    return this._multiFactorSignatureService;
  }

  /**
   * High-level Reservation service for UTXO reservation operations.
   */
  get reservations(): ReservationService {
    this.ensureOpen();
    if (!this._reservationService) {
      this._reservationService = new ReservationService(this.reservationsApi);
    }
    return this._reservationService;
  }

  /**
   * High-level Score service for address scoring operations.
   */
  get scores(): ScoreService {
    this.ensureOpen();
    if (!this._scoreService) {
      this._scoreService = new ScoreService(this.scoresApi);
    }
    return this._scoreService;
  }

  /**
   * High-level Staking service for staking operations.
   */
  get staking(): StakingService {
    this.ensureOpen();
    if (!this._stakingService) {
      this._stakingService = new StakingService(this.stakingApi);
    }
    return this._stakingService;
  }

  /**
   * High-level User Device service for user device management.
   */
  get userDevices(): UserDeviceService {
    this.ensureOpen();
    if (!this._userDeviceService) {
      this._userDeviceService = new UserDeviceService(this.userDeviceApi);
    }
    return this._userDeviceService;
  }

  /**
   * High-level Webhook Call service for webhook call operations.
   */
  get webhookCalls(): WebhookCallService {
    this.ensureOpen();
    if (!this._webhookCallService) {
      this._webhookCallService = new WebhookCallService(this.webhookCallsApi);
    }
    return this._webhookCallService;
  }
}
