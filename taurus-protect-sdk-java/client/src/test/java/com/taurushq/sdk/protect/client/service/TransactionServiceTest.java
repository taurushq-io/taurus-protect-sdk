package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.TransactionExportResult;
import com.taurushq.sdk.protect.client.model.TransactionResult;
import com.taurushq.sdk.protect.client.testutil.StubTransport;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

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
     * both endpoints and Go/TS expose them, while this SDK hardcoded null. Through the
     * transport stub both filters are seen on the wire.
     */
    @Test
    void getTransactionsWithChainFilters_sendsBothFilters() throws Exception {
        StubTransport stub = StubTransport.replying("{}");
        new TransactionService(stub.client(), apiExceptionMapper)
                .getTransactions(null, null, null, null, "ETH", "mainnet", 0, 0);
        assertEquals("ETH", stub.only().param("blockchain"));
        assertEquals("mainnet", stub.only().param("network"));
        assertEquals("20", stub.only().param("limit"));
    }

    @Test
    void getTransactionsWithChainFilters_throwsOnNegativeOffset() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactions(null, null, null, null, "ETH", "mainnet", 50, -1));
    }

    @Test
    void getTransactions_throwsAboveTheMaximumPageSize() {
        IllegalArgumentException e = assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactions(null, null, null, null, 101, 0));
        assertTrue(e.getMessage().contains("limit") && e.getMessage().contains("100"), e.getMessage());
    }

    @Test
    void getTransactions_zeroLimitSendsTheDefault() throws Exception {
        StubTransport stub = StubTransport.replying("{\"result\":[],\"totalItems\":\"45\"}");
        TransactionResult page = new TransactionService(stub.client(), apiExceptionMapper)
                .getTransactions(null, null, null, null, 0, 0);
        assertEquals("20", stub.only().param("limit"));
        assertNull(stub.only().param("offset"), "offset 0 is the first page and is not sent");
        assertEquals(20, page.getPagination().getLimit());
        assertEquals(45L, page.getPagination().getTotalItems());
    }

    @Test
    void getTransactions_plusRows() throws Exception {
        // offset + rows: two rows at offset 40 of 45 make the next offset 42, not 60.
        StubTransport stub = StubTransport.replying(
                "{\"result\":[{\"id\":\"1\"},{\"id\":\"2\"}],\"totalItems\":\"45\"}");
        TransactionResult page = new TransactionService(stub.client(), apiExceptionMapper)
                .getTransactions(null, null, null, null, 20, 40);
        assertEquals("40", stub.only().param("offset"));
        assertEquals(2, page.getTransactions().size());
        assertEquals(42L, page.getPagination().getNextOffset());
        assertTrue(page.getPagination().hasMore());
    }

    @Test
    void exportTransactions_isLimitOnly() throws Exception {
        // The server ignores an export offset, so none is ever sent; the limit has no SDK
        // maximum, and no format is sent unless the caller asks for one.
        StubTransport stub = StubTransport.replying("{\"result\":\"[]\",\"totalItems\":\"7\"}");
        TransactionExportResult export = new TransactionService(stub.client(), apiExceptionMapper)
                .exportTransactions(null, null, null, null, 5000);
        assertEquals("5000", stub.only().param("limit"));
        assertNull(stub.only().param("offset"));
        assertNull(stub.only().param("format"));
        assertEquals("[]", export.getContent());
        assertEquals(7L, export.getTotalItems());
    }

    @Test
    void exportTransactions_emptyReplyIsAnEmptyResult() throws Exception {
        // {} used to come back as a null String.
        StubTransport stub = StubTransport.replying("{}");
        TransactionExportResult export = new TransactionService(stub.client(), apiExceptionMapper)
                .exportTransactions(null, null, null, null, "ETH", "mainnet", "csv", 0);
        assertEquals("20", stub.only().param("limit"));
        assertEquals("csv", stub.only().param("format"));
        assertEquals("", export.getContent());
        assertEquals(0L, export.getTotalItems());
    }

    @Test
    void exportTransactions_throwsOnNegativeLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.exportTransactions(null, null, null, null, -1));
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
    void getTransactionsByAddress_throwsOnNegativeLimit() {
        assertThrows(IllegalArgumentException.class, () ->
                transactionService.getTransactionsByAddress("0x123", -1, 0));
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
