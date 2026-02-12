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
