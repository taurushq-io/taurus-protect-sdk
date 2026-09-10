package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;

class TransactionServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private TransactionService transactionService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        transactionService = new TransactionService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new TransactionService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new TransactionService(apiClient, null));
    }

    /**
     * The chain-filter overloads exist because the API accepts blockchain and network on
     * both endpoints and Go/TS expose them, while this SDK hardcoded null. Without a
     * network stub (the project forbids Mockito) the reachable assertion is that the new
     * arity validates its arguments exactly like the 6-arg form, which also pins the
     * overload's existence and parameter order at compile time.
     */
    @Test
    void getTransactionsWithChainFilters_throwsOnNonPositiveLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactions(null, null, null, null, "ETH", "mainnet", 0, 0));
    }

    @Test
    void getTransactionsWithChainFilters_throwsOnNegativeOffset() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactions(null, null, null, null, "ETH", "mainnet", 50, -1));
    }

    @Test
    void exportTransactionsWithChainFilters_throwsOnNonPositiveLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.exportTransactions(null, null, null, null, "ETH", "mainnet", 0, 0));
    }

    @Test
    void exportTransactionsWithChainFilters_throwsOnNegativeOffset() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.exportTransactions(null, null, null, null, "ETH", "mainnet", 50, -1));
    }

    @Test
    void getTransactionById_throwsOnZeroId() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactionById(0));
    }

    @Test
    void getTransactionById_throwsOnNegativeId() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactionById(-1));
    }

    @Test
    void getTransactions_throwsOnZeroLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactions(null, null, null, null, 0, 0));
    }

    @Test
    void getTransactions_throwsOnNegativeOffset() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactions(null, null, null, null, 10, -1));
    }

    @Test
    void getTransactionsByAddress_throwsOnNullAddress() {
        assertThrows(NullPointerException.class, () ->
                transactionService.getTransactionsByAddress(null, 10, 0));
    }

    @Test
    void getTransactionsByAddress_throwsOnEmptyAddress() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactionsByAddress("", 10, 0));
    }

    @Test
    void getTransactionsByAddress_throwsOnZeroLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactionsByAddress("0x123", 0, 0));
    }

    @Test
    void getTransactionByHash_throwsOnNullHash() {
        assertThrows(NullPointerException.class, () ->
                transactionService.getTransactionByHash(null));
    }

    @Test
    void getTransactionByHash_throwsOnEmptyHash() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactionByHash(""));
    }
}
