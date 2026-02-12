package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for WhitelistedAddressService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class WhitelistedAddressIntegrationTest {

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
    void getWhitelistedAddress() throws ApiException, WhitelistException {
        // First get a list to find a valid ID
        List<SignedWhitelistedAddressEnvelope> addresses = client.getWhitelistedAddressService()
                .getWhitelistedAddresses(1, 0);

        if (addresses.isEmpty()) {
            System.out.println("No whitelisted addresses available for testing");
            return;
        }

        long id = addresses.get(0).getId();
        WhitelistedAddress wla = client.getWhitelistedAddressService().getWhitelistedAddress(id);

        System.out.println("Whitelisted address:");
        System.out.println("  Blockchain: " + wla.getBlockchain());
        System.out.println("  Network: " + wla.getNetwork());
        System.out.println("  Address: " + wla.getAddress());

        assertNotNull(wla);
    }

    @Test
    void listWhitelistedAddresses() throws ApiException, WhitelistException {
        List<SignedWhitelistedAddressEnvelope> addresses = client.getWhitelistedAddressService()
                .getWhitelistedAddresses(10, 0);

        System.out.println("Found " + addresses.size() + " whitelisted addresses");
        for (SignedWhitelistedAddressEnvelope envelope : addresses) {
            WhitelistedAddress wla = envelope.getWhitelistedAddress();
            System.out.println("  " + wla.getBlockchain() + "/" + wla.getNetwork() + ": " + wla.getAddress());
        }

        assertNotNull(addresses);
    }

    @Test
    void listWhitelistedAddressesByBlockchain() throws ApiException, WhitelistException {
        List<SignedWhitelistedAddressEnvelope> ethAddresses = client.getWhitelistedAddressService()
                .getWhitelistedAddresses(10, 0, "ETH");

        System.out.println("Found " + ethAddresses.size() + " ETH whitelisted addresses");
        for (SignedWhitelistedAddressEnvelope envelope : ethAddresses) {
            WhitelistedAddress wla = envelope.getWhitelistedAddress();
            System.out.println("  Address: " + wla.getAddress());
            System.out.println("  Blockchain: " + wla.getBlockchain());
        }

        assertNotNull(ethAddresses);
    }

    @Test
    void listWhitelistedAddressesByBlockchainAndNetwork() throws ApiException, WhitelistException {
        List<SignedWhitelistedAddressEnvelope> mainnetAddresses = client.getWhitelistedAddressService()
                .getWhitelistedAddresses(10, 0, "ETH", "mainnet");

        System.out.println("Found " + mainnetAddresses.size() + " ETH mainnet whitelisted addresses");

        assertNotNull(mainnetAddresses);
    }

    @Test
    void paginateAllWhitelistedAddresses() throws ApiException, WhitelistException {
        int limit = 50;
        int offset = 0;
        int totalCount = 0;

        List<SignedWhitelistedAddressEnvelope> addresses;
        do {
            addresses = client.getWhitelistedAddressService()
                    .getWhitelistedAddresses(limit, offset);

            for (SignedWhitelistedAddressEnvelope envelope : addresses) {
                WhitelistedAddress wla = envelope.getWhitelistedAddress();
                System.out.printf("Blockchain: %s, Network: %s, Address: %s%n",
                        wla.getBlockchain(),
                        wla.getNetwork(),
                        wla.getAddress());
                totalCount++;
            }

            offset += addresses.size();

            // Safety limit for tests
            if (offset > 2000) {
                System.out.println("Stopping pagination test at 2000 items");
                break;
            }
        } while (addresses.size() == limit);

        System.out.println("Total whitelisted addresses scanned: " + totalCount);
        assertTrue(totalCount >= 0);
    }
}
