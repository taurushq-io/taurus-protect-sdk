package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class WalletServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private WalletService walletService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        walletService = new WalletService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new WalletService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new WalletService(apiClient, null));
    }

    @Test
    void createWallet_throwsOnNullRequest() {
        assertThrows(NullPointerException.class, () ->
                walletService.createWallet(null));
    }

    @Test
    void createWallet_throwsOnNullBlockchain() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.createWallet(null, "mainnet", "Test", false));
    }

    @Test
    void createWallet_throwsOnEmptyBlockchain() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.createWallet("", "mainnet", "Test", false));
    }

    @Test
    void createWallet_throwsOnNullNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.createWallet("ETH", null, "Test", false));
    }

    @Test
    void createWallet_throwsOnEmptyWalletName() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.createWallet("ETH", "mainnet", "", false));
    }

    @Test
    void getWallet_throwsOnZeroWalletId() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.getWallet(0));
    }

    @Test
    void getWallet_throwsOnNegativeWalletId() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.getWallet(-1));
    }

    @Test
    void getWallets_throwsOnZeroLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.getWallets(0, 0));
    }

    @Test
    void getWallets_throwsOnNegativeOffset() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.getWallets(10, -1));
    }

    @Test
    void getWalletsByName_throwsOnEmptyName() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.getWalletsByName("", 10, 0));
    }

    @Test
    void getWalletsByName_throwsOnZeroLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.getWalletsByName("wallet", 0, 0));
    }

    @Test
    void createWalletAttribute_throwsOnZeroWalletId() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.createWalletAttribute(0, "key", "value"));
    }

    @Test
    void getWalletBalanceHistory_throwsOnZeroWalletId() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.getWalletBalanceHistory(0, 24));
    }

    @Test
    void getWalletBalanceHistory_throwsOnZeroInterval() {
        assertThrows(IllegalArgumentException.class, () ->
                walletService.getWalletBalanceHistory(1, 0));
    }
}
