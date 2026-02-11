package com.taurushq.sdk.protect.client.integration;

import com.taurushq.sdk.protect.client.ProtectClient;
import com.taurushq.sdk.protect.client.testutil.TestHelper;
import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.Wallet;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInstance;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.*;

/**
 * Integration tests for AddressService.
 */
@TestInstance(TestInstance.Lifecycle.PER_CLASS)
class AddressIntegrationTest {

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
    void getAddress() throws ApiException {
        // Get first wallet to find an address
        List<Wallet> wallets = client.getWalletService().getWallets(1, 0);
        if (wallets.isEmpty()) {
            System.out.println("No wallets available for testing");
            return;
        }

        // Get addresses from the wallet
        List<Address> addresses = client.getAddressService().getAddresses(wallets.get(0).getId(), 1, 0);
        if (addresses.isEmpty()) {
            System.out.println("No addresses available for testing");
            return;
        }

        long addressId = addresses.get(0).getId();
        Address address = client.getAddressService().getAddress(addressId);

        System.out.println("Address details:");
        System.out.println("  ID: " + address.getId());
        System.out.println("  Address: " + address.getAddress());
        System.out.println("  Currency: " + address.getCurrency());
        System.out.println("  Available funds: " + address.getBalance().getAvailableConfirmed());

        assertNotNull(address);
        assertEquals(addressId, address.getId());
    }

    @Test
    void listAddresses() throws ApiException {
        // Get first wallet
        List<Wallet> wallets = client.getWalletService().getWallets(1, 0);
        if (wallets.isEmpty()) {
            System.out.println("No wallets available for testing");
            return;
        }

        long walletId = wallets.get(0).getId();
        List<Address> addresses = client.getAddressService().getAddresses(walletId, 10, 0);

        System.out.println("Found " + addresses.size() + " addresses for wallet " + walletId);
        for (Address a : addresses) {
            System.out.println("  " + a.getCurrency() + " " + a.getAddress());
        }

        assertNotNull(addresses);
    }

    @Test
    void paginateAddresses() throws ApiException {
        // Get first wallet
        List<Wallet> wallets = client.getWalletService().getWallets(1, 0);
        if (wallets.isEmpty()) {
            System.out.println("No wallets available for testing");
            return;
        }

        long walletId = wallets.get(0).getId();
        int limit = 10;
        int offset = 0;
        List<Address> allAddresses = new ArrayList<>();

        List<Address> addresses;
        do {
            addresses = client.getAddressService().getAddresses(walletId, limit, offset);
            allAddresses.addAll(addresses);
            offset += limit;
        } while (!addresses.isEmpty() && allAddresses.size() < 100);

        System.out.println("Found " + allAddresses.size() + " total addresses for wallet " + walletId);

        assertNotNull(allAddresses);
    }
}
