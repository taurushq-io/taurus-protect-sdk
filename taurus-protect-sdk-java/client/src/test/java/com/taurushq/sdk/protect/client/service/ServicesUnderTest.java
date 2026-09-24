package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.testutil.StubTransport;
import com.taurushq.sdk.protect.openapi.ApiClient;

import java.security.GeneralSecurityException;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.Collections;
import java.util.List;

/**
 * Real services wired to one {@link StubTransport}, so a test drives the generated client
 * end to end. A reply without rows needs no rules container, so nothing here is verified
 * against a real key; tests that need verification build their own fixtures.
 */
final class ServicesUnderTest {

    private static final List<PublicKey> SUPER_ADMIN_KEYS = Collections.singletonList(newKey());

    private final ApiClient client;
    private final ApiExceptionMapper mapper = new ApiExceptionMapper();

    ServicesUnderTest(final StubTransport stub) {
        this.client = stub.client();
    }

    private static PublicKey newKey() {
        try {
            KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
            generator.initialize(new ECGenParameterSpec("secp256r1"));
            return generator.generateKeyPair().getPublic();
        } catch (GeneralSecurityException e) {
            throw new IllegalStateException(e);
        }
    }

    private RulesContainerCache cache() {
        return new RulesContainerCache(governance());
    }

    GovernanceRuleService governance() {
        return new GovernanceRuleService(client, mapper, SUPER_ADMIN_KEYS, 1);
    }

    WalletService wallets() {
        return new WalletService(client, mapper);
    }

    AddressService addresses() {
        return new AddressService(client, mapper, cache());
    }

    TransactionService transactions() {
        return new TransactionService(client, mapper);
    }

    UserService users() {
        return new UserService(client, mapper);
    }

    GroupService groups() {
        return new GroupService(client, mapper);
    }

    FeePayerService feePayers() {
        return new FeePayerService(client, mapper);
    }

    ActionService actions() {
        return new ActionService(client, mapper);
    }

    WhitelistedAddressService whitelistedAddresses() {
        return new WhitelistedAddressService(client, mapper, SUPER_ADMIN_KEYS, 1);
    }

    WhitelistedAssetService whitelistedAssets() {
        return new WhitelistedAssetService(client, mapper, SUPER_ADMIN_KEYS, 1);
    }

    RequestService requests() {
        return new RequestService(client, mapper);
    }

    ChangeService changes() {
        return new ChangeService(client, mapper);
    }

    AuditService audit() {
        return new AuditService(client, mapper);
    }

    BusinessRuleService businessRules() {
        return new BusinessRuleService(client, mapper);
    }

    StakingService staking() {
        return new StakingService(client, mapper);
    }

    BalanceService balances() {
        return new BalanceService(client, mapper);
    }

    AssetService assets() {
        RulesContainerCache cache = cache();
        return new AssetService(client, mapper, cache, new AddressService(client, mapper, cache),
                whitelistedAddresses());
    }

    ReservationService reservations() {
        return new ReservationService(client, mapper);
    }

    FiatService fiat() {
        return new FiatService(client, mapper);
    }

    WebhookService webhooks() {
        return new WebhookService(client, mapper);
    }

    WebhookCallsService webhookCalls() {
        return new WebhookCallsService(client, mapper);
    }

    PriceService prices() {
        return new PriceService(client, mapper, cache());
    }

    EarnService earn() {
        return new EarnService(client, mapper);
    }

    FeeService fees() {
        return new FeeService(client, mapper);
    }

    TokenMetadataService tokenMetadata() {
        return new TokenMetadataService(client, mapper);
    }

    TaurusNetworkLendingService lending() {
        return new TaurusNetworkLendingService(client, mapper);
    }

    TaurusNetworkPledgeService pledges() {
        return new TaurusNetworkPledgeService(client, mapper);
    }

    TaurusNetworkSettlementService settlements() {
        return new TaurusNetworkSettlementService(client, mapper);
    }

    TaurusNetworkSharingService sharing() {
        return new TaurusNetworkSharingService(client, mapper);
    }
}
