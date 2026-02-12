package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.BusinessRule;
import com.taurushq.sdk.protect.client.model.BusinessRuleResult;
import com.taurushq.sdk.protect.client.model.ExchangeCounterparty;
import com.taurushq.sdk.protect.client.model.FeePayer;
import com.taurushq.sdk.protect.client.model.FiatProvider;
import com.taurushq.sdk.protect.client.model.Job;
import com.taurushq.sdk.protect.client.model.Webhook;
import com.taurushq.sdk.protect.client.model.WebhookCallResult;
import com.taurushq.sdk.protect.client.model.WebhookResult;
import com.taurushq.sdk.protect.client.model.WhitelistedContractAddressResult;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingAgreementResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingOfferResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Participant;
import com.taurushq.sdk.protect.client.model.taurusnetwork.PledgeResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.SettlementResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.SharedAddressResult;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for extended services that previously lacked test coverage.
 *
 * <p>Tests read-only operations for: BusinessRuleService, WebhookService,
 * WebhookCallsService, FeePayerService, ExchangeService, FiatService,
 * JobService, ContractWhitelistingService, and TaurusNetwork
 * (participants, pledges, lending, settlements, sharing).
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class ExtendedServicesIntegrationTest {

    private ProtectClient client;

    @BeforeAll
    void setup() throws Exception {
        TestHelper.skipIfNotEnabled();
        client = TestHelper.getTestClient(1);
    }

    @AfterAll
    void teardown() {
        if (client != null) {
            client.close();
        }
    }

    // =========================================================================
    // BusinessRuleService
    // =========================================================================

    @Test
    void listBusinessRules() throws ApiException {
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, 50);
        BusinessRuleResult result = client.getBusinessRuleService().getBusinessRules(cursor);

        assertNotNull(result);
        List<BusinessRule> rules = result.getRules();
        assertNotNull(rules);
        System.out.println("Found " + rules.size() + " business rules");
        for (BusinessRule rule : rules.subList(0, Math.min(5, rules.size()))) {
            System.out.println("  Rule: ID=" + rule.getId());
        }
    }

    // =========================================================================
    // WebhookService
    // =========================================================================

    @Test
    void listWebhooks() throws ApiException {
        WebhookResult result = client.getWebhookService().getWebhooks(null, null, null);

        assertNotNull(result);
        List<Webhook> webhooks = result.getWebhooks();
        assertNotNull(webhooks);
        System.out.println("Found " + webhooks.size() + " webhooks");
        for (Webhook wh : webhooks.subList(0, Math.min(5, webhooks.size()))) {
            System.out.println("  Webhook: ID=" + wh.getId() + ", URL=" + wh.getUrl());
        }
    }

    // =========================================================================
    // WebhookCallsService
    // =========================================================================

    @Test
    void listWebhookCalls() throws ApiException {
        WebhookCallResult result = client.getWebhookCallsService()
                .getWebhookCalls(null, null, null, null, null);

        assertNotNull(result);
        System.out.println("Found " + result.getCalls().size() + " webhook calls");
    }

    // =========================================================================
    // FeePayerService
    // =========================================================================

    @Test
    void listFeePayers() throws ApiException {
        List<FeePayer> feePayers = client.getFeePayerService().getFeePayers();

        assertNotNull(feePayers);
        System.out.println("Found " + feePayers.size() + " fee payers");
        for (FeePayer fp : feePayers.subList(0, Math.min(5, feePayers.size()))) {
            System.out.println("  FeePayer: ID=" + fp.getId());
        }
    }

    // =========================================================================
    // ExchangeService
    // =========================================================================

    @Test
    void listExchangeCounterparties() throws ApiException {
        List<ExchangeCounterparty> counterparties = client.getExchangeService()
                .getExchangeCounterparties();

        assertNotNull(counterparties);
        System.out.println("Found " + counterparties.size() + " exchange counterparties");
        for (ExchangeCounterparty cp : counterparties.subList(0, Math.min(5, counterparties.size()))) {
            System.out.println("  Exchange: " + cp.getName());
        }
    }

    // =========================================================================
    // FiatService
    // =========================================================================

    @Test
    void listFiatProviders() throws ApiException {
        List<FiatProvider> providers = client.getFiatService().getFiatProviders();

        assertNotNull(providers);
        System.out.println("Found " + providers.size() + " fiat providers");
        for (FiatProvider provider : providers.subList(0, Math.min(5, providers.size()))) {
            System.out.println("  Provider: " + provider.getProvider());
        }
    }

    // =========================================================================
    // JobService
    // =========================================================================

    @Test
    void listJobs() {
        try {
            List<Job> jobs = client.getJobService().getJobs();

            assertNotNull(jobs);
            System.out.println("Found " + jobs.size() + " jobs");
            for (Job job : jobs.subList(0, Math.min(5, jobs.size()))) {
                System.out.println("  Job: " + job.getName());
            }
        } catch (ApiException e) {
            // JobService may require 'tgvalidatord' role not available in test credentials
            System.out.println("Jobs not available: " + e.getMessage());
        }
    }

    // =========================================================================
    // TaurusNetwork - Participants
    // =========================================================================

    @Test
    void listTaurusNetworkParticipants() {
        try {
            List<Participant> participants = client.taurusNetwork().participants()
                    .list(null, null);

            assertNotNull(participants);
            System.out.println("Found " + participants.size() + " TaurusNetwork participants");
            for (Participant p : participants.subList(0, Math.min(5, participants.size()))) {
                System.out.println("  Participant: ID=" + p.getId() + ", Name=" + p.getName());
            }
        } catch (ApiException e) {
            // TaurusNetwork may not be enabled in all environments
            System.out.println("TaurusNetwork participants not available: " + e.getMessage());
        }
    }

    // =========================================================================
    // TaurusNetwork - Pledges
    // =========================================================================

    @Test
    void listTaurusNetworkPledges() {
        try {
            PledgeResult result = client.taurusNetwork().pledges()
                    .list(null, null, null, null, null, null);

            assertNotNull(result);
            System.out.println("Found " + result.getPledges().size() + " TaurusNetwork pledges");
        } catch (ApiException e) {
            // TaurusNetwork may not be enabled in all environments
            System.out.println("TaurusNetwork pledges not available: " + e.getMessage());
        }
    }

    // =========================================================================
    // TaurusNetwork - Lending
    // =========================================================================

    @Test
    void listTaurusNetworkLendingOffers() {
        try {
            LendingOfferResult result = client.taurusNetwork().lending()
                    .getLendingOffers(null, null, null, null, null);

            assertNotNull(result);
            System.out.println("Found " + result.getOffers().size() + " TaurusNetwork lending offers");
        } catch (ApiException e) {
            System.out.println("TaurusNetwork lending offers not available: " + e.getMessage());
        }
    }

    @Test
    void listTaurusNetworkLendingAgreements() {
        try {
            LendingAgreementResult result = client.taurusNetwork().lending()
                    .getLendingAgreements(null, null);

            assertNotNull(result);
            System.out.println("Found " + result.getAgreements().size() + " TaurusNetwork lending agreements");
        } catch (ApiException e) {
            System.out.println("TaurusNetwork lending agreements not available: " + e.getMessage());
        }
    }

    // =========================================================================
    // TaurusNetwork - Settlements
    // =========================================================================

    @Test
    void listTaurusNetworkSettlements() {
        try {
            SettlementResult result = client.taurusNetwork().settlements()
                    .getSettlements(null, null, null, null);

            assertNotNull(result);
            System.out.println("Found " + result.getSettlements().size() + " TaurusNetwork settlements");
        } catch (ApiException e) {
            System.out.println("TaurusNetwork settlements not available: " + e.getMessage());
        }
    }

    // =========================================================================
    // TaurusNetwork - Sharing
    // =========================================================================

    @Test
    void listTaurusNetworkSharedAddresses() {
        try {
            SharedAddressResult result = client.taurusNetwork().sharing()
                    .listSharedAddresses(null, null, null, null, null, null, null, null);

            assertNotNull(result);
            System.out.println("Found " + result.getSharedAddresses().size() + " TaurusNetwork shared addresses");
        } catch (ApiException e) {
            System.out.println("TaurusNetwork shared addresses not available: " + e.getMessage());
        }
    }

    // =========================================================================
    // ContractWhitelistingService
    // =========================================================================

    @Test
    void listWhitelistedContracts() throws ApiException {
        WhitelistedContractAddressResult result = client.getContractWhitelistingService()
                .getWhitelistedContracts(null, null, null, null, null, null);

        assertNotNull(result);
        System.out.println("Found " + result.getContracts().size() + " whitelisted contracts");
    }
}
