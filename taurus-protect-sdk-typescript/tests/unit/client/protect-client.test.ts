/**
 * Unit tests for ProtectClient.
 *
 * These tests verify:
 * - Client creation and configuration validation
 * - Lazy initialization of APIs and services
 * - Client lifecycle (open/close)
 * - TaurusNetwork namespace functionality
 * - Configuration accessors
 */

import { ProtectClient, ProtectClientConfig } from "../../../src/client";
import { ConfigurationError } from "../../../src/errors";

// A valid P-256 public key for testing (test key with no production value)
const TEST_SUPER_ADMIN_KEY_PEM = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvW
L+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==
-----END PUBLIC KEY-----`;

// Helper to create valid config
function createValidConfig(
  overrides?: Partial<ProtectClientConfig>
): ProtectClientConfig {
  return {
    host: "https://protect.example.com",
    apiKey: "test-api-key-12345",
    apiSecret: "aabbccdd11223344556677889900aabbccddeeff",
    superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
    ...overrides,
  };
}

describe("ProtectClient", () => {
  let client: ProtectClient;

  afterEach(() => {
    if (client && !client.isClosed) {
      client.close();
    }
  });

  // ===== Client Creation Tests =====

  describe("create", () => {
    describe("valid configuration", () => {
      it("should create client with valid configuration", () => {
        client = ProtectClient.create(createValidConfig());
        expect(client).toBeDefined();
        expect(client.isClosed).toBe(false);
      });

      it("should accept host with trailing slash", () => {
        client = ProtectClient.create(
          createValidConfig({ host: "https://example.com/" })
        );
        expect(client.host).toBe("https://example.com/");
      });

      it("should accept host with multiple trailing slashes", () => {
        client = ProtectClient.create(
          createValidConfig({ host: "https://example.com///" })
        );
        expect(client).toBeDefined();
      });

      it("should accept http URLs", () => {
        client = ProtectClient.create(
          createValidConfig({ host: "http://localhost:8080" })
        );
        expect(client.host).toBe("http://localhost:8080");
      });

      it("should accept URLs with ports", () => {
        client = ProtectClient.create(
          createValidConfig({ host: "https://protect.example.com:443" })
        );
        expect(client).toBeDefined();
      });

      it("should accept URLs with paths", () => {
        client = ProtectClient.create(
          createValidConfig({ host: "https://protect.example.com/api/v1" })
        );
        expect(client).toBeDefined();
      });

      it("should reject empty superAdminKeysPem array", () => {
        expect(() =>
          ProtectClient.create(
            createValidConfig({
              superAdminKeysPem: [],
            })
          )
        ).toThrow(ConfigurationError);
        expect(() =>
          ProtectClient.create(
            createValidConfig({
              superAdminKeysPem: [],
            })
          )
        ).toThrow("superAdminKeysPem is required");
      });

      it("should accept superAdminKeysPem with valid keys", () => {
        client = ProtectClient.create(
          createValidConfig({
            superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
          })
        );
        expect(client).toBeDefined();
      });

      it("should accept minValidSignatures configuration", () => {
        client = ProtectClient.create(
          createValidConfig({ minValidSignatures: 3 })
        );
        expect(client.minValidSignatures).toBe(3);
      });

      it("should accept rulesCacheTtlMs configuration", () => {
        client = ProtectClient.create(
          createValidConfig({ rulesCacheTtlMs: 600000 })
        );
        expect(client.rulesCacheTtlMs).toBe(600000);
      });

      it("should accept timeout configuration", () => {
        client = ProtectClient.create(createValidConfig({ timeout: 60000 }));
        expect(client).toBeDefined();
      });

      it("should accept lowercase hex in apiSecret", () => {
        client = ProtectClient.create(
          createValidConfig({ apiSecret: "aabbccdd11223344" })
        );
        expect(client).toBeDefined();
      });

      it("should accept uppercase hex in apiSecret", () => {
        client = ProtectClient.create(
          createValidConfig({ apiSecret: "AABBCCDD11223344" })
        );
        expect(client).toBeDefined();
      });

      it("should accept mixed case hex in apiSecret", () => {
        client = ProtectClient.create(
          createValidConfig({ apiSecret: "AaBbCcDd11223344" })
        );
        expect(client).toBeDefined();
      });
    });

    describe("host validation", () => {
      it("should throw ConfigurationError when host is empty", () => {
        expect(() => ProtectClient.create(createValidConfig({ host: "" }))).toThrow(
          ConfigurationError
        );
      });

      it("should throw ConfigurationError when host is missing", () => {
        const config = {
          apiKey: "key",
          apiSecret: "aabbccdd",
        } as ProtectClientConfig;
        expect(() => ProtectClient.create(config)).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError when host URL is invalid", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ host: "not-a-valid-url" }))
        ).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError with message for invalid URL", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ host: "invalid" }))
        ).toThrow("Invalid host URL");
      });

      it("should throw ConfigurationError for whitespace-only host", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ host: "   " }))
        ).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError for host with only protocol", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ host: "https://" }))
        ).toThrow(ConfigurationError);
      });
    });

    describe("apiKey validation", () => {
      it("should throw ConfigurationError when apiKey is empty", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ apiKey: "" }))
        ).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError when apiKey is missing", () => {
        const config = {
          host: "https://example.com",
          apiSecret: "aabbccdd",
        } as ProtectClientConfig;
        expect(() => ProtectClient.create(config)).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError with message for empty apiKey", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ apiKey: "" }))
        ).toThrow("apiKey is required");
      });
    });

    describe("apiSecret validation", () => {
      it("should throw ConfigurationError when apiSecret is empty", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ apiSecret: "" }))
        ).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError when apiSecret is missing", () => {
        const config = {
          host: "https://example.com",
          apiKey: "key",
        } as ProtectClientConfig;
        expect(() => ProtectClient.create(config)).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError when apiSecret is not valid hex", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ apiSecret: "not-hex-string!" }))
        ).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError for apiSecret with spaces", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ apiSecret: "aa bb cc dd" }))
        ).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError for apiSecret with special characters", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ apiSecret: "aabbccdd!" }))
        ).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError with message for non-hex apiSecret", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ apiSecret: "xyz123" }))
        ).toThrow("apiSecret must be a valid hexadecimal string");
      });

      it("should throw ConfigurationError for apiSecret with only letters g-z", () => {
        expect(() =>
          ProtectClient.create(createValidConfig({ apiSecret: "ghijklmnop" }))
        ).toThrow(ConfigurationError);
      });
    });
  });

  // ===== Configuration Accessors Tests =====

  describe("configuration accessors", () => {
    describe("host", () => {
      it("should return configured host", () => {
        client = ProtectClient.create(createValidConfig());
        expect(client.host).toBe("https://protect.example.com");
      });

      it("should return host with path if provided", () => {
        client = ProtectClient.create(
          createValidConfig({ host: "https://protect.example.com/api" })
        );
        expect(client.host).toBe("https://protect.example.com/api");
      });
    });

    describe("superAdminKeysPem", () => {
      it("should return configured keys", () => {
        client = ProtectClient.create(createValidConfig());
        expect(client.superAdminKeysPem).toEqual([TEST_SUPER_ADMIN_KEY_PEM]);
      });

      it("should return all configured keys", () => {
        const keys = [TEST_SUPER_ADMIN_KEY_PEM, TEST_SUPER_ADMIN_KEY_PEM];
        client = ProtectClient.create(
          createValidConfig({ superAdminKeysPem: keys })
        );
        expect(client.superAdminKeysPem).toEqual(keys);
      });
    });

    describe("minValidSignatures", () => {
      it("should return configured minValidSignatures", () => {
        client = ProtectClient.create(
          createValidConfig({ minValidSignatures: 3 })
        );
        expect(client.minValidSignatures).toBe(3);
      });

      it("should return default minValidSignatures when not set", () => {
        client = ProtectClient.create(createValidConfig());
        expect(client.minValidSignatures).toBe(1);
      });

      it("should throw when explicitly set to 0", () => {
        expect(() =>
          ProtectClient.create(
            createValidConfig({ minValidSignatures: 0 })
          )
        ).toThrow("minValidSignatures must be greater than 0 when specified");
      });
    });

    describe("rulesCacheTtlMs", () => {
      it("should return configured rulesCacheTtlMs", () => {
        client = ProtectClient.create(
          createValidConfig({ rulesCacheTtlMs: 600000 })
        );
        expect(client.rulesCacheTtlMs).toBe(600000);
      });

      it("should return default rulesCacheTtlMs when not set", () => {
        client = ProtectClient.create(createValidConfig());
        expect(client.rulesCacheTtlMs).toBe(300000);
      });

      it("should return 0 when explicitly set to 0", () => {
        client = ProtectClient.create(createValidConfig({ rulesCacheTtlMs: 0 }));
        expect(client.rulesCacheTtlMs).toBe(0);
      });
    });
  });

  // ===== Low-Level API Accessors Tests =====

  describe("low-level API accessors", () => {
    const apiGetters = [
      "actionsApi",
      "addressWhitelistingApi",
      "addressesApi",
      "airGapApi",
      "assetsApi",
      "auditApi",
      "authenticationApi",
      "authenticationHMACApi",
      "authenticationOIDCApi",
      "authenticationSAMLApi",
      "balancesApi",
      "blockchainApi",
      "businessRulesApi",
      "changesApi",
      "configApi",
      "contractWhitelistingApi",
      "currenciesApi",
      "exchangeApi",
      "feeApi",
      "feePayersApi",
      "fiatApi",
      "governanceRulesApi",
      "groupsApi",
      "healthApi",
      "jobsApi",
      "multiFactorSignatureApi",
      "pricesApi",
      "requestsApi",
      "requestsADAApi",
      "requestsALGOApi",
      "requestsContractsApi",
      "requestsCosmosApi",
      "requestsDOTApi",
      "requestsFTMApi",
      "requestsHederaApi",
      "requestsICPApi",
      "requestsMinaApi",
      "requestsNEARApi",
      "requestsSOLApi",
      "requestsXLMApi",
      "requestsXTZApi",
      "reservationsApi",
      "restrictedVisibilityGroupsApi",
      "scimApi",
      "scoresApi",
      "stakingApi",
      "statisticsApi",
      "stewardApi",
      "tagsApi",
      "tokenMetadataApi",
      "transactionsApi",
      "userDeviceApi",
      "usersApi",
      "walletsApi",
      "webhookCallsApi",
      "webhooksApi",
    ];

    beforeEach(() => {
      client = ProtectClient.create(createValidConfig());
    });

    test.each(apiGetters)(
      "%s should be lazily initialized and return same instance",
      (apiName) => {
        const api1 = (client as unknown as Record<string, unknown>)[apiName];
        expect(api1).toBeDefined();
        const api2 = (client as unknown as Record<string, unknown>)[apiName];
        expect(api2).toBe(api1); // Same instance
      }
    );

    test.each(apiGetters)("%s should throw after client closed", (apiName) => {
      client.close();
      expect(() => (client as unknown as Record<string, unknown>)[apiName]).toThrow(
        ConfigurationError
      );
    });

    it("should have 56 API getters available", () => {
      expect(apiGetters.length).toBe(56);
    });
  });

  // ===== High-Level Service Accessors Tests =====

  describe("high-level service accessors", () => {
    // List of service getters actually implemented on ProtectClient
    const serviceGetters = [
      "wallets",
      "addresses",
      "requests",
      "transactions",
      "balances",
      "currencies",
      "health",
      "jobs",
      "users",
      "groups",
      "visibilityGroups",
      "tags",
      "webhooks",
      "audits",
      "governanceRules",
      "whitelistedAddresses",
      "whitelistedAssets",
      "statistics",
      "configService",
      "assets",
      "prices",
      "fees",
      "feePayers",
      "exchanges",
      "airGap",
      "tokenMetadata",
    ];

    beforeEach(() => {
      client = ProtectClient.create(createValidConfig());
    });

    test.each(serviceGetters)(
      "%s should be lazily initialized and return same instance",
      (serviceName) => {
        const svc1 = (client as unknown as Record<string, unknown>)[serviceName];
        expect(svc1).toBeDefined();
        const svc2 = (client as unknown as Record<string, unknown>)[serviceName];
        expect(svc2).toBe(svc1); // Same instance
      }
    );

    test.each(serviceGetters)(
      "%s should throw after client closed",
      (serviceName) => {
        client.close();
        expect(
          () => (client as unknown as Record<string, unknown>)[serviceName]
        ).toThrow(ConfigurationError);
      }
    );

    it("should have the expected number of service getters", () => {
      // Note: staking is only exposed as stakingApi, not as high-level service
      expect(serviceGetters.length).toBe(26);
    });
  });

  // ===== TaurusNetwork Namespace Tests =====

  describe("taurusNetwork namespace", () => {
    beforeEach(() => {
      client = ProtectClient.create(createValidConfig());
    });

    describe("namespace access", () => {
      it("should return TaurusNetworkNamespace", () => {
        expect(client.taurusNetwork).toBeDefined();
      });

      it("should return same TaurusNetworkNamespace instance", () => {
        const tn1 = client.taurusNetwork;
        const tn2 = client.taurusNetwork;
        expect(tn2).toBe(tn1);
      });

      it("should throw after client closed", () => {
        client.close();
        expect(() => client.taurusNetwork).toThrow(ConfigurationError);
      });

      it("should throw with descriptive message after client closed", () => {
        client.close();
        expect(() => client.taurusNetwork).toThrow("Client has been closed");
      });
    });

    describe("lendingApi", () => {
      it("should lazily initialize lendingApi", () => {
        const api1 = client.taurusNetwork.lendingApi;
        const api2 = client.taurusNetwork.lendingApi;
        expect(api1).toBeDefined();
        expect(api2).toBe(api1);
      });

      it("should throw after client closed", () => {
        const taurusNetwork = client.taurusNetwork;
        client.close();
        expect(() => taurusNetwork.lendingApi).toThrow(ConfigurationError);
      });
    });

    describe("participantApi", () => {
      it("should lazily initialize participantApi", () => {
        const api1 = client.taurusNetwork.participantApi;
        const api2 = client.taurusNetwork.participantApi;
        expect(api1).toBeDefined();
        expect(api2).toBe(api1);
      });

      it("should throw after client closed", () => {
        const taurusNetwork = client.taurusNetwork;
        client.close();
        expect(() => taurusNetwork.participantApi).toThrow(ConfigurationError);
      });
    });

    describe("pledgeApi", () => {
      it("should lazily initialize pledgeApi", () => {
        const api1 = client.taurusNetwork.pledgeApi;
        const api2 = client.taurusNetwork.pledgeApi;
        expect(api1).toBeDefined();
        expect(api2).toBe(api1);
      });

      it("should throw after client closed", () => {
        const taurusNetwork = client.taurusNetwork;
        client.close();
        expect(() => taurusNetwork.pledgeApi).toThrow(ConfigurationError);
      });
    });

    describe("settlementApi", () => {
      it("should lazily initialize settlementApi", () => {
        const api1 = client.taurusNetwork.settlementApi;
        const api2 = client.taurusNetwork.settlementApi;
        expect(api1).toBeDefined();
        expect(api2).toBe(api1);
      });

      it("should throw after client closed", () => {
        const taurusNetwork = client.taurusNetwork;
        client.close();
        expect(() => taurusNetwork.settlementApi).toThrow(ConfigurationError);
      });
    });

    describe("sharedAddressAssetApi", () => {
      it("should lazily initialize sharedAddressAssetApi", () => {
        const api1 = client.taurusNetwork.sharedAddressAssetApi;
        const api2 = client.taurusNetwork.sharedAddressAssetApi;
        expect(api1).toBeDefined();
        expect(api2).toBe(api1);
      });

      it("should throw after client closed", () => {
        const taurusNetwork = client.taurusNetwork;
        client.close();
        expect(() => taurusNetwork.sharedAddressAssetApi).toThrow(
          ConfigurationError
        );
      });
    });
  });

  // ===== Client Lifecycle Tests =====

  describe("client lifecycle", () => {
    describe("isClosed property", () => {
      it("should have isClosed = false initially", () => {
        client = ProtectClient.create(createValidConfig());
        expect(client.isClosed).toBe(false);
      });

      it("should set isClosed = true after close()", () => {
        client = ProtectClient.create(createValidConfig());
        client.close();
        expect(client.isClosed).toBe(true);
      });
    });

    describe("close method", () => {
      it("should be idempotent (multiple close calls safe)", () => {
        client = ProtectClient.create(createValidConfig());
        expect(() => {
          client.close();
          client.close();
          client.close();
        }).not.toThrow();
        expect(client.isClosed).toBe(true);
      });

      it("should clear cached API instances on close", () => {
        client = ProtectClient.create(createValidConfig());
        const api = client.walletsApi; // Initialize
        expect(api).toBeDefined();
        client.close();
        expect(() => client.walletsApi).toThrow(ConfigurationError);
      });

      it("should clear cached service instances on close", () => {
        client = ProtectClient.create(createValidConfig());
        const svc = client.wallets; // Initialize
        expect(svc).toBeDefined();
        client.close();
        expect(() => client.wallets).toThrow(ConfigurationError);
      });

      it("should clear TaurusNetwork namespace on close", () => {
        client = ProtectClient.create(createValidConfig());
        const tn = client.taurusNetwork; // Initialize
        expect(tn).toBeDefined();
        client.close();
        expect(() => client.taurusNetwork).toThrow(ConfigurationError);
      });

      it("should clear TaurusNetwork APIs after close", () => {
        client = ProtectClient.create(createValidConfig());
        const tn = client.taurusNetwork;
        const api = tn.lendingApi; // Initialize
        expect(api).toBeDefined();
        client.close();
        expect(() => tn.lendingApi).toThrow(ConfigurationError);
      });
    });

    describe("access after close", () => {
      it("should throw ConfigurationError when accessing APIs after close", () => {
        client = ProtectClient.create(createValidConfig());
        client.close();
        expect(() => client.walletsApi).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError when accessing services after close", () => {
        client = ProtectClient.create(createValidConfig());
        client.close();
        expect(() => client.wallets).toThrow(ConfigurationError);
      });

      it("should throw ConfigurationError when accessing taurusNetwork after close", () => {
        client = ProtectClient.create(createValidConfig());
        client.close();
        expect(() => client.taurusNetwork).toThrow(ConfigurationError);
      });

      it("should throw with descriptive error message when accessing after close", () => {
        client = ProtectClient.create(createValidConfig());
        client.close();
        expect(() => client.walletsApi).toThrow("Client has been closed");
      });
    });
  });

  // ===== ensureOpen Behavior Tests =====

  describe("ensureOpen behavior", () => {
    beforeEach(() => {
      client = ProtectClient.create(createValidConfig());
    });

    it("should allow API access on open client", () => {
      expect(() => client.walletsApi).not.toThrow();
    });

    it("should allow service access on open client", () => {
      expect(() => client.wallets).not.toThrow();
    });

    it("should allow taurusNetwork access on open client", () => {
      expect(() => client.taurusNetwork).not.toThrow();
    });

    it("should allow multiple API accesses on open client", () => {
      expect(() => {
        client.walletsApi;
        client.addressesApi;
        client.requestsApi;
      }).not.toThrow();
    });

    it("should allow multiple service accesses on open client", () => {
      expect(() => {
        client.wallets;
        client.addresses;
        client.requests;
      }).not.toThrow();
    });
  });

  // ===== Middleware Configuration Tests =====

  describe("middleware configuration", () => {
    it("should accept empty middleware array", () => {
      client = ProtectClient.create(createValidConfig({ middleware: [] }));
      expect(client).toBeDefined();
    });

    it("should accept custom middleware", () => {
      const customMiddleware = {
        pre: async (context: { init: RequestInit; url: string }) => context,
        post: async (context: { init: RequestInit; url: string; response: Response }) => context.response,
      };
      client = ProtectClient.create(
        createValidConfig({ middleware: [customMiddleware] })
      );
      expect(client).toBeDefined();
    });
  });

  // ===== API Instance Type Tests =====

  describe("API instance types", () => {
    beforeEach(() => {
      client = ProtectClient.create(createValidConfig());
    });

    it("should return WalletsApi instance with expected methods", () => {
      const api = client.walletsApi;
      expect(typeof api.walletServiceGetWalletsV2).toBe("function");
    });

    it("should return AddressesApi instance with expected methods", () => {
      const api = client.addressesApi;
      expect(typeof api.walletServiceGetAddresses).toBe("function");
    });

    it("should return RequestsApi instance with expected methods", () => {
      const api = client.requestsApi;
      expect(typeof api.requestServiceGetRequest).toBe("function");
      expect(typeof api.requestServiceApproveRequests).toBe("function");
    });

    it("should return TransactionsApi instance with expected methods", () => {
      const api = client.transactionsApi;
      expect(typeof api.transactionServiceGetTransactions).toBe("function");
    });

    it("should return HealthApi instance with expected methods", () => {
      const api = client.healthApi;
      expect(typeof api.healthServiceGetHealthChecks).toBe("function");
    });
  });

  // ===== Service Instance Type Tests =====

  describe("service instance types", () => {
    beforeEach(() => {
      client = ProtectClient.create(createValidConfig());
    });

    it("should return WalletService instance with expected methods", () => {
      const svc = client.wallets;
      expect(typeof svc.list).toBe("function");
      expect(typeof svc.get).toBe("function");
    });

    it("should return AddressService instance with expected methods", () => {
      const svc = client.addresses;
      expect(typeof svc.list).toBe("function");
      expect(typeof svc.get).toBe("function");
    });

    it("should return RequestService instance with expected methods", () => {
      const svc = client.requests;
      expect(typeof svc.list).toBe("function");
      expect(typeof svc.get).toBe("function");
      expect(typeof svc.approveRequest).toBe("function");
      expect(typeof svc.rejectRequest).toBe("function");
    });

    it("should return TransactionService instance with expected methods", () => {
      const svc = client.transactions;
      expect(typeof svc.list).toBe("function");
      expect(typeof svc.get).toBe("function");
    });

    it("should return HealthService instance with expected methods", () => {
      const svc = client.health;
      expect(typeof svc.check).toBe("function");
    });

    it("should return GovernanceRuleService instance with expected methods", () => {
      const svc = client.governanceRules;
      expect(typeof svc.get).toBe("function");
    });
  });

  // ===== Multiple Client Instances Tests =====

  describe("multiple client instances", () => {
    let client2: ProtectClient;

    afterEach(() => {
      if (client2 && !client2.isClosed) {
        client2.close();
      }
    });

    it("should create independent client instances", () => {
      client = ProtectClient.create(createValidConfig());
      client2 = ProtectClient.create(
        createValidConfig({ host: "https://other.example.com" })
      );

      expect(client.host).toBe("https://protect.example.com");
      expect(client2.host).toBe("https://other.example.com");
    });

    it("should have independent lifecycle", () => {
      client = ProtectClient.create(createValidConfig());
      client2 = ProtectClient.create(createValidConfig());

      client.close();

      expect(client.isClosed).toBe(true);
      expect(client2.isClosed).toBe(false);
    });

    it("should have independent API instances", () => {
      client = ProtectClient.create(createValidConfig());
      client2 = ProtectClient.create(createValidConfig());

      const api1 = client.walletsApi;
      const api2 = client2.walletsApi;

      expect(api1).not.toBe(api2);
    });

    it("should have independent service instances", () => {
      client = ProtectClient.create(createValidConfig());
      client2 = ProtectClient.create(createValidConfig());

      const svc1 = client.wallets;
      const svc2 = client2.wallets;

      expect(svc1).not.toBe(svc2);
    });
  });

  // ===== Edge Cases Tests =====

  describe("edge cases", () => {
    it("should handle very long apiKey", () => {
      const longKey = "a".repeat(1000);
      client = ProtectClient.create(createValidConfig({ apiKey: longKey }));
      expect(client).toBeDefined();
    });

    it("should handle very long apiSecret (valid hex)", () => {
      const longSecret = "a".repeat(1000);
      client = ProtectClient.create(createValidConfig({ apiSecret: longSecret }));
      expect(client).toBeDefined();
    });

    it("should handle apiSecret with only digits", () => {
      client = ProtectClient.create(
        createValidConfig({ apiSecret: "1234567890" })
      );
      expect(client).toBeDefined();
    });

    it("should handle minValidSignatures of large value", () => {
      client = ProtectClient.create(
        createValidConfig({ minValidSignatures: 1000000 })
      );
      expect(client.minValidSignatures).toBe(1000000);
    });

    it("should handle rulesCacheTtlMs of large value", () => {
      client = ProtectClient.create(
        createValidConfig({ rulesCacheTtlMs: 999999999 })
      );
      expect(client.rulesCacheTtlMs).toBe(999999999);
    });
  });

  // ===== Address Service Rules Cache Tests =====

  describe("addresses service with mandatory rules cache", () => {
    it("should create AddressService with rules cache when superAdminKeysPem provided", () => {
      client = ProtectClient.create(
        createValidConfig({
          superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
        })
      );
      const svc = client.addresses;
      expect(svc).toBeDefined();
    });

    it("should throw ConfigurationError when creating client without superAdminKeysPem", () => {
      expect(
        () => ProtectClient.create(createValidConfig({ superAdminKeysPem: undefined }))
      ).toThrow(ConfigurationError);
      expect(
        () => ProtectClient.create(createValidConfig({ superAdminKeysPem: [] }))
      ).toThrow(ConfigurationError);
    });
  });
});
