package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Currency;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for CurrencyService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class CurrencyIntegrationTest {

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

    @Test
    void listCurrencies() throws ApiException {
        List<Currency> currencies = client.getCurrencyService().getCurrencies();

        System.out.println("Found " + currencies.size() + " currencies");
        for (Currency c : currencies) {
            System.out.println("  " + c.getId() + ": " + c.getName() + " (" + c.getSymbol() + ")");
        }

        assertNotNull(currencies);
        assertFalse(currencies.isEmpty(), "Should have at least one currency");
    }

    @Test
    void getCurrencyById() throws ApiException {
        // First get a list to find a valid currency with blockchain and network
        List<Currency> currencies = client.getCurrencyService().getCurrencies();

        if (currencies.isEmpty()) {
            System.out.println("No currencies available for testing");
            return;
        }

        // Find a currency with blockchain and network set (required by API)
        Currency source = currencies.stream()
                .filter(c -> c.getBlockchain() != null && c.getNetwork() != null)
                .findFirst()
                .orElse(null);

        if (source == null) {
            System.out.println("No currency with blockchain and network available");
            return;
        }

        // Look up the currency using blockchain and network
        Currency currency = client.getCurrencyService()
                .getCurrencyByBlockchain(source.getBlockchain(), source.getNetwork());

        System.out.println("Currency details:");
        System.out.println("  ID: " + currency.getId());
        System.out.println("  Name: " + currency.getName());
        System.out.println("  Symbol: " + currency.getSymbol());
        System.out.println("  Blockchain: " + currency.getBlockchain());
        System.out.println("  Network: " + currency.getNetwork());
        System.out.println("  Decimals: " + currency.getDecimals());

        assertNotNull(currency);
        assertEquals(source.getBlockchain(), currency.getBlockchain());
        assertEquals(source.getNetwork(), currency.getNetwork());
    }

    @Test
    void getCurrencyByBlockchain() throws ApiException {
        // First get a list to find a valid currency with blockchain/network
        List<Currency> currencies = client.getCurrencyService().getCurrencies();

        if (currencies.isEmpty()) {
            System.out.println("No currencies available for testing");
            return;
        }

        // Find a currency with blockchain and network set
        Currency source = currencies.stream()
                .filter(c -> c.getBlockchain() != null && c.getNetwork() != null)
                .findFirst()
                .orElse(null);

        if (source == null) {
            System.out.println("No currency with blockchain and network available");
            return;
        }

        Currency currency = client.getCurrencyService()
                .getCurrencyByBlockchain(source.getBlockchain(), source.getNetwork());

        System.out.println("Currency by blockchain/network:");
        System.out.println("  Blockchain: " + currency.getBlockchain());
        System.out.println("  Network: " + currency.getNetwork());
        System.out.println("  Name: " + currency.getName());

        assertNotNull(currency);
        assertEquals(source.getBlockchain(), currency.getBlockchain());
        assertEquals(source.getNetwork(), currency.getNetwork());
    }

    @Test
    void getBaseCurrency() throws ApiException {
        Currency baseCurrency = client.getCurrencyService().getBaseCurrency();

        if (baseCurrency == null) {
            System.out.println("No base currency configured");
            return;
        }

        System.out.println("Base currency:");
        System.out.println("  ID: " + baseCurrency.getId());
        System.out.println("  Name: " + baseCurrency.getName());
        System.out.println("  Symbol: " + baseCurrency.getSymbol());

        assertNotNull(baseCurrency.getId());
    }
}
